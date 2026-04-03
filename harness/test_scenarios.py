"""Test scenarios for the Go CLI with mock server."""

import json
import os
import subprocess
import time

import pytest

from harness.mock_server import registry, Scenario


def run_go_cli(binary, base_url, input_text, env_overrides=None, timeout=30):
    """Run the go-code CLI with given input and return output."""
    env = {
        "ANTHROPIC_API_KEY": "test",
        "ANTHROPIC_BASE_URL": base_url,
    }
    if env_overrides:
        env.update(env_overrides)

    process = subprocess.Popen(
        [str(binary)],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env={**os.environ, **env},
        text=True,
    )

    stdout, stderr = process.communicate(input=input_text, timeout=timeout)
    return stdout, stderr, process.returncode


def set_mock_scenario(scenario_name: str):
    """Configure the mock server's default scenario before starting."""
    import harness.mock_server.app as app_module
    
    # Store the default scenario name in a module-level variable that app.py can check
    app_module._default_scenario = scenario_name


# Patch the app.py to use the configurable default scenario
import harness.mock_server.app as app_module
_original_get_scenario = app_module.registry.get_scenario

def _patched_get_scenario(self, name):
    # Check if there's a default scenario set
    default_scenario = getattr(app_module, '_default_scenario', None)
    if default_scenario:
        return _original_get_scenario(self, default_scenario)
    return _original_get_scenario(self, name)

app_module.registry.get_scenario = _patched_get_scenario.__get__(
    app_module.registry, type(app_module.registry)
)


class TestStreamingText:
    """Tests for streaming text responses."""

    def test_streaming_text(self, mock_server, go_binary):
        """Test that streaming text response is received from CLI."""
        # The default scenario is "streaming_text"
        input_text = "Hello\n"
        
        stdout, stderr, returncode = run_go_cli(
            go_binary, 
            mock_server, 
            input_text,
            timeout=30
        )

        # Verify output contains streamed text
        assert "Hello" in stdout or "Claude" in stdout, f"Expected streaming text in output, got: {stdout}"


class TestToolRoundtrip:
    """Tests for tool roundtrip functionality."""

    def test_tool_roundtrip_read(self, mock_server, go_binary, test_file):
        """Test that Read tool works via CLI."""
        # Configure mock server to use tool_use_read scenario
        set_mock_scenario("tool_use_read")
        
        # The mock server will respond with a tool_use for Read
        input_text = f"Read the file at {test_file}\n"
        
        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            input_text,
            timeout=30
        )

        # Verify the output mentions file content or the tool was called
        assert "Hello World" in stdout or "test.txt" in stdout, f"Expected file read in output, got: {stdout}"

    def test_tool_roundtrip_bash(self, mock_server, go_binary):
        """Test that Bash tool works via CLI."""
        # Configure mock server to use tool_use_bash scenario
        set_mock_scenario("tool_use_bash")
        
        input_text = "run echo hello\n"
        
        stdout, stderr, returncode = run_go_cli(
            go_binary,
            mock_server,
            input_text,
            timeout=30
        )

        # Verify the output mentions the command result
        assert "hello" in stdout.lower() or "echo" in stdout, f"Expected bash output in output, got: {stdout}"