## ADDED Requirements

### Requirement: Sessions persist complete conversation state
The system SHALL persist user messages, assistant messages, tool uses, tool results, model, timestamps, token usage when available, and final status in a JSONL session file.

#### Scenario: Tool session saved
- **WHEN** an agent run includes a tool call and final response
- **THEN** the saved session contains the user input, assistant tool request, tool result, and final assistant text

#### Scenario: Interrupted session
- **WHEN** a session is interrupted before normal completion
- **THEN** the saved session records the interruption status and recoverable history

### Requirement: Resume restores usable context
The system MUST allow users to list and resume prior sessions.

#### Scenario: Resume by session id
- **WHEN** the user runs `/resume <session-id>`
- **THEN** the system loads the prior conversation history and continues subsequent turns in that context

### Requirement: Replay supports debugging
The system SHALL provide replay tooling that can inspect a saved session and reproduce or summarize the agent/tool sequence without calling a real provider.

#### Scenario: Replay latest session
- **WHEN** the user runs the replay command for the latest session
- **THEN** the system prints the sequence of requests, tool calls, results, errors, and final status

