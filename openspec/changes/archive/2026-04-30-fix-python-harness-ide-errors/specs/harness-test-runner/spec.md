## ADDED Requirements

### Requirement: Pytest can discover and run harness tests
The pytest test runner SHALL successfully discover and execute all harness tests.

#### Scenario: Pytest runs from project root
- **WHEN** a developer runs `pytest harness/` from project root
- **THEN** pytest discovers all test files
- **AND** no ModuleNotFoundError occurs during test collection

#### Scenario: Conftest loads successfully
- **WHEN** pytest loads `harness/conftest.py`
- **THEN** all imports in conftest.py resolve successfully
- **AND** no ImportError is raised

#### Scenario: All harness tests pass
- **WHEN** pytest executes all harness tests
- **THEN** all tests complete (pass or expected skip)
- **AND** no test fails due to import errors

### Requirement: Harness has setup documentation
The harness module SHALL have documentation explaining how to run tests.

#### Scenario: README exists
- **WHEN** a developer opens `harness/README.md`
- **THEN** they find setup instructions for running tests

#### Scenario: Setup instructions are complete
- **WHEN** a developer follows the README instructions
- **THEN** they can successfully run `pytest harness/`
