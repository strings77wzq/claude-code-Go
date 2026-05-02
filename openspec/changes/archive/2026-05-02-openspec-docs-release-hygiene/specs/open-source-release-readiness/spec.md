## ADDED Requirements

### Requirement: Install smoke test gates published versions
The project MUST run an install smoke test for the target release artifact before publishing or recommending a version.

#### Scenario: Binary is built for release
- **WHEN** a release binary or package is produced
- **THEN** the smoke test verifies installation, `--help`, doctor or offline status, and one deterministic non-provider command

### Requirement: Release checklist includes spec and docs validation
The release checklist SHALL include strict OpenSpec validation, docs truth checks, harness gate status, and known-risk notes.

#### Scenario: Release checklist is incomplete
- **WHEN** required validation evidence is missing
- **THEN** the release is not marked ready
