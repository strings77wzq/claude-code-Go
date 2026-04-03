## ADDED Requirements

### Requirement: Hero section clarity
The homepage hero section SHALL present a clean, non-redundant brand identity.

#### Scenario: No visual redundancy
- **WHEN** a user views the English homepage hero
- **THEN** the project name "claude-code-Go" appears once, without repeating "Claude Code in Go"

#### Scenario: Chinese hero section
- **WHEN** a user views the Chinese homepage hero
- **THEN** the tagline "模型提供智能，harness 提供可靠性" is prominently displayed

### Requirement: CTA section
The homepage SHALL include a call-to-action section below the feature cards.

#### Scenario: English CTA
- **WHEN** a user scrolls to the bottom of the English homepage
- **THEN** they see: "Ready to try it?" with links to Get Started and GitHub

#### Scenario: Chinese CTA
- **WHEN** a user scrolls to the bottom of the Chinese homepage
- **THEN** they see: "准备开始？" with links to 快速开始 and 查看源码

### Requirement: Use cases section
The homepage SHALL show what users can do with the project.

#### Scenario: Use cases display
- **WHEN** a user views the homepage
- **THEN** they see 3-4 use case examples: code review, refactoring, debugging, boilerplate generation
