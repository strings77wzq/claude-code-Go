## 1. Token Counting Infrastructure

- [ ] 1.1 Add Go tokenization dependency (tiktoken-go or equivalent cl100k_base encoder)
- [ ] 1.2 Implement `CountTokens(messages []api.Message) int` replacing the 4-chars/token heuristic
- [ ] 1.3 Add token counting tests comparing against known reference values
- [ ] 1.4 Integrate token counting into `CompactIfNeeded` and `ShouldCompact`
- [ ] 1.5 Verify compaction threshold behavior with precise counting

## 2. Tool Tier System

- [ ] 2.1 Define `ToolTier` type (Core, Extension, MCP) in `internal/tool/tool.go`
- [ ] 2.2 Add optional `Tier() ToolTier` method to Tool interface
- [ ] 2.3 Annotate each builtin tool with tier: Read/Write/Edit/Bash/Grep=Core, Glob/Diff/Tree/WebFetch/TodoWrite/NotebookEdit=Extension
- [ ] 2.4 Add `GetDefinitionsByTier(tiers ...ToolTier) []ToolDefinition` to Registry
- [ ] 2.5 Implement progressive disclosure logic in `buildRequest()`: turn 1 sends all, subsequent turns send core only; on-demand re-include extension tools when model requests them
- [ ] 2.6 Track which extension tools the model has seen in the session state

## 3. Anthropic Prompt Caching

- [ ] 3.1 Add cache-related fields to `api.ApiRequest` (cache control headers)
- [ ] 3.2 Modify request builder to inject cache breakpoints at system prompt and tool definitions
- [ ] 3.3 Track cache invalidation triggers: tool registry change, model switch, system prompt change
- [ ] 3.4 Add cache hit/miss metrics to API response parsing
- [ ] 3.5 Record cache metrics in session trace (`traceCacheMetrics`)

## 4. Integration and Verification

- [ ] 4.1 Run `go test -count=1 ./internal/agent` and confirm all tests pass
- [ ] 4.2 Run `go test -count=1 ./internal/api` and confirm all tests pass
- [ ] 4.3 Run `go test -count=1 ./internal/tool` and confirm all tests pass
- [ ] 4.4 Run full test suite `go test ./...` and confirm 26/26 pass
- [ ] 4.5 Run `go vet ./...` and confirm clean
- [ ] 4.6 Run `gofmt -l .` and confirm clean
- [ ] 4.7 Build binary and run `./bin/go-code doctor --offline` confirm PASS
- [ ] 4.8 Run harness `./scripts/run-harness.sh` and confirm 48/48 pass
- [ ] 4.9 OpenSpec validate and verify artifacts complete
