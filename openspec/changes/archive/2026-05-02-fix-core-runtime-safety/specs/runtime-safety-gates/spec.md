## ADDED Requirements

### Requirement: Tool execution is panic-contained
The system MUST recover from panics raised by registered tools and return a structured tool error without terminating the agent process.

#### Scenario: Tool panics during execution
- **WHEN** a registered tool panics while handling an invocation
- **THEN** the agent receives a tool error result containing the tool name and a safe failure message
- **AND** the process continues running

### Requirement: Runtime requests have explicit terminal states
The system SHALL track each agent request through started, completed, cancelled, or failed states.

#### Scenario: Request is cancelled
- **WHEN** a user cancels an in-flight agent request
- **THEN** the request context is cancelled
- **AND** the UI or CLI records a cancelled terminal state instead of remaining loading

### Requirement: Runtime safety events are traceable
The system SHALL emit trace events for request cancellation, tool panic recovery, permission denial, and provider failure.

#### Scenario: Tool panic is recovered
- **WHEN** a tool panic is converted into a structured tool error
- **THEN** the trace records the event type, tool name, request id, and redacted error summary

#### Scenario: Request is cancelled
- **WHEN** a request context is cancelled before completion
- **THEN** the trace records a runtime event subtype for cancellation and the request terminal state is cancelled
