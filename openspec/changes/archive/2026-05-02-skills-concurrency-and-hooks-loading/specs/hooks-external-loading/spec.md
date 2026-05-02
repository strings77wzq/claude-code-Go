## ADDED Requirements

### Requirement: External hooks loaded from user directory
The system SHALL load hook definitions from `~/.go-code/hooks/` directory. Each valid JSON file SHALL define a hook with `name`, `type` (pre or post), and `command` fields. Invalid files SHALL be skipped with warnings.

#### Scenario: Valid hook file loaded
- **WHEN** `~/.go-code/hooks/my-hook.json` contains a valid hook definition
- **THEN** the hook is registered and executes on matching tool events

#### Scenario: Invalid hook file skipped
- **WHEN** `~/.go-code/hooks/bad.json` contains invalid JSON
- **THEN** the file is skipped, a warning is generated, and other hooks continue loading

#### Scenario: No hooks directory
- **WHEN** `~/.go-code/hooks/` does not exist
- **THEN** the agent starts normally with only built-in hooks

### Requirement: PostExecute errors are logged
The system SHALL log PostExecute hook errors using structured logging (`slog`) instead of silently discarding them.

#### Scenario: PostExecute hook fails
- **WHEN** a post-execute hook returns an error
- **THEN** the error is logged at WARN level with the hook name and error message
- **AND** subsequent hooks continue executing
