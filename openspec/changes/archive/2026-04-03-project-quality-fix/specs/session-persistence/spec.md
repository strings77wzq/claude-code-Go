## ADDED Requirements

### Requirement: Session persistence
The system SHALL persist conversation sessions in JSONL format for later replay and debugging.

#### Scenario: Save session
- **WHEN** a conversation ends
- **THEN** the session is saved as a JSONL file with timestamp in the filename

#### Scenario: Session format
- **WHEN** a session file is examined
- **THEN** each line is a valid JSON object with type, role, content, and timestamp fields

#### Scenario: Load session
- **WHEN** a session file is loaded
- **THEN** the conversation history is reconstructed in order

### Requirement: Session metadata
Each session file SHALL include metadata about the conversation.

#### Scenario: Session metadata
- **WHEN** a session is saved
- **THEN** it includes: session_id, model, start_time, end_time, turn_count, token_usage
