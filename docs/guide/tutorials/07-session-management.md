# Tutorial 7: Session Management

Learn to manage conversation sessions.

## What are Sessions?

Sessions are persistent conversation histories. They allow you to:
- Resume conversations after restart
- Review past work
- Share context with team members

## Session Storage

Sessions are stored in `~/.go-code/sessions/`:

```
~/.go-code/sessions/
├── 2026-04-07-143022-abc123.jsonl
├── 2026-04-06-091544-xyz789.jsonl
└── 2026-04-05-185533-def456.jsonl
```

## Session Commands

### List Sessions

```
> /sessions

Recent sessions:
1. 2026-04-07 14:30:22 (abc123) - 15 messages
2. 2026-04-06 09:15:44 (xyz789) - 8 messages
3. 2026-04-05 18:55:33 (def456) - 42 messages

Use /resume <id> to resume a session
```

### Resume Session

```
> /resume abc123

Resumed session from 2026-04-07 14:30:22
Last message: "What's the status of the auth module?"

> Continue from where we left off
```

### Save Session

Sessions auto-save, but you can explicitly save:

```
> /save

Session saved: 2026-04-07-143022-abc123.jsonl
```

### Clear Session

Start fresh:

```
> /clear

⚠️  This will clear the current session.
Continue? [y/N] y

Session cleared. Starting fresh.
```

## Session Format

Sessions use JSONL format:

```jsonl
{"type":"session_meta","session_id":"abc123","created_at":"2026-04-07T14:30:22Z"}
{"type":"message","role":"user","content":"Hello"}
{"type":"message","role":"assistant","content":"Hi! How can I help?"}
{"type":"tool_call","tool":"Read","args":{"file_path":"main.go"}}
{"type":"tool_result","result":"package main..."}
```

## Session Compaction

When sessions get too long:

```
> /compact

Compacting session...
- 50 messages before
- 5 messages after (summary included)

Session compacted successfully
```

## Best Practices

1. **Name important sessions**: Add meaningful notes
2. **Compact regularly**: Keep sessions focused
3. **Clear between tasks**: Start fresh for new topics
4. **Backup sessions**: Copy important sessions elsewhere

## Next Steps

- [Tutorial 8: Error Handling](08-error-handling.md)
- [Architecture: Session Persistence](../../architecture/session-persistence.md)
