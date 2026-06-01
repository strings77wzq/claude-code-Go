## Context

Anthropic's extended thinking feature allows Claude to spend additional tokens on internal reasoning before producing a response. Thinking tokens are billed separately and are invisible in the final answer but can be retrieved via API for debugging. The feature uses a `thinking` parameter with a `budget_tokens` value and returns thinking blocks alongside regular content blocks.

## Goals / Non-Goals

**Goals:**
- Add configurable thinking budget to agent (default: 0 = disabled, typical: 1024-4096)
- Parse thinking blocks from API responses
- Store thinking in session traces for observability
- Handle thinking-related stop reasons correctly

**Non-Goals:**
- Multi-model thinking support (Anthropic-specific initially)
- Streaming thinking content (batch mode only)
- Dynamic thinking budget adjustment based on task complexity
- Thinking content display in TUI (trace-only for now)

## Decisions

### 1. Thinking budget: agent-level config with sensible defaults

The thinking budget is set per-Agent via `SetThinkingBudget(tokens int)`. Default is 0 (disabled). Recommended for complex tasks: 1600-4096 tokens.

**Rationale**: Thinking has real cost (thinking tokens are billed at output rates). Making it opt-in with explicit configuration prevents surprise costs.

### 2. Thinking blocks in trace: full content, not truncated

Thinking content is stored in full in session traces (JSONL format). This enables:
- Debugging agent reasoning
- Replay with full context
- Quality analysis of reasoning patterns

**Rationale**: Thinking is invisible to users in final output but critical for debugging agent behavior.

### 3. Thinking stop_reason handling

When stop_reason is `thinking`, the agent sends the thinking block as a tool_result to continue the conversation (Anthropic API expects this pattern for multi-turn thinking).

## Risks / Trade-offs

- [Cost] → Thinking tokens are billed at output rates (~$15/MTok). Default disabled; explicit opt-in
- [Latency] → Thinking adds latency proportional to budget. TUI shows "thinking..." status
- [Provider lock-in] → Thinking is Anthropic-specific. OpenAI-compatible providers silently skip
