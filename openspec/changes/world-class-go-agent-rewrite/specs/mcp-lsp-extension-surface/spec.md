## ADDED Requirements

### Requirement: MCP servers are configurable and permission-aware
The system SHALL load configured MCP servers, register their tools, and apply the same permission policy used for built-in tools.

#### Scenario: MCP tool registration
- **WHEN** a configured MCP server starts and lists tools
- **THEN** the system registers namespaced MCP tools in the tool registry

#### Scenario: MCP tool requiring approval
- **WHEN** an MCP tool requests a write-like or external side effect
- **THEN** the system evaluates the operation through the permission policy before execution

### Requirement: Hooks have documented lifecycle guarantees
The system SHALL document and test pre-tool and post-tool hook execution order, error behavior, and trace integration.

#### Scenario: Pre-hook blocks execution
- **WHEN** a pre-tool hook returns an error
- **THEN** the target tool is not executed and the agent receives a tool error result

### Requirement: LSP features are exposed as optional capabilities
The system MUST expose LSP diagnostics, symbols, definitions, references, and hover only when an LSP server is configured and healthy.

#### Scenario: LSP unavailable
- **WHEN** no LSP server is configured
- **THEN** LSP commands or tools report unavailable status without failing the core agent workflow

### Requirement: Skills are discoverable and executable
The system SHALL load user-defined skills, list them, validate required fields, and execute selected skill prompts through the agent.

#### Scenario: List valid skills
- **WHEN** the user runs `/skills`
- **THEN** the system lists valid skills and skips invalid skill files with a non-fatal warning

