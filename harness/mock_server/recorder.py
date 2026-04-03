"""Request recording utilities for the mock Anthropic API server."""

from typing import Any


class RequestRecorder:
    """Records all incoming requests for later inspection."""

    def __init__(self):
        self._requests: list[dict[str, Any]] = []

    def record(self, request: dict[str, Any]) -> None:
        """Record an incoming request."""
        self._requests.append(request)

    def get_requests(self) -> list[dict[str, Any]]:
        """Get all recorded requests."""
        return list(self._requests)

    def get_last_request(self) -> dict[str, Any] | None:
        """Get the most recent request."""
        if self._requests:
            return self._requests[-1]
        return None

    def clear(self) -> None:
        """Clear all recorded requests."""
        self._requests.clear()


# Global recorder instance
recorder = RequestRecorder()