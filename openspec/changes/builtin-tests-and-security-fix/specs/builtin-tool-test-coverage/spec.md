## ADDED Requirements

### Requirement: All builtin tools have comprehensive unit tests

The system SHALL have unit tests for every tool in `internal/tool/builtin` covering happy path, error conditions, timeout, output truncation, and permission gating. The package SHALL achieve at least 80% statement test coverage.

#### Scenario: Bash tool test coverage
- **WHEN** `go test -cover ./internal/tool/builtin` is run
- **THEN** coverage is at least 80% and all tests pass

#### Scenario: Glob tool returns matching files
- **WHEN** `Glob.Execute` is called with a valid pattern against a temp directory
- **THEN** it returns the list of matching file paths

#### Scenario: Grep tool finds matching lines
- **WHEN** `Grep.Execute` is called with a pattern against a temp file
- **THEN** it returns lines containing the pattern

#### Scenario: Tree tool returns directory structure
- **WHEN** `Tree.Execute` is called on a temp directory with nested files
- **THEN** it returns a tree-formatted listing of the directory

#### Scenario: WebFetch tool fetches URL content
- **WHEN** `WebFetch.Execute` is called with a valid URL pointing to a mock HTTP server
- **THEN** it returns the response body content

#### Scenario: Diff tool shows file differences
- **WHEN** `Diff.Execute` is called with two different temp files
- **THEN** it returns a unified diff showing the differences

#### Scenario: Notebook tool operates on notebook cells
- **WHEN** `NotebookEdit.Execute` is called with a valid notebook path and cell edit parameters
- **THEN** it returns the modified notebook content

#### Scenario: Todo tool manages task list
- **WHEN** `TodoWrite.Execute` is called with todo items
- **THEN** it returns the formatted todo list

#### Scenario: Path validation utility rejects blocked paths
- **WHEN** `ResolvePath` or `ValidatePath` is called with a system path like `/dev/sda`
- **THEN** it returns an error indicating the path is blocked

#### Scenario: Path validation utility accepts workspace paths
- **WHEN** `ResolvePath` or `ValidatePath` is called with a path within the working directory
- **THEN** it returns the resolved absolute path without error
