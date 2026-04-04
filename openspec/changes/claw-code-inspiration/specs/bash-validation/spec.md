## ADDED Requirements

### Requirement: Bash command validation
The system SHALL validate Bash commands before execution.

#### Scenario: Read-only commands auto-allowed
- **WHEN** a read-only command is executed (ls, cat, grep, find, wc, head, tail, echo, pwd, tree)
- **THEN** it is auto-allowed without user approval

#### Scenario: Dangerous commands blocked
- **WHEN** a dangerous command is detected (rm -rf /, curl | bash, sudo, dd, mkfs)
- **THEN** the command is blocked with a warning message

#### Scenario: Write commands validated
- **WHEN** a write command is executed (sed -i, awk with output, tee)
- **THEN** the target paths are validated against workspace boundaries

#### Scenario: Path injection prevented
- **WHEN** a command contains path traversal attempts
- **THEN** the command is blocked
