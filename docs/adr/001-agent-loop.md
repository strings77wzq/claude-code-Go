# ADR-001: Agent Loop Design

## Status
Accepted

## Context
The agent loop is the core of claude-code-Go. It determines how the AI:
1. Receives user input
2. Decides which tools to use
3. Executes tools
4. Processes results
5. Generates responses

## Decision
We implemented a **Stop-Reason Driven State Machine** with the following states:

```
UserInput → Think → ToolSelection → Execute → Observe → Respond
                ↑                              ↓
                └────────── Loopback ←─────────┘
```

### State Transitions

| Current State | Trigger | Next State |
|--------------|---------|------------|
| UserInput | User sends message | Think |
| Think | Model generates tool_use | ToolSelection |
| Think | Model generates text | Respond |
| ToolSelection | Tool selected | Execute |
| Execute | Tool completes | Observe |
| Observe | Tool result sent to model | Think |
| Respond | Response shown | UserInput |

### Stop Reasons

The model can return these stop reasons:

- `end_turn`: Conversation complete, wait for user
- `tool_use`: Model wants to use a tool
- `max_tokens`: Context limit reached, need compaction
- `stop_sequence`: Special stop token encountered

### Consequences

**Positive:**
- Predictable behavior
- Easy to debug
- Clear state tracking
- Supports both single-turn and multi-turn

**Negative:**
- Complex state management
- Need to handle all edge cases
- Potential for infinite loops (mitigated with max_turns)

## Alternatives Considered

1. **Event-Driven Architecture**: More flexible but harder to reason about
2. **Coroutine-Based**: Could use Go channels, but added complexity
3. **Recursive Loop**: Simpler but harder to track state

## Related Decisions

- ADR-003: Context Management
- ADR-005: Session Persistence
