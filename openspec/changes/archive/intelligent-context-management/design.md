## Context

The current agent sends all 11 tool definitions on every API request via `buildRequest()`. System prompt + tool definitions consume ~3.5K tokens per request. For a typical 10-turn session, that's 35K tokens wasted on redundant transmission. Anthropic's Prompt Caching offers up to 90% cost reduction on cached content. Additionally, the `EstimateTokens` function uses `totalChars / 4` — a heuristic that can be off by 30-50% compared to real tokenization.

## Goals / Non-Goals

**Goals:**
- Add Anthropic Prompt Caching for system prompt and tool definitions with cache breakpoints
- Implement 3-tier tool categorization (core: always visible; extension: sent on-demand; MCP: registered only)
- Replace character-based token estimation with accurate token counting
- Maintain backward compatibility — no behavior change, only efficiency improvement
- Track cache hit/miss metrics in session trace for observability

**Non-Goals:**
- Multi-provider caching support (focus on Anthropic first; OpenAI-compatible providers can be added later)
- Dynamic tool selection based on conversation context (future enhancement)
- Cache persistence across sessions
- Token-level streaming optimization

## Decisions

### 1. Cache architecture: ephemeral breakpoints on system + tools

Place cache breakpoints at the system prompt and before the first tool definition block. The messages array (which changes every turn) is NOT cached. This mirrors Anthropic's recommended pattern:

```
[system] ← cache breakpoint here
[tool_defs...] ← cache breakpoint here  
[messages...] ← NOT cached (changes every turn)
```

**Rationale**: This gives maximum cache reuse. System prompt rarely changes, tool definitions change only when MCP tools register/unregister. Messages inherently change every turn. Alternative considered: caching individual tools separately — rejected as over-complex for v1.

**Alternative considered**: Include last N messages in cache. Rejected because messages are the primary thing that changes between turns — caching them would cause frequent cache invalidations.

### 2. Tool tiers: static categorization

| Tier | Tools | When Sent | Token Budget |
|------|-------|-----------|--------------|
| Core | Read, Write, Edit, Bash, Grep | Always | ~1.5K |
| Extension | Glob, Diff, Tree, WebFetch, TodoWrite | On-demand via model request | ~1.2K |
| MCP | mcp__* | When registered | variable |

The model can request extension tools by name. If the response contains `tool_use` for an extension tool not yet disclosed, the agent sends the tool definition on the next request. MCP tools are auto-disclosed when their server is connected.

**Rationale**: Read/Write/Edit/Bash/Grep cover 90%+ of interactions. Keeping only these 5 in every request saves ~1.8K tokens/request. Extension tools can be added back on-demand — the model already knows about them from the full list sent on turn 1.

**Alternative considered**: Dynamic tool selection using embeddings. Rejected as premature optimization — static tiering is simpler, more predictable, and sufficient for the tool count (11).

### 3. Token counting: pure Go implementation

Use a Go port of the cl100k_base tokenizer used by Claude (same encoding as GPT-4). The `tiktoken-go` library provides this with no CGo dependency.

**Alternative considered**: Calling Anthropic's token-counting endpoint. Rejected because it adds network latency and a required API call for a purely local operation.

### 4. Cache invalidation

Cache is invalidated when:
- Tool registry changes (MCP server connects/disconnects)
- Model switches (different models have different system prompts)
- System prompt changes

Invalidation triggers a "cache miss" on the next request, after which caching resumes normally.

## Risks / Trade-offs

- [Cache not supported by all providers] → Feature-gate behind provider capability check; OpenAI-compatible providers silently skip caching
- [Extension tool not in cache causes extra round-trip] → Send extension tools on turn 1, then drop them; model can always request them back
- [Token counting library adds binary size] → Estimated ~200KB increase; acceptable for a Go binary
- [Tool tier misclassification] → Tier assignments are based on usage analysis and can be adjusted; core tier is conservative
