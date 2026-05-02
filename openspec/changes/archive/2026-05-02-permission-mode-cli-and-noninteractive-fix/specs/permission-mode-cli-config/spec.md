## ADDED Requirements

### Requirement: CLI permission mode flag
The system SHALL accept a `--permission-mode` flag with values `read-only`, `workspace-write`, or `danger-full-access`. The mode SHALL default to `workspace-write` when not specified. Invalid values SHALL produce a clear error message listing valid options.

#### Scenario: Valid mode specified
- **WHEN** user runs `go-code --permission-mode read-only`
- **THEN** the agent starts with `ReadOnly` permission mode

#### Scenario: Default mode
- **WHEN** user runs `go-code` without `--permission-mode`
- **THEN** the agent starts with `WorkspaceWrite` permission mode

#### Scenario: Invalid mode
- **WHEN** user runs `go-code --permission-mode full-access`
- **THEN** the process exits with a non-zero code and an error message listing `read-only`, `workspace-write`, `danger-full-access`

### Requirement: Environment variable fallback
The system SHALL load the permission mode from the `GO_CODE_PERMISSION_MODE` environment variable when the CLI flag is not set. The flag takes priority over the env var.

#### Scenario: Env var sets mode
- **WHEN** `GO_CODE_PERMISSION_MODE=read-only` is set and no `--permission-mode` flag is given
- **THEN** the agent starts with `ReadOnly` permission mode

#### Scenario: Flag overrides env var
- **WHEN** `GO_CODE_PERMISSION_MODE=read-only` is set AND `--permission-mode danger-full-access` is given
- **THEN** the agent starts with `DangerFullAccess` permission mode
