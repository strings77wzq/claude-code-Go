## Why

MCP, LSP, hooks, and skills are strategically important extension surfaces, but the roadmap audit found that their security, lifecycle, diagnostics, and reuse boundaries are not yet strong enough for a release candidate. Productizing these boundaries after core runtime safety keeps extension power from becoming an uncontrolled trust or reliability risk.

## What Changes

- Define a reusable extension lifecycle contract for startup, health, timeout, shutdown, and diagnostics.
- Harden MCP server execution with command allowlisting, environment scrubbing, identity reporting, and per-call timeouts.
- Make LSP, hooks, and skills report actionable diagnostics without breaking the core agent loop.
- Keep provider profiles separate from reusable transport implementations.
- Add offline doctor/TUI/replay diagnostics for extension readiness.

## Capabilities

### New Capabilities
- `extension-runtime-diagnostics`: Defines shared diagnostic contracts for extension startup, health, failures, and offline readiness.

### Modified Capabilities
- `mcp-lsp-extension-surface`: Adds hardened lifecycle, timeout, permission, and diagnostic requirements for MCP/LSP/hooks/skills.
- `provider-model-system`: Separates model/provider profiles from transport behavior and capability reporting.

## Impact

- Affected code: `internal/tool/mcp`, `internal/lsp`, `internal/hooks`, `internal/skills`, `internal/provider`, `internal/session`, `cmd/go-code`, `pkg/tui`.
- Affected tests: mock MCP server tests, LSP unavailable/healthy tests, invalid hook/skill tests, provider profile tests, doctor/replay diagnostics tests.
- No new dependencies unless a local fixture helper is already available in the repo or standard library.
