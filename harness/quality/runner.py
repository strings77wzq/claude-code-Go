"""Scenario result evaluation for manifest-driven quality gates."""

from __future__ import annotations

import json
import re
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

from harness.quality.manifest import ScenarioManifest

REDACTED = "[REDACTED]"


@dataclass
class ScenarioResult:
    name: str
    passed: bool
    duration_ms: int
    tool_count: int
    permission_decisions: list[str] = field(default_factory=list)
    trace_path: str | None = None
    failure_reason: str = ""
    latency_budget_violated: bool = False
    stdout: str = ""
    stderr: str = ""

    def to_dict(self) -> dict[str, Any]:
        return _redact(
            {
                "name": self.name,
                "passed": self.passed,
                "duration_ms": self.duration_ms,
                "tool_count": self.tool_count,
                "permission_decisions": self.permission_decisions,
                "trace_path": self.trace_path,
                "failure_reason": self.failure_reason,
                "latency_budget_violated": self.latency_budget_violated,
                "stdout": self.stdout,
                "stderr": self.stderr,
            }
        )


class ScenarioRunner:
    def __init__(self, manifest: ScenarioManifest):
        self.manifest = manifest

    def evaluate(
        self,
        *,
        stdout: str,
        stderr: str,
        returncode: int,
        duration_ms: int,
        trace_path: Path | None,
    ) -> ScenarioResult:
        failures: list[str] = []
        trace_events = _load_trace(trace_path)
        event_types = [event.get("type") for event in trace_events]
        tool_count = sum(1 for event in trace_events if event.get("type") == "tool")
        permission_decisions = [
            str(event.get("decision"))
            for event in trace_events
            if event.get("type") == "permission" and event.get("decision") is not None
        ]

        if returncode != 0:
            failures.append(f"returncode={returncode}")
        for expected in self.manifest.assertions.stdout_contains:
            if expected not in stdout:
                failures.append(f"stdout missing {expected!r}")
        for event_type in self.manifest.trace.required_events:
            if event_type not in event_types:
                failures.append(f"trace missing event {event_type!r}")
        if self.manifest.trace.require_redaction_status:
            for event in trace_events:
                if event.get("type") in self.manifest.trace.required_events and event.get("redaction_status") != "applied":
                    failures.append(f"trace event {event.get('type')} missing redaction_status=applied")
        if self.manifest.budgets.max_tools is not None and tool_count > self.manifest.budgets.max_tools:
            failures.append(f"tool budget exceeded: {tool_count}>{self.manifest.budgets.max_tools}")

        latency_budget_violated = (
            self.manifest.budgets.latency_ms is not None
            and duration_ms > self.manifest.budgets.latency_ms
        )
        if latency_budget_violated:
            failures.append(f"latency budget exceeded: {duration_ms}>{self.manifest.budgets.latency_ms}")

        return ScenarioResult(
            name=self.manifest.name,
            passed=not failures,
            duration_ms=duration_ms,
            tool_count=tool_count,
            permission_decisions=permission_decisions,
            trace_path=str(trace_path) if trace_path else None,
            failure_reason="; ".join(failures),
            latency_budget_violated=latency_budget_violated,
            stdout=stdout,
            stderr=stderr,
        )


def _load_trace(trace_path: Path | None) -> list[dict[str, Any]]:
    if trace_path is None or not trace_path.exists():
        return []
    events: list[dict[str, Any]] = []
    for line in trace_path.read_text().splitlines():
        if line.strip():
            events.append(json.loads(line))
    return events


def _redact(value: Any) -> Any:
    if isinstance(value, dict):
        return {key: _redact(item) for key, item in value.items()}
    if isinstance(value, list):
        return [_redact(item) for item in value]
    if isinstance(value, str):
        value = re.sub(r"(?i)bearer\s+[A-Za-z0-9._~+/=-]+", "Bearer " + REDACTED, value)
        value = re.sub(r"sk-[A-Za-z0-9][A-Za-z0-9._-]{8,}", REDACTED, value)
        value = re.sub(r"(?i)\b(api[_-]?key|authorization|password|secret|token)[_:=.-][A-Za-z0-9._~+/=-]+", REDACTED, value)
        return value
    return value
