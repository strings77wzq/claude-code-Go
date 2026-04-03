"""Latency monitoring for measuring Go CLI response times."""

import time
from typing import Optional, Dict, Any


class LatencyMonitor:
    """Measures and tracks latency metrics for CLI responses.
    
    Tracks:
    - TTFT (Time To First Token): Time from request start to first token
    - Total latency: Time from request start to completion
    - Tokens per second: Throughput calculation
    
    Usage:
        monitor = LatencyMonitor()
        monitor.start_request()
        # ... run CLI command ...
        monitor.record_first_token()
        # ... more output ...
        monitor.record_complete()
        metrics = monitor.get_metrics()
    """
    
    def __init__(self):
        """Initialize the latency monitor."""
        self._request_start_ms: Optional[float] = None
        self._first_token_ms: Optional[float] = None
        self._complete_ms: Optional[float] = None
        self._token_count: int = 0
    
    def start_request(self) -> None:
        """Record the request start timestamp."""
        self._request_start_ms = time.perf_counter() * 1000
        self._first_token_ms = None
        self._complete_ms = None
        self._token_count = 0
    
    def record_first_token(self) -> None:
        """Record the first token timestamp (TTFT milestone)."""
        if self._request_start_ms is not None:
            self._first_token_ms = time.perf_counter() * 1000
    
    def record_complete(self) -> None:
        """Record the completion timestamp."""
        if self._request_start_ms is not None:
            self._complete_ms = time.perf_counter() * 1000
    
    def add_tokens(self, count: int) -> None:
        """Add token count for throughput calculation.
        
        Args:
            count: Number of tokens received.
        """
        self._token_count += count
    
    def get_metrics(self) -> Dict[str, Any]:
        """Get the calculated latency metrics.
        
        Returns:
            Dictionary with:
                - ttft_ms: Time to first token in milliseconds (float)
                - total_ms: Total request time in milliseconds (float)
                - tokens_per_second: Token throughput (float)
                - token_count: Total tokens received (int)
        """
        ttft_ms = 0.0
        total_ms = 0.0
        tokens_per_second = 0.0
        
        if self._request_start_ms is not None:
            # Calculate TTFT
            if self._first_token_ms is not None:
                ttft_ms = self._first_token_ms - self._request_start_ms
            
            # Calculate total latency
            if self._complete_ms is not None:
                total_ms = self._complete_ms - self._request_start_ms
            
            # Calculate tokens per second
            if total_ms > 0 and self._token_count > 0:
                tokens_per_second = (self._token_count / total_ms) * 1000
        
        return {
            "ttft_ms": ttft_ms,
            "total_ms": total_ms,
            "tokens_per_second": tokens_per_second,
            "token_count": self._token_count
        }
    
    def reset(self) -> None:
        """Reset all metrics to start fresh."""
        self._request_start_ms = None
        self._first_token_ms = None
        self._complete_ms = None
        self._token_count = 0
    
    def __enter__(self):
        """Context manager entry - starts request timing."""
        self.start_request()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit - records completion."""
        if exc_type is None:
            self.record_complete()
        return False


def measure_subprocess_latency(cmd: list, **kwargs) -> Dict[str, Any]:
    """Measure latency of a subprocess call.
    
    This is a convenience function to wrap a subprocess call and measure timing.
    
    Args:
        cmd: Command list to execute (passed to subprocess.Popen)
        **kwargs: Additional arguments passed to subprocess.Popen
        
    Returns:
        Dictionary with metrics from LatencyMonitor plus returncode and stdout.
    """
    import subprocess
    
    monitor = LatencyMonitor()
    monitor.start_request()
    
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, **kwargs)
    
    # Read output line by line, tracking first token
    stdout_lines = []
    first_token_recorded = False
    
    for line in iter(process.stdout.readline, b''):
        if line:
            if not first_token_recorded:
                monitor.record_first_token()
                first_token_recorded = True
            stdout_lines.append(line.decode('utf-8', errors='replace'))
    
    process.wait()
    monitor.record_complete()
    
    # Estimate token count (rough approximation)
    full_output = ''.join(stdout_lines)
    # Rough estimate: ~4 characters per token on average
    estimated_tokens = len(full_output) // 4
    monitor.add_tokens(estimated_tokens)
    
    return {
        **monitor.get_metrics(),
        "returncode": process.returncode,
        "stdout": full_output,
        "stderr": process.stderr.read().decode('utf-8', errors='replace')
    }