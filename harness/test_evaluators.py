"""Pytest tests for harness evaluators and replay utilities."""

import pytest
import time
import json
import tempfile
import os

from harness.evaluators.output_quality import (
    TextCompletenessEvaluator,
    ToolCallCorrectnessEvaluator,
)
from harness.evaluators.latency import LatencyMonitor
from harness.replay.session_replay import SessionReplayer
from harness.replay.trace_analyzer import TraceAnalyzer


class TestTextCompletenessEvaluator:
    """Tests for TextCompletenessEvaluator."""
    
    def test_exact_match(self):
        """Test exact text match returns complete."""
        expected = "Hello world"
        evaluator = TextCompletenessEvaluator(expected)
        result = evaluator.evaluate("Hello world")
        
        assert result["complete"] is True
        assert result["missing_ratio"] == 0.0
        assert result["missing_parts"] == []
    
    def test_partial_match(self):
        """Test partial text match returns incomplete with missing parts."""
        expected = "Hello world from testing"
        evaluator = TextCompletenessEvaluator(expected)
        result = evaluator.evaluate("Hello world")
        
        assert result["complete"] is False
        assert result["missing_ratio"] > 0.0
        assert len(result["missing_parts"]) > 0
    
    def test_empty_actual_returns_incomplete(self):
        """Test empty actual text returns incomplete."""
        expected = "Some expected text"
        evaluator = TextCompletenessEvaluator(expected)
        result = evaluator.evaluate("")
        
        assert result["complete"] is False
        assert result["missing_ratio"] == 1.0
    
    def test_empty_expected_returns_complete(self):
        """Test empty expected text returns complete."""
        evaluator = TextCompletenessEvaluator("")
        result = evaluator.evaluate("Any text")
        
        assert result["complete"] is True
        assert result["missing_ratio"] == 0.0
    
    def test_substring_match(self):
        """Test substring within larger text returns complete."""
        expected = "important info"
        evaluator = TextCompletenessEvaluator(expected)
        result = evaluator.evaluate("Here is some important info for you")
        
        assert result["complete"] is True
        assert result["missing_ratio"] == 0.0
    
    def test_word_level_ratio(self):
        """Test word-level ratio calculation."""
        expected = "The quick brown fox"
        evaluator = TextCompletenessEvaluator(expected)
        result = evaluator.evaluate("The quick")
        
        assert result["complete"] is False
        assert "brown" in result["missing_parts"] or "fox" in result["missing_parts"]


class TestToolCallCorrectnessEvaluator:
    """Tests for ToolCallCorrectnessEvaluator."""
    
    def test_exact_match(self):
        """Test exact tool call match returns correct."""
        expected_calls = [
            {"name": "bash", "params": {"command": "ls -la"}}
        ]
        evaluator = ToolCallCorrectnessEvaluator(expected_calls)
        
        actual_calls = [
            {"name": "bash", "params": {"command": "ls -la"}}
        ]
        result = evaluator.evaluate(actual_calls)
        
        assert result["correct"] is True
        assert result["mismatches"] == []
    
    def test_missing_tool_call(self):
        """Test missing expected tool call returns mismatch."""
        expected_calls = [
            {"name": "bash", "params": {"command": "ls"}},
            {"name": "read", "params": {"file_path": "/test.txt"}}
        ]
        evaluator = ToolCallCorrectnessEvaluator(expected_calls)
        
        actual_calls = [
            {"name": "bash", "params": {"command": "ls"}}
        ]
        result = evaluator.evaluate(actual_calls)
        
        assert result["correct"] is False
        assert len(result["mismatches"]) == 1
        assert result["mismatches"][0]["type"] == "missing"
    
    def test_extra_tool_call(self):
        """Test extra actual tool call returns mismatch."""
        expected_calls = [
            {"name": "bash", "params": {"command": "ls"}}
        ]
        evaluator = ToolCallCorrectnessEvaluator(expected_calls)
        
        actual_calls = [
            {"name": "bash", "params": {"command": "ls"}},
            {"name": "read", "params": {"file_path": "/test.txt"}}
        ]
        result = evaluator.evaluate(actual_calls)
        
        assert result["correct"] is False
        assert len(result["mismatches"]) == 1
        assert result["mismatches"][0]["type"] == "extra"
    
    def test_wrong_params(self):
        """Test wrong parameters returns mismatch."""
        expected_calls = [
            {"name": "bash", "params": {"command": "ls -la"}}
        ]
        evaluator = ToolCallCorrectnessEvaluator(expected_calls)
        
        actual_calls = [
            {"name": "bash", "params": {"command": "pwd"}}
        ]
        result = evaluator.evaluate(actual_calls)
        
        assert result["correct"] is False
        assert len(result["mismatches"]) > 0
    
    def test_empty_calls_match(self):
        """Test empty expected and actual calls match."""
        evaluator = ToolCallCorrectnessEvaluator([])
        result = evaluator.evaluate([])
        
        assert result["correct"] is True
    
    def test_param_containment(self):
        """Test parameter containment matching."""
        expected_calls = [
            {"name": "read", "params": {"file_path": "/test"}}
        ]
        evaluator = ToolCallCorrectnessEvaluator(expected_calls)
        
        actual_calls = [
            {"name": "read", "params": {"file_path": "/test.txt"}}
        ]
        result = evaluator.evaluate(actual_calls)
        
        # Should match because actual contains expected param value
        assert result["correct"] is True


class TestLatencyMonitor:
    """Tests for LatencyMonitor."""
    
    def test_basic_timing(self):
        """Test basic timing measurements."""
        monitor = LatencyMonitor()
        
        monitor.start_request()
        time.sleep(0.05)  # 50ms
        monitor.record_first_token()
        time.sleep(0.05)  # Another 50ms
        monitor.record_complete()
        
        metrics = monitor.get_metrics()
        
        assert metrics["ttft_ms"] >= 40  # Allow some tolerance
        assert metrics["total_ms"] >= 90
        assert metrics["ttft_ms"] < metrics["total_ms"]
    
    def test_token_count(self):
        """Test token count tracking."""
        monitor = LatencyMonitor()
        
        monitor.start_request()
        monitor.add_tokens(100)
        monitor.record_complete()
        
        metrics = monitor.get_metrics()
        
        assert metrics["token_count"] == 100
        assert metrics["tokens_per_second"] > 0
    
    def test_no_timing_recorded(self):
        """Test metrics when no timing recorded."""
        monitor = LatencyMonitor()
        metrics = monitor.get_metrics()
        
        assert metrics["ttft_ms"] == 0.0
        assert metrics["total_ms"] == 0.0
        assert metrics["tokens_per_second"] == 0.0
        assert metrics["token_count"] == 0
    
    def test_context_manager(self):
        """Test context manager usage."""
        with LatencyMonitor() as monitor:
            time.sleep(0.02)
        
        metrics = monitor.get_metrics()
        assert metrics["total_ms"] > 0
    
    def test_reset(self):
        """Test reset clears all metrics."""
        monitor = LatencyMonitor()
        
        monitor.start_request()
        monitor.record_first_token()
        monitor.add_tokens(50)
        monitor.reset()
        
        metrics = monitor.get_metrics()
        assert metrics["ttft_ms"] == 0.0
        assert metrics["token_count"] == 0


class TestSessionReplayer:
    """Tests for SessionReplayer."""
    
    def test_load_from_lines(self):
        """Test loading session from JSON lines."""
        lines = [
            '{"type": "session_meta", "session_id": "abc123", "created_at_ms": 123456}',
            '{"type": "message", "role": "user", "content": "Hello"}',
            '{"type": "message", "role": "assistant", "content": "Hi there", "tool_uses": [{"name": "bash", "input": {"command": "ls"}}]}',
            '{"type": "message", "role": "tool", "tool_use_id": "tool1", "content": "file1\\nfile2"}'
        ]
        
        replayer = SessionReplayer()
        replayer.load_from_lines(lines)
        
        summary = replayer.get_summary()
        
        assert summary["total_messages"] == 3
        assert summary["user_messages"] == 1
        assert summary["assistant_messages"] == 1
        assert summary["tool_calls"] == 1
    
    def test_replay_callback(self):
        """Test replay calls callback for each message."""
        lines = [
            '{"type": "message", "role": "user", "content": "Hello"}',
            '{"type": "message", "role": "assistant", "content": "Hi"}'
        ]
        
        replayer = SessionReplayer()
        replayer.load_from_lines(lines)
        
        received = []
        replayer.replay(on_message=lambda m: received.append(m))
        
        assert len(received) == 2
    
    def test_filter_by_role(self):
        """Test filtering messages by role."""
        lines = [
            '{"type": "message", "role": "user", "content": "Hello"}',
            '{"type": "message", "role": "assistant", "content": "Hi"}',
            '{"type": "message", "role": "user", "content": "Again"}'
        ]
        
        replayer = SessionReplayer()
        replayer.load_from_lines(lines)
        
        user_msgs = replayer.filter_by_role("user")
        
        assert len(user_msgs) == 2
    
    def test_get_tool_calls(self):
        """Test extracting tool calls from messages."""
        lines = [
            '{"type": "message", "role": "assistant", "content": "Working...", "tool_uses": [{"name": "bash", "input": {"command": "ls"}, "id": "call1"}]}'
        ]
        
        replayer = SessionReplayer()
        replayer.load_from_lines(lines)
        
        tool_calls = replayer.get_tool_calls()
        
        assert len(tool_calls) == 1
        assert tool_calls[0]["name"] == "bash"
        assert tool_calls[0]["params"]["command"] == "ls"


class TestTraceAnalyzer:
    """Tests for TraceAnalyzer."""
    
    def test_tool_call_patterns(self):
        """Test extracting tool call patterns."""
        events = [
            {"type": "tool_call", "name": "bash", "timestamp": 100},
            {"type": "tool_call", "name": "read", "timestamp": 200},
            {"type": "tool_call", "name": "bash", "timestamp": 300},
        ]
        
        analyzer = TraceAnalyzer(events)
        patterns = analyzer.get_tool_call_patterns()
        
        assert patterns["bash"] == 2
        assert patterns["read"] == 1
    
    def test_error_patterns(self):
        """Test extracting error patterns."""
        events = [
            {"type": "api_call", "error": "Connection timeout", "timestamp": 100},
            {"type": "api_call", "error": "Connection timeout", "timestamp": 200},
            {"type": "api_call", "error": "Permission denied", "timestamp": 300},
        ]
        
        analyzer = TraceAnalyzer(events)
        patterns = analyzer.get_error_patterns()
        
        assert len(patterns) == 2
        timeout_err = next(p for p in patterns if p["error_type"] == "timeout")
        assert timeout_err["count"] == 2
    
    def test_timeline(self):
        """Test timeline generation."""
        events = [
            {"type": "start", "timestamp": 100},
            {"type": "tool_call", "name": "bash", "timestamp": 200},
            {"type": "end", "timestamp": 300},
        ]
        
        analyzer = TraceAnalyzer(events)
        timeline = analyzer.get_timeline()
        
        assert len(timeline) == 3
        assert timeline[0]["type"] == "start"
        assert timeline[1]["name"] == "bash"
    
    def test_generate_report(self):
        """Test report generation."""
        events = [
            {"type": "tool_call", "name": "bash", "timestamp": 100},
            {"type": "api_call", "error": "Timeout", "timestamp": 200},
        ]
        
        analyzer = TraceAnalyzer(events)
        report = analyzer.generate_report()
        
        assert "TRACE ANALYSIS REPORT" in report
        assert "bash" in report
        assert "Timeout" in report or "error" in report.lower()


if __name__ == "__main__":
    pytest.main([__file__, "-v"])