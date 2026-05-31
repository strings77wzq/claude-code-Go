## MODIFIED Requirements

### Requirement: Bash semantic validation

The system SHALL validate Bash commands at the semantic level. Write indicators SHALL only include operators that actually produce file modifications (redirection operators `>`, `>>`, command substitution `$(` and backtick, and the `tee` command). Command composition operators (`|`, `;`, `&&`, `||`) SHALL NOT be classified as write indicators.

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

#### Scenario: Piped read-only commands are correctly classified
- **WHEN** a command uses only pipes with read-only commands (e.g., `ls | grep foo | wc -l`)
- **THEN** the command is classified as read-only, not as a write operation

#### Scenario: Command chains with read-only commands are correctly classified
- **WHEN** a command uses `&&` or `;` to chain read-only commands
- **THEN** the command is not falsely classified as containing write indicators

### Requirement: hasWriteArguments has no dead code paths

The `hasWriteArguments` function SHALL NOT contain unreachable conditional branches. All code paths SHALL be reachable with valid input. The function SHALL use exact field matching against write commands.

#### Scenario: Write command detected in arguments
- **WHEN** a read-only base command has a write command as an argument (e.g., `echo cp`)
- **THEN** the write command is detected via exact field match

#### Scenario: No dead code paths exist
- **WHEN** `hasWriteArguments` is analyzed with static analysis
- **THEN** no unreachable code branches are reported
