## MODIFIED Requirements

### Requirement: Bash semantic validation
The system SHALL validate Bash commands at the semantic level before shell execution. The validation behavior MUST be covered by deterministic regression tests for read-only verification, destructive command detection, sed/awk write targets, redirects, subshells, command chaining, and workspace path boundaries.

#### Scenario: Read-only verification
- **WHEN** a command uses only supported read-only operations and workspace-safe paths
- **THEN** it is verified to not modify the filesystem

#### Scenario: Destructive command warning
- **WHEN** a destructive command is detected (rm, mv, cp overwrite, privilege escalation, device write, process kill, or remote script execution)
- **THEN** validation rejects the command with a reason that identifies the unsafe pattern

#### Scenario: sed/awk write validation
- **WHEN** sed -i, sed output redirection, or awk output redirection is used
- **THEN** the target path is extracted and validated against workspace boundaries

#### Scenario: Command semantics analysis
- **WHEN** a command contains pipes, redirects, subshells, command substitution, or command chaining
- **THEN** the full command structure is analyzed for safety before the command is allowed

#### Scenario: Workspace escape blocked
- **WHEN** a command references a write target or validated path outside the configured working directory
- **THEN** validation rejects the command and reports the workspace boundary violation

#### Scenario: Dangerous subshell blocked
- **WHEN** command substitution or a subshell contains destructive command content
- **THEN** validation rejects the command before execution
