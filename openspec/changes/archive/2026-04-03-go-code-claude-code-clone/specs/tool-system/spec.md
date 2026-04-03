## ADDED Requirements

### Requirement: Tool interface
The system SHALL define a uniform Tool interface that all tools implement.

#### Scenario: Tool contract
- **WHEN** a tool is registered
- **THEN** it provides Name(), Description(), InputSchema(), RequiresPermission(), and Execute() methods

### Requirement: Tool registry
The system SHALL provide a central registry for tool discovery and execution.

#### Scenario: Register tool
- **WHEN** a tool is registered with a unique name
- **THEN** it is available for lookup and execution

#### Scenario: Duplicate registration
- **WHEN** a tool is registered with a name that already exists
- **THEN** an error is returned

#### Scenario: Execute tool
- **WHEN** a tool is executed by name
- **THEN** the tool's Execute method is called with the provided input and the result is returned

#### Scenario: Execute unknown tool
- **WHEN** a non-existent tool name is used for execution
- **THEN** an error result is returned without panicking

#### Scenario: Execute tool with panic
- **WHEN** a tool's Execute method panics
- **THEN** the panic is recovered and an error result is returned

### Requirement: Built-in Bash tool
The system SHALL provide a Bash tool for executing shell commands.

#### Scenario: Execute command
- **WHEN** the Bash tool is called with a valid command
- **THEN** the command is executed and stdout+stderr output is returned

#### Scenario: Command timeout
- **WHEN** a command exceeds the timeout (default 120s)
- **THEN** the process is forcefully killed and a timeout error is returned

#### Scenario: Output truncation
- **WHEN** command output exceeds 100KB
- **THEN** the output is truncated with a "... (output truncated)" suffix

#### Scenario: Non-zero exit code
- **WHEN** a command exits with non-zero code
- **THEN** the output includes the exit code but the result is not marked as error

### Requirement: Built-in Read tool
The system SHALL provide a Read tool for reading file contents.

#### Scenario: Read file with line numbers
- **WHEN** the Read tool is called with a valid file path
- **THEN** the file content is returned with line numbers in "N: content" format

#### Scenario: Read with offset and limit
- **WHEN** the Read tool is called with offset and limit parameters
- **THEN** only the specified range of lines is returned

#### Scenario: Read non-existent file
- **WHEN** the Read tool is called with a non-existent file path
- **THEN** an error is returned

### Requirement: Built-in Write tool
The system SHALL provide a Write tool for creating or overwriting files.

#### Scenario: Create new file
- **WHEN** the Write tool is called with a file path and content
- **THEN** the file is created with the specified content, parent directories created if needed

#### Scenario: Overwrite existing file
- **WHEN** the Write tool is called with an existing file path
- **THEN** the file content is replaced

### Requirement: Built-in Edit tool
The system SHALL provide an Edit tool for precise string replacement in files.

#### Scenario: Unique replacement
- **WHEN** the Edit tool is called with an old_string that appears exactly once
- **THEN** the old_string is replaced with new_string

#### Scenario: Non-unique match
- **WHEN** the Edit tool is called with an old_string that appears multiple times and replace_all is false
- **THEN** an error is returned indicating the number of matches

#### Scenario: Replace all
- **WHEN** the Edit tool is called with replace_all=true
- **THEN** all occurrences of old_string are replaced with new_string

#### Scenario: String not found
- **WHEN** the Edit tool is called with an old_string not found in the file
- **THEN** an error is returned

### Requirement: Built-in Glob tool
The system SHALL provide a Glob tool for filename pattern matching.

#### Scenario: Pattern matching
- **WHEN** the Glob tool is called with a pattern like "**/*.go"
- **THEN** all matching file paths are returned

### Requirement: Built-in Grep tool
The system SHALL provide a Grep tool for regex content search.

#### Scenario: Content search
- **WHEN** the Grep tool is called with a regex pattern
- **THEN** matching lines are returned in "file:line:content" format

#### Scenario: Include/exclude filters
- **WHEN** the Grep tool is called with include or exclude patterns
- **THEN** only matching/non-matching files are searched
