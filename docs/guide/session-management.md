---
title: Session Management
description: How to list, resume, and manage saved sessions in claude-code-Go
---

# Session Management

claude-code-Go automatically saves your conversation sessions to disk, allowing you to review past conversations and resume interrupted sessions.

## How Sessions Are Saved

Sessions are saved in **JSONL** (JSON Lines) format in the session directory (default: `~/.claude-code-go/sessions/`).

Each session file contains normalized trace lines:
- `meta`: session ID, model, start timestamp, and initial status.
- `request` / `response`: provider request and response summaries.
- `tool`: tool execution name, input, output, and duration.
- `permission`: permission decision, sanitized summary, and timestamp.
- `message`: recoverable user and assistant messages.
- `status`: final status, turn count, and token usage.

Example file structure:

```jsonl
{"type":"meta","session_id":"sess_123","model":"claude-sonnet-4-6-20251001","start_time_ms":1234567890000,"end_time_ms":0,"turn_count":0,"input_tokens":0,"output_tokens":0,"status":"running"}
{"type":"request","model":"claude-sonnet-4-6-20251001","messages_count":1,"timestamp_ms":1234567891000}
{"type":"message","role":"user","content":"hello","timestamp_ms":1234567890000}
{"type":"message","role":"assistant","content":"Hello! How can I help you?","timestamp_ms":1234567895000}
{"type":"status","status":"completed","turn_count":1,"input_tokens":1000,"output_tokens":500,"timestamp_ms":1234567990000}
```

## Session File Location

By default, sessions are stored in:

```
~/.claude-code-go/sessions/
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

## Replaying a Session

Use `go-code replay` to inspect a saved session without calling a provider:

```bash
go-code replay latest
go-code replay sess_123
go-code replay ~/.claude-code-go/sessions/session-1234567890.jsonl
```

Replay prints the sequence of requests, tool calls, permission decisions, messages, errors, and final status. It is intended for debugging and issue reports.

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
| `status` | `running`, `completed`, `failed`, or `max_turns` |

## Programmatic Access

You can also access sessions programmatically using the `session` package:

```go
import "github.com/strings77wzq/claude-code-Go/internal/session"

// List all sessions
sessions, err := session.ListSessions("~/.claude-code-go/sessions")

// Load a specific session
sess, messages, err := session.LoadSession("~/.claude-code-go/sessions/session-1234567890.jsonl")

// Replay a session trace
events, err := session.ReplaySessionFile("~/.claude-code-go/sessions/session-1234567890.jsonl")
```

The `SessionInfo` struct contains:
- `ID`: Session identifier
- `FilePath`: Full path to the session file
- `StartTime`: Session start time
- `EndTime`: Session end time
- `TurnCount`: Number of turns
- `Model`: Model used
