## ADDED Requirements

### Requirement: Additional built-in tools
The system SHALL provide additional built-in tools beyond the Phase 1 set.

#### Scenario: Diff tool
- **WHEN** the Diff tool is called with two file paths
- **THEN** it returns a unified diff output

#### Scenario: Tree tool
- **WHEN** the Tree tool is called with a directory path
- **THEN** it returns a directory tree visualization

#### Scenario: WebFetch tool
- **WHEN** the WebFetch tool is called with a URL
- **THEN** it fetches the page content and returns readable text
