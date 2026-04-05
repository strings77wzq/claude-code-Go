## ADDED Requirements

### Requirement: Hard timeout
The system SHALL enforce a hard timeout for API requests.

#### Scenario: Default timeout
- **WHEN** an API request is made
- **THEN** it times out after 5 minutes with a clear error message

#### Scenario: Real-time timer in TUI
- **WHEN** waiting for API response
- **THEN** the TUI shows elapsed time: "Waiting... (2.3s)"

#### Scenario: Timeout error
- **WHEN** the API request exceeds the timeout
- **THEN** a "Request timed out (5 minutes)" error is displayed and the user can continue
