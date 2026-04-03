## ADDED Requirements

### Requirement: Token estimation
The system SHALL estimate token usage for conversation history.

#### Scenario: Estimate message tokens
- **WHEN** the token count of conversation history is requested
- **THEN** an estimate is returned based on character count (approximately 4 chars per token)

### Requirement: Auto compaction
The system SHALL automatically compact conversation history when approaching the context window limit.

#### Scenario: Compact at threshold
- **WHEN** the estimated token count exceeds 80% of the model's context window
- **THEN** old messages are summarized and replaced with a summary message

#### Scenario: Preserve recent messages
- **WHEN** compaction is triggered
- **THEN** the first user message and the most recent 10 turns are preserved verbatim

### Requirement: Manual compaction
The system SHALL support manual compaction via user command.

#### Scenario: User triggers compaction
- **WHEN** the user issues a compact command
- **THEN** the conversation history is compacted regardless of token count
