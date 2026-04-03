## ADDED Requirements

### Requirement: MCP server configuration
The system SHALL support MCP server configuration via settings files.

#### Scenario: Configure MCP server
- **WHEN** a settings.json file contains an "mcpServers" section
- **THEN** each server's command, args, and env are parsed and used for connection

#### Scenario: Environment variable interpolation
- **WHEN** an env value contains "${ENV_VAR}" syntax
- **THEN** it is replaced with the actual environment variable value

### Requirement: Stdio transport
The system SHALL implement stdio transport for MCP server communication.

#### Scenario: Start MCP server process
- **WHEN** an MCP server is initialized
- **THEN** a subprocess is started with the configured command and args

#### Scenario: JSON-RPC communication
- **WHEN** a JSON-RPC request is sent to an MCP server
- **THEN** the request is written to stdin and the response is read from stdout

### Requirement: MCP tool discovery
The system SHALL discover tools from connected MCP servers.

#### Scenario: List tools from server
- **WHEN** an MCP server is connected
- **THEN** the tools/list method is called and returned tools are registered

#### Scenario: Tool naming convention
- **WHEN** an MCP tool is registered
- **THEN** its name follows the format "mcp__<serverName>__<toolName>"

### Requirement: MCP tool adapter
The system SHALL adapt MCP remote tools to the local Tool interface.

#### Scenario: Execute MCP tool
- **WHEN** an MCP tool is executed
- **THEN** a JSON-RPC tools/call request is sent to the server and the result is returned

#### Scenario: MCP tool requires permission
- **WHEN** an MCP tool's RequiresPermission() is called
- **THEN** it always returns true (external tools always require approval)

### Requirement: MCP server lifecycle
The system SHALL manage MCP server process lifecycle.

#### Scenario: Graceful shutdown
- **WHEN** the application exits
- **THEN** all MCP server subprocesses are gracefully terminated

#### Scenario: Server connection failure
- **WHEN** an MCP server fails to start
- **THEN** an error is logged and other servers continue to work
