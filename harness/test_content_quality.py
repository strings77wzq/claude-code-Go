"""Tests for content-level quality checks."""

from pathlib import Path

from harness.workflow.content_check import ContentQualityChecker, ContentCheckResult, DocIssue


class TestContentQualityChecker:
    def test_passes_when_all_refs_valid(self, tmp_path: Path):
        (tmp_path / "internal").mkdir(parents=True, exist_ok=True)
        (tmp_path / "internal" / "agent").mkdir(parents=True, exist_ok=True)
        (tmp_path / "internal" / "agent" / "loop.go").write_text("package agent")
        (tmp_path / "CONTRIBUTING.md").write_text(
            "# Contributing\n\n"
            "The core loop is in `internal/agent/loop.go`.\n\n"
            "```bash\n"
            "go test ./...\n"
            "go build ./...\n"
            "go vet ./...\n"
            "gofmt -l .\n"
            "```\n"
        )
        (tmp_path / "CLAUDE.md").write_text("# Project rules")
        (tmp_path / "ARCHITECTURE.md").write_text(
            "# Architecture\n\n"
            "Main entry at `cmd/go-code/`.\n"
            "See `CLAUDE.md` for project rules.\n"
        )

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        assert result.passed
        assert len(result.errors) == 0

    def test_detects_missing_file_reference(self, tmp_path: Path):
        (tmp_path / "ARCHITECTURE.md").write_text(
            "# Architecture\n\nSee `internal/missing/file.go` for details.\n"
        )

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        assert not result.passed
        assert any("missing/file.go" in e.message for e in result.errors)

    def test_ignores_urls_and_stdlib(self, tmp_path: Path):
        (tmp_path / "ARCHITECTURE.md").write_text(
            "# Architecture\n\n"
            "See https://github.com/example/repo for more.\n"
            "Requires `github.com/foo/bar` module.\n"
            "Uses Go `v1.24.0`.\n"
        )
        (tmp_path / "CONTRIBUTING.md").write_text(
            "# Contributing\n\n"
            "```bash\n"
            "go test ./...\n"
            "go build ./...\n"
            "go vet ./...\n"
            "gofmt -l .\n"
            "```\n"
        )

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        # URLs and stdlib paths should not trigger errors
        assert result.passed

    def test_detects_missing_required_commands(self, tmp_path: Path):
        (tmp_path / "CONTRIBUTING.md").write_text(
            "# Contributing\n\n## Development Setup\n\nNo commands here\n"
        )

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        assert not result.passed
        assert any("go test" in e.message for e in result.errors)

    def test_warns_on_smart_quotes(self, tmp_path: Path):
        (tmp_path / "CONTRIBUTING.md").write_text(
            "# Contributing\n\n"
            "```bash\n"
            "echo ‘hello’\n"
            "go test ./...\n"
            "```\n"
        )

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        assert any("Smart quotes" in i.message for i in result.issues)

    def test_checks_both_architecture_and_contributing(self, tmp_path: Path):
        (tmp_path / "CONTRIBUTING.md").write_text("minimal")
        (tmp_path / "ARCHITECTURE.md").write_text("minimal")

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        assert "CONTRIBUTING.md" in result.checked_files
        assert "ARCHITECTURE.md" in result.checked_files

    def test_result_render_output(self, tmp_path: Path):
        (tmp_path / "ARCHITECTURE.md").write_text("# OK\n")
        (tmp_path / "CONTRIBUTING.md").write_text("# OK\n")

        checker = ContentQualityChecker()
        result = checker.check(tmp_path)

        output = result.render()
        assert "CONTENT QUALITY CHECK" in output
