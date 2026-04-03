## ADDED Requirements

### Requirement: Auto-update
The system SHALL support checking for and downloading updates.

#### Scenario: Update check
- **WHEN** a user types `/update`
- **THEN** the system checks GitHub Releases for the latest version

#### Scenario: Self-update
- **WHEN** a newer version is available
- **THEN** the system downloads and replaces the current binary
