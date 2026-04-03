## ADDED Requirements

### Requirement: JSONL session persistence
The system SHALL persist conversation history in JSONL format.

#### Scenario: Save session metadata
- **WHEN** a new session is created
- **THEN** a session_meta line with session_id, version, and timestamps is written

#### Scenario: Append messages
- **WHEN** a message is added to the conversation
- **THEN** a message line with role and content blocks is appended to the JSONL file

#### Scenario: Atomic writes
- **WHEN** session data is written
- **THEN** it is written to a temp file first, then renamed to ensure atomicity

### Requirement: Session loading
The system SHALL load existing sessions from JSONL files.

#### Scenario: Load existing session
- **WHEN** a session file exists
- **THEN** all messages are reconstructed from the JSONL lines

#### Scenario: Handle corrupted session
- **WHEN** a session file contains invalid JSON lines
- **THEN** valid lines are loaded and invalid lines are skipped with a warning

### Requirement: Session resume
The system SHALL support resuming a previous conversation.

#### Scenario: Resume by session ID
- **WHEN** a user requests to resume a specific session
- **THEN** the conversation history is restored and the agent can continue

### Requirement: Log rotation
The system SHALL rotate session files when they exceed a size threshold.

#### Scenario: Rotate large session
- **WHEN** a session file exceeds 256KB
- **THEN** old messages are compacted and the file is rotated
