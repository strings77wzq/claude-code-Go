# release-state-governance Specification

## Purpose
TBD - created by archiving change v03-extension-productization. Update Purpose after archive.
## Requirements
### Requirement: Active OpenSpec changes reflect real work in progress
The project SHALL keep completed, obsolete, or parked OpenSpec changes out of the active implementation lane before starting a new milestone.

#### Scenario: Completed change is still active
- **WHEN** `openspec list` shows a change as complete
- **THEN** the maintainer either archives it, records why it remains active, or moves follow-up work into a new change

#### Scenario: Parked change has no tasks
- **WHEN** an active change has proposal-only scope and no actionable tasks
- **THEN** the maintainer labels it as parked in the project status or converts it into a complete spec-driven change before implementation

### Requirement: Release evidence is captured before public status changes
The project SHALL capture verification evidence before README, roadmap, PARITY.md, or release notes mark a milestone as verified.

#### Scenario: Milestone is marked verified
- **WHEN** docs claim a milestone is verified
- **THEN** the evidence includes Go tests, Python harness results, OpenSpec validation, docs build status, and known-risk notes

### Requirement: Generated documentation artifacts are handled intentionally
The project SHALL separate docs source edits from generated `docs/.vitepress/dist` output unless a release task explicitly requires generated artifacts.

#### Scenario: Docs source changes are reviewed
- **WHEN** docs source files change
- **THEN** generated dist changes are either excluded from the review or included with a release-specific rationale

### Requirement: Missing steering documents are recorded explicitly
The project SHALL record when expected steering documents such as `CLAUDE.md`, `task.md`, or `TASK.md` are absent so future planning does not assume they exist.

#### Scenario: Project progress audit runs
- **WHEN** a progress audit checks for steering documents
- **THEN** the result names which steering documents were found and which were absent

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
