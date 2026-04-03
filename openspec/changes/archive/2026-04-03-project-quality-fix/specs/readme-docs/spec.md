## ADDED Requirements

### Requirement: README accuracy
The README SHALL accurately reflect the actual project state.

#### Scenario: Project name matches repository
- **WHEN** a user views the README
- **THEN** the project name is "claude-code-Go" matching the repository name

#### Scenario: Clone URL is correct
- **WHEN** a user copies the clone command
- **THEN** the URL is `https://github.com/strings77wzq/claude-code-Go`

#### Scenario: Build commands are current
- **WHEN** a user follows the build instructions
- **THEN** all commands work without errors (vitepress, not mkdocs)

#### Scenario: Project structure is accurate
- **WHEN** a user views the project structure diagram
- **THEN** it reflects the actual directory layout including `docs/en/` and `docs/zh/`

### Requirement: README completeness
The README SHALL include elements expected of a professional open source project.

#### Scenario: Star encouragement
- **WHEN** a user views the README
- **THEN** there is a polite request to star the project if found useful

#### Scenario: Demo section
- **WHEN** a user views the README
- **THEN** there is a placeholder or description for a demo GIF/screenshot

#### Scenario: License is visible
- **WHEN** a user views the README
- **THEN** the MIT license is mentioned with a link to LICENSE file
