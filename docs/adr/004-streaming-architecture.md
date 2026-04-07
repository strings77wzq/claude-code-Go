# ADR-004: Streaming Architecture

## Status
Accepted

## Context
Users expect real-time feedback when:
1. AI is generating responses
2. Tools are executing
3. Files are being processed

We need streaming that:
1. Works with SSE (Server-Sent Events)
2. Handles network interruptions
3. Shows progress in UI
4. Is testable

## Decision
We implemented **Custom SSE Parser** with:

### SSE Event Format

```
event: content_block_delta
data: {"type": "text_delta", "text": "Hello"}

event: content_block_stop
data: {}
```

### Parser Design

```go
type SSEParser struct {
    reader *bufio.Reader
}

func (p *SSEParser) NextEvent() (*Event, error) {
    for {
        line, err := p.reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
        
        // Parse event: and data: lines
        if strings.HasPrefix(line, "event: ") {
            eventType = strings.TrimPrefix(line, "event: ")
        }
        if strings.HasPrefix(line, "data: ") {
            eventData = strings.TrimPrefix(line, "data: ")
        }
        
        // Empty line means end of event
        if line == "\n" {
            return &Event{Type: eventType, Data: eventData}, nil
        }
    }
}
```

### Zero Dependencies

We parse SSE manually (no external libraries) for:
- Smaller binary size
- Full control over parsing
- Easier debugging

### Consequences

**Positive:**
- Fast streaming
- Works with all providers
- Small binary size
- Full control

**Negative:**
- Need to maintain parser
- Edge cases in SSE spec
- Must handle reconnections

## Alternatives Considered

1. **External SSE Library**: Added dependencies
2. **WebSockets**: Overkill for this use case
3. **Long Polling**: Not real-time enough

## Related Decisions

- ADR-001: Agent Loop (streaming during Respond state)
