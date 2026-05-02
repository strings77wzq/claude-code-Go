## ADDED Requirements

### Requirement: Extensions report shared diagnostics
The system SHALL represent MCP, LSP, hook, skill, and provider-profile readiness issues using a shared diagnostic structure.

#### Scenario: Invalid skill is discovered
- **WHEN** the skills loader finds an invalid skill file
- **THEN** the diagnostic includes the component, severity, stable code, summary, and redacted file metadata

### Requirement: Offline extension diagnostics are available
The system MUST report extension readiness without requiring live provider credentials.

#### Scenario: Offline doctor checks extensions
- **WHEN** the user runs an offline doctor or status check
- **THEN** the output reports MCP, LSP, hooks, skills, and provider profile readiness from local configuration

### Requirement: Extension diagnostics are replayable
The system SHALL include extension diagnostics in session trace and replay output when they affect available tools or commands.

#### Scenario: MCP server startup times out
- **WHEN** an MCP server fails to start before its timeout
- **THEN** replay output includes a redacted diagnostic explaining the unavailable server
