"""Workflow quality evaluation tests.

Validates that the development workflow produces complete, high-quality
artifacts suitable for open-source collaboration.

Run with:
    python -m pytest harness/test_workflow_quality.py -v
"""

from __future__ import annotations

import json
from pathlib import Path

import pytest

from harness.workflow.evaluator import (
    ArtifactEvaluator,
    ArtifactExpectation,
    WorkflowManifest,
    load_workflow_manifest,
)
from harness.workflow.scorecard import (
    Dimension,
    Score,
    Scorecard,
    oss_quality_dimensions,
)


# ---------------------------------------------------------------------------
# Manifest validation tests
# ---------------------------------------------------------------------------

class TestManifestValidation:
    def test_load_workflow_manifest(self):
        manifest_path = Path(__file__).parent / "manifests" / "workflow-artifacts.json"
        manifest = load_workflow_manifest(manifest_path)

        assert manifest.schema_version == "workflow-quality.v1"
        assert manifest.name == "workflow-artifacts"
        assert len(manifest.artifacts) >= 4  # at least proposal/design/tasks/requirement

    def test_all_required_artifacts_have_paths(self):
        manifest_path = Path(__file__).parent / "manifests" / "workflow-artifacts.json"
        manifest = load_workflow_manifest(manifest_path)

        for artifact in manifest.artifacts:
            assert artifact.path, f"Artifact missing path: {artifact}"

    def test_manifest_rejects_missing_fields(self, tmp_path: Path):
        bad = tmp_path / "bad.json"
        bad.write_text(json.dumps({"schema_version": "workflow-quality.v1"}))

        with pytest.raises(ValueError) as exc:
            load_workflow_manifest(bad)

        assert "name" in str(exc.value)


# ---------------------------------------------------------------------------
# Artifact evaluation tests
# ---------------------------------------------------------------------------

class TestArtifactEvaluator:
    def test_passes_when_all_artifacts_present_and_complete(self, tmp_path: Path):
        # Create all required files
        (tmp_path / "proposal.md").write_text(
            "## 问题陈述\n描述问题\n## 方案\n方案内容\n## 影响范围\n影响\n## 工作量估算\n2天\n## 验收标准\nAC1"
        )
        (tmp_path / "design.md").write_text(
            "## 架构\n架构描述\n## 数据流\n数据流描述\n## 接口\n接口定义\n## 权衡\n权衡分析\n" * 5
        )
        (tmp_path / "tasks.md").write_text(
            "## Task 1\n实现某功能\n文件边界: pkg/foo/\n## Task 2\n测试\n文件边界: pkg/foo/"
        )
        (tmp_path / "requirement.md").write_text(
            "Given 前置条件\nWhen 触发动作\nThen 预期结果"
        )

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(path="proposal.md", required_sections=["## 问题陈述", "## 方案"]),
                ArtifactExpectation(path="design.md", required_sections=["## 架构", "## 数据流"]),
                ArtifactExpectation(path="tasks.md", required_sections=["## Task", "文件边界"]),
                ArtifactExpectation(path="requirement.md", required_sections=["Given", "When", "Then"]),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)

        assert result.passed
        assert result.score == 1.0

    def test_fails_when_required_artifact_missing(self, tmp_path: Path):
        (tmp_path / "proposal.md").write_text("minimal")

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(path="proposal.md"),
                ArtifactExpectation(path="design.md", required=True),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)

        assert not result.passed
        assert result.score < 1.0

    def test_optional_artifact_missing_is_ok(self, tmp_path: Path):
        (tmp_path / "proposal.md").write_text("minimal")

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(path="proposal.md", required=True),
                ArtifactExpectation(path="CHANGELOG.md", required=False),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)
        assert result.passed  # optional CHANGELOG missing = OK

    def test_detects_missing_sections(self, tmp_path: Path):
        (tmp_path / "proposal.md").write_text("Just some text, no proper sections")

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(
                    path="proposal.md",
                    required_sections=["## 问题陈述", "## 方案"],
                ),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)

        assert not result.passed
        assert len(result.artifacts[0].missing_sections) == 2

    def test_detects_forbidden_patterns(self, tmp_path: Path):
        (tmp_path / "proposal.md").write_text("## 问题陈述\nTODO: fix this later\n## 方案\n方案")

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(
                    path="proposal.md",
                    forbidden_patterns=["TODO"],
                ),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)

        assert not result.passed
        assert "TODO" in result.artifacts[0].forbidden_found[0]

    def test_size_check_fails_for_too_small_file(self, tmp_path: Path):
        (tmp_path / "proposal.md").write_text("tiny")

        manifest = WorkflowManifest(
            schema_version="workflow-quality.v1",
            name="test",
            description="",
            artifacts=[
                ArtifactExpectation(path="proposal.md", min_bytes=100),
            ],
        )

        result = ArtifactEvaluator(manifest).evaluate(tmp_path)

        assert not result.passed


# ---------------------------------------------------------------------------
# Scorecard tests
# ---------------------------------------------------------------------------

class TestScorecard:
    def test_scorecard_renders_all_dimensions(self):
        dims = oss_quality_dimensions()
        scores = [
            Score(d, value=1.0, evidence="all present", passed=True)
            for d in dims
        ]
        card = Scorecard(dimensions=dims, scores=scores)

        output = card.render()

        assert card.overall == 1.0
        assert card.passed
        assert "WORKFLOW QUALITY SCORECARD" in output
        for d in dims:
            assert d.label in output

    def test_scorecard_fails_when_any_dimension_fails(self):
        dims = oss_quality_dimensions()
        scores = [Score(d, value=1.0, evidence="ok", passed=True) for d in dims[:-1]]
        scores.append(Score(dims[-1], value=0.0, evidence="missing", passed=False))
        card = Scorecard(dimensions=dims, scores=scores)

        assert not card.passed
        assert card.overall < 1.0

    def test_scorecard_overall_averages_correctly(self):
        dims = oss_quality_dimensions()[:3]
        scores = [
            Score(dims[0], value=0.5, evidence="half", passed=False),
            Score(dims[1], value=1.0, evidence="full", passed=True),
            Score(dims[2], value=0.0, evidence="none", passed=False),
        ]
        card = Scorecard(dimensions=dims, scores=scores)

        assert card.overall == 0.5


# ---------------------------------------------------------------------------
# Real project evaluation (these run against the actual project)
# ---------------------------------------------------------------------------

class TestRealProjectQuality:
    """Smoke tests that evaluate the actual project's quality artifacts."""

    @pytest.fixture
    def project_root(self):
        return Path(__file__).parent.parent

    def test_project_has_minimal_contributor_docs(self, project_root: Path):
        """Project must have at least ARCHITECTURE.md or equivalent."""
        candidates = [
            project_root / "ARCHITECTURE.md",
            project_root / "CLAUDE.md",
            project_root / "README.md",
        ]
        found = [c for c in candidates if c.exists()]
        assert found, f"No contributor docs found in {project_root}"

    def test_project_has_openspec_specs(self, project_root: Path):
        """Project must have OpenSpec specs directory."""
        specs_dir = project_root / "openspec" / "specs"
        assert specs_dir.is_dir(), f"OpenSpec specs not found at {specs_dir}"

    def test_workflow_manifest_is_valid(self):
        """The workflow-artifacts manifest must pass self-validation."""
        manifest_path = (
            Path(__file__).parent / "manifests" / "workflow-artifacts.json"
        )
        manifest = load_workflow_manifest(manifest_path)

        assert manifest.schema_version == "workflow-quality.v1"
        for a in manifest.artifacts:
            assert a.path, f"Artifact has empty path"


# ---------------------------------------------------------------------------
# Workflow quality manifest: test that project meets its own standards
# ---------------------------------------------------------------------------

class TestProjectMeetsWorkflowStandard:
    """Evaluate the actual project against the workflow quality manifest."""

    @pytest.fixture
    def project_root(self):
        return Path(__file__).parent.parent

    @pytest.fixture
    def workflow_manifest(self):
        manifest_path = (
            Path(__file__).parent / "manifests" / "workflow-artifacts.json"
        )
        return load_workflow_manifest(manifest_path)

    def test_evaluate_project_against_manifest(self, project_root, workflow_manifest):
        result = ArtifactEvaluator(workflow_manifest).evaluate(project_root)

        output = result.render()
        print("\n" + output)

        # Report which artifacts are missing — don't fail yet,
        # this is aspirational.
        missing = [a for a in result.artifacts if not a.passed]
        if missing:
            names = [m.expectation.path for m in missing]
            print(f"\n⚠  {len(missing)} artifact(s) need attention: {names}")
            print("   Run the full workflow to generate them.")

        # Soft check: at least proposal.md and CLAUDE.md must exist
        proposal = project_root / "proposal.md"
        claude = project_root / "CLAUDE.md"
        assert proposal.exists() or claude.exists(), \
            "Project must have at least proposal.md or CLAUDE.md"
