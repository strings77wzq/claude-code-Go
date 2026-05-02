"""Normalized comparison evidence for local agent trials."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class EvidenceRecord:
    agent: str
    scenario: str
    outcome: str
    duration_ms: int
    tool_count: int
    source: str
    notes: str = ""


class ComparisonReport:
    def __init__(self, records: list[EvidenceRecord]):
        self.records = records

    def render(self) -> str:
        lines = ["Agent comparison evidence"]
        for record in self.records:
            lines.append(
                f"- agent={record.agent} scenario={record.scenario} outcome={record.outcome} "
                f"duration_ms={record.duration_ms} tool_count={record.tool_count} source={record.source}"
            )
            if record.notes:
                lines.append(f"  notes={record.notes}")
        return "\n".join(lines)
