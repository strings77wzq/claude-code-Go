## ADDED Requirements

### Requirement: Prompt Caching reduces API cost and latency

The system SHALL support Anthropic Prompt Caching by placing cache breakpoints on the system prompt and tool definition block. Cached content SHALL be reused across consecutive turns within a session. The system SHALL track cache hit/miss metrics and include them in session trace output.

#### Scenario: System prompt is cached across turns
- **WHEN** the agent makes multiple API requests in a session
- **THEN** the system prompt is sent with a cache breakpoint and reused across turns

#### Scenario: Tool definitions are cached
- **WHEN** tool definitions are sent to the model
- **THEN** a cache breakpoint is placed before the tool definitions block

#### Scenario: Cache metrics are recorded
- **WHEN** an API response includes cache metrics (cache_creation_input_tokens, cache_read_input_tokens)
- **THEN** the metrics are recorded in the session trace

#### Scenario: Cache invalidates on tool registry change
- **WHEN** a new MCP server connects and registers tools
- **THEN** the next API request sends fresh (uncached) tool definitions
