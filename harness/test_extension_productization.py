"""Harness scenarios for v0.3 extension productization evidence."""

import os
import subprocess
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def run_checked(args, *, env=None, timeout=30):
    completed = subprocess.run(
        args,
        cwd=PROJECT_ROOT,
        env={**os.environ, **(env or {})},
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        timeout=timeout,
    )
    assert completed.returncode == 0, completed.stderr + completed.stdout
    return completed.stdout + completed.stderr


class TestExtensionProductization:
    """Release-gate scenarios for MCP, LSP, and replay evidence."""

    def test_mcp_registration_go_evidence(self, go_binary):
        """MCP fixture registration creates a namespaced tool entry."""
        output = run_checked([
            "go",
            "test",
            "./internal/tool/mcp",
            "-run",
            "TestMcpManagerRegistersToolsFromFixture",
            "-count=1",
        ])
        assert "ok" in output

    def test_mcp_permission_denial_go_evidence(self, go_binary):
        """MCP permission denial returns a structured result without execution."""
        output = run_checked([
            "go",
            "test",
            "./internal/agent",
            "-run",
            "TestAgent_DeniedMCPToolReturnsStructuredResultWithoutExecuting",
            "-count=1",
        ])
        assert "ok" in output

    def test_lsp_unavailable_doctor_scenario(self, go_binary, tmp_path):
        """Offline doctor reports LSP unavailable without failing core diagnostics."""
        env = {
            "HOME": str(tmp_path),
            "ANTHROPIC_API_KEY": "sk-ant-test",
        }
        output = run_checked([str(go_binary), "doctor", "--offline"], env=env)
        assert "[SKIP] lsp: not configured" in output
        assert "[PASS] tools" in output

    def test_replay_evidence_output_scenario(self, go_binary, tmp_path):
        """Replay evidence mode summarizes extension, permission, and final status."""
        session_file = tmp_path / "session-evidence.jsonl"
        session_file.write_text(
            "\n".join([
                '{"type":"meta","session_id":"sess_harness","model":"test-model","start_time_ms":1,"status":"running"}',
                '{"type":"message","role":"user","content":"check extension state","timestamp_ms":2}',
                '{"type":"extension","name":"lsp","event":"health_check","status":"unavailable","reason":"not configured","timestamp_ms":3}',
                '{"type":"permission","tool":"mcp__fixture__write","decision":"Deny","summary":"blocked","timestamp_ms":4}',
                '{"type":"status","status":"completed","turn_count":1,"input_tokens":3,"output_tokens":4,"timestamp_ms":5}',
                "",
            ])
        )

        output = run_checked([str(go_binary), "replay", "--evidence", str(session_file)])
        assert "Replay evidence:" in output
        assert "extension lsp health_check status=unavailable" in output
        assert "permission mcp__fixture__write decision=Deny" in output
        assert "status completed" in output
