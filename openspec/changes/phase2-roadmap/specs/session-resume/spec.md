## ADDED Requirements

### Requirement: Session resume
The system SHALL support resuming a previous conversation session.

#### Scenario: Resume command
- **WHEN** a user types `/resume <session-id>`
- **THEN** the session history is loaded from JSONL and the agent continues

#### Scenario: List sessions
- **WHEN** a user types `/sessions`
- **THEN** available sessions are listed with ID, date, and turn count
