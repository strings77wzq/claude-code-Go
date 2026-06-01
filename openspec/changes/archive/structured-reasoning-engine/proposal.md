## Why

The current agent loop is purely reactive: model responds → execute tools → repeat. There is no explicit reasoning phase where the model can plan before acting. Modern frontier agents (Claude with extended thinking, OpenAI o1/o3) use structured reasoning with `<thinking>` blocks to plan, analyze trade-offs, verify assumptions, and self-critique before committing to tool calls. Adding structured reasoning improves reliability for complex multi-step tasks without changing the fundamental loop structure.

## What Changes

- **Extended thinking budget**: Support Anthropic's `thinking` parameter with configurable token budget for reasoning
- **Thinking block parsing**: Parse and preserve `<thinking>` blocks in API responses for trace/session visibility
- **Reasoning-aware loop**: Track whether the model is in a reasoning phase and handle stop_reason transitions accordingly
- **Reasoning trace events**: Record thinking blocks in session traces for debuggability and replay

## Capabilities

### New Capabilities

- `extended-thinking`: Configurable thinking budget for API requests, thinking block parsing in responses, reasoning trace events in session logs

### Modified Capabilities

None.

## Impact

- `internal/api/types.go` — new `Thinking` config struct, `ThinkingBlock` content type
- `internal/api/client.go` — inject `thinking` parameter into requests
- `internal/agent/loop.go` — handle `thinking` stop_reason, pass thinking blocks to trace
- `internal/session/` — append thinking events to trace
- No tool changes, no breaking API changes
