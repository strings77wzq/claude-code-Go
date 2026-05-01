## ADDED Requirements

### Requirement: CI runs Python harness tests with editable install
The CI pipeline SHALL install the harness package in editable mode before running tests.

#### Scenario: Editable install step exists
- **WHEN** the CI python-harness job runs
- **THEN** it executes `pip install -e .` after installing requirements

#### Scenario: Tests run from project root
- **WHEN** the CI python-harness job runs tests
- **THEN** pytest is executed from the project root directory
- **AND** the test path includes `harness/`

#### Scenario: CI fails on test failures
- **WHEN** any harness test fails
- **THEN** the CI job reports failure
- **AND** `continue-on-error` is NOT set for the python-harness job
