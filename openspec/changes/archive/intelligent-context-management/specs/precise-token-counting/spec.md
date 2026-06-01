## ADDED Requirements

### Requirement: Token counting uses accurate tokenization

The system SHALL count tokens using a real tokenizer (cl100k_base encoding) instead of the 4-characters-per-token heuristic. Token counts SHALL be accurate within 5% of Anthropic's reported token counts.

#### Scenario: Token count matches actual usage
- **WHEN** a message list is tokenized
- **THEN** the counted tokens are within 5% of the count reported by the API response `usage.input_tokens`

#### Scenario: Compaction triggers at correct threshold
- **WHEN** the accurate token count exceeds the compaction threshold (80% of context window)
- **THEN** compaction is triggered

#### Scenario: Tokenizer handles tool content blocks
- **WHEN** messages contain tool_use and tool_result content blocks
- **THEN** the tokenizer correctly accounts for their JSON structure tokens
