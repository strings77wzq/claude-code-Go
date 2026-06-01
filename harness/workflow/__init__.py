"""Workflow quality evaluation framework.

Evaluates whether a development workflow produces complete, high-quality artifacts
that enable open-source collaboration and codebase understanding.

Layers:
    1. evaluator  — structural artifact checks (files exist, sections present)
    2. scorecard  — 6‑dimension quality scoring for open‑source excellence
    3. content_check — content‑level validation (file refs, commands, copy‑paste safety)
"""

from harness.workflow.evaluator import (
    ArtifactExpectation,
    ArtifactEvaluator,
    WorkflowManifest,
    WorkflowResult,
    load_workflow_manifest,
)
from harness.workflow.scorecard import (
    Scorecard,
    Dimension,
    Score,
    render_scorecard,
)
from harness.workflow.content_check import (
    ContentQualityChecker,
    ContentCheckResult,
    DocIssue,
)

__all__ = [
    "ArtifactExpectation",
    "ArtifactEvaluator",
    "WorkflowManifest",
    "WorkflowResult",
    "load_workflow_manifest",
    "Scorecard",
    "Dimension",
    "Score",
    "render_scorecard",
    "ContentQualityChecker",
    "ContentCheckResult",
    "DocIssue",
]
