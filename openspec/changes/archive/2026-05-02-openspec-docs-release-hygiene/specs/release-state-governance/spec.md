## ADDED Requirements

### Requirement: Release state transitions require evidence
The project MUST require validation evidence before moving a release state from local/dev to published or recommended-use status.

#### Scenario: Release candidate is promoted
- **WHEN** maintainers promote a release candidate
- **THEN** the release evidence includes tests, OpenSpec validation, harness status, docs checks, and install smoke results

#### Scenario: Release state matrix is reviewed
- **WHEN** a maintainer checks release readiness
- **THEN** the release state matrix lists required evidence for local/dev, release-candidate, and published states

### Requirement: Known gaps are published with release evidence
Release notes SHALL identify known gaps that affect agent usability, safety, or compatibility.

#### Scenario: Harness comparison is incomplete
- **WHEN** a release is published before full external-agent comparison is complete
- **THEN** the release notes identify the comparison as incomplete
