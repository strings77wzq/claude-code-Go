## ADDED Requirements

### Requirement: Model list command
The system SHALL support listing available models via `/models` command.

#### Scenario: List models
- **WHEN** a user types `/models`
- **THEN** available models are displayed, grouped by provider (Anthropic, Tencent Coding Plan)

#### Scenario: Default model indicator
- **WHEN** a model is listed as default
- **THEN** it is marked with "(default)" indicator
