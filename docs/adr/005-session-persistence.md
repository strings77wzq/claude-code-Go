# ADR-005: Session Persistence

## Status
Accepted

## Context
Users expect to:
1. Resume conversations after closing the app
2. Review past sessions
3. Share sessions with others
4. Not lose work on crashes

We need persistence that:
1. Is human-readable
2. Is append-only for safety
3. Supports large files
4. Is easy to parse

## Decision
We use **JSONL (JSON Lines) Format**:

```jsonl
{"type": "session_meta", "session_id": "abc123", "created_at_ms": 123456}
{"type": "message", "role": "user", "content": "Hello"}
{"type": "message", "role": "assistant", "content": "Hi there!"}
{"type": "compaction", "compacted_message_count": 10}
```

### Why JSONL?

| Format | Pros | Cons |
|--------|------|------|
| **JSONL** | Append-only, human-readable, streamable | Larger than binary |
| SQLite | Compact, queryable | Binary, harder to debug |
| Protobuf | Compact, fast | Binary, harder to debug |
| JSON | Human-readable | Must rewrite entire file |

### Storage Location

```
~/.go-code/
├── sessions/
│   ├── 2026-04-07-abc123.jsonl
│   └── 2026-04-06-xyz789.jsonl
└── settings.json
```

### Session Recovery

On startup:
1. List all session files
2. Sort by modification time
3. Offer to resume latest
4. Load session on user request

### Consequences

**Positive:**
- Human-readable sessions
- Can `tail -f` for debugging
- Git-friendly format
- Easy to parse

**Negative:**
- Larger than binary formats
- No built-in querying
- Need to implement compaction

## Alternatives Considered

1. **SQLite**: Better for querying but harder to debug
2. **Single JSON File**: Must rewrite entire file on each save
3. **Binary Format**: Smaller but opaque

## Related Decisions

- ADR-001: Agent Loop (sessions saved after each turn)
- ADR-003: Context Management (compaction events in JSONL)
