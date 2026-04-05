## ADDED Requirements

### Requirement: File boundary guards
The system SHALL enforce file operation boundaries.

#### Scenario: Binary file protection
- **WHEN** a write operation targets a binary file (.exe, .bin, .so, .dylib, .dll)
- **THEN** the operation is blocked with a warning

#### Scenario: File size limit
- **WHEN** a write operation exceeds 10MB
- **THEN** the operation is blocked

#### Scenario: Symlink escape prevention
- **WHEN** a file path contains a symlink that escapes the workspace
- **THEN** the operation is blocked
