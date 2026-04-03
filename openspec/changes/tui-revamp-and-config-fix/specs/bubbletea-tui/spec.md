## ADDED Requirements

### Requirement: Bubbletea TUI
The system SHALL provide a new TUI based on bubbletea framework.

#### Scenario: Stream rendering
- **WHEN** the agent streams text
- **THEN** it renders character by character with smooth output

#### Scenario: Tool call display
- **WHEN** a tool is called
- **THEN** it shows the tool name, input, and result with status icons

#### Scenario: Legacy fallback
- **WHEN** the user runs with --legacy-repl flag
- **THEN** the old bufio-based REPL is used
