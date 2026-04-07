## ADDED Requirements

### Requirement: CI runs Python harness with editable install
The CI pipeline SHALL install and run Python harness tests correctly.

#### Scenario: Editable install step
- **WHEN** CI runs the python-harness job
- **THEN** it executes `pip install -e .` after requirements

#### Scenario: Tests from project root
- **WHEN** CI runs harness tests
- **THEN** pytest is executed from project root
- **AND** the test path includes `harness/`

#### Scenario: CI fails on test failure
- **WHEN** any harness test fails
- **THEN** the CI job reports failure
- **AND** `continue-on-error` is NOT set
