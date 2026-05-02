## Why

The roadmap audit identified several release-blocking runtime safety risks: TUI request cancellation can outlive its context, non-interactive permission prompts can hang or degrade unpredictably, tool panic recovery is incomplete, and resumed sessions can lose coherent trace context. These defects sit on the critical path for a usable local agent because they affect every real trial run.

## What Changes

- Make TUI and non-interactive agent runs use a deterministic request lifecycle with explicit cancellation, completion, and error states.
- Make permission behavior fail closed when approval cannot be collected interactively.
- Normalize permission mode hierarchy so higher-trust modes satisfy lower requirements without bypassing explicit deny policies.
- Add a runtime safety gate around tool execution so panics become structured tool errors and trace events.
- Version session trace events and preserve resumable context across TUI and CLI flows.

## Capabilities

### New Capabilities
- `runtime-safety-gates`: Defines crash containment, request lifecycle, and validation gates for agent runtime execution.

### Modified Capabilities
- `permission-and-sandbox-flow`: Adds fail-closed non-interactive semantics and an explicit permission decision matrix.
- `session-trace-and-replay`: Adds versioned trace envelopes and coherent resume evidence for interrupted or cancelled runs.
- `smooth-cli-tui-experience`: Adds deterministic TUI loading, cancellation, and error-state behavior.

## Impact

- Affected code: `pkg/tui`, `pkg/tty`, `internal/agent`, `internal/tool`, `internal/permission`, `internal/session`, command composition under `cmd/go-code`.
- Affected tests: permission policy unit tests, TUI command lifecycle tests, tool registry panic tests, session trace/replay tests.
- No new external dependencies are expected.
