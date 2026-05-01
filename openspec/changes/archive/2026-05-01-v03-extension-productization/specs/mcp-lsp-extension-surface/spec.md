## ADDED Requirements

### Requirement: Extension status is diagnosable offline
The system SHALL report MCP, LSP, hooks, and skills status without requiring live provider credentials.

#### Scenario: Offline doctor checks extension readiness
- **WHEN** the user runs `go-code doctor --offline`
- **THEN** the output reports whether MCP config exists, whether configured servers are skipped or unavailable, whether LSP is configured, and whether skills/hooks directories are readable

### Requirement: MCP tool registration is testable with a local mock server
The system SHALL support deterministic MCP registration tests using a local mock server or fixture transport.

#### Scenario: Mock MCP server exposes a tool
- **WHEN** a local MCP fixture starts and lists one tool
- **THEN** the tool registry includes the namespaced MCP tool
- **AND** the tool result is traceable in session output

### Requirement: MCP tools use the same permission model as built-ins
The system MUST apply permission checks to MCP tools that read, write, execute commands, or access external resources.

#### Scenario: MCP tool requests a write-like action
- **WHEN** an MCP tool invocation is classified as write-like
- **THEN** the permission policy evaluates the action before execution
- **AND** denial is returned as a tool result without crashing the agent loop

### Requirement: LSP capabilities are only advertised when healthy
The system MUST expose LSP diagnostics, symbols, definitions, references, and hover only when an LSP server is configured and passes a health check.

#### Scenario: LSP is not configured
- **WHEN** no LSP server is configured
- **THEN** LSP commands or tools report unavailable status
- **AND** the core prompt, tools, and TUI continue to work

#### Scenario: LSP health check succeeds
- **WHEN** an LSP server is configured and responds to initialization
- **THEN** LSP capabilities are listed as available
- **AND** diagnostics include the server identity and workspace root

### Requirement: Hooks and skills failures are non-fatal by default
The system SHALL surface invalid hooks or skills as warnings unless the configured lifecycle requires blocking execution.

#### Scenario: Invalid skill file is discovered
- **WHEN** the skills loader reads an invalid skill file
- **THEN** the valid skills remain available
- **AND** the invalid file is reported with an actionable warning
