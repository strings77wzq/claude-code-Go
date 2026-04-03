## ADDED Requirements

### Requirement: Multi-source configuration
The system SHALL load configuration from multiple sources with defined priority.

#### Scenario: Priority chain
- **WHEN** configuration is loaded
- **THEN** the priority order is: CLI args > environment variables > project config > user config > defaults

### Requirement: Environment variable support
The system SHALL read configuration from environment variables.

#### Scenario: API key from environment
- **WHEN** ANTHROPIC_API_KEY is set
- **THEN** the API key is read from the environment variable

#### Scenario: Base URL from environment
- **WHEN** ANTHROPIC_BASE_URL is set
- **THEN** the base URL is read from the environment variable

### Requirement: Config file loading
The system SHALL load configuration from JSON settings files.

#### Scenario: User config file
- **WHEN** ~/.go-code/settings.json exists
- **THEN** its values are loaded as the user-level configuration

#### Scenario: Project config file
- **WHEN** ./.go-code/settings.json exists
- **THEN** its values override the user-level configuration

### Requirement: Default values
The system SHALL provide sensible default values for all configuration fields.

#### Scenario: Default configuration
- **WHEN** no configuration sources are present
- **THEN** defaults are used: BaseURL="https://api.anthropic.com", Model="claude-sonnet-4-20250514", MaxTokens=8192

### Requirement: API key validation
The system SHALL validate that an API key is present before starting.

#### Scenario: Missing API key
- **WHEN** no API key is found in any configuration source
- **THEN** an error is returned indicating the API key is required
