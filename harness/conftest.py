"""Pytest fixtures for the Go CLI test harness."""

import os
import tempfile
from pathlib import Path

import pytest

from harness.mock_server.server import MockServer


@pytest.fixture(scope="session")
def go_binary() -> Path:
    """Return path to the compiled Go binary."""
    binary = Path("/home/strin/go/src/devLearn/aiLab/claude_code_Go/bin/go-code")
    if not binary.exists():
        pytest.skip(f"Go binary not found at {binary}")
    return binary


@pytest.fixture(scope="function")
def mock_server():
    """Start MockServer, yield base_url, stop on teardown."""
    server = MockServer()
    base_url = server.start()
    yield base_url
    server.stop()


@pytest.fixture(scope="function")
def test_workspace():
    """Create a temp directory with test files for tool operations."""
    with tempfile.TemporaryDirectory() as tmpdir:
        yield Path(tmpdir)


@pytest.fixture(scope="function")
def test_file(test_workspace):
    """Create a test file with known content in the workspace."""
    file_path = test_workspace / "test.txt"
    file_path.write_text("Hello World")
    return file_path


@pytest.fixture(scope="function")
def unique_content_file(test_workspace):
    """Create a file with unique content for edit tests."""
    file_path = test_workspace / "unique.txt"
    file_path.write_text("line 1\nline 2 unique content\nline 3\n")
    return file_path


@pytest.fixture(scope="function")
def repeated_content_file(test_workspace):
    """Create a file with repeated content for non-unique edit tests."""
    file_path = test_workspace / "repeated.txt"
    file_path.write_text("line 1\nrepeated line\nline 3\nrepeated line\nline 5\n")
    return file_path