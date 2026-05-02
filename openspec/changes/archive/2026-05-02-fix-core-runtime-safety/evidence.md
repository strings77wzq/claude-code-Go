## Verification Evidence

- `go test ./internal/permission ./internal/tool ./internal/tool/init ./internal/session ./internal/agent ./internal/runstate ./pkg/tui ./cmd/go-code -count=1`
  - Result: pass.
- `go test -race ./pkg/tui -run 'TestModel_QuitCancelsActiveRequestAndIgnoresLateOutput|TestModel_StreamMessageKeepsDrainingUntilDone' -count=1`
  - Result: pass.
- `go test ./...`
  - Result: pass.
- `openspec validate fix-core-runtime-safety --strict --json --no-interactive`
  - Result: valid.

## Residual Risks

- Runtime trace `request_id` currently uses the session id for agent-level runtime events. This is stable and replayable, but a future enhancement can add per-turn request ids if the UI needs concurrent request correlation.
