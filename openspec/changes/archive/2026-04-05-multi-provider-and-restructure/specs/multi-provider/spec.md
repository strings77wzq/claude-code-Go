## ADDED Requirements

### Requirement: Multi-provider support
The system SHALL support multiple API providers through a unified interface.

#### Scenario: Provider selection
- **WHEN** a model name is configured
- **THEN** the correct provider is automatically selected based on model name prefix

#### Scenario: Anthropic models
- **WHEN** the model name starts with "claude-"
- **THEN** the Anthropic provider is used

#### Scenario: OpenAI-compatible models
- **WHEN** the model name starts with "gpt-", "deepseek-", "qwen-", "glm-"
- **THEN** the OpenAI-compatible provider is used
