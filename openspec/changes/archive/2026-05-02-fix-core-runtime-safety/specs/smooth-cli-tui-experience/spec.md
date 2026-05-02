## ADDED Requirements

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
