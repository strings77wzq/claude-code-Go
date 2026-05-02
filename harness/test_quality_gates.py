"""Manifest-driven agent quality gate tests."""

import json
from pathlib import Path

import pytest

from harness.quality.comparison import ComparisonReport, EvidenceRecord
from harness.quality.manifest import ManifestValidationError, load_manifest
from harness.quality.runner import ScenarioRunner


def test_load_valid_manifest(tmp_path: Path):
    manifest_path = tmp_path / "inspect.json"
    manifest_path.write_text(
        json.dumps(
            {
                "schema_version": "agent-quality.v1",
                "name": "repo-inspection",
                "prompt": "Inspect this repo",
                "workspace": {"files": {"README.md": "# Demo"}},
                "allowed_tools": ["Read", "Grep"],
                "assertions": {"stdout_contains": ["Demo"]},
                "trace": {"required_events": ["request"], "require_redaction_status": True},
                "budgets": {"latency_ms": 1000, "max_tools": 3},
            }
        )
    )

    manifest = load_manifest(manifest_path)

    assert manifest.name == "repo-inspection"
    assert manifest.schema_version == "agent-quality.v1"
    assert manifest.budgets.latency_ms == 1000


def test_builtin_manifests_are_valid():
    manifest_dir = Path(__file__).parent / "manifests"
    manifests = sorted(manifest_dir.glob("*.json"))

    assert manifests
    for manifest_path in manifests:
        manifest = load_manifest(manifest_path)
        assert manifest.name
        assert manifest.prompt


def test_invalid_manifest_reports_missing_fields(tmp_path: Path):
    manifest_path = tmp_path / "bad.json"
    manifest_path.write_text(json.dumps({"schema_version": "agent-quality.v1"}))

    with pytest.raises(ManifestValidationError) as exc:
        load_manifest(manifest_path)

    assert "name" in str(exc.value)
    assert "prompt" in str(exc.value)


def test_runner_records_trace_and_latency_budget_failure(tmp_path: Path):
    trace_path = tmp_path / "session.jsonl"
    trace_path.write_text(
        "\n".join(
            [
                json.dumps({"type": "request", "schema_version": "trace.v1", "redaction_status": "applied"}),
                json.dumps({"type": "permission", "decision": "Deny", "schema_version": "trace.v1", "redaction_status": "applied"}),
            ]
        )
        + "\n"
    )
    manifest_path = tmp_path / "permission.json"
    manifest_path.write_text(
        json.dumps(
            {
                "schema_version": "agent-quality.v1",
                "name": "permission-denial",
                "prompt": "Try denied command",
                "allowed_tools": ["Bash"],
                "trace": {"required_events": ["request", "permission"], "require_redaction_status": True},
                "budgets": {"latency_ms": 1},
            }
        )
    )
    manifest = load_manifest(manifest_path)

    result = ScenarioRunner(manifest).evaluate(
        stdout="Permission denied",
        stderr="",
        returncode=0,
        duration_ms=5,
        trace_path=trace_path,
    )

    assert result.passed is False
    assert result.latency_budget_violated is True
    assert result.permission_decisions == ["Deny"]
    assert result.trace_path == str(trace_path)


def test_evidence_redacts_secret_values(tmp_path: Path):
    manifest_path = tmp_path / "redact.json"
    manifest_path.write_text(
        json.dumps(
            {
                "schema_version": "agent-quality.v1",
                "name": "redaction",
                "prompt": "Check redaction",
                "trace": {"required_events": []},
            }
        )
    )
    result = ScenarioRunner(load_manifest(manifest_path)).evaluate(
        stdout="token=secret-token",
        stderr="Authorization: Bearer abc123",
        returncode=0,
        duration_ms=1,
        trace_path=None,
    )

    encoded = json.dumps(result.to_dict())
    assert "secret-token" not in encoded
    assert "abc123" not in encoded
    assert "[REDACTED]" in encoded


def test_comparison_report_labels_evidence_sources():
    report = ComparisonReport(
        [
            EvidenceRecord(agent="go-code", scenario="safe-edit", outcome="pass", duration_ms=100, tool_count=2, source="measured"),
            EvidenceRecord(agent="codex", scenario="safe-edit", outcome="pass", duration_ms=90, tool_count=3, source="manual"),
        ]
    )

    text = report.render()

    assert "go-code" in text
    assert "source=measured" in text
    assert "codex" in text
    assert "source=manual" in text
