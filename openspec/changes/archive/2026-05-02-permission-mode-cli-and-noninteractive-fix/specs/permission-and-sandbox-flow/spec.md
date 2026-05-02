## MODIFIED Requirements

### Requirement: Safe permission defaults
The system SHALL start new users in a safe permission mode that allows read-only exploration and asks before file writes, edits, shell execution, network fetches, or destructive operations. The starting mode SHALL be `workspace-write` by default and SHALL be overridable via `--permission-mode` flag or `GO_CODE_PERMISSION_MODE` environment variable.

#### Scenario: First write attempt
- **WHEN** the agent attempts to write a file during a new session
- **THEN** the system prompts for approval before executing the write

#### Scenario: Read-only operation
- **WHEN** the agent reads a file inside the workspace
- **THEN** the system allows the operation without prompting

#### Scenario: CLI flag overrides default
- **WHEN** user starts agent with `--permission-mode danger-full-access`
- **THEN** the agent runs with `DangerFullAccess` mode and skips all permission prompts

#### Scenario: Invalid mode rejected
- **WHEN** user specifies an unrecognized permission mode value
- **THEN** the agent exits with a non-zero code and lists valid values

### Requirement: Non-interactive approvals fail closed
The system MUST deny operations that require approval when no interactive approval channel is available. Denial SHALL be immediate (no blocking on stdin) and SHALL include the tool name, required permission, and remediation suggestion.

#### Scenario: Write requires approval in CI
- **WHEN** the agent attempts a write operation in a non-interactive run
- **THEN** the system denies the operation with an agent-visible permission result
- **AND** the system does not block waiting for stdin

#### Scenario: Denial includes remediation hint
- **WHEN** a permission is denied in non-interactive mode
- **THEN** the error message includes a suggestion to re-run with `--permission-mode danger-full-access`

#### Scenario: Non-interactive with elevated mode succeeds
- **WHEN** agent runs non-interactively with `--permission-mode danger-full-access`
- **THEN** operations that would normally prompt are allowed without blocking
