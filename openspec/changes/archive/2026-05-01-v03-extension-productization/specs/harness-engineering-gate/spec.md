## ADDED Requirements

### Requirement: Harness covers extension trajectories
The Python harness SHALL include deterministic scenarios for MCP registration, MCP permission denial, LSP unavailable behavior, and replay evidence output before v0.3 extension support is marked verified.

#### Scenario: MCP registration harness scenario runs
- **WHEN** the harness starts a local MCP fixture
- **THEN** `go-code` registers the fixture tool and records the tool lifecycle without live credentials

#### Scenario: LSP unavailable harness scenario runs
- **WHEN** the harness runs without an LSP server
- **THEN** `go-code` reports LSP unavailable without failing the prompt workflow

### Requirement: Release gate includes extension diagnostics
The release gate SHALL run extension diagnostics in addition to Go tests, harness tests, docs build, and OpenSpec validation.

#### Scenario: Maintainer checks v0.3 readiness
- **WHEN** the maintainer runs the v0.3 verification commands
- **THEN** the result includes `go test ./...`, `./scripts/run-harness.sh`, docs build, `go-code doctor --offline`, and OpenSpec validation

### Requirement: Feature verification updates PARITY.md
Every extension workflow marked verified in PARITY.md MUST name the test, harness scenario, or smoke check proving it.

#### Scenario: MCP status changes to verified
- **WHEN** PARITY.md marks MCP as verified
- **THEN** the row names the unit test or harness scenario that proves registration, permission behavior, and diagnostics
