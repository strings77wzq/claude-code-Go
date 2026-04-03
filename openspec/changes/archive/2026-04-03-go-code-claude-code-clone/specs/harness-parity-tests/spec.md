## ADDED Requirements

### Requirement: Streaming text parity test
The harness SHALL verify that streaming text output works correctly end-to-end.

#### Scenario: Streaming text delivery
- **WHEN** the Go CLI sends a text query to the mock server
- **THEN** the output contains the complete streamed text with no corruption

### Requirement: Tool roundtrip parity test
The harness SHALL verify that tool call → execution → result → response cycles work correctly.

#### Scenario: Read file roundtrip
- **WHEN** the CLI is asked to read a file
- **THEN** the tool_use block is generated, the Read tool executes, and the model summarizes the content

#### Scenario: Bash command roundtrip
- **WHEN** the CLI is asked to run a command
- **THEN** the Bash tool executes (with auto-approval in test mode) and the model reports the output

### Requirement: Permission flow parity test
The harness SHALL verify that permission allow/deny flows work correctly.

#### Scenario: Permission allowed
- **WHEN** a tool requiring permission is auto-approved in test mode
- **THEN** the tool executes and the result is returned to the model

#### Scenario: Permission denied
- **WHEN** a tool requiring permission is denied in test mode
- **THEN** an error result is returned to the model and the model adapts its response

### Requirement: MCP integration parity test
The harness SHALL verify that MCP server integration works correctly.

#### Scenario: MCP tool discovery
- **WHEN** an MCP server is configured
- **THEN** its tools are discovered and registered with the mcp__ prefix

#### Scenario: MCP tool execution
- **WHEN** an MCP tool is executed
- **THEN** the JSON-RPC call succeeds and the result is returned

### Requirement: Edit tool parity test
The harness SHALL verify that the Edit tool's uniqueness check works correctly.

#### Scenario: Unique edit succeeds
- **WHEN** the Edit tool is called with a unique old_string
- **THEN** the file is modified correctly

#### Scenario: Non-unique edit fails
- **WHEN** the Edit tool is called with a non-unique old_string
- **THEN** an error is returned indicating multiple matches
