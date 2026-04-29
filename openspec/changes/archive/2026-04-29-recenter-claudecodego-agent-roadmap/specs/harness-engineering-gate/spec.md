## ADDED Requirements

### Requirement: Harness Must Gate Supported Claims
The project SHALL require deterministic verification before marking an agent workflow, provider, or tool behavior as supported in public docs.

#### Scenario: Documentation claims support
- **WHEN** docs claim a workflow or provider is supported
- **THEN** the claim links to or names at least one Go test, harness scenario, parity row, or manual verification note

#### Scenario: A feature lacks verification
- **WHEN** a feature exists in code but lacks tests or harness coverage
- **THEN** public docs mark it as partial, experimental, or planned rather than fully supported

### Requirement: Harness Must Cover Core Agent Trajectories
The Python harness SHALL include deterministic scenarios for provider streaming, tool calls, file edits, bash permission decisions, session persistence, replay, recoverable provider failures, and model/profile selection.

#### Scenario: Core runtime changes
- **WHEN** code changes the agent loop, provider adapter, permission layer, session layer, or built-in tools
- **THEN** the relevant harness scenarios run locally and in CI before completion is claimed

### Requirement: Mock Providers Must Avoid Real API Dependency
The harness SHALL verify expected request/response behavior through mock providers by default, without requiring live API keys, external network calls, or paid token usage.

#### Scenario: Contributor runs harness locally
- **WHEN** a contributor runs `./scripts/run-harness.sh`
- **THEN** the command completes using local mock provider scenarios and no real provider credentials

### Requirement: Release Verification Must Be One Command
The project SHALL provide a documented release verification path that runs Go tests and Python harness tests with a small number of commands suitable for CI and contributors.

#### Scenario: Maintainer checks release readiness
- **WHEN** a maintainer prepares a release or public milestone
- **THEN** `go test ./...` and `./scripts/run-harness.sh` pass or the release notes include explicit known failures
