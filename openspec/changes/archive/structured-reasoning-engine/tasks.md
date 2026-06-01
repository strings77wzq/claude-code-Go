## 1. API Types

- [ ] 1.1 Add `Thinking` config struct to `internal/api/types.go`
- [ ] 1.2 Add `thinking` field to `ApiRequest`
- [ ] 1.3 Add `ThinkingBlock` content block type for responses

## 2. Agent Integration

- [ ] 2.1 Add `thinkingBudget` field to Agent struct
- [ ] 2.2 Add `SetThinkingBudget(tokens int)` method
- [ ] 2.3 Inject `thinking` parameter into `buildRequest()` when budget > 0
- [ ] 2.4 Handle `thinking` stop_reason in the agent loop

## 3. Trace Integration

- [ ] 3.1 Add `traceThinking` method to Agent
- [ ] 3.2 Add `AppendTraceThinking` to session package
- [ ] 3.3 Record thinking blocks in trace when received

## 4. Verification

- [ ] 4.1 Run `go test -count=1 ./internal/agent` and confirm all tests pass
- [ ] 4.2 Run `go test -count=1 ./internal/api` and confirm all tests pass
- [ ] 4.3 Run full test suite `go test ./...` and confirm 26/26 pass
- [ ] 4.4 Run `go vet ./...` and `gofmt -l .` confirm clean
