## ADDED Requirements

### Requirement: File logging
The system SHALL write structured logs to a file.

#### Scenario: Log file creation
- **WHEN** the application starts
- **THEN** a log file is created at `~/.go-code/go-code.log`

#### Scenario: Structured log format
- **WHEN** a log entry is written
- **THEN** it includes timestamp, level, message, and context fields in JSON format

#### Scenario: Key events logged
- **WHEN** key events occur (API request, tool execution, error, session start/end)
- **THEN** they are logged with relevant context (model, tool name, duration, etc.)
