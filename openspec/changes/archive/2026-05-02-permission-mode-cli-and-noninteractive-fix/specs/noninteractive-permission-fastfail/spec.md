## ADDED Requirements

### Requirement: Non-interactive fast-fail
When stdin is not a terminal (non-interactive mode, e.g., `go-code -p "..."`), the system SHALL immediately deny any permission prompt instead of blocking on stdin read. The denial SHALL include the tool name, the required permission, and a suggestion for how to grant it.

#### Scenario: Non-interactive prompt denial
- **WHEN** agent runs in non-interactive mode (`-p`) AND a tool requires user approval
- **THEN** the permission is denied immediately with a message containing the tool name and required permission level

#### Scenario: Non-interactive with danger-full-access
- **WHEN** agent runs in non-interactive mode with `--permission-mode danger-full-access`
- **THEN** all tools are allowed without prompting

#### Scenario: Interactive mode prompt
- **WHEN** agent runs in interactive mode (stdin is a terminal)
- **THEN** permission prompts display normally and wait for user input

### Requirement: Clear error messaging
The non-interactive denial message SHALL include a suggestion to re-run with `--permission-mode danger-full-access` to grant all permissions.

#### Scenario: Denial message content
- **WHEN** a non-interactive permission prompt is denied
- **THEN** the error message includes "Re-run with --permission-mode danger-full-access"
