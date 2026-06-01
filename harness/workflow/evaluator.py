"""Artifact-level evaluation for workflow quality.

Tests whether a workflow run produces complete, well-structured artifacts
that serve both the developer and future contributors.
"""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


# ---------------------------------------------------------------------------
# Data types
# ---------------------------------------------------------------------------

@dataclass(frozen=True)
class ArtifactExpectation:
    """What we require from a single workflow artifact file."""

    path: str                # relative path from project root
    required: bool = True
    min_bytes: int = 0       # file must be at least this many bytes
    required_sections: list[str] = field(default_factory=list)
    forbidden_patterns: list[str] = field(default_factory=list)  # regex


@dataclass(frozen=True)
class WorkflowManifest:
    schema_version: str
    name: str
    description: str
    artifacts: list[ArtifactExpectation]
    quality_indicators: list[str] = field(default_factory=list)


@dataclass
class ArtifactResult:
    expectation: ArtifactExpectation
    found: bool
    size_bytes: int
    missing_sections: list[str] = field(default_factory=list)
    forbidden_found: list[str] = field(default_factory=list)

    @property
    def passed(self) -> bool:
        if self.expectation.required and not self.found:
            return False
        if not self.found:
            return True  # optional, not present = OK
        if self.size_bytes < self.expectation.min_bytes:
            return False
        if self.missing_sections:
            return False
        if self.forbidden_found:
            return False
        return True


@dataclass
class WorkflowResult:
    manifest: WorkflowManifest
    artifacts: list[ArtifactResult]
    project_root: Path

    @property
    def passed(self) -> bool:
        return all(a.passed for a in self.artifacts)

    @property
    def score(self) -> float:
        """0.0 – 1.0 where 1.0 = all required artifacts fully complete."""
        required = [a for a in self.artifacts if a.expectation.required]
        if not required:
            return 1.0
        return sum(1.0 for a in required if a.passed) / len(required)

    def render(self) -> str:
        lines = [
            f"Workflow: {self.manifest.name}",
            f"  Score: {self.score:.0%}  ({'PASS' if self.passed else 'FAIL'})",
            f"  Root:  {self.project_root}",
            "",
        ]
        for a in self.artifacts:
            status = "✅" if a.passed else "❌"
            lines.append(f"  {status} {a.expectation.path}")
            if not a.found:
                lines.append(f"       MISSING (required={a.expectation.required})")
            else:
                lines.append(f"       {a.size_bytes} bytes")
                if a.missing_sections:
                    lines.append(f"       missing sections: {a.missing_sections}")
                if a.forbidden_found:
                    lines.append(f"       forbidden patterns: {a.forbidden_found}")
        return "\n".join(lines)


# ---------------------------------------------------------------------------
# Manifest loading
# ---------------------------------------------------------------------------

def load_workflow_manifest(path: Path) -> WorkflowManifest:
    data = json.loads(path.read_text())
    _validate_fields(data, ["schema_version", "name", "description", "artifacts"])

    if data["schema_version"] != "workflow-quality.v1":
        raise ValueError(f"unsupported schema_version: {data['schema_version']}")

    artifacts = []
    for entry in data["artifacts"]:
        artifacts.append(
            ArtifactExpectation(
                path=entry["path"],
                required=entry.get("required", True),
                min_bytes=entry.get("min_bytes", 0),
                required_sections=entry.get("required_sections", []),
                forbidden_patterns=entry.get("forbidden_patterns", []),
            )
        )

    return WorkflowManifest(
        schema_version=data["schema_version"],
        name=data["name"],
        description=data["description"],
        artifacts=artifacts,
        quality_indicators=data.get("quality_indicators", []),
    )


# ---------------------------------------------------------------------------
# Evaluation
# ---------------------------------------------------------------------------

class ArtifactEvaluator:
    """Evaluates a project directory against a WorkflowManifest."""

    def __init__(self, manifest: WorkflowManifest):
        self.manifest = manifest

    def evaluate(self, project_root: Path) -> WorkflowResult:
        results: list[ArtifactResult] = []
        for exp in self.manifest.artifacts:
            results.append(self._check_artifact(project_root, exp))
        return WorkflowResult(
            manifest=self.manifest,
            artifacts=results,
            project_root=project_root,
        )

    def _check_artifact(self, root: Path, exp: ArtifactExpectation) -> ArtifactResult:
        file_path = root / exp.path
        if not file_path.exists():
            return ArtifactResult(
                expectation=exp,
                found=False,
                size_bytes=0,
                missing_sections=list(exp.required_sections),
            )

        content = file_path.read_text(encoding="utf-8", errors="replace")
        size = len(content.encode("utf-8"))

        missing_sections = [
            s for s in exp.required_sections
            if s not in content
        ]

        import re
        forbidden_found = [
            p for p in exp.forbidden_patterns
            if re.search(p, content)
        ]

        return ArtifactResult(
            expectation=exp,
            found=True,
            size_bytes=size,
            missing_sections=missing_sections,
            forbidden_found=forbidden_found,
        )


def _validate_fields(data: dict[str, Any], fields: list[str]) -> None:
    missing = [f for f in fields if f not in data]
    if missing:
        raise ValueError(f"missing required fields: {missing}")
