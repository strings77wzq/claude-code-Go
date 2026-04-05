## ADDED Requirements

### Requirement: ANTHROPIC_MODEL environment variable
The system SHALL read ANTHROPIC_MODEL environment variable.

#### Scenario: Model from environment
- **WHEN** ANTHROPIC_MODEL is set
- **THEN** it overrides the default model
