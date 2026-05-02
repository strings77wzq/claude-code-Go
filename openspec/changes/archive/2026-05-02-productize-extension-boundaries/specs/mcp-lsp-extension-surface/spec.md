## ADDED Requirements

### Requirement: MCP launch is policy-validated
The system MUST validate configured MCP server command, arguments, working directory, and environment before starting the server.

#### Scenario: MCP command is not allowed
- **WHEN** an MCP server configuration references a command outside the allowed launch policy
- **THEN** the server is not started
- **AND** a diagnostic explains the blocked launch

#### Scenario: Existing config omits launch policy
- **WHEN** an existing MCP config does not declare a launch policy
- **THEN** the system applies the documented default allowlist, cwd constraints, and env inheritance rules

### Requirement: MCP lifecycle has bounded timeouts
The system MUST enforce startup, list-tools, tool-call, and shutdown timeouts for MCP servers.

#### Scenario: MCP tool call exceeds timeout
- **WHEN** an MCP tool call exceeds its configured timeout
- **THEN** the call is cancelled
- **AND** the agent receives a structured tool error result

#### Scenario: List tools blocks
- **WHEN** an MCP server does not respond to list-tools before the configured timeout
- **THEN** the server is marked unavailable and core agent startup continues

### Requirement: MCP tools are safety-classified
The system MUST classify MCP tools into the shared permission action matrix before execution.

#### Scenario: MCP tool is write-like
- **WHEN** an MCP tool is classified as write-like
- **THEN** the shared permission policy evaluates it as a side-effecting operation before execution

### Requirement: MCP environments are scrubbed before trace output
The system MUST redact configured MCP environment variables that contain secrets before writing diagnostics, trace, or replay output.

#### Scenario: MCP config includes a token
- **WHEN** MCP diagnostics include environment metadata
- **THEN** token-like values are replaced with a redacted marker

### Requirement: Optional extensions degrade gracefully
The system SHALL keep the core prompt, built-in tools, and TUI usable when MCP, LSP, hooks, or skills are unavailable.

#### Scenario: LSP server is not configured
- **WHEN** no LSP server is configured
- **THEN** LSP capabilities are marked unavailable
- **AND** the core agent workflow continues
