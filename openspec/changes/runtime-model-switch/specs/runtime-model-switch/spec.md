## ADDED Requirements

### Requirement: Runtime model switching
The system SHALL support switching models at runtime via `/model` command.

#### Scenario: Display current model
- **WHEN** a user types `/model` without arguments
- **THEN** the current model name is displayed

#### Scenario: Switch model
- **WHEN** a user types `/model <model-name>`
- **THEN** the model is switched and a confirmation is shown

#### Scenario: Invalid model name
- **WHEN** a user switches to an invalid model
- **THEN** the switch succeeds locally, but the next API call will show an error
