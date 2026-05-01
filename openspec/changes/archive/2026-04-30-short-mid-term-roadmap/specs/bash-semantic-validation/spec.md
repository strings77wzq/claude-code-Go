## ADDED Requirements

### Requirement: Bash semantic validation
The system SHALL validate Bash commands at the semantic level.

#### Scenario: Read-only verification
- **WHEN** a read-only command is executed
- **THEN** it is verified to not modify the filesystem

#### Scenario: Destructive command warning
- **WHEN** a destructive command is detected (rm, mv, cp overwrite)
- **THEN** a warning is shown and explicit approval is required

#### Scenario: sed/awk write validation
- **WHEN** sed -i or awk with output redirection is used
- **THEN** the target path is validated against workspace boundaries

#### Scenario: Command semantics analysis
- **WHEN** a command contains pipes, redirects, or subshells
- **THEN** the full command chain is analyzed for safety
