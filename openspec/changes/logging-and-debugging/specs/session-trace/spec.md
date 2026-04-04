## ADDED Requirements

### Requirement: Session trace
The system SHALL save complete session traces for debugging.

#### Scenario: Trace saved to session file
- **WHEN** a session ends
- **THEN** the session file contains: API requests, responses, tool calls, errors

#### Scenario: Trace command
- **WHEN** a user types `/trace last`
- **THEN** the last session trace is displayed

#### Scenario: Export command
- **WHEN** a user types `/export session`
- **THEN** the current session is exported as JSON
