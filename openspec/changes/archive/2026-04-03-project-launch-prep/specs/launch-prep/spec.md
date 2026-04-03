## ADDED Requirements

### Requirement: Launch preparation script
The project SHALL provide a launch.sh script that prepares the codebase for GitHub deployment.

#### Scenario: Script execution
- **WHEN** a user runs `./launch.sh`
- **THEN** non-project directories are removed, go.sum is generated, binary is built, tests pass, and git is initialized

#### Scenario: Cleanup
- **WHEN** launch.sh runs
- **THEN** owncode-analysis/, claw-code-parity/, test/, bin/, .pytest_cache/ directories are removed

### Requirement: Git configuration files
The project SHALL include .gitignore and LICENSE files for GitHub readiness.

#### Scenario: Gitignore excludes
- **WHEN** git status is run
- **THEN** bin/, *.exe, __pycache__/, .pytest_cache/, venv/, .vscode/, .idea/, owncode-analysis/, claw-code-parity/, test/ are excluded

#### Scenario: License present
- **WHEN** a user checks the LICENSE file
- **THEN** MIT License text is present with current year and project name
