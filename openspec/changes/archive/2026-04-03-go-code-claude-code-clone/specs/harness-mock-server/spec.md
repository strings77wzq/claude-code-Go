## ADDED Requirements

### Requirement: Mock Anthropic API
The harness SHALL provide a mock Anthropic API service that implements the Messages API protocol.

#### Scenario: Mock streaming text response
- **WHEN** a streaming request is sent to the mock server
- **THEN** the server returns SSE events (message_start, content_block_delta, message_delta) simulating a text response

#### Scenario: Mock tool_use response
- **WHEN** a request is sent that should trigger tool use
- **THEN** the server returns SSE events including a tool_use content block

#### Scenario: Mock multi-turn response
- **WHEN** multiple requests are sent in sequence
- **THEN** the mock server returns different responses based on the conversation context

### Requirement: Request recording
The mock server SHALL record all incoming requests for assertion in tests.

#### Scenario: Record requests
- **WHEN** the Go CLI sends requests to the mock server
- **THEN** each request is stored and can be queried by test scenarios

### Requirement: Configurable scenarios
The mock server SHALL support pre-configured response scenarios.

#### Scenario: Configure response sequence
- **WHEN** a test scenario is configured
- **THEN** the mock server returns the specified sequence of responses

### Requirement: Server lifecycle
The mock server SHALL support starting and stopping for test isolation.

#### Scenario: Start and stop
- **WHEN** a test starts the mock server
- **THEN** it binds to an available port and provides the base URL for test configuration
