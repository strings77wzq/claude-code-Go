"""CLI entry point for the mock Anthropic API server.

Run with: python -m harness.mock_server
"""

import sys
import time

from .server import MockServer


def main():
    """Start the mock server and print the base URL."""
    server = MockServer()

    try:
        base_url = server.start()
        print(f"MOCK_ANTHROPIC_BASE_URL={base_url}")
        print(f"Server running at {base_url}")
        print("Press Ctrl+C to stop...")

        # Keep the server running
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nStopping mock server...")
    finally:
        server.stop()
        print("Mock server stopped.")


if __name__ == "__main__":
    main()