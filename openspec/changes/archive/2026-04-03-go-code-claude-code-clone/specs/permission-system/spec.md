## ADDED Requirements

### Requirement: Permission modes
The system SHALL support three permission modes for controlling tool execution.

#### Scenario: ReadOnly mode
- **WHEN** the active mode is ReadOnly
- **THEN** only read-only tools (Read, Glob, Grep) are auto-allowed; all others require explicit approval

#### Scenario: WorkspaceWrite mode
- **WHEN** the active mode is WorkspaceWrite
- **THEN** read-only and file-write tools (Read, Write, Edit, Glob, Grep) are auto-allowed; Bash requires approval

#### Scenario: DangerFullAccess mode
- **WHEN** the active mode is DangerFullAccess
- **THEN** all tools including Bash are auto-allowed

### Requirement: Permission evaluation
The system SHALL evaluate permissions before each tool execution.

#### Scenario: Auto-allow read-only tools
- **WHEN** a tool with RequiresPermission()=false is executed
- **THEN** the tool is allowed without user interaction

#### Scenario: Ask for write tools
- **WHEN** a tool with RequiresPermission()=true is executed and no rules match
- **THEN** the user is prompted to approve or deny

### Requirement: Interactive permission prompt
The system SHALL prompt the user for permission decisions.

#### Scenario: User approves
- **WHEN** the user enters "y" or "yes"
- **THEN** the tool execution is allowed for this single call

#### Scenario: User denies
- **WHEN** the user enters "n" or "no"
- **THEN** the tool execution is denied and an error result is returned to the model

#### Scenario: User approves always
- **WHEN** the user enters "a" or "always"
- **THEN** the tool execution is allowed and the decision is remembered for the session

### Requirement: Session memory
The system SHALL remember "always" decisions for the current session.

#### Scenario: Remembered decision
- **WHEN** a tool was previously approved with "always"
- **THEN** subsequent calls to the same tool are auto-allowed without prompting

### Requirement: Rule-based permissions
The system SHALL support glob-based rules for fine-grained permission control.

#### Scenario: Bash command rule
- **WHEN** a rule "bash(git:*)" is configured
- **THEN** "git commit" is auto-allowed but "rm -rf /" still requires approval

#### Scenario: File path rule
- **WHEN** a rule "read(/tmp:*)" is configured
- **THEN** reading files under /tmp/ is auto-allowed

### Requirement: Deny rules override
The system SHALL prioritize deny rules over all other permission decisions.

#### Scenario: Deny overrides allow
- **WHEN** a deny rule matches a tool call that would otherwise be allowed
- **THEN** the tool call is denied
