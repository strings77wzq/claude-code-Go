"""Scenario manifest loading and validation."""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


class ManifestValidationError(ValueError):
    """Raised when a scenario manifest is not valid."""


@dataclass(frozen=True)
class WorkspaceSpec:
    files: dict[str, str] = field(default_factory=dict)


@dataclass(frozen=True)
class AssertionsSpec:
    stdout_contains: list[str] = field(default_factory=list)


@dataclass(frozen=True)
class TraceSpec:
    required_events: list[str] = field(default_factory=list)
    require_redaction_status: bool = False


@dataclass(frozen=True)
class BudgetSpec:
    latency_ms: int | None = None
    max_tools: int | None = None


@dataclass(frozen=True)
class ScenarioManifest:
    schema_version: str
    name: str
    prompt: str
    workspace: WorkspaceSpec = field(default_factory=WorkspaceSpec)
    allowed_tools: list[str] = field(default_factory=list)
    assertions: AssertionsSpec = field(default_factory=AssertionsSpec)
    trace: TraceSpec = field(default_factory=TraceSpec)
    budgets: BudgetSpec = field(default_factory=BudgetSpec)


def load_manifest(path: Path) -> ScenarioManifest:
    data = json.loads(path.read_text())
    _validate_required(data, ["schema_version", "name", "prompt"])
    if data["schema_version"] != "agent-quality.v1":
        raise ManifestValidationError(f"unsupported schema_version: {data['schema_version']}")

    workspace = data.get("workspace", {})
    assertions = data.get("assertions", {})
    trace = data.get("trace", {})
    budgets = data.get("budgets", {})

    return ScenarioManifest(
        schema_version=data["schema_version"],
        name=data["name"],
        prompt=data["prompt"],
        workspace=WorkspaceSpec(files=dict(workspace.get("files", {}))),
        allowed_tools=list(data.get("allowed_tools", [])),
        assertions=AssertionsSpec(stdout_contains=list(assertions.get("stdout_contains", []))),
        trace=TraceSpec(
            required_events=list(trace.get("required_events", [])),
            require_redaction_status=bool(trace.get("require_redaction_status", False)),
        ),
        budgets=BudgetSpec(
            latency_ms=_optional_int(budgets.get("latency_ms")),
            max_tools=_optional_int(budgets.get("max_tools")),
        ),
    )


def _validate_required(data: dict[str, Any], fields: list[str]) -> None:
    missing = [field for field in fields if not data.get(field)]
    if missing:
        raise ManifestValidationError("missing required fields: " + ", ".join(missing))


def _optional_int(value: Any) -> int | None:
    if value is None:
        return None
    return int(value)
