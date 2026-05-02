## ADDED Requirements

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
