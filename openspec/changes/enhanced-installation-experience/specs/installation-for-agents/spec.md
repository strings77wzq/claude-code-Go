## ADDED Requirements

### Requirement: Installation Guide for LLM Agents
A dedicated markdown file SHALL exist at `docs/guide/installation-for-agents.md` providing step-by-step instructions for AI agents.

#### Scenario: Agent reads the guide
- **WHEN** an AI agent fetches the installation-for-agents.md file
- **THEN** it contains clear, structured steps to install and configure go-code for the user

#### Scenario: Agent asks user about provider
- **WHEN** an agent follows the guide
- **THEN** it first asks the user which LLM provider they want to use

#### Scenario: Agent completes setup
- **WHEN** an agent follows all steps
- **THEN** go-code is installed, configured with an API key, and verified
