## ADDED Requirements

### Requirement: PermissionEnforcer module
The system SHALL have an independent permission enforcement module.

#### Scenario: Tool-level permission labels
- **WHEN** a tool is registered
- **THEN** it declares its required permission level (ReadOnly, WorkspaceWrite, DangerFullAccess)

#### Scenario: Enforcer evaluation
- **WHEN** a tool execution is requested
- **THEN** the PermissionEnforcer evaluates: tool-level label → bash validation → file boundary → policy
