## MODIFIED Requirements

### Requirement: Accurate README
The README SHALL accurately reflect the current project state.

#### Scenario: Default model
- **WHEN** a user reads the README
- **THEN** the default model is claude-sonnet-4-6-20251001

#### Scenario: Project structure
- **WHEN** a user reads the project structure
- **THEN** it includes provider/, cost/, lsp/, scripts/ directories

#### Scenario: Tool count
- **WHEN** a user reads the Features section
- **THEN** it shows 10 built-in tools including TodoWrite

#### Scenario: Supported providers
- **WHEN** a user reads the Supported Providers section
- **THEN** it includes Anthropic, OpenAI, DeepSeek, Qwen, GLM, Tencent Cloud
