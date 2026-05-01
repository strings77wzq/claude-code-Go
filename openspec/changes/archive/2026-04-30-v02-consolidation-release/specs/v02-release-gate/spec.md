## ADDED Requirements

### Requirement: All Go tests pass
The full `go test ./...` suite SHALL pass with zero failures for all packages in the repository.

#### Scenario: Clean test run
- **WHEN** `go test ./...` is executed from the repository root
- **THEN** the exit code is 0
- **AND** every package reports `ok` (not `?` for missing tests, except intentionally untestable entry points)

### Requirement: Python harness scenarios pass
The `./scripts/run-harness.sh` script SHALL build the Go binary, start local mock servers, run all harness scenarios, and exit with code 0.

#### Scenario: Harness green
- **WHEN** `./scripts/run-harness.sh` is executed from the repository root
- **THEN** all pytest scenarios pass
- **AND** the script exits with code 0

### Requirement: Docs build passes
The documentation site SHALL build without errors.

#### Scenario: Docs build clean
- **WHEN** the docs build command is executed
- **THEN** the build completes without broken links or fatal errors
- **AND** the exit code is 0

### Requirement: Changelog entry exists
The CHANGELOG.md SHALL contain a v0.2 entry summarizing changes since the last release, with links to the relevant OpenSpec change and parity evidence.

#### Scenario: Changelog is current
- **WHEN** CHANGELOG.md is inspected
- **THEN** a v0.2 section exists with date, summary, and links
- **AND** the entry references the v0.2 mandatory workflow verification results

### Requirement: PARITY.md v0.2 mandatory workflows are all verified
All 8 v0.2 mandatory workflows defined in PARITY.md SHALL have status `verified` with linked evidence.

#### Scenario: v0.2 workflow verification complete
- **WHEN** PARITY.md is inspected
- **THEN** all 8 mandatory v0.2 workflow rows have status `verified`
- **AND** each row includes an evidence link (test file path or smoke check doc)
