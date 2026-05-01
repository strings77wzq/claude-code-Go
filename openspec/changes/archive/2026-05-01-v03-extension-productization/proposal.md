## Why

The v0.2 runtime is now testable and the core harness is green, but the project still has active completed changes, stale audit contradictions, and partially productized extension surfaces. The next best move is to close the v0.2 evidence loop and turn MCP, LSP, hooks, skills, and replay from "code exists" into documented, gated v0.3 product capabilities.

## What Changes

- Archive or explicitly resolve completed OpenSpec changes so active work reflects the real roadmap.
- Add a v0.3 extension productization gate that requires working configuration, diagnostics, documentation, and deterministic tests before MCP/LSP/features are called supported.
- Expose MCP and LSP status through user-facing commands or doctor checks, including clear unavailable/error states.
- Strengthen session trace and replay as a debugging workflow for extension/tool runs.
- Align README, roadmap, PARITY.md, docs, and Chinese docs with the verified v0.2 baseline and planned v0.3 scope.
- Keep enterprise and content-marketing proposals parked until core extension productization is real.

## Capabilities

### New Capabilities
- `release-state-governance`: OpenSpec, docs, generated artifacts, and release evidence stay in a coherent state before new feature work starts.

### Modified Capabilities
- `mcp-lsp-extension-surface`: MCP, LSP, hooks, and skills move from partial/internal surfaces to permission-aware, documented, diagnosable product capabilities.
- `session-trace-and-replay`: Replay and trace requirements expand to cover extension/tool debugging and release evidence.
- `docs-spec-alignment`: Public docs must distinguish verified v0.2 support, v0.3 planned work, experimental extension surfaces, and parked business/marketing proposals.
- `harness-engineering-gate`: Harness gates must cover v0.3 extension trajectories, not only the v0.2 prompt/tool happy path.

## Impact

- Affected code: `cmd/go-code`, `internal/tool/mcp`, `internal/lsp`, `internal/hooks`, `internal/skills`, `internal/session`, `internal/logger`, `pkg/tui`, `pkg/tty`.
- Affected docs/specs: `README.md`, `PARITY.md`, `CHANGELOG.md`, `docs/`, `docs/zh/`, `openspec/specs/`, `openspec/changes/`.
- Affected workflows: `go test ./...`, `./scripts/run-harness.sh`, docs build, `go-code doctor --offline`, OpenSpec validation, release notes.
- No new runtime dependency is required unless a specific MCP/LSP smoke fixture cannot be represented with existing local mocks.
