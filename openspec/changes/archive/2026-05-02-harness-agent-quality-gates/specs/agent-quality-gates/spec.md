## ADDED Requirements

### Requirement: Agent quality scenarios are manifest-driven
The system SHALL define agent quality scenarios in reviewable manifests that include task prompt, workspace setup, allowed tools, assertions, trace expectations, and budgets.

#### Scenario: Maintainer reviews a quality scenario
- **WHEN** a maintainer opens a scenario manifest
- **THEN** the task intent, setup, expected outcomes, allowed tools, and budgets are visible without reading runner code

#### Scenario: Invalid manifest is rejected
- **WHEN** a scenario manifest has an unsupported schema version or missing required field
- **THEN** the harness rejects it with an actionable validation error

### Requirement: Scenario results include product evidence
The harness MUST record pass/fail status, duration, tool count, permission decisions, trace file path, and failure reason for each scenario.

#### Scenario: Scenario fails a trace assertion
- **WHEN** a scenario completes but required trace evidence is missing
- **THEN** the scenario result fails with the missing assertion identified

#### Scenario: Evidence contains a secret
- **WHEN** scenario output, trace, or environment metadata contains token-like values
- **THEN** the persisted evidence replaces them with a redacted marker

### Requirement: Scenarios run in isolated environments
The harness MUST run each scenario with a temporary HOME, temporary workspace, and explicit environment allowlist.

#### Scenario: Host environment has credentials
- **WHEN** the host process has provider or Git credentials in its environment
- **THEN** the scenario run does not inherit those variables unless the manifest explicitly allows them

### Requirement: Release gates cover real agent workflows
The release gate SHALL include scenarios for repository inspection, safe edit, test execution, failure recovery, and user-facing explanation.

#### Scenario: Safe edit workflow
- **WHEN** the agent completes a safe edit scenario
- **THEN** the expected file changes and tests are present
- **AND** the trace records the tools and permission decisions used
