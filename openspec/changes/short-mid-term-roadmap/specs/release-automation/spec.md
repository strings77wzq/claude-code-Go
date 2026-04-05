## ADDED Requirements

### Requirement: Release automation
The project SHALL support automated releases via GoReleaser.

#### Scenario: Automated release
- **WHEN** a new tag is pushed
- **THEN** binaries are built for all platforms and published to GitHub Releases
