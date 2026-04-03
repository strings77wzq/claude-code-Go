"""Session replay utilities for playing back recorded sessions from JSONL."""

import json
from typing import List, Dict, Any, Callable, Optional


class SessionReplayer:
    """Replays session data from JSONL format.
    
    JSONL format expected (each line is a JSON object):
    - {"type": "session_meta", "session_id": "...", "created_at_ms": ...}
    - {"type": "message", "role": "user", "content": "..."}
    - {"type": "message", "role": "assistant", "content": "...", "tool_uses": [...]}
    - {"type": "message", "role": "tool", "tool_use_id": "...", "content": "..."}
    - {"type": "compaction", " compacted_message_count": ...}
    
    The replayer handles any valid JSONL format and extracts message data.
    """
    
    def __init__(self):
        """Initialize the session replayer."""
        self._sessions: List[Dict[str, Any]] = []
        self._messages: List[Dict[str, Any]] = []
    
    def load_from_jsonl(self, filepath: str) -> List[Dict[str, Any]]:
        """Load session data from a JSONL file.
        
        Args:
            filepath: Path to the JSONL file to load.
            
        Returns:
            List of parsed JSON objects (session records).
        """
        self._sessions = []
        self._messages = []
        
        with open(filepath, 'r', encoding='utf-8') as f:
            for line_num, line in enumerate(f, 1):
                line = line.strip()
                if not line:
                    continue
                
                try:
                    record = json.loads(line)
                    self._sessions.append(record)
                    
                    # Extract messages for summary
                    if record.get("type") == "message":
                        self._messages.append(record)
                except json.JSONDecodeError as e:
                    # Skip malformed lines
                    continue
        
        return self._sessions
    
    def load_from_lines(self, lines: List[str]) -> List[Dict[str, Any]]:
        """Load session data from a list of JSON strings.
        
        Args:
            lines: List of JSON strings (one per line).
            
        Returns:
            List of parsed JSON objects.
        """
        self._sessions = []
        self._messages = []
        
        for line in lines:
            line = line.strip()
            if not line:
                continue
            
            try:
                record = json.loads(line)
                self._sessions.append(record)
                
                if record.get("type") == "message":
                    self._messages.append(record)
            except json.JSONDecodeError:
                continue
        
        return self._sessions
    
    def replay(self, messages: Optional[List[Dict[str, Any]]] = None, 
               on_message: Optional[Callable[[Dict[str, Any]], None]] = None) -> None:
        """Iterate through messages, calling the callback for each.
        
        Args:
            messages: Optional list of messages to replay. If None, uses loaded messages.
            on_message: Callback function called for each message with the message dict.
        """
        msgs = messages if messages is not None else self._messages
        
        if on_message is None:
            return
        
        for msg in msgs:
            on_message(msg)
    
    def get_messages(self) -> List[Dict[str, Any]]:
        """Get all loaded messages.
        
        Returns:
            List of message dictionaries.
        """
        return self._messages
    
    def get_summary(self) -> Dict[str, Any]:
        """Get a summary of the loaded session.
        
        Returns:
            Dictionary with:
                - total_messages: Total number of message records
                - user_messages: Number of user messages
                - assistant_messages: Number of assistant messages
                - tool_calls: Number of tool uses in assistant messages
                - total_tokens: Estimated total tokens (rough count)
        """
        user_messages = 0
        assistant_messages = 0
        tool_calls = 0
        total_tokens = 0
        
        for msg in self._messages:
            role = msg.get("role", "")
            content = msg.get("content", "")
            
            if role == "user":
                user_messages += 1
                # Rough estimate: ~4 chars per token
                total_tokens += len(content) // 4
            elif role == "assistant":
                assistant_messages += 1
                total_tokens += len(content) // 4
                
                # Count tool uses
                tool_uses = msg.get("tool_uses", [])
                if tool_uses:
                    tool_calls += len(tool_uses)
            elif role == "tool":
                # Tool response messages
                tool_calls += 1
                total_tokens += len(content) // 4
        
        return {
            "total_messages": len(self._messages),
            "user_messages": user_messages,
            "assistant_messages": assistant_messages,
            "tool_calls": tool_calls,
            "total_tokens": total_tokens
        }
    
    def filter_by_role(self, role: str) -> List[Dict[str, Any]]:
        """Filter messages by role.
        
        Args:
            role: The role to filter by (e.g., "user", "assistant", "tool").
            
        Returns:
            List of messages with the specified role.
        """
        return [msg for msg in self._messages if msg.get("role") == role]
    
    def get_tool_calls(self) -> List[Dict[str, Any]]:
        """Get all tool calls from assistant messages.
        
        Returns:
            List of tool call dictionaries with name and params.
        """
        tool_calls = []
        
        for msg in self._messages:
            if msg.get("role") == "assistant":
                tool_uses = msg.get("tool_uses", [])
                for tool_use in tool_uses:
                    tool_calls.append({
                        "name": tool_use.get("name", ""),
                        "params": tool_use.get("input", {}),
                        "id": tool_use.get("id", "")
                    })
        
        return tool_calls