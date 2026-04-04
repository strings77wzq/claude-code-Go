## ADDED Requirements

### Requirement: Non-interactive mode
The system SHALL support running with a single prompt via `-p` flag.

#### Scenario: Single prompt
- **WHEN** a user runs `go-code -p "explain this code"`
- **THEN** the agent processes the prompt and prints the result, then exits

#### Scenario: JSON output
- **WHEN** a user runs `go-code -p "prompt" -f json`
- **THEN** the output is formatted as JSON

#### Scenario: Quiet mode
- **WHEN** a user runs `go-code -p "prompt" -q`
- **THEN** no spinner or status messages are shown

#### Scenario: Auto-approve permissions
- **WHEN** running in non-interactive mode
- **THEN** all tool permissions are automatically approved
