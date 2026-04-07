# Architecture Deep Dive

Technical deep dive into claude-code-Go's architecture.

## Agent Loop State Machine

The agent loop is the core decision-making engine. It's a state machine with 6 states:

```
UserInput → Think → Act → Observe → [Loop to Think] → Respond
```

### State Details

**UserInput**
- Wait for user message
- Parse commands (/help, /mode, etc.)
- Queue message for processing

**Think**
- Analyze context
- Decide next action
- Generate tool calls or response

**Act**
- Execute tool
- Validate inputs
- Handle errors

**Observe**
- Process tool results
- Update context
- Decide if more actions needed

**Respond**
- Generate final response
- Stream output
- Return to UserInput

## Tool Registry Design

The tool registry manages all available tools:

```go
type Registry struct {
    tools map[string]Definition
    mu    sync.RWMutex
}
```

### Registration

Tools register at startup:
```go
func init() {
    tool.Register(Definition{
        Name: "Read",
        Handler: readHandler,
    })
}
```

### Execution

Tools execute with context:
```go
result, err := registry.Execute(ctx, "Read", args)
```

## Permission Model

The 3-tier permission system:

```go
type Enforcer struct {
    mode  PermissionMode
    rules []GlobRule
}
```

### Rule Evaluation

Rules evaluated in order:
1. Check explicit denials
2. Check explicit allows
3. Apply default mode behavior

## Context Management

Context compaction algorithm:

```go
func (c *Context) Compact() {
    // 1. Summarize old messages
    summary := c.Summarize(c.Messages[:threshold])
    
    // 2. Keep recent messages
    recent := c.Messages[threshold:]
    
    // 3. Reconstruct context
    c.Messages = append([]Message{summary}, recent...)
}
```

## Streaming Architecture

Custom SSE parser:

```go
func (p *Parser) NextEvent() (*Event, error) {
    for {
        line := p.reader.ReadString('\n')
        
        if strings.HasPrefix(line, "event: ") {
            event.Type = parseEventType(line)
        }
        
        if strings.HasPrefix(line, "data: ") {
            event.Data = parseData(line)
        }
        
        if line == "\n" {
            return event, nil
        }
    }
}
```

## Session Persistence

JSONL format for sessions:

```jsonl
{"type":"meta","session_id":"abc123"}
{"type":"message","role":"user","content":"Hello"}
{"type":"tool_call","tool":"Read","args":{"file":"main.go"}}
```

### Storage Strategy

- Append-only for durability
- Async writes for performance
- Compaction for size management

## Related Documents

- [ADR-001: Agent Loop Design](../adr/001-agent-loop.md)
- [ADR-002: Permission Model](../adr/002-permission-model.md)
- [ADR-003: Context Management](../adr/003-context-management.md)
- [ADR-004: Streaming Architecture](../adr/004-streaming-architecture.md)
- [ADR-005: Session Persistence](../adr/005-session-persistence.md)
