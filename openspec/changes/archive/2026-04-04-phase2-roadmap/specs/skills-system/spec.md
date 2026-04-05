## ADDED Requirements

### Requirement: Skills system
The system SHALL support custom skills defined via YAML configuration files.

#### Scenario: Skill definition
- **WHEN** a skill YAML is placed in `.go-code/skills/`
- **THEN** it is loaded and available as a slash command (e.g., `/review-pr`)

#### Scenario: Skill execution
- **WHEN** a user invokes a skill command
- **THEN** the skill's prompt is sent to the LLM as system context for the session
