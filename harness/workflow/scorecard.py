"""Multi-dimensional scorecard for workflow quality.

Evaluates a project across quality dimensions that predict
open-source excellence and contributor attraction.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from pathlib import Path


@dataclass(frozen=True)
class Dimension:
    """A single quality dimension with target and measurement method."""
    key: str
    label: str
    description: str
    target: str                     # human-readable target
    check: str                      # how to measure ("file_exists", "section_present", "cli_output", "coverage_pct")
    params: dict = field(default_factory=dict)


@dataclass
class Score:
    dimension: Dimension
    value: float                    # 0.0 – 1.0
    evidence: str                   # what we found
    passed: bool

    def render(self) -> str:
        icon = "✅" if self.passed else "❌"
        return f"  {icon} {self.dimension.label}: {self.value:.0%}  ({self.evidence})"


@dataclass
class Scorecard:
    dimensions: list[Dimension]
    scores: list[Score]

    @property
    def overall(self) -> float:
        if not self.scores:
            return 0.0
        return sum(s.value for s in self.scores) / len(self.scores)

    @property
    def passed(self) -> bool:
        return all(s.passed for s in self.scores)

    def render(self) -> str:
        lines = [
            "=" * 60,
            f"  WORKFLOW QUALITY SCORECARD  —  {self.overall:.0%}  {'✅' if self.passed else '❌'}",
            "=" * 60,
            "",
        ]
        for s in self.scores:
            lines.append(s.render())
        lines.append("")
        lines.append(f"  Overall: {self.overall:.0%}  ({'PASS' if self.passed else 'FAIL'})")
        return "\n".join(lines)


# ---------------------------------------------------------------------------
# Built-in dimensions for open-source workflow quality
# ---------------------------------------------------------------------------

def oss_quality_dimensions() -> list[Dimension]:
    """The 6 dimensions that predict open-source project excellence."""
    return [
        Dimension(
            key="design_artifacts",
            label="Design Artifacts",
            description="Proposal, design doc, requirements, and tasks are present and complete",
            target="All 4 artifacts present with required sections",
            check="artifact_manifest",
            params={"manifest": "workflow-artifacts.json"},
        ),
        Dimension(
            key="contributor_onboarding",
            label="Contributor Onboarding",
            description="ARCHITECTURE.md and CONTRIBUTING.md help new contributors start in < 30 min",
            target="Both files present, ARCHITECTURE.md ≥ 500 words, CONTRIBUTING.md has PR process",
            check="contributor_docs",
            params={
                "files": ["ARCHITECTURE.md", "CONTRIBUTING.md"],
                "architecture_min_words": 500,
                "contributing_required_sections": ["## PR", "## Development", "## Testing"],
            },
        ),
        Dimension(
            key="code_quality_gates",
            label="Code Quality Gates",
            description="Lint, build, test, coverage all pass automatically",
            target="Lint clean, build passes, test ≥ 80% coverage, no race conditions",
            check="quality_gate_output",
            params={
                "required_checks": ["lint", "build", "test", "coverage", "race"],
            },
        ),
        Dimension(
            key="process_compliance",
            label="Process Compliance",
            description="Commits follow conventional commits, PRs have complete descriptions",
            target="100% conventional commits, every PR links to proposal",
            check="git_history",
            params={
                "conventional_commit_pattern": r"^(feat|fix|refactor|docs|test|chore|perf|ci)(\(.+\))?:.+",
            },
        ),
        Dimension(
            key="review_thoroughness",
            label="Review Thoroughness",
            description="Every change reviewed for security, performance, and correctness",
            target="Security + performance + code review on every merge",
            check="review_evidence",
            params={
                "required_dimensions": ["security", "performance", "code-quality"],
            },
        ),
        Dimension(
            key="design_traceability",
            label="Design Traceability",
            description="Every implementation decision traces back to a design artifact",
            target="tasks.md → design.md → proposal.md chain is unbroken",
            check="traceability_chain",
            params={
                "required_links": ["tasks→design", "design→proposal"],
            },
        ),
    ]


def render_scorecard(card: Scorecard) -> str:
    return card.render()
