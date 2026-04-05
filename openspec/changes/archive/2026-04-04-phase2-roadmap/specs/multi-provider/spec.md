## ADDED Requirements

### Requirement: Multi-provider support
The system SHALL support multiple LLM providers through a unified interface.

#### Scenario: Provider interface
- **WHEN** a new provider is implemented
- **THEN** it implements the Provider interface (Name, SendMessage, SendMessageStream)

#### Scenario: Provider selection
- **WHEN** a user configures a provider in settings.json
- **THEN** the agent uses that provider for all API calls

#### Scenario: Supported providers
- **WHEN** the system starts
- **THEN** it supports Anthropic, OpenAI, and OpenAI-compatible providers
