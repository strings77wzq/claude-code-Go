## Purpose
Define the CLI and TUI experience required for task-oriented entrypoints, consistent slash-command behavior, and concise user-facing errors.
## Requirements
### Requirement: CLI supports task-oriented entrypoints
The system SHALL provide consistent command-line entrypoints for setup, doctor, interactive TUI, single prompt execution, JSON output, quiet mode, version, and help.

#### Scenario: Single prompt JSON output
- **WHEN** the user runs `go-code -p "hello" -f json`
- **THEN** the command returns valid JSON containing the assistant result and exits without launching the TUI

#### Scenario: Help lists available entrypoints
- **WHEN** the user runs `go-code --help`
- **THEN** the output lists setup, doctor, prompt mode, output format, quiet mode, debug mode, and interactive mode

### Requirement: TUI and REPL commands remain consistent
The system MUST route TUI and legacy REPL commands through a shared command layer so advertised commands behave the same way in both surfaces.

#### Scenario: Sessions command in default TUI
- **WHEN** the user enters `/sessions` in the default TUI
- **THEN** the system lists available sessions or prints an actionable empty-state message

#### Scenario: Unknown command
- **WHEN** the user enters an unsupported slash command
- **THEN** the system prints a clear unknown-command message and points to `/help`

### Requirement: User-facing errors are smooth and specific
The system SHALL convert common failures into concise messages that name the failing subsystem and next action.

#### Scenario: Provider authentication failure
- **WHEN** the provider returns an authentication error
- **THEN** the UI reports invalid credentials and points to setup or config docs without dumping raw response bodies

#### Scenario: Long connection delay
- **WHEN** a provider request takes longer than the configured threshold
- **THEN** the UI shows elapsed time and a clear connection status without freezing input rendering

### Requirement: TUI loading state terminates deterministically
The TUI MUST leave loading state when an agent request completes, fails, is denied by permission policy, or is cancelled.

#### Scenario: Provider returns an error
- **WHEN** the active provider returns an error during a TUI request
- **THEN** the TUI displays the failure state
- **AND** the loading indicator stops

### Requirement: TUI cancellation does not leak request work
The TUI MUST cancel in-flight request work when the user cancels or exits the request flow.

#### Scenario: User cancels generation
- **WHEN** the user cancels an in-flight generation from the TUI
- **THEN** the request context is cancelled
- **AND** no further assistant output is appended for that cancelled request
