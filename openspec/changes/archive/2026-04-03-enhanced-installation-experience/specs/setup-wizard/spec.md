## ADDED Requirements

### Requirement: Setup Wizard CLI
The system SHALL provide a `--setup` flag that launches an interactive configuration wizard.

#### Scenario: User runs go-code --setup
- **WHEN** user runs `go-code --setup`
- **THEN** an interactive wizard starts asking for provider, API key, and model

#### Scenario: User selects provider
- **WHEN** user is prompted to select a provider
- **THEN** they can choose from Anthropic, OpenAI, or Custom

#### Scenario: User enters API key
- **WHEN** user enters an API key
- **THEN** the system validates the format (Anthropic: sk-ant-, OpenAI: sk-, Custom: non-empty)

#### Scenario: User skips API key
- **WHEN** user chooses to skip entering an API key
- **THEN** the wizard completes without writing a config file and informs the user how to configure later

#### Scenario: Config is written
- **WHEN** user completes the wizard with valid input
- **THEN** `~/.go-code/settings.json` is created with the provider, API key, and model

### Requirement: install.sh Integration
The install.sh script SHALL call `go-code --setup` after downloading the binary.

#### Scenario: install.sh calls setup
- **WHEN** install.sh finishes downloading and installing the binary
- **THEN** it calls `go-code --setup` to start the configuration wizard

#### Scenario: setup call fails
- **WHEN** `go-code --setup` fails or is not found
- **THEN** install.sh prints manual configuration instructions as fallback
