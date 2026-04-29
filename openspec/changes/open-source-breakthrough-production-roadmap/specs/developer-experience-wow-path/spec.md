## ADDED Requirements

### Requirement: First Success Must Be Achievable In Minutes
The project SHALL provide a first-run path that lets a new user build or install, run a health check, execute a local demo or prompt, and inspect output within five minutes on a normal development machine.

#### Scenario: New user follows quickstart
- **WHEN** a new user follows the primary quickstart from a clean checkout
- **THEN** they can reach a visible successful agent output without reading architecture docs first

### Requirement: Offline Demo Must Avoid Paid Provider Dependency
The project SHALL provide an offline-capable demo path using mock provider or harness infrastructure so users can evaluate behavior without API keys or token spend.

#### Scenario: User has no API key
- **WHEN** a user runs the documented offline demo path
- **THEN** the project demonstrates provider streaming, at least one tool interaction, and session/replay evidence without contacting a live provider

### Requirement: Doctor Must Guide The User To Success
The project SHALL make `doctor` or equivalent health checks the first troubleshooting surface for config, provider, tools, permissions, sessions, docs, and environment readiness.

#### Scenario: Health check finds a problem
- **WHEN** doctor detects missing config, invalid provider, unavailable tool, bad session path, or docs mismatch
- **THEN** it prints a specific fix and points to the relevant docs

### Requirement: Replay Must Make Agent Behavior Inspectable
The project SHALL make session replay or trace inspection part of the developer experience so users can understand what the agent did.

#### Scenario: Demo completes
- **WHEN** a demo or prompt run finishes
- **THEN** the user can replay or inspect the session to see messages, tool calls, permission decisions, and final output

### Requirement: Quickstart Must Support Global And China-Friendly Providers
The project SHALL include first-run configuration paths for Anthropic/OpenAI-compatible providers and China-friendly model families such as DeepSeek and MiMo when their profiles are implemented.

#### Scenario: User chooses DeepSeek or MiMo
- **WHEN** a user follows provider setup docs for DeepSeek or MiMo
- **THEN** they see model IDs, environment variables, base URL behavior, compatibility status, and verification commands
