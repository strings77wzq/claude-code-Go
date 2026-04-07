## ADDED Requirements

### Requirement: GitHub Security scanning enabled
The repository SHALL have automated security scanning via GitHub Actions.

#### Scenario: CodeQL for Go
- **WHEN** code is pushed to main or a PR is opened
- **THEN** CodeQL analysis runs for Go code
- **AND** vulnerabilities are reported in Security tab

#### Scenario: CodeQL for Python
- **WHEN** code is pushed to main or a PR is opened
- **THEN** CodeQL analysis runs for Python code
- **AND** vulnerabilities are reported in Security tab

#### Scenario: Dependency review
- **WHEN** a PR updates dependencies
- **THEN** dependency review action runs
- **AND** vulnerable dependencies are flagged
