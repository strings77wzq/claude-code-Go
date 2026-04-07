## ADDED Requirements

### Requirement: IDE can resolve harness imports
The IDE (VSCode/GoLand) SHALL resolve harness imports without showing errors.

#### Scenario: VSCode settings exist
- **WHEN** a developer opens `.vscode/settings.json`
- **THEN** they see Python path configuration for harness

#### Scenario: Counter import resolves
- **WHEN** IDE analyzes `harness/replay/trace_analyzer.py`
- **THEN** the `Counter` import from `collections` shows no error

#### Scenario: Harness imports resolve
- **WHEN** IDE analyzes `harness/conftest.py`
- **THEN** imports like `from harness.mock_server.server import MockServer` show no error

#### Scenario: Mock server imports resolve
- **WHEN** IDE analyzes `harness/mock_server/app.py`
- **THEN** relative imports like `from .scenarios import registry` show no error
