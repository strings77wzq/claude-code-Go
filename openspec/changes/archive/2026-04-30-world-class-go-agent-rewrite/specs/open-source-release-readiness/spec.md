## ADDED Requirements

### Requirement: CI gates release-quality changes
The project SHALL run formatting, vet/static checks, Go tests, harness tests, docs build, and OpenSpec validation in CI.

#### Scenario: Pull request changes runtime code
- **WHEN** a pull request modifies runtime code
- **THEN** CI runs the runtime test suite and parity harness before merge

### Requirement: Releases are reproducible
The project MUST document and automate building signed or checksummed binaries for supported platforms.

#### Scenario: Release artifact generated
- **WHEN** a release workflow runs for a tag
- **THEN** it produces platform binaries, checksums, release notes, and installation instructions

### Requirement: Community contribution path is clear
The project SHALL provide issue templates, pull request template, contribution guide, coding standards, test commands, and first-good-issue guidance.

#### Scenario: New contributor opens contribution guide
- **WHEN** a contributor reads the guide
- **THEN** they can set up the project, run tests, understand package boundaries, and choose a starter task

### Requirement: Roadmap is honest and actionable
The project MUST maintain a roadmap that separates now, next, later, and not-planned work.

#### Scenario: Unsupported feature request
- **WHEN** users ask for a feature outside the current scope
- **THEN** maintainers can point to roadmap status and acceptance criteria
