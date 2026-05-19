## Why

The current agent sends all 11 tool definitions (~3K tokens) on every API request, uses a crude 4-characters-per-token heuristic for context management, and does not leverage Anthropic's Prompt Caching — one of the most cost-effective optimizations available (90% cost reduction on cached prompt tokens). In modern agent design, context is the most expensive resource. These three optimizations together can reduce API costs by 30-50% while lowering latency and enabling longer effective conversations.

## What Changes

- **Anthropic Prompt Caching**: Add cache breakpoints to system prompt and tool definitions so they are cached and reused across turns. Cache read hits cost 90% less than base input tokens.
- **Progressive Tool Disclosure**: Categorize 11 tools into 3 tiers — core (5 always-sent), extension (5 sent on demand), MCP (sent when registered). Core tools cover 90%+ of typical interactions.
- **Precise Token Counting**: Replace the `EstimateTokens` 4 chars/token heuristic with accurate counting using Anthropic's tokenization or a Go-native tiktoken port.
- **Cache-aware request builder**: Track cache breakpoints per request, detect when underlying tool/message state invalidates cached prefixes.

## Capabilities

### New Capabilities

- `prompt-caching`: Anthropic Prompt Caching support — cache breakpoints on system prompt and tool definitions, cache-aware request building, cache hit/miss metrics in trace
- `progressive-tool-disclosure`: Three-tier tool categorization (core/extension/mcp) with lazy tool definition loading, reducing per-request token overhead by 40-60%
- `precise-token-counting`: Replace 4-chars-per-token heuristic with accurate token counting, enabling precise context window management

### Modified Capabilities

None — all are new capabilities that don't change existing spec-level behavior.

## Impact

- `internal/api/client.go` — cache breakpoint injection, cache-aware request headers
- `internal/agent/loop.go` — request builder modified to use progressive tool disclosure
- `internal/tool/registry.go` — tier-based tool registration and filtered definition retrieval
- `internal/agent/compact.go` — `EstimateTokens` replaced with precise counter
- `internal/tool/builtin/*.go` — each tool annotated with tier metadata
- New dependency: Go tokenization library or Anthropic token-counting integration
