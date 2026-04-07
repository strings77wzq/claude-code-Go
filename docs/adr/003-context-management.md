# ADR-003: Context Management

## Status
Accepted

## Context
LLMs have limited context windows (100K-200K tokens). We need to:
1. Fit conversations within token limits
2. Preserve important information
3. Support long-running sessions
4. Optimize token usage

## Decision
We implemented **Intelligent Context Management** with:

### Token Estimation

We estimate token count using:
```go
func EstimateTokens(text string) int {
    // ~4 characters per token for English
    return len(text) / 4
}
```

### Context Budget

We reserve tokens for:
- System prompt: ~500 tokens
- Recent messages: ~80% of remaining
- Tool results: ~20% of remaining

### Automatic Compaction

When context approaches limit:

1. **Summarization**: Old messages summarized
2. **Key Points Extraction**: Important facts preserved
3. **File References**: Keep file paths, summarize content

```go
func (c *Context) Compact() error {
    if c.TokenCount() < c.MaxTokens * 0.8 {
        return nil // No need to compact
    }
    
    // Summarize oldest messages
    summary := c.SummarizeMessages(c.Messages[:10])
    c.Messages = append([]Message{summary}, c.Messages[10:]...)
    
    return nil
}
```

### Manual Compaction

Users can trigger compaction:
- `/compact` - Summarize current session
- `/clear` - Clear all context (start fresh)

### Consequences

**Positive:**
- Sessions can run for hours
- Important info preserved
- Token usage optimized

**Negative:**
- Summarization may lose nuance
- Complex implementation
- Need to balance memory vs tokens

## Alternatives Considered

1. **Sliding Window**: Simple but loses old context
2. **Full History**: Hits token limits quickly
3. **External Vector Store**: Too complex for local tool

## Related Decisions

- ADR-001: Agent Loop (compaction happens between turns)
- ADR-005: Session Persistence (compact sessions saved to disk)
