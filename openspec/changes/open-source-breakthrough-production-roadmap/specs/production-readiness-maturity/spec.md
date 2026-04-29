## ADDED Requirements

### Requirement: Maturity Levels Must Be Defined
The project SHALL define maturity levels for demo, usable, reliable, production-grade, and community-leading status.

#### Scenario: Release status is displayed
- **WHEN** README, docs, or release notes describe project maturity
- **THEN** they use the defined maturity levels and explain the evidence behind the current label

### Requirement: Production-Grade Claims Must Require Gates
The project SHALL require build, test, harness, docs, security, permission, packaging, and known-risk gates before claiming production-grade status.

#### Scenario: Maintainer prepares a production-grade release
- **WHEN** a release is labeled production-grade
- **THEN** the release evidence includes Go tests, Python harness results, docs build, security review or scan, permission/sandbox verification, packaged artifacts, and known-risk notes

### Requirement: Reliability Must Include Failure Behavior
The project SHALL treat graceful failure, actionable errors, recovery, and rollback notes as part of production readiness.

#### Scenario: Provider call fails
- **WHEN** provider auth, rate limit, timeout, streaming, or model errors occur
- **THEN** the user receives actionable errors and the failure mode is covered by tests or harness scenarios

#### Scenario: Tool execution is denied
- **WHEN** permission or sandbox policy blocks a tool call
- **THEN** the user receives a clear explanation and the session records the denial for replay/debugging

### Requirement: Security-Sensitive Surfaces Must Have Explicit Review
The project SHALL require explicit review and verification for shell execution, file writes, workspace boundaries, secrets handling, update downloads, telemetry, and MCP/tool extensions.

#### Scenario: A security-sensitive feature changes
- **WHEN** a change touches shell, filesystem, secrets, update, telemetry, or external tool execution behavior
- **THEN** the task includes security-focused tests or review notes before completion

### Requirement: Performance Baselines Must Be Tracked
The project SHALL track startup time, prompt-mode latency under mock provider, harness runtime, and memory-sensitive behavior before using performance claims publicly.

#### Scenario: Performance claim is added
- **WHEN** docs or release notes claim speed, low overhead, or production performance
- **THEN** the claim references a reproducible benchmark command and environment
