## ADDED Requirements

### Requirement: Output quality evaluation
The harness SHALL evaluate the quality of model output responses.

#### Scenario: Evaluate text completeness
- **WHEN** a streaming response completes
- **THEN** the evaluator verifies that the full expected text was received without truncation

#### Scenario: Evaluate tool call correctness
- **WHEN** the model generates tool calls
- **THEN** the evaluator verifies the tool name, parameters, and call order match expectations

### Requirement: Latency monitoring
The harness SHALL monitor and report latency metrics for API interactions.

#### Scenario: Measure response latency
- **WHEN** an API call is made
- **THEN** the time from request to first token and to completion is recorded

### Requirement: Token usage evaluation
The harness SHALL evaluate token usage accuracy and efficiency.

#### Scenario: Track token consumption
- **WHEN** a conversation completes
- **THEN** total input and output tokens are reported and compared against estimates
