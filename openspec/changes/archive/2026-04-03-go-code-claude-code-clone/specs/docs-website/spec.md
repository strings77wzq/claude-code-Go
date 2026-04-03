## ADDED Requirements

### Requirement: MkDocs website
The project SHALL provide a documentation website built with MkDocs.

#### Scenario: Local development
- **WHEN** `make docs` is run
- **THEN** the MkDocs dev server starts and the site is accessible at localhost:8000

#### Scenario: Site structure
- **WHEN** the site is viewed
- **THEN** it includes Home, Guide (Install, Quick Start, Config), Architecture (Overview, Agent Loop, Tools), and Harness sections

### Requirement: Installation guide
The documentation SHALL provide clear installation instructions.

#### Scenario: Go installation
- **WHEN** a user reads the install guide
- **THEN** they can install the CLI via `go install` or download a pre-built binary

### Requirement: Quick start guide
The documentation SHALL provide a quick start tutorial.

#### Scenario: First run
- **WHEN** a user follows the quick start
- **THEN** they can set up their API key, run the CLI, and have their first conversation

### Requirement: Architecture documentation
The documentation SHALL explain the system architecture.

#### Scenario: Architecture overview
- **WHEN** a user reads the architecture docs
- **THEN** they understand the Agent Loop, Tool System, Permission System, and MCP Integration

### Requirement: GitHub Pages deployment
The documentation SHALL be deployable to GitHub Pages.

#### Scenario: Deploy to GitHub Pages
- **WHEN** the CI workflow runs on main branch
- **THEN** the MkDocs site is built and deployed to GitHub Pages
