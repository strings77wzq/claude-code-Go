## ADDED Requirements

### Requirement: Safe permission defaults
The system SHALL start new users in a safe permission mode that allows read-only exploration and asks before file writes, edits, shell execution, network fetches, or destructive operations.

#### Scenario: First write attempt
- **WHEN** the agent attempts to write a file during a new session
- **THEN** the system prompts for approval before executing the write

#### Scenario: Read-only operation
- **WHEN** the agent reads a file inside the workspace
- **THEN** the system allows the operation without prompting

### Requirement: Approval prompt decisions are enforceable
The system MUST distinguish allow, deny, allow-once, and allow-for-session decisions.

#### Scenario: Denied shell command
- **WHEN** the user denies a proposed shell command
- **THEN** the command is not executed and a tool result explaining the denial is returned to the agent

#### Scenario: Remembered session approval
- **WHEN** the user approves a matching operation for the session
- **THEN** the system executes subsequent matching operations without repeating the prompt

### Requirement: Workspace boundary and command validation
The system MUST validate file paths and shell commands before execution according to the active permission mode.

#### Scenario: Path escapes workspace
- **WHEN** a write or edit resolves outside the configured workspace
- **THEN** the system blocks the operation even if the model requested it

#### Scenario: Destructive command in safe mode
- **WHEN** the agent proposes a command matching a destructive pattern in safe mode
- **THEN** the system requires explicit approval or denies the command according to policy

### Requirement: Permission audit trace
The system SHALL record permission decisions in the session trace without storing secrets.

#### Scenario: Audited approval
- **WHEN** the user approves a shell command
- **THEN** the trace records the tool name, normalized decision, command summary, and timestamp

