## 1. Manifest Contract

- [x] 1.1 Define a scenario manifest schema covering prompt, setup, allowed tools, assertions, trace expectations, and budgets.
- [x] 1.2 Add parser and validation tests for valid and invalid manifests.
- [x] 1.3 Document the manifest fields in harness developer docs.
- [x] 1.4 Define manifest file location, schema version, and golden invalid-case fixtures.

## 2. Runner Evidence

- [x] 2.1 Record pass/fail status, duration, tool count, permission decisions, trace path, and failure reason for each scenario.
- [x] 2.2 Add trace assertion support for required event types and redaction status.
- [x] 2.3 Add latency budget reporting that distinguishes functional failures from budget violations.
- [x] 2.4 Ensure evidence output is deterministic enough for CI artifacts.
- [x] 2.5 Run each scenario with temporary HOME, workspace, and explicit environment allowlist.
- [x] 2.6 Add evidence artifact redaction tests.

## 3. Scenario Coverage

- [x] 3.1 Convert existing shallow harness cases into manifests.
- [x] 3.2 Add repository inspection scenario.
- [x] 3.3 Add safe edit plus test execution scenario.
- [x] 3.4 Add permission denial and recovery scenario.
- [x] 3.5 Add provider/tool failure recovery scenario.
- [x] 3.6 Add user-facing explanation scenario.

## 4. Comparison Workflow

- [x] 4.1 Define normalized evidence fields for manually supplied external-agent runs.
- [x] 4.2 Add report generation that labels measured, manual, and inferred evidence separately.
- [x] 4.3 Document a local trial procedure for comparing this agent with Codex- and Claude-style agents.

## 5. CI And Release Integration

- [x] 5.1 Add a deterministic harness command suitable for CI.
- [x] 5.2 Store or print release evidence paths for scenario reports and traces.
- [x] 5.3 Add release checklist entries requiring harness gate results.

## 6. Verification

- [x] 6.1 Run harness parser, runner, and report tests.
- [x] 6.2 Run deterministic harness scenarios locally.
- [x] 6.3 Run `go test ./...`.
- [x] 6.4 Run `openspec validate harness-agent-quality-gates --strict`.

## Evidence

- `python -m pytest harness/test_quality_gates.py -q`
- `python -m pytest harness/ -q`
- `./scripts/run-harness.sh`
- `go test ./...`
- `openspec validate harness-agent-quality-gates --strict --json --no-interactive`

## Residual Risks

- Manifest evaluation currently validates supplied run evidence rather than executing every manifest end-to-end through a live provider; deterministic CI remains mock-provider based.
- External Codex/Claude-style comparison is normalized and labeled, but competitor runs are intentionally manual evidence to avoid credential and maintenance coupling.
