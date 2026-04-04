## ADDED Requirements

### Requirement: Debug mode
The system SHALL support a debug mode via `--debug` flag.

#### Scenario: Debug flag enables verbose logging
- **WHEN** the user runs `go-code --debug`
- **THEN** logs are written to both file and stderr

#### Scenario: Debug status bar in TUI
- **WHEN** debug mode is enabled
- **THEN** the TUI shows a debug status bar with API latency, token usage, and tool execution details

#### Scenario: HTTP trace mode
- **WHEN** the user runs `go-code --debug --trace-http`
- **THEN** complete HTTP request/response bodies are logged
