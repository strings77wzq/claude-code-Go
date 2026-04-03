"""Trace analysis utilities for examining session execution patterns."""

from typing import List, Dict, Any, Optional
from collections import defaultdict, Counter


class TraceAnalyzer:
    """Analyzes trace events from session execution.
    
    Trace events include tool calls, API responses, errors, and other
    runtime events. Provides pattern detection and reporting.
    """
    
    def __init__(self, events: List[Dict[str, Any]]):
        """Initialize with trace events.
        
        Args:
            events: List of trace event dictionaries. Each event should have
                   at least a "type" field, and optionally "name", "timestamp",
                   "error", "tool_name", etc.
        """
        self.events = events
        self._sorted_events: Optional[List[Dict[str, Any]]] = None
    
    def _ensure_sorted(self) -> None:
        """Ensure events are sorted by timestamp."""
        if self._sorted_events is None:
            self._sorted_events = sorted(
                self.events, 
                key=lambda e: e.get("timestamp", e.get("timestamp_ms", 0))
            )
    
    def get_tool_call_patterns(self) -> Dict[str, int]:
        """Get tool call frequency patterns.
        
        Returns:
            Dictionary mapping tool names to call counts.
        """
        tool_counts: Dict[str, int] = defaultdict(int)
        
        for event in self.events:
            # Check various event types for tool calls
            if event.get("type") == "tool_call":
                tool_name = event.get("name", event.get("tool_name", "unknown"))
                tool_counts[tool_name] += 1
            elif event.get("type") == "tool_use":
                tool_name = event.get("tool_name", "unknown")
                tool_counts[tool_name] += 1
            elif event.get("type") == "message" and event.get("role") == "assistant":
                # Check for tool_uses in assistant messages
                tool_uses = event.get("tool_uses", [])
                for tool_use in tool_uses:
                    tool_name = tool_use.get("name", "unknown")
                    tool_counts[tool_name] += 1
        
        return dict(tool_counts)
    
    def get_error_patterns(self) -> List[Dict[str, Any]]:
        """Get error patterns with counts and examples.
        
        Returns:
            List of dictionaries with error_type, count, and examples.
        """
        error_groups: Dict[str, Dict[str, Any]] = {}
        
        for event in self.events:
            error = event.get("error")
            if error:
                # Determine error type
                if isinstance(error, str):
                    # Extract error type from error message
                    if "timeout" in error.lower():
                        error_type = "timeout"
                    elif "connection" in error.lower():
                        error_type = "connection"
                    elif "permission" in error.lower():
                        error_type = "permission"
                    elif "not found" in error.lower():
                        error_type = "not_found"
                    elif "invalid" in error.lower():
                        error_type = "invalid"
                    else:
                        error_type = "general"
                else:
                    error_type = str(type(error).__name__)
                
                if error_type not in error_groups:
                    error_groups[error_type] = {
                        "error_type": error_type,
                        "count": 0,
                        "examples": []
                    }
                
                error_groups[error_type]["count"] += 1
                
                # Add up to 3 examples
                if len(error_groups[error_type]["examples"]) < 3:
                    error_groups[error_type]["examples"].append(
                        str(error)[:200]  # Truncate long errors
                    )
        
        return list(error_groups.values())
    
    def get_timeline(self) -> List[Dict[str, Any]]:
        """Get chronological list of events with timestamps.
        
        Returns:
            List of events sorted by timestamp, each with type, timestamp, and other fields.
        """
        self._ensure_sorted()
        
        timeline = []
        for event in self._sorted_events:
            timeline_entry = {
                "type": event.get("type", "unknown"),
                "timestamp": event.get("timestamp", event.get("timestamp_ms", 0)),
            }
            
            # Add relevant fields based on event type
            if event.get("type") == "tool_call":
                timeline_entry["tool_name"] = event.get("name", event.get("tool_name", ""))
            elif event.get("type") == "message":
                timeline_entry["role"] = event.get("role", "")
            
            if "error" in event:
                timeline_entry["has_error"] = True
            
            timeline.append(timeline_entry)
        
        return timeline
    
    def get_api_call_stats(self) -> Dict[str, Any]:
        """Get statistics about API calls in the trace.
        
        Returns:
            Dictionary with API call statistics.
        """
        api_calls = [e for e in self.events if e.get("type") == "api_call"]
        
        return {
            "total_api_calls": len(api_calls),
            "successful_calls": len([c for c in api_calls if not c.get("error")]),
            "failed_calls": len([c for c in api_calls if c.get("error")]),
        }
    
    def generate_report(self) -> str:
        """Generate a human-readable report of the trace analysis.
        
        Returns:
            String containing the formatted report.
        """
        lines = []
        lines.append("=" * 60)
        lines.append("TRACE ANALYSIS REPORT")
        lines.append("=" * 60)
        lines.append("")
        
        # Basic stats
        lines.append(f"Total Events: {len(self.events)}")
        lines.append("")
        
        # Tool call patterns
        tool_patterns = self.get_tool_call_patterns()
        if tool_patterns:
            lines.append("TOOL CALL PATTERNS:")
            lines.append("-" * 40)
            for tool_name, count in sorted(tool_patterns.items(), key=lambda x: -x[1]):
                lines.append(f"  {tool_name}: {count} calls")
            lines.append("")
        
        # Error patterns
        error_patterns = self.get_error_patterns()
        if error_patterns:
            lines.append("ERROR PATTERNS:")
            lines.append("-" * 40)
            for err in error_patterns:
                lines.append(f"  {err['error_type']}: {err['count']} occurrences")
                for example in err["examples"][:2]:
                    lines.append(f"    Example: {example[:80]}...")
            lines.append("")
        
        # Timeline summary
        timeline = self.get_timeline()
        if timeline:
            lines.append("TIMELINE SUMMARY:")
            lines.append("-" * 40)
            lines.append(f"  Total events: {len(timeline)}")
            if timeline[0].get("timestamp"):
                lines.append(f"  First event: {timeline[0]['timestamp']}")
            if timeline[-1].get("timestamp"):
                lines.append(f"  Last event: {timeline[-1]['timestamp']}")
            lines.append("")
        
        # API stats
        api_stats = self.get_api_call_stats()
        if api_stats["total_api_calls"] > 0:
            lines.append("API CALL STATISTICS:")
            lines.append("-" * 40)
            lines.append(f"  Total: {api_stats['total_api_calls']}")
            lines.append(f"  Successful: {api_stats['successful_calls']}")
            lines.append(f"  Failed: {api_stats['failed_calls']}")
            lines.append("")
        
        lines.append("=" * 60)
        
        return "\n".join(lines)
    
    def get_session_summary(self) -> Dict[str, Any]:
        """Get a summary dictionary of the trace.
        
        Returns:
            Dictionary with summary statistics.
        """
        return {
            "total_events": len(self.events),
            "tool_call_patterns": self.get_tool_call_patterns(),
            "error_count": len(self.get_error_patterns()),
            "api_stats": self.get_api_call_stats(),
            "timeline_length": len(self.get_timeline())
        }