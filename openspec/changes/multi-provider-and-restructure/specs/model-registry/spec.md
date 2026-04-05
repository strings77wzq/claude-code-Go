## ADDED Requirements

### Requirement: Model registry
The system SHALL maintain a registry of supported models with pricing information.

#### Scenario: Model list
- **WHEN** a user runs `/models`
- **THEN** they see all supported models grouped by provider

#### Scenario: Default model
- **WHEN** no model is specified
- **THEN** the default model is claude-sonnet-4-6-20251001
