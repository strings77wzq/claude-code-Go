## ADDED Requirements

### Requirement: English homepage
The English homepage SHALL present the project as a professional product.

#### Scenario: Hero section
- **WHEN** a user visits the English homepage
- **THEN** they see: project name "claude-code-Go", tagline "Claude Code in Go — AI-powered coding assistant", one-click install command, and GitHub stars badge

#### Scenario: Feature cards
- **WHEN** a user scrolls the English homepage
- **THEN** they see 6 feature cards: Agent Loop, 6 Built-in Tools, Permission System, MCP Integration, SSE Streaming, Context Management

#### Scenario: Install command
- **WHEN** a user views the install command
- **THEN** it is displayed in a copyable code block: `go install github.com/strings77wzq/claude-code-Go@latest`

### Requirement: Chinese homepage
The Chinese homepage SHALL present the same content in natural, professional Chinese.

#### Scenario: Chinese hero section
- **WHEN** a user visits the Chinese homepage
- **THEN** they see: project name "claude-code-Go", Chinese tagline, Chinese install instructions

#### Scenario: Chinese feature cards
- **WHEN** a user scrolls the Chinese homepage
- **THEN** all 6 feature cards are in professional Chinese technical writing style

### Requirement: English documentation
The English documentation SHALL cover all major aspects of the project.

#### Scenario: Guide pages
- **WHEN** a user navigates to Guide sections
- **THEN** Installation, Quick Start, Configuration pages are available in English

#### Scenario: Architecture pages
- **WHEN** a user navigates to Architecture sections
- **THEN** Overview, Agent Loop, Tools pages are available in English

### Requirement: Chinese documentation
The Chinese documentation SHALL provide accurate, professional translations.

#### Scenario: Chinese guide pages
- **WHEN** a user navigates to Chinese Guide sections
- **THEN** Installation, Quick Start, Configuration pages are available in professional Chinese

#### Scenario: Chinese architecture pages
- **WHEN** a user navigates to Chinese Architecture sections
- **THEN** Overview, Agent Loop, Tools pages are available in professional Chinese
