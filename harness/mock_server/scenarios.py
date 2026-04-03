"""Scenario definitions for the mock Anthropic API server."""

from dataclasses import dataclass, field
from typing import Any


@dataclass
class Scenario:
    """A scenario defines a sequence of responses to return."""
    name: str
    messages: list[dict[str, Any]] = field(default_factory=list)


# Pre-built scenario definitions
SCENARIOS: dict[str, Scenario] = {
    "streaming_text": Scenario(
        name="streaming_text",
        messages=[
            {
                "type": "text",
                "text": "Hello! I'm Claude, an AI assistant created by Anthropic."
            }
        ]
    ),
    "tool_use_read": Scenario(
        name="tool_use_read",
        messages=[
            # First turn: return tool_use for Read
            {
                "type": "tool_use",
                "id": "toolu_001",
                "name": "Read",
                "input": {"file_path": "/test.txt"}
            },
            # Second turn: return text after tool_result
            {
                "type": "text",
                "text": "The file contains: Hello World"
            }
        ]
    ),
    "tool_use_bash": Scenario(
        name="tool_use_bash",
        messages=[
            # First turn: return tool_use for Bash
            {
                "type": "tool_use",
                "id": "toolu_002",
                "name": "Bash",
                "input": {"command": "ls -la"}
            },
            # Second turn: return text after tool_result
            {
                "type": "text",
                "text": "total 4\ndrwxr-xr-x 1 user user 4096 Apr  3 10:00 .\nddrwxr-xr-x   1 user user 4096 Apr  3 10:00 .."
            }
        ]
    ),
    "multi_turn": Scenario(
        name="multi_turn",
        messages=[
            # Turn 1
            {
                "type": "text",
                "text": "First response"
            },
            # Turn 2
            {
                "type": "text",
                "text": "Second response"
            },
            # Turn 3
            {
                "type": "text",
                "text": "Third response"
            }
        ]
    ),
}


class ScenarioRegistry:
    """Registry for managing scenarios and tracking turn counts."""

    def __init__(self):
        self._scenarios: dict[str, Scenario] = {}
        self._turn_counts: dict[str, int] = {}
        # Initialize with pre-built scenarios
        for name, scenario in SCENARIOS.items():
            self._scenarios[name] = scenario

    def get_scenario(self, name: str) -> Scenario | None:
        """Get a scenario by name."""
        return self._scenarios.get(name)

    def record_request(self, name: str, request: dict[str, Any]) -> None:
        """Track that a request was made for a scenario.
        
        This increments the turn counter for that scenario.
        """
        if name not in self._turn_counts:
            self._turn_counts[name] = 0
        self._turn_counts[name] += 1

    def get_turn_count(self, name: str) -> int:
        """Get the current turn count for a scenario."""
        return self._turn_counts.get(name, 0)

    def reset_turn_count(self, name: str) -> None:
        """Reset the turn count for a scenario."""
        self._turn_counts[name] = 0

    def register_scenario(self, scenario: Scenario) -> None:
        """Register a new scenario."""
        self._scenarios[scenario.name] = scenario
        self._turn_counts[scenario.name] = 0


# Global registry instance
registry = ScenarioRegistry()