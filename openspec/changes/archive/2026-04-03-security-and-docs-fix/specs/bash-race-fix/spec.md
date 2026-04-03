## ADDED Requirements

### Requirement: Bash tool race condition fix
The Bash tool SHALL use proper synchronization for concurrent output access.

#### Scenario: Command execution
- **WHEN** a command is executed
- **THEN** the output goroutine completes before the result is read, using sync.WaitGroup
