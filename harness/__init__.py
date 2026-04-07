"""Python harness for claude-code-Go testing infrastructure.

This package provides testing utilities for the claude-code-Go project:
- Mock Anthropic API server for testing without API costs
- Output quality evaluators for response validation
- Session replay utilities for debugging
- Integration tests for parity verification

Usage:
    from harness.mock_server import MockServer
    from harness.evaluators import TextCompletenessEvaluator
"""

__version__ = "0.1.0"
__all__ = [
    "mock_server",
    "evaluators",
    "replay",
]
