---
title: Agent Loop Architecture
description: Understanding the agent loop as a state machine — stop_reason dispatch, safety mechanisms, and design trade-offs
---

# Agent Loop Architecture

The agent loop is the **runtime engine** of the system. While the entry point establishes the dependency graph at startup, the agent loop executes the actual conversation cycles. Viewing it as a **state machine** rather than a procedural loop reveals the core design decisions that make the system robust and predictable.

## State Machine Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Agent Loop State Machine                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│     ┌──────────┐                                                   │
│     │ thinking │                                                   │
│     └────┬─────┘                                                   │
│          │ API response received                                    │
│          ▼                                                          │
│  ┌───────────────────────────────────────┐                        │
│  │         stop_reason dispatch          │                        │
│  └───────────────────────────────────────┘                        │
│          │                                                         │
│    ┌─────┴─────┬───────────┬────────────┐                       │
│    ▼           ▼           ▼            ▼                        │
│ ┌──────┐  ┌────────┐  ┌──────────┐  ┌──────────┐                 │
│ │end_  │  │max_    │  │tool_use  │  │unknown   │                 │
│ │turn  │  │tokens  │  │          │  │          │                 │
│ └──┬───┘  └────┬───┘  └────┬─────┘  └────┬─────┘                 │
│    │           │           │             │                        │
│    ▼           ▼           ▼             ▼                        │
│  Return      Return      Execute      Return                      │
│  response   +warning   tools +                               │
│                           continue                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## State Definitions

The agent loop operates in these states:

| State | Meaning | Next Action |
|-------|---------|--------------|
| **thinking** | Awaiting API response | Transition based on `stop_reason` |
| **tool_use** | Model requested tool execution | Execute tools, add results, continue loop |
| **end_turn** | Model completed task | Return response to user |
| **max_tokens** | Response truncated | Return partial response with warning |
| **unknown** | Unexpected stop reason | Return whatever content exists |

## Why a State Machine?

### Procedural vs. State Machine Thinking

A **procedural** view treats the loop as:
```
1. Send request
2. Get response
3. Check stop_reason
4. If tool_use: execute and loop
5. Else: return
```

This misses the deeper structure. A **state machine** view reveals:

- Each `stop_reason` is a **state transition** — the system moves from "thinking" to a specific outcome state
- The loop isn't just repeating; it's **cycling through well-defined states**
- Safety limits (MAX_TURNS) are **circuit breakers** that halt the state machine

### Design Trade-offs

Why use a state machine instead of alternatives?

| Approach | Why Not Used |
|----------|---------------|
| **Event-driven** | Over-engineered for this use case; the API already provides discrete responses |
| **Reactive/stream-based** | SSE provides streaming, but we need to reason about complete responses, not individual tokens |
| **Coroutine-based** | Go's goroutines are available, but the loop logic is simpler as explicit state transitions |

The state machine is the **minimal complexity** solution that captures the needed behavior: the API tells us what happened (`stop_reason`), and we respond accordingly.

## Stop Reason Dispatch

The `stop_reason` field is the **protocol contract** between the API and the agent. It tells the agent what the model did and what to do next:

### end_turn / stop_sequence

```
Model believes task is complete → Return response to user
```

This is the **happy path**. The model decided it has enough information and produced a final answer. The agent returns the text content directly.

### max_tokens

```
Model hit token limit → Return partial response with warning
```

This is a **degraded but safe** path. The response is truncated, so the user gets partial results plus a warning. They can continue the conversation to get the rest.

### tool_use

```
Model wants to call tools → Execute tools, add results, continue loop
```

This is the **iterative path**. The model identified that it needs to take action (read a file, run a command, search for something). The agent:
1. Executes all requested tools
2. Adds results to conversation history
3. Continues the loop to get the next response

### unknown / default

```
Unexpected stop_reason → Return whatever content exists
```

This is a **defensive fallback**. We don't crash on unexpected values; we return what we have and let the user decide what to do.

## Circuit Breaker — MAX_TURNS

```go
const MaxTurns = 50
```

MAX_TURNS is a **circuit breaker** in the electrical safety sense — it's designed to prevent catastrophic failure:

### What it prevents

- **Infinite loops**: Model calling tools that produce results that trigger more tool calls endlessly
- **Oscillation**: Model going back and forth between same approaches without progress
- **Resource exhaustion**: Running out of API tokens, memory, or time

### How it works

The loop enforces a maximum of 50 iterations (each iteration = one API call + potentially many tool executions). After 50 turns, the agent stops and returns a message indicating the limit was reached.

### Why 50?

This is an **empirical choice**:
- 50 turns is enough for complex tasks (analyze codebase, make changes, verify)
- It's beyond what a reasonable conversation would need
- It provides a safety net without being too restrictive

### What this design choice means

The circuit breaker assumes that **reasonable tasks complete within 50 turns**. If a task legitimately needs more, the design says: break the task into smaller pieces, or accept that the current approach isn't working.

## Context Window Optimization — History Management

The conversation history grows with each turn. The API has a **context window limit** (e.g., 200K tokens). The history management system optimizes this:

### Compaction Strategy

```
┌─────────────────────────────────────────────────────┐
│              History Compaction                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Turn 1:  [user: "..."]                             │
│           [assistant: "..."]                        │
│                                                     │
│  Turn 2:  [user: "..."]                             │
│           [assistant: "tool_use: read file"]        │
│           [tool_result: "file contents..."]        │
│           [assistant: "..."]                        │
│                                                     │
│  ...                                                │
│                                                     │
│  After compaction:                                  │
│  ─────────────────                                  │
│  [user: "original request"]                        │
│  [system: "summary of middle turns"]               │
│  [assistant: "current response"]                   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Why not keep everything?

- **Cost**: More tokens = more API costs
- **Performance**: Larger contexts = slower API responses
- **Model attention**: Extremely long contexts can reduce model focus on the current task

### Compaction triggers

Compaction happens when:
1. Token count approaches the limit
2. Before each API request (to ensure the request fits)

### What gets compacted

The system summarizes or removes older turns while preserving:
- The original user request (context)
- The most recent conversation (working memory)
- Tool definitions (always needed)

## Security Checkpoint — Permission Gate

Before any tool executes, it passes through the **permission gate**:

```
┌─────────────────────────────────────────────────────┐
│              Permission Gate                         │
├─────────────────────────────────────────────────────┤
│                                                     │
│  tool_request ──▶ checkPermission() ──▶ decision    │
│                      │                              │
│                      ▼                              │
│              ┌──────────────────┐                   │
│              │ PermissionPolicy │                   │
│              │  - Allow         │                   │
│              │  - Deny          │                   │
│              │  - Ask (prompt)  │                   │
│              └──────────────────┘                   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Why this design?

- **Defense in depth**: Not all tools are dangerous; permission check is a security layer
- **User control**: Users can allow/deny specific operations
- **Auditability**: Permission decisions are logged

### Tool classification

- **No permission needed**: Read, Glob, Grep (read-only, low-risk)
- **Requires permission**: Bash, Write, Edit (can modify system)

The permission policy evaluates each tool call based on:
1. Whether the tool requires permission
2. The user's configured permission level
3. The specific operation being attempted

## Crash Recovery — Session Persistence

Each turn completes with session persistence:

```go
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int) {
    // Save to ~/.go-code/sessions/
}
```

### What gets saved

- Session ID and timestamps
- Model used
- Turn count and token usage
- **Full conversation history**

### Why save after each turn?

```
┌─────────────────────────────────────────────────────┐
│            Session Persistence Flow                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Turn 1:  request → response → save                 │
│  Turn 2:  request → response → save                 │
│  Turn 3:  crash!                                    │
│                                                     │
│  Recovery: Load last session → resume from turn 3  │
│                                                     │
└─────────────────────────────────────────────────────┘
```

This design provides **crash recovery** without complexity:
- If the process dies mid-turn, the next run can resume
- No complex transaction log needed
- Sessions are small enough to save quickly

### Recovery considerations

- Sessions saved to disk survive process crash
- On next run, user can continue or start fresh
- Old sessions accumulate (cleanup is a future enhancement)

## Architectural Summary

The agent loop demonstrates several architectural principles:

| Principle | Implementation |
|-----------|----------------|
| **State Machine** | `stop_reason` as state transition trigger |
| **Circuit Breaker** | MAX_TURNS prevents infinite loops |
| **Context Optimization** | History compaction before each request |
| **Defense in Depth** | Permission gate before tool execution |
| **Crash Recovery** | Session persistence after each turn |
| **Fail Gracefully** | Unknown stop_reason returns partial response |

### Why these choices work

1. **State machine**: Matches the API's discrete response model
2. **Circuit breaker**: Provides safety without complexity
3. **History compaction**: Keeps costs predictable, performance high
4. **Permission gate**: Balances security with usability
5. **Session persistence**: Enables recovery without transaction complexity

The agent loop is the **operational core** — it takes the initialized components from the entry point and makes them do useful work through well-defined state transitions.

---

<div class="nav-prev-next">

- [Entry Point Architecture](entry-point.md) ←
- → [Tool System Overview](../tools/overview.md)

</div>