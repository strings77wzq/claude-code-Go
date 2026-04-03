## MODIFIED Requirements

### Requirement: Natural Chinese documentation
All Chinese documentation SHALL read naturally to native Chinese developers.

#### Scenario: No translation artifacts
- **WHEN** a native Chinese developer reads the documentation
- **THEN** the text reads naturally, not like machine translation

#### Scenario: Technical terms preserved
- **WHEN** technical terms appear in Chinese docs
- **THEN** the following are kept in English: Roadmap, Skills, Harness, Agent Loop, Provider, MCP, Hooks, Token
