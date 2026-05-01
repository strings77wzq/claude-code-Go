# session-trace-and-replay Specification

## Purpose
Defines session persistence, replay, and trace evidence behavior for agent runs and extension debugging.
## Requirements
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

### Requirement: Replay includes extension and permission events
Replay output SHALL include enough structured information to debug MCP, LSP, hook, skill, and permission behavior without calling a real provider.

#### Scenario: MCP tool is denied
- **WHEN** a saved session includes an MCP tool call denied by policy
- **THEN** replay shows the tool name, permission decision, denial reason, and final agent-visible result

#### Scenario: LSP capability is unavailable
- **WHEN** a saved session attempts an unavailable LSP operation
- **THEN** replay shows the unavailable capability, configuration state, and non-fatal result

### Requirement: Trace output redacts sensitive values
Trace and replay output MUST redact API keys, provider tokens, authorization headers, and configured secrets.

#### Scenario: Provider configuration contains a secret
- **WHEN** replay or trace output includes provider configuration
- **THEN** secret values are replaced with a redacted marker

### Requirement: Replay has a release-evidence mode
The replay command SHALL provide a concise mode suitable for attaching to release evidence or issue reports.

#### Scenario: Maintainer collects replay evidence
- **WHEN** the maintainer runs replay in evidence mode
- **THEN** the output summarizes prompts, tool calls, permission decisions, extension availability, errors, and final status
- **AND** it omits raw provider secrets and excessive token payloads

