## ADDED Requirements

### Requirement: Tools are categorized into tiers for progressive disclosure

The system SHALL categorize all registered tools into three tiers: core, extension, and MCP. Core tools SHALL always be included in API requests. Extension tools SHALL be sent on request turn 1 then omitted; the model may request them back. MCP tools SHALL be included when their server is connected and registered.

#### Scenario: Core tools always sent
- **WHEN** an API request is built
- **THEN** all core-tier tools (Read, Write, Edit, Bash, Grep) are included in the tool definitions

#### Scenario: Extension tools sent on first turn only
- **WHEN** the first API request of a session is built
- **THEN** core AND extension tools are included
- **WHEN** subsequent API requests are built and the model has not requested extension tools
- **THEN** only core tools are included

#### Scenario: Model requests extension tool
- **WHEN** the model attempts to use an extension tool not in the current request
- **THEN** the next API request includes that specific extension tool definition

#### Scenario: MCP tools included when registered
- **WHEN** an MCP server registers its tools
- **THEN** those tools are included in the tool definitions for all subsequent requests

### Requirement: Tool tier metadata is declarative

Each builtin tool SHALL declare its tier (core or extension) via a method on the Tool interface. MCP tools SHALL be automatically assigned to the MCP tier.

#### Scenario: Tool declares its tier
- **WHEN** a builtin tool's Tier() method is called
- **THEN** it returns the tier constant matching its intended tier assignment
