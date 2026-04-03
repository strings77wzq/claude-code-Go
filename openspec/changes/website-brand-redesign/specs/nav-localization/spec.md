## ADDED Requirements

### Requirement: Navigation bar follows locale
The navigation bar SHALL display content in the current locale language.

#### Scenario: English nav
- **WHEN** a user views the English site
- **THEN** nav shows: Guide, Architecture, Core Code, Tools, MCP, GitHub

#### Scenario: Chinese nav
- **WHEN** a user views the Chinese site
- **THEN** nav shows: 指南, 架构, 核心代码, 工具, MCP, GitHub

#### Scenario: Nav switches on language toggle
- **WHEN** a user toggles language
- **THEN** the nav bar updates to the selected language
