## ADDED Requirements

### Requirement: Path validation for file operations
All file operation tools SHALL validate that the file path is within the working directory.

#### Scenario: Path inside working directory
- **WHEN** a file tool is called with a path inside the working directory
- **THEN** the operation proceeds normally

#### Scenario: Path outside working directory
- **WHEN** a file tool is called with a path outside the working directory (e.g., /etc/passwd)
- **THEN** an error is returned: "path outside working directory"
