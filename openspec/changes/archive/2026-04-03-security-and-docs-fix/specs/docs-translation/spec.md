## ADDED Requirements

### Requirement: Chinese translation fixes
All Chinese documentation SHALL use correct technical terminology.

#### Scenario: Core terms preserved in English
- **WHEN** a user reads Chinese documentation
- **THEN** the following terms are preserved in English: Harness, Skills, Agent, Provider, Token

#### Scenario: No machine-translation artifacts
- **WHEN** a user reads Chinese documentation
- **THEN** the text reads naturally in Chinese, not like machine translation
