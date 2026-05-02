## 1. Diagnostic Foundation

- [x] 1.1 Define a shared extension diagnostic type with severity, stable code, component, summary, detail, retryability, and redacted metadata.
- [x] 1.2 Add formatting helpers for doctor, TUI, trace, and replay consumers.
- [x] 1.3 Add golden tests for diagnostic formatting and redaction.

## 2. MCP Security And Lifecycle

- [x] 2.1 Inventory current MCP config parsing, process launch, tool registration, and call paths.
- [x] 2.2 Add MCP launch validation for command allowlist, working directory, args, and env handling.
- [x] 2.3 Scrub MCP environment and server metadata before diagnostics or trace output.
- [x] 2.4 Enforce startup, list-tools, tool-call, and shutdown timeouts.
- [x] 2.5 Route MCP tool calls through shared permission policy and trace constructors.
- [x] 2.6 Add deterministic mock MCP server tests for registration, timeout, denial, and shutdown.
- [x] 2.7 Define MCP launch policy config format, default allowlist, cwd constraints, env inheritance rules, and compatibility behavior for existing configs.
- [x] 2.8 Make MCP transport request/response APIs context-aware so list-tools and tool-call timeouts do not depend on an unbounded blocking read.
- [x] 2.9 Map MCP tool safety classification into the shared permission action matrix.

## 3. LSP, Hooks, And Skills Diagnostics

- [x] 3.1 Convert LSP unavailable and unhealthy states into shared diagnostics.
- [x] 3.2 Convert invalid hook configuration and hook execution failures into shared diagnostics.
- [x] 3.3 Convert invalid skill files and skill loading warnings into shared diagnostics.
- [x] 3.4 Ensure optional extension failures do not fail core agent startup unless configured as blocking.

## 4. Provider Profile Reuse

- [x] 4.1 Separate provider profile metadata from transport request and streaming code.
- [x] 4.2 Add tests proving multiple profiles can reuse the same transport.
- [x] 4.3 Surface provider profile and model capability diagnostics in doctor and trace summaries.

## 5. User-Facing Surfaces

- [x] 5.1 Add offline doctor output for MCP, LSP, hooks, skills, and provider profile readiness.
- [x] 5.2 Add TUI status or command output for extension availability.
- [x] 5.3 Add replay output for extension diagnostics that affected a session.

## 6. Verification

- [x] 6.1 Run targeted MCP, LSP, hooks, skills, provider, doctor, and replay tests.
- [x] 6.2 Run `go test ./...`.
- [x] 6.3 Run `openspec validate productize-extension-boundaries --strict`.

## Evidence

- `go test ./internal/diagnostic ./internal/tool/mcp ./internal/lsp ./internal/hooks ./internal/skills ./internal/provider/registry ./cmd/go-code -count=1`
- `go test ./...`
- `openspec validate productize-extension-boundaries --strict --json --no-interactive`

## Residual Risks

- MCP launch policy is conservative and defaults to explicit allowlist behavior; users with custom commands may need documented launch policy entries.
- Extension diagnostics are now shared across doctor, trace, and replay surfaces; richer TUI presentation can be improved in a later UX-focused change.
