---
title: Agent Loop Deep Dive
description: The state machine nature of the agent loop — stop_reason transitions, history management, and the permission gate
---

# Agent Loop Deep Dive

The agent loop is the heart of go-code's execution engine. But describing it as a "loop" obscures its true nature: it is a **state machine** where each iteration represents a transition between well-defined states, governed by the model's response through the `stop_reason` mechanism.

Understanding the agent loop as a state machine — rather than a simple iteration — reveals why certain design decisions were made and how they contribute to reliability.

## The State Machine View

At any point in execution, the agent is in one of a defined set of states. The transitions between these states are determined by the model's response:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Agent State Machine                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│         ┌──────────────────────────────────────────────────────────────┐   │
│         │                                                              │   │
│         │    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌───────┐ │   │
│         │    │ REQUEST │───▶│ WAITING │───▶│ RESPONSE│───▶│DECISION│ │   │
│         │    │  STATE  │    │  STATE  │    │  STATE  │    │ STATE  │ │   │
│         │    └─────────┘    └─────────┘    └─────────┘    └────┬────┘ │   │
│         │         ▲              │              │               │      │   │
│         │         │              │              │               │      │   │
│         │         └──────────────┴──────────────┴───────────────┘      │   │
│         │                              │                                  │   │
│         │                              │                                  │   │
│         │              ┌───────────────┴───────────────┐                 │   │
│         │              │                               │                 │   │
│         │              ▼                               ▼                 │   │
│         │    ┌─────────────────┐           ┌─────────────────────────┐     │   │
│         │    │  end_turn       │           │  tool_use               │     │   │
│         │    │  (TERMINAL)    │           │  (CONTINUE)             │     │   │
│         │    │                 │           │                         │     │   │
│         │    │ Return result  │           │  1. Permission check   │     │   │
│         │    │ to user        │           │  2. Tool execution     │     │   │
│         │    └─────────────────┘           │  3. Add results to     │     │   │
│         │                                    │  history              │     │   │
│         │                                    │  4. Transition back   │     │   │
│         │                                    │  to REQUEST state     │     │   │
│         │                                    └─────────────────────────┘     │   │
│         │                                                              │      │   │
│         └──────────────────────────────────────────────────────────────┘      │   │
│                                                                              │
│    Error Transitions:                                                       │
│    ┌──────────────────────────────────────────────────────────────────┐     │
│    │  max_tokens ──▶ Return truncated response with warning          │     │
│    │  max_turns  ──▶ Return with loop termination warning            │     │
│    │  error      ──▶ Return error to user                             │     │
│    └──────────────────────────────────────────────────────────────────┘     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Stop Reason as State Transition

The `stop_reason` field in the API response is not merely a status code — it is the **state transition signal** that determines the agent's next action. Each possible value represents a different state transition:

### end_turn: Terminal State

When `stop_reason` is `end_turn` (or `stop_sequence`), the model has completed its task. The response contains the final message to the user, and the agent loop terminates. This is a **terminal state** — no further transitions occur.

This state represents the model's confidence that the user's request has been fulfilled. The harness does not second-guess this decision — it trusts the model's assessment and returns the result.

### tool_use: Continue State

When `stop_reason` is `tool_use`, the model has requested to execute one or more tools to gather additional information or perform actions. The agent transitions to tool execution:

1. **Permission Gate**: Each requested tool passes through the permission system. If any tool is denied, its result is an error, but the loop continues.

2. **Parallel Execution**: Multiple tool requests are executed in parallel. The model can request several operations simultaneously, and the harness coordinates their execution.

3. **Result Capture**: Tool results (or errors) are captured and added to the conversation history.

4. **State Transition**: After tool execution, the agent transitions back to the REQUEST state, where the updated history is used to build the next API request.

### max_tokens: Truncation State

When `stop_reason` is `max_tokens`, the response was truncated because the maximum token limit was reached. This is not a failure — it's a boundary condition. The harness returns the partial response with a warning, allowing the user to decide whether to continue.

### max_turns: Loop Limit State

When the agent reaches the maximum number of iterations (MAX_TURNS = 50) without completing, it terminates with a warning. This is a **safety bound** — see the next section for why this exists.

## Why MAX_TURNS Exists

The model has no intrinsic understanding of time, resource constraints, or iteration counts. Left to run indefinitely, an agent could:

- Call the same tool repeatedly in an infinite loop
- Oscillate between ineffective strategies
- Consume unbounded computational resources

MAX_TURNS provides a **hard safety bound**. At 50 iterations, even a misbehaving agent will terminate. This number was chosen based on empirical observation: complex tasks rarely require more than 50 tool invocations, while simpler tasks complete well before this limit.

The value is configurable, allowing users who need more iterations to increase the limit while still maintaining *some* bound on execution time.

### What Happens at MAX_TURNS

When the limit is reached, the agent:
1. Stops making new API requests
2. Returns the accumulated response to the user
3. Includes a warning message indicating the loop was terminated

The user can then decide whether to continue the task manually or adjust the MAX_TURNS setting.

## History Management for Long Conversations

The conversation history is not merely a log — it is the **context** that allows the model to understand the ongoing task. Managing this history is one of the harness's most critical responsibilities.

### The Token Budget Problem

Every message in the history consumes tokens from the model's context window. A 50-turn conversation with 6 tool calls per turn could easily exceed 100,000 tokens — more than most models support.

Without management, the agent would either:
- Fail when hitting the token limit
- Truncate critical context needed for the task

### Three-Tier Compaction

go-code implements a three-tier approach to history management:

**Tier 1: Full Retention** — Early messages in the conversation are retained in full. These establish the task context and initial reasoning.

**Tier 2: Selective Compression** — Messages from the middle of the conversation are compressed into summaries that preserve key information (tool calls, their results, important decisions) while removing verbose details.

**Tier 3: Pruning** — If the conversation exceeds a certain length, older messages are pruned entirely, keeping only the most recent context.

This tiered approach ensures the model always has access to:
- The original task description
- Recent tool executions and their results
- The current state of work

While discarding information that is no longer relevant to the task at hand.

### Why Compaction Matters

Without compaction, the system would be limited to approximately 8-10 tool invocations before hitting token limits. With compaction, conversations can span 50+ iterations while maintaining enough context for the model to complete complex tasks.

Compaction is lossy — some information is inevitably discarded. But this is a necessary trade-off. The alternative (no compaction, limited conversation length) would be far more restrictive.

## The Permission Gate as Safety Checkpoint

Every tool execution passes through the permission system before execution. This is not optional — it is a **mandatory checkpoint** that guards dangerous operations.

### Why a Gate, Not a Filter

A filter would allow dangerous operations to pass through, relying on the model to avoid them. A gate actively intercepts and blocks operations that don't meet the configured policy.

This distinction is critical because the model operates with the user's full privileges. Without a gate, a misbehaving model could:
- Delete all files in the user's directory
- Execute malicious shell commands
- Exfiltrate sensitive data

### The Gate Architecture

The permission gate operates at the tool call level:

1. **Tier Check**: The requested tool is evaluated against the user's permission tier (ReadOnly, WorkspaceWrite, DangerFullAccess).

2. **Path Matching**: For file operations, glob patterns determine whether the specific path is allowed.

3. **Operation Evaluation**: Some operations (like deletion) are treated as inherently dangerous and may require additional approval.

4. **User Prompt**: If human-in-the-loop mode is enabled and the operation requires approval, the user is prompted before execution proceeds.

5. **Decision Enforcement**: The gate returns either approval or denial. If denied, the tool is not executed, and an error result is returned to the model.

### What the Model Sees

The model doesn't see the permission gate — it only sees tool results. If a tool is denied, the model receives an error result indicating the operation was blocked. The model must then decide how to proceed: try a different approach, request permission, or report the issue to the user.

## Streaming and Real-Time Feedback

The agent loop supports streaming output, allowing the user to see the model's response as it's generated. This is not merely a UX feature — it provides **early indication of progress**.

### Why Streaming Matters

Without streaming, the user would wait in silence until the complete response arrives — potentially many seconds for long outputs. With streaming, the user can:
- See that the request was received and is being processed
- Observe partial progress on long responses
- Interrupt the operation if something appears wrong

### How Streaming Works

The API client receives tokens incrementally and invokes a callback for each token delta. The harness uses this callback to:
- Stream text output to the terminal
- Display progress indicators during tool execution
- Provide real-time feedback to the user

The loop itself continues to process the full response — streaming is purely a display optimization.

## Execution Flow in State Machine Terms

```
1. INITIALIZE
   - Create new session history
   - Add user message to history
   - Transition to: REQUEST

2. REQUEST
   - Apply compaction if needed
   - Build API request (messages + tools)
   - Send request to API
   - Transition to: WAITING

3. WAITING
   - Receive streaming response
   - Buffer complete response
   - Transition to: RESPONSE

4. RESPONSE
   - Parse stop_reason
   - Extract content blocks
   - Transition to: DECISION

5. DECISION (state branching)
   - end_turn    → TERMINATE
   - tool_use    → TOOL_EXECUTION
   - max_tokens  → TERMINATE (with warning)
   - max_turns   → TERMINATE (with warning)
   - error       → TERMINATE (with error)

6. TOOL_EXECUTION
   - For each tool_use block:
     a. Permission gate check
     b. Execute tool (with timeout)
     c. Capture result (or error)
   - Add all results to history
   - Transition to: REQUEST

7. TERMINATE
   - Format final response
   - Return to user
   - End of execution
```

## Related Documentation

- [Architecture Overview](overview.md) — System components and data flow
- [Design Philosophy](design-philosophy.md) — Intelligence/reliability separation
- [Tool System](tools.md) — Tool implementations and permission configuration