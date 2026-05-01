## Why

Four overlapping active changes create implementation confusion, while the core agent engine—already clean and tested—can't ship because 9 packages lack test coverage, the provider registry is stale, and public docs over-claim unsupported features. This change consolidates the roadmap, fills the critical gaps, and delivers a credible v0.2 release.

## What Changes

- Archive `recenter-claudecodego-agent-roadmap`, `open-source-breakthrough-production-roadmap`, and `premium-project-upgrade`—extracting strategic insights into project specs, removing overlapping task lists.
- Rescope `world-class-go-agent-rewrite` sections 8–11 to focus on v0.2 essentials: test coverage, provider alignment, docs truth, and release gate. Defer MCP/LSP productization to v0.3.
- Add test coverage to 9 packages with zero tests: `internal/lsp`, `internal/tool`, `internal/tool/init`, `internal/tool/mcp`, `internal/telemetry`, `internal/update`, `internal/provider/anthropic`, `internal/provider/openai`, `pkg/tui`.
- Update model registry from stale names (`deepseek-chat`, `deepseek-reasoner`) to current names (`deepseek-v4-pro`, `deepseek-v4-flash`). Add MiMo-V2.5 profile. Support dynamic unknown-model passthrough.
- Remove or relabel unsupported claims in README, docs, and PARITY.md. Replace placeholder demo GIF with honest status. Sync English and Chinese docs.
- Deliver v0.2 release gate: `go test ./...` passes all packages, harness scenarios pass, docs build passes, PARITY.md reflects verified status only.

## Capabilities

### New Capabilities

- `test-coverage-gap-fill`: At least one happy-path and one error-path test for each of the 9 zero-test packages, with `go test ./...` passing on all.
- `provider-registry-alignment`: Model registry updated to current model names (DeepSeek v4 series, MiMo-V2.5), unknown-model passthrough, and config-driven model discovery.
- `docs-truth-alignment`: README, docs/zh, and PARITY.md stripped of unsupported claims. Placeholder demo, testimonials, and metrics removed or clearly labeled.
- `v02-release-gate`: All Go tests pass, harness scenarios pass, docs build passes, PARITY.md reflects verified-only status, and a changelog entry exists.

### Modified Capabilities

None. `openspec/specs/` is empty—this change introduces new capability contracts.

## Impact

- Affected code: `internal/lsp/`, `internal/tool/`, `internal/tool/init/`, `internal/tool/mcp/`, `internal/telemetry/`, `internal/update/`, `internal/provider/anthropic/`, `internal/provider/openai/`, `internal/provider/registry/`, `pkg/tui/`.
- Affected harness: `harness/test_scenarios.py`, `harness/mock_server/scenarios.py`.
- Affected docs: `README.md`, `PARITY.md`, `docs/`, `docs/zh/`, `CHANGELOG.md`.
- Affected project: 3 changes archived, 1 change rescoped. OpenSpec change count goes from 4 to 1.
