## 1. Regression Evidence

- [x] 1.1 Add a failing regression test or harness fixture for TUI request cancellation leaving loading state or appending late output.
- [x] 1.2 Add permission policy tests for stdin closed, piped input, non-TTY, and CI-like non-interactive execution.
- [x] 1.3 Add a tool registry test with a panic-producing test tool.
- [x] 1.4 Add trace fixture expectations for cancellation, permission denial, and recovered tool panic events.

## 2. Permission Policy

- [x] 2.1 Define the permission mode hierarchy and action matrix in a shared policy package.
- [x] 2.2 Replace equality-based mode checks with hierarchy-aware policy evaluation.
- [x] 2.3 Implement fail-closed non-interactive approval behavior with stable reason codes.
- [x] 2.4 Wire the shared policy into built-in file, shell, network, and edit-like tools.
- [x] 2.5 Document user-visible permission mode behavior in the relevant CLI/TUI help text.
- [x] 2.6 Implement `EvaluateDetailed` with stable reason constants and migrate agent permission trace output to include reason codes.

## 3. Runtime Lifecycle

- [x] 3.1 Introduce a shared request lifecycle state type for CLI and TUI runs.
- [x] 3.2 Fix TUI request context ownership so cancellation happens after async command completion or explicit user cancellation.
- [x] 3.3 Ensure provider errors, permission denials, and cancelled requests stop loading state.
- [x] 3.4 Add goroutine leak or race-sensitive tests for cancellation paths where practical.
- [x] 3.5 Persist the active TUI request cancel function, keep stream commands draining until terminal state, and discard late output for cancelled request ids.

## 4. Tool Panic Containment

- [x] 4.1 Wrap registry tool execution with panic recovery.
- [x] 4.2 Return structured agent-visible tool errors for recovered panics.
- [x] 4.3 Emit redacted trace events for recovered panics.
- [x] 4.4 Remove or resolve stale duplicate tool registry paths that bypass the recovery wrapper.

## 5. Trace And Resume

- [x] 5.1 Add versioned trace event constructors with shared redaction.
- [x] 5.2 Migrate permission, runtime, and tool error events to the constructors.
- [x] 5.3 Preserve cancelled and failed terminal events in session history.
- [x] 5.4 Verify CLI and TUI resume behavior with persisted tool results and terminal events.
- [x] 5.5 Add `request_id` and runtime event subtype fields to trace envelopes where available.

## 6. Verification

- [x] 6.1 Run targeted tests for permission, TUI lifecycle, tool registry, and session trace packages.
- [x] 6.2 Run `go test ./...`.
- [x] 6.3 Run `openspec validate fix-core-runtime-safety --strict`.
- [x] 6.4 Record residual risks or follow-up tasks in the change notes before implementation completion.

## Evidence

- `go test ./internal/tool -run TestRegistryExecuteRecoversPanickingTool -count=1`
- `go test ./internal/permission -count=1`
- `go test ./internal/agent ./internal/tool ./internal/session ./pkg/tui ./internal/runstate -count=1`
- `go test ./...`
- `openspec validate fix-core-runtime-safety --strict --json --no-interactive`

## Residual Risks

- `internal/runstate` is introduced as the shared lifecycle contract and covered by tests; broader CLI/TUI lifecycle consolidation can continue during app composition cleanup.
- Non-interactive approval now fails closed through the default prompter and stdin EOF handling; explicit CLI permission-mode selection remains outside this change.
