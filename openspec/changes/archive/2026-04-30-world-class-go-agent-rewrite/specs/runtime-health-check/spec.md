## ADDED Requirements

### Requirement: Doctor validates local readiness
The system SHALL provide a `doctor` command that validates whether the local installation can run a basic agent session.

#### Scenario: Healthy installation
- **WHEN** the user runs `go-code doctor` with valid configuration and accessible runtime directories
- **THEN** the command exits successfully and reports passing checks for binary version, config source, provider/model configuration, session directory, tool availability, and documentation links

#### Scenario: Missing API key
- **WHEN** the user runs `go-code doctor` without a usable provider API key
- **THEN** the command exits with a non-zero status and prints the exact configuration source to fix

### Requirement: Doctor provides actionable remediation
The system MUST report failed checks with concrete next commands or documentation paths.

#### Scenario: Invalid session directory
- **WHEN** the configured session directory cannot be created or written
- **THEN** the doctor output identifies the path, the failing operation, and a remediation command or configuration key

#### Scenario: Provider probe disabled
- **WHEN** the user runs doctor in offline or no-network mode
- **THEN** local checks still run and provider checks are marked as skipped rather than failed

