## ADDED Requirements

### Requirement: Makefile build targets
The project SHALL provide a Makefile with standard build targets.

#### Scenario: Build Go binary
- **WHEN** `make build` is run
- **THEN** the Go CLI binary is compiled and placed in a standard output directory

#### Scenario: Run tests
- **WHEN** `make test` is run
- **THEN** Go unit tests and Python harness tests are executed

#### Scenario: Start docs server
- **WHEN** `make docs` is run
- **THEN** the MkDocs development server starts

### Requirement: GitHub Actions CI
The project SHALL provide CI automation via GitHub Actions.

#### Scenario: Go test CI
- **WHEN** a push or PR is made
- **THEN** Go tests (`go test ./...`) and race detection (`go test -race ./...`) run

#### Scenario: Python harness CI
- **WHEN** a push or PR is made
- **THEN** Python harness tests (pytest) run against the built Go binary

#### Scenario: Build verification
- **WHEN** a push or PR is made
- **THEN** `go build` succeeds for the target platform

### Requirement: Multi-platform builds
The project SHALL support building for multiple platforms.

#### Scenario: Cross-platform build
- **WHEN** `make build-all` is run
- **THEN** binaries are produced for linux/amd64, darwin/amd64, and darwin/arm64
