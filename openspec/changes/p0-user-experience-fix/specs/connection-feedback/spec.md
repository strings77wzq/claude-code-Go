## ADDED Requirements

### Requirement: Connection timeout feedback
The TUI SHALL show progressive feedback while waiting for API response.

#### Scenario: Normal wait (0-3s)
- **WHEN** the API request is in progress for less than 3 seconds
- **THEN** only the spinner animation is shown

#### Scenario: Connecting message (3-30s)
- **WHEN** the API request takes more than 3 seconds
- **THEN** a "Connecting to API..." message is displayed below the spinner

#### Scenario: Still connecting message (30s-5min)
- **WHEN** the API request takes more than 30 seconds
- **THEN** the message changes to "Still connecting... check your network or API key"

#### Scenario: Timeout error (5min)
- **WHEN** the API request exceeds 5 minutes
- **THEN** a timeout error is displayed with instructions to check API key and network
