## Why

The project needs evidence that the agent is actually useful compared with Codex- and Claude-style coding agents, not only that unit tests pass. The roadmap audit found that harness scenarios are shallow, latency is not tracked, and release evidence does not yet prove real-world task quality.

## What Changes

- Introduce scenario manifests with task intent, workspace setup, allowed tools, assertions, latency budgets, and trace requirements.
- Add agent quality gates for real coding workflows: inspect, edit, test, recover, and explain.
- Record pass/fail, latency, tool count, permission decisions, and trace evidence for each scenario.
- Add comparison-oriented workflows that let maintainers evaluate this agent against other local coding agents without hardcoding competitor internals.

## Capabilities

### New Capabilities
- `agent-quality-gates`: Defines product-quality scenario gates for agent behavior, latency, traceability, and release evidence.

### Modified Capabilities
- `parity-harness`: Adds manifest-driven scenario execution and comparison evidence.
- `harness-engineering-gate`: Adds latency, trace, and real-workflow assertions beyond binary pass/fail.

## Impact

- Affected code: `harness`, `internal/session`, `internal/tool`, `cmd/go-code`, CI scripts, release evidence docs.
- Affected tests: harness manifest parser tests, scenario runner tests, trace assertion tests, latency budget checks.
- No new provider credentials should be required for deterministic offline harness scenarios.
