"""MockServer class for managing the FastAPI server lifecycle."""

import socket
import subprocess
import sys
import time
from typing import Any


class MockServer:
    """Wrapper for the FastAPI mock Anthropic API server."""

    def __init__(self):
        self._process: subprocess.Popen | None = None
        self._base_url: str = ""
        self._port: int = 0

    @property
    def base_url(self) -> str:
        """Return the base URL of the running server."""
        return self._base_url

    def _find_free_port(self) -> int:
        """Find a free port on localhost."""
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind(("127.0.0.1", 0))
            s.listen(1)
            port = s.getsockname()[1]
        return port

    def start(self, host: str = "127.0.0.1") -> str:
        """Start the mock server on a dynamic port.

        Args:
            host: The host to bind to (default: 127.0.0.1)

        Returns:
            The base URL of the running server
        """
        self._port = self._find_free_port()
        self._base_url = f"http://{host}:{self._port}"

        # Start the server as a subprocess using uvicorn
        # We use subprocess to avoid complex threading issues
        self._process = subprocess.Popen(
            [
                sys.executable, "-m", "uvicorn",
                "harness.mock_server.app:app",
                "--host", host,
                "--port", str(self._port),
                "--log-level", "error",
            ],
            cwd="/home/strin/go/src/devLearn/aiLab/claude_code_Go",
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        # Wait for the server to be ready
        max_retries = 30
        for i in range(max_retries):
            try:
                import httpx
                response = httpx.get(f"{self._base_url}/health", timeout=1.0)
                if response.status_code == 200:
                    break
            except Exception:
                pass
            time.sleep(0.1)
        else:
            self.stop()
            raise RuntimeError("Mock server failed to start")

        return self._base_url

    def stop(self) -> None:
        """Stop the mock server."""
        if self._process:
            self._process.terminate()
            try:
                self._process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self._process.kill()
                self._process.wait()
            self._process = None
        self._base_url = ""
        self._port = 0