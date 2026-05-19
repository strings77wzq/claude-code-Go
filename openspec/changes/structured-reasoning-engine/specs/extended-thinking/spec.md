## ADDED Requirements

### Requirement: Agent supports configurable extended thinking

The system SHALL support Anthropic's extended thinking with a configurable token budget. When the thinking budget is non-zero, API requests SHALL include the `thinking` parameter with `budget_tokens`. When the budget is zero, thinking SHALL be disabled and no `thinking` parameter is sent.

#### Scenario: Thinking enabled with budget
- **WHEN** the agent's thinking budget is set to 1600
- **THEN** the API request includes `{"thinking": {"type": "enabled", "budget_tokens": 1600}}`

#### Scenario: Thinking disabled by default
- **WHEN** a new agent is created without setting a thinking budget
- **THEN** the API request does NOT include a `thinking` parameter

#### Scenario: Thinking blocks are stored in trace
- **WHEN** an API response contains thinking content blocks
- **THEN** the thinking content is appended to the session trace as a `thinking` event

#### Scenario: Thinking stop_reason triggers continue
- **WHEN** the API response has stop_reason "thinking"
- **THEN** the agent continues the loop with thinking content as a tool_result block
