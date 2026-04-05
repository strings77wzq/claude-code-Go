## ADDED Requirements

### Requirement: OpenAI-compatible provider
The system SHALL support OpenAI API format for compatible models.

#### Scenario: OpenAI format request
- **WHEN** sending a message via OpenAI provider
- **THEN** the request uses POST /v1/chat/completions with Authorization: Bearer header

#### Scenario: OpenAI streaming response
- **WHEN** streaming is enabled
- **THEN** the response is parsed from SSE format with data: {"choices":[{"delta":{"content":"..."}}]}
