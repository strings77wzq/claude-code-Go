## ADDED Requirements

### Requirement: Agent Loop
The system SHALL implement a "think → act → observe" loop that drives autonomous task completion.

#### Scenario: Text-only response
- **WHEN** the model returns a text response with stop_reason "end_turn"
- **THEN** the loop exits and returns the text to the user

#### Scenario: Tool use response
- **WHEN** the model returns a tool_use block with stop_reason "tool_use"
- **THEN** the system executes the tool, feeds the result back to the model, and continues the loop

### Requirement: Stop reason dispatch
The system SHALL dispatch behavior based on the stop_reason field from the API response.

#### Scenario: end_turn
- **WHEN** stop_reason is "end_turn"
- **THEN** the loop terminates and returns the assistant's text

#### Scenario: tool_use
- **WHEN** stop_reason is "tool_use"
- **THEN** tools are executed, results added to history, and the loop continues

#### Scenario: max_tokens
- **WHEN** stop_reason is "max_tokens"
- **THEN** the loop terminates with a truncation warning message

#### Scenario: Unknown stop_reason
- **WHEN** stop_reason is an unrecognized value
- **THEN** the loop safely terminates and returns any available text

### Requirement: Maximum turns limit
The system SHALL enforce a maximum number of loop iterations to prevent infinite loops.

#### Scenario: Max turns reached
- **WHEN** the loop exceeds 50 iterations
- **THEN** the loop terminates with a "max turns reached" message

### Requirement: Streaming output callback
The system SHALL support a callback function for streaming text output during API calls.

#### Scenario: Real-time text display
- **WHEN** the API streams text_delta events
- **THEN** each delta is passed to the outputCallback function for immediate display

### Requirement: Message history management
The system SHALL maintain conversation history with strict user/assistant alternation.

#### Scenario: History accumulation
- **WHEN** the agent processes multiple turns
- **THEN** user messages, assistant messages, and tool results are accumulated in correct alternation order
