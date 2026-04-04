## ADDED Requirements

### Requirement: Tool description constraints
The system SHALL enforce tool description and schema constraints.

#### Scenario: Description length limit
- **WHEN** a tool is registered with a description longer than 250 characters
- **THEN** the registration fails with a validation error

#### Scenario: InputSchema required
- **WHEN** a tool is registered without an InputSchema
- **THEN** the registration fails with a validation error
