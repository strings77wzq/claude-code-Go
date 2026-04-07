## ADDED Requirements

### Requirement: Git repository excludes build artifacts
The repository SHALL have a comprehensive .gitignore that excludes build artifacts.

#### Scenario: Python artifacts excluded
- **WHEN** a developer runs `git status` after `pip install -e .`
- **THEN** no `*.egg-info/` directories appear as untracked

#### Scenario: IDE configs excluded
- **WHEN** a developer opens the project in VSCode or GoLand
- **THEN** no `.vscode/` or `.idea/` files appear as untracked

#### Scenario: Go build artifacts excluded
- **WHEN** a developer runs `go build`
- **THEN** no `bin/` or `dist/` files appear as untracked
