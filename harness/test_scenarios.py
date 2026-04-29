"""Test scenarios for the Go CLI with mock server."""

import os
import subprocess


def run_go_cli(binary, base_url, input_text, env_overrides=None, timeout=30):
    """Run the go-code CLI with given input and return output."""
    env = {
        "ANTHROPIC_API_KEY": "test",
        "ANTHROPIC_BASE_URL": base_url,
        "ANTHROPIC_MODEL": "claude-sonnet-4-6-20251001",
        "LLM_PROVIDER": "anthropic",
    }
    if env_overrides:
        env.update(env_overrides)

    completed = subprocess.run(
        [str(binary), "-p", input_text, "-q", "-f", "json"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env={**os.environ, **env, "HOME": os.environ.get("HOME", "")},
        text=True,
        timeout=timeout,
    )
    return completed.stdout, completed.stderr, completed.returncode


class TestStreamingText:
    """Tests for streaming text responses."""

    def test_streaming_text(self, mock_server, go_binary):
        """Test that streaming text response is received from CLI."""
        # The default scenario is "streaming_text"
        input_text = "Hello"

        stdout, stderr, returncode = run_go_cli(
            go_binary, 
            mock_server, 
            input_text,
            timeout=30
        )

        # Verify output contains streamed text
        assert returncode == 0, stderr
        assert "Hello" in stdout or "Claude" in stdout, f"Expected streaming text in output, got: {stdout}"


class TestToolRoundtrip:
    """Tests for tool roundtrip functionality."""

    def test_tool_roundtrip_read(self, mock_server, go_binary, test_file):
        """Test that Read tool works via CLI."""
        # The mock server will respond with a tool_use for Read
        input_text = f"Read the file at {test_file}"

        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            input_text,
            timeout=30
        )

        # Verify the output mentions file content or the tool was called
        assert returncode == 0, stderr
        assert "Hello World" in stdout or "file contains" in stdout, f"Expected file read in output, got: {stdout}"

    def test_tool_roundtrip_bash(self, mock_server, go_binary):
        """Test that Bash tool works via CLI."""
        input_text = "run bash echo hello"

        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            input_text,
            timeout=30
        )

        assert returncode == 0, stderr
        assert "total" in stdout.lower() or "command" in stdout.lower(), f"Expected bash scenario output, got: {stdout}"

    def test_permission_denial_scenario(self, mock_server, go_binary):
        """Test that permission denial is represented as a deterministic scenario."""
        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            "Please test permission denial for rm -rf",
            timeout=30,
        )

        assert returncode == 0, stderr
        assert "denied" in stdout.lower() or "permission" in stdout.lower(), stdout

    def test_malformed_stream_does_not_panic(self, mock_server, go_binary):
        """Malformed stream events should fail cleanly, not panic."""
        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            "malformed stream please",
            timeout=30,
        )

        combined = stdout + stderr
        assert "panic" not in combined.lower(), combined
