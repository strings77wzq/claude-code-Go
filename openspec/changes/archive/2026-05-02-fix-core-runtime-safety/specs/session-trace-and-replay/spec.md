## ADDED Requirements

### Requirement: Trace events use versioned envelopes
The system MUST write session trace events using a versioned envelope that includes event type, timestamp, request id when available, and redaction status.

#### Scenario: Permission denial is traced
- **WHEN** a permission policy denies a tool invocation
- **THEN** the trace event includes the schema version, event type, tool name, decision reason, and redaction status

#### Scenario: Request id is available
- **WHEN** a runtime event is associated with a specific request id
- **THEN** the trace event includes that request id in the versioned envelope

### Requirement: Resume preserves coherent request history
The system SHALL preserve enough prior message and tool result history for resumed TUI and CLI sessions to continue coherently.

#### Scenario: Resume after cancelled TUI request
- **WHEN** a user resumes a session after cancelling an in-flight TUI request
- **THEN** the resumed session includes the prior conversation and a cancelled terminal event

### Requirement: Trace redaction covers runtime boundaries
Trace output MUST redact provider tokens, authorization headers, embedded credentials in URLs, configured secret environment variables, and raw provider request or response bodies unless explicitly safe.

#### Scenario: Provider request contains an authorization header
- **WHEN** a provider request is summarized in trace output
- **THEN** the authorization value is replaced with a redacted marker
