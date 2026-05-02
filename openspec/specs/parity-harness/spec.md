## Purpose
Define the deterministic parity harness requirements used to prove core Claude Code-style workflows, CI gates, and documented parity status.
## Requirements
### Requirement: Deterministic mock provider scenarios
The system SHALL include deterministic mock-provider scenarios for streaming text, tool use, edit flows, bash flows, permission denial, retries, context pressure, and malformed provider responses.

#### Scenario: Tool-use loop
- **WHEN** the parity harness runs a scenario where the mock provider requests a file read and then returns a final answer
- **THEN** the agent executes the read, sends the tool result, and completes with the expected final answer

#### Scenario: Malformed stream event
- **WHEN** the mock provider emits a malformed streaming event
- **THEN** the agent returns a classified error or recovery result without panicking

### Requirement: CI runs parity gates
The system MUST run the deterministic parity harness in CI for pull requests that affect runtime, provider, tool, permission, session, or TUI code.

#### Scenario: Harness failure
- **WHEN** a parity scenario fails in CI
- **THEN** the workflow fails and publishes enough logs to identify the failing scenario

### Requirement: Parity status is documented
The system SHALL maintain a parity matrix that maps important Claude Code-style workflows to status, tests, docs, and known gaps.

#### Scenario: Unsupported workflow
- **WHEN** a workflow is not implemented
- **THEN** the parity matrix marks it as unsupported or planned instead of implying support

### Requirement: Harness supports normalized comparison evidence
The harness SHALL support normalized evidence records that allow maintainers to compare this agent's behavior with manually supplied runs from other coding agents.

#### Scenario: External agent evidence is imported
- **WHEN** a maintainer supplies a normalized external run record
- **THEN** the harness report compares outcome, duration, tool count, and notes without requiring external credentials

### Requirement: Comparison reports avoid unsupported claims
The harness MUST distinguish measured local evidence from manual notes or inferred comparisons.

#### Scenario: Report includes manual competitor notes
- **WHEN** a comparison report includes manually entered notes about another agent
- **THEN** the notes are labeled as manual evidence rather than automated measurement
