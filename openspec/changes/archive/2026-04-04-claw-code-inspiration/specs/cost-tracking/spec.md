## ADDED Requirements

### Requirement: Cost tracking
The system SHALL track and estimate API costs.

#### Scenario: Per-request cost estimation
- **WHEN** an API response is received
- **THEN** the cost is estimated based on token usage and model pricing

#### Scenario: Session cost summary
- **WHEN** a session ends
- **THEN** the total cost is displayed to the user
