## ADDED Requirements

### Requirement: Session replay
The harness SHALL support replaying recorded sessions for debugging.

#### Scenario: Replay from JSONL
- **WHEN** a session file is loaded for replay
- **THEN** each message is displayed in sequence with timing information

#### Scenario: Replay with mock server
- **WHEN** a session is replayed against the mock server
- **THEN** the same tool calls are generated and the same responses are produced

### Requirement: Trace analysis
The harness SHALL analyze execution traces to identify issues.

#### Scenario: Identify tool call patterns
- **WHEN** a trace is analyzed
- **THEN** the frequency, order, and parameters of tool calls are reported

#### Scenario: Identify error patterns
- **WHEN** a trace contains errors
- **THEN** the error types, frequency, and context are identified
