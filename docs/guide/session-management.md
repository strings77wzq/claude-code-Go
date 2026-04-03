---
title: Session Management
description: How to list, resume, and manage saved sessions in claude-code-Go
---

# Session Management

claude-code-Go automatically saves your conversation sessions to disk, allowing you to review past conversations and resume interrupted sessions.

## How Sessions Are Saved

Sessions are saved in **JSONL** (JSON Lines) format in the session directory (default: `~/.go-code/sessions/`).

Each session file contains:
- **Line 1**: Session metadata (session ID, model, timestamps, turn count, token usage)
- **Lines 2+**: Individual messages (role, content, timestamp)

Example file structure:

```jsonl
{"type":"meta","session_id":"sess_123","model":"claude-sonnet-4-20250514","start_time_ms":1234567890000,"end_time_ms":1234567990000,"turn_count":5,"input_tokens":1000,"output_tokens":500}
{"type":"message","role":"user","content":"hello","timestamp_ms":1234567890000}
{"type":"message","role":"assistant","content":"Hello! How can I help you?","timestamp_ms":1234567895000}
```

## Session File Location

By default, sessions are stored in:

```
~/.go-code/sessions/
```

Each session file is named: `session-{timestamp}.jsonl`

The timestamp is the Unix epoch of the session start time.

## Listing Sessions

To view all saved sessions, use the `/sessions` command in the REPL:

```
/sessions
```

This will display a list of sessions sorted by most recent first, showing:
- Session ID
- Model used
- Start/end times
- Number of turns

## Resuming a Session

To resume a previous session, use the `/resume` command with the session ID:

```
/resume <session_id>
```

For example:
```
/resume sess_123
```

This will load all messages from the session file and restore your conversation state, allowing you to continue where you left off.

## Session Metadata

Each session stores the following metadata:

| Field | Description |
|-------|-------------|
| `session_id` | Unique identifier for the session |
| `model` | Claude model used (e.g., claude-sonnet-4-20250514) |
| `start_time_ms` | Session start time (Unix milliseconds) |
| `end_time_ms` | Session end time (Unix milliseconds) |
| `turn_count` | Number of conversation turns |
| `input_tokens` | Total input tokens used |
| `output_tokens` | Total output tokens generated |

## Programmatic Access

You can also access sessions programmatically using the `session` package:

```go
import "github.com/strings77wzq/claude-code-Go/internal/session"

// List all sessions
sessions, err := session.ListSessions("~/.go-code/sessions")

// Load a specific session
sess, messages, err := session.LoadSession("~/.go-code/sessions/session-1234567890.jsonl")
```

The `SessionInfo` struct contains:
- `ID`: Session identifier
- `FilePath`: Full path to the session file
- `StartTime`: Session start time
- `EndTime`: Session end time
- `TurnCount`: Number of turns
- `Model`: Model used