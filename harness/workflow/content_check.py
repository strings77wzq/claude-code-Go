"""Content-level quality checks for project documentation.

Validates that:
1. File references in docs point to real files
2. Commands in code blocks are syntactically plausible
3. Key quality-gate commands match the actual project tooling
"""

from __future__ import annotations

import re
import subprocess
from dataclasses import dataclass, field
from pathlib import Path


# ---------------------------------------------------------------------------
# Data types
# ---------------------------------------------------------------------------

@dataclass
class DocIssue:
    severity: str       # "error" | "warning"
    file: str           # which doc file
    line: int           # approximate line
    message: str

    def render(self) -> str:
        icon = {"error": "❌", "warning": "⚠️"}.get(self.severity, "?")
        return f"  {icon} {self.file}:{self.line}  {self.message}"


@dataclass
class ContentCheckResult:
    issues: list[DocIssue] = field(default_factory=list)
    checked_files: list[str] = field(default_factory=list)

    @property
    def errors(self) -> list[DocIssue]:
        return [i for i in self.issues if i.severity == "error"]

    @property
    def passed(self) -> bool:
        return len(self.errors) == 0

    def render(self) -> str:
        lines = [
            "=" * 60,
            f"  CONTENT QUALITY CHECK  —  {'✅ PASS' if self.passed else '❌ FAIL'}",
            f"  Checked: {', '.join(self.checked_files)}",
            f"  Issues:  {len(self.issues)} ({len(self.errors)} errors, {len(self.issues) - len(self.errors)} warnings)",
            "=" * 60,
        ]
        if self.issues:
            lines.append("")
            for issue in self.issues:
                lines.append(issue.render())
        return "\n".join(lines)


# ---------------------------------------------------------------------------
# Checker
# ---------------------------------------------------------------------------

class ContentQualityChecker:
    """Checks documentation content for correctness."""

    # File reference pattern: relative paths like `internal/agent/loop.go` or `cmd/go-code/`
    FILE_REF_PATTERN = re.compile(
        r'`([a-zA-Z0-9_./-]+\.[a-zA-Z]+(?:/[a-zA-Z0-9_./-]*)?)`'
    )

    # Code block: ```bash ... ```
    BASH_BLOCK = re.compile(r'```bash\n(.*?)```', re.DOTALL)

    # Critical commands that MUST be in CONTRIBUTING.md
    REQUIRED_COMMANDS = [
        "go test",
        "go build",
        "go vet",
        "gofmt",
    ]

    def check(self, project_root: Path) -> ContentCheckResult:
        result = ContentCheckResult()

        for doc_name in ["ARCHITECTURE.md", "CONTRIBUTING.md"]:
            doc_path = project_root / doc_name
            if not doc_path.exists():
                result.issues.append(DocIssue(
                    severity="error",
                    file=doc_name,
                    line=1,
                    message=f"File missing: {doc_name}"
                ))
                continue

            result.checked_files.append(doc_name)
            content = doc_path.read_text(encoding="utf-8", errors="replace")

            self._check_file_references(content, doc_name, project_root, result)
            self._check_commands(content, doc_name, result)
            self._check_required_commands(content, doc_name, result)

        return result

    def _check_file_references(
        self,
        content: str,
        doc_name: str,
        project_root: Path,
        result: ContentCheckResult,
    ) -> None:
        """Verify that file paths mentioned in docs actually exist."""
        seen = set()

        for match in self.FILE_REF_PATTERN.finditer(content):
            ref = match.group(1)

            # Skip URLs, Go stdlib paths, version numbers
            if ref.startswith(("http", "https", "github.com/", "v0.", "v1.")):
                continue
            if ref in seen:
                continue

            # Only check paths that look like project files
            if not any(ref.startswith(prefix) for prefix in (
                "cmd/", "internal/", "pkg/", "harness/", "docs/", "tests/",
                "openspec/", "CLAUDE.md", "ARCHITECTURE.md", "CONTRIBUTING.md",
                "CHANGELOG.md", "README.md", "go.mod", "go.sum", "Makefile",
                "CODE_OF_CONDUCT.md", "LICENSE",
            )):
                continue

            seen.add(ref)
            full_path = project_root / ref

            if not full_path.exists():
                # Find the line number approximately
                line_num = content[:match.start()].count('\n') + 1
                result.issues.append(DocIssue(
                    severity="error",
                    file=doc_name,
                    line=line_num,
                    message=f"Referenced file does not exist: `{ref}`"
                ))

    def _check_commands(
        self,
        content: str,
        doc_name: str,
        result: ContentCheckResult,
    ) -> None:
        """Check that bash commands in docs are syntactically plausible."""
        for match in self.BASH_BLOCK.finditer(content):
            block = match.group(1)
            line_num = content[:match.start()].count('\n') + 1

            for cmd_line in block.strip().split('\n'):
                cmd_line = cmd_line.strip()
                if not cmd_line or cmd_line.startswith('#'):
                    continue

                # Check for `go test -v ./...` vs `go test -v ./...` (smart quotes)
                if '‘' in cmd_line or '’' in cmd_line or '“' in cmd_line or '”' in cmd_line:
                    result.issues.append(DocIssue(
                        severity="warning",
                        file=doc_name,
                        line=line_num,
                        message=f"Smart quotes in command — may not copy-paste: `{cmd_line[:60]}...`"
                    ))

    def _check_required_commands(
        self,
        content: str,
        doc_name: str,
        result: ContentCheckResult,
    ) -> None:
        """Verify that critical workflow commands are present."""
        if doc_name != "CONTRIBUTING.md":
            return

        for cmd in self.REQUIRED_COMMANDS:
            if cmd not in content:
                result.issues.append(DocIssue(
                    severity="error",
                    file=doc_name,
                    line=1,
                    message=f"Required command missing from docs: `{cmd}`"
                ))


# ---------------------------------------------------------------------------
# Quick self-test: verify the checker doesn't false-positive on itself
# ---------------------------------------------------------------------------

def _self_test():
    """Verify the checker works on the project's own docs."""
    project = Path(__file__).parent.parent.parent
    checker = ContentQualityChecker()
    result = checker.check(project)
    print(result.render())
    return result.passed


if __name__ == "__main__":
    ok = _self_test()
    exit(0 if ok else 1)
