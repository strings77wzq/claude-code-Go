## Purpose
Define the permission and sandbox behavior required to keep local agent execution safe, enforceable, workspace-bounded, and auditable.
## Requirements
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

### Requirement: Non-interactive approvals fail closed
The system MUST deny operations that require approval when no interactive approval channel is available.

#### Scenario: Write requires approval in CI
- **WHEN** the agent attempts a write operation in a non-interactive run
- **THEN** the system denies the operation with an agent-visible permission result
- **AND** the system does not block waiting for stdin

### Requirement: Permission mode hierarchy is explicit
The system MUST evaluate tool requirements against a documented permission mode hierarchy and action matrix.

#### Scenario: Higher mode satisfies lower requirement
- **WHEN** a tool requires workspace-write permission and the active mode is a higher-trust mode
- **THEN** the policy evaluates the operation using the hierarchy instead of simple string equality

#### Scenario: Explicit deny still wins
- **WHEN** a remembered or policy-level deny matches an operation
- **THEN** the operation is denied even if the active permission mode would otherwise allow it

### Requirement: Permission decisions include stable reasons
The system SHALL return stable reason codes for allow, deny, and prompt-required decisions.

#### Scenario: Approval cannot be collected
- **WHEN** an operation requires approval but the runtime is non-interactive
- **THEN** the decision reason identifies that approval was unavailable

#### Scenario: Active mode is insufficient
- **WHEN** an operation requires a higher permission mode than the active mode provides
- **THEN** the decision reason identifies insufficient mode
