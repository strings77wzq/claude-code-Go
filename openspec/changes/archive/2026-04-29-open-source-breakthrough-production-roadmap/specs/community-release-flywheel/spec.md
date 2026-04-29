## ADDED Requirements

### Requirement: Contributor Lanes Must Be Clear
The project SHALL define contributor lanes for docs, harness scenarios, provider profiles, built-in tools, permissions, TUI polish, and release infrastructure.

#### Scenario: New contributor looks for work
- **WHEN** a contributor opens contributing docs or issue labels
- **THEN** they can choose a lane with scope, acceptance criteria, and verification commands

### Requirement: Issues And PRs Must Capture Reproducibility
The project SHALL provide issue and pull request templates that collect version, OS, config shape, provider profile, command output, reproduction steps, verification commands, and known risks.

#### Scenario: User files a bug
- **WHEN** a bug report is created
- **THEN** the template asks for enough information to reproduce agent, provider, permission, or harness failures

#### Scenario: Contributor opens a PR
- **WHEN** a PR is opened
- **THEN** the template asks which tests/harness/docs checks ran and whether security-sensitive surfaces changed

### Requirement: Releases Must Include Evidence
The project SHALL publish release notes with feature changes, maturity level, verification results, known gaps, migration notes, and contributor credits.

#### Scenario: Release is published
- **WHEN** a GitHub release or changelog entry is created
- **THEN** it includes evidence for tests/harness/docs and clearly lists not-tested or partial areas

### Requirement: Benchmarks Must Be Reproducible
The project SHALL publish benchmark and showcase methodology before using benchmark results as community-facing claims.

#### Scenario: Benchmark result is advertised
- **WHEN** docs, README, or social copy cite benchmark results
- **THEN** the repository includes fixtures, commands, environment notes, scoring criteria, and limitations

### Requirement: Roadmap Must Support Community Trust
The project SHALL keep roadmap items linked to specs, tasks, status, and known blockers so users understand what is stable and what is experimental.

#### Scenario: User reads roadmap
- **WHEN** a user opens roadmap docs
- **THEN** each major item has status, priority, next step, and relationship to production readiness
