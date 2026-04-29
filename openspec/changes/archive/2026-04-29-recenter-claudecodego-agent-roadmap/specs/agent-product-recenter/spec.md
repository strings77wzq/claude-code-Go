## ADDED Requirements

### Requirement: Product Identity Must Be Explicit
The project SHALL define itself as a Go-first Claude Code-style coding agent runtime with a Python verification harness, and SHALL describe Claude Code / Claw Code influence as observable workflow and harness inspiration unless legally usable source is present.

#### Scenario: User reads the public project overview
- **WHEN** a user opens the README or primary docs homepage
- **THEN** the project identity states the Go runtime and Python harness split without claiming private Claude Code source compatibility

#### Scenario: Contributor evaluates scope
- **WHEN** a contributor reads architecture or roadmap docs
- **THEN** they can distinguish supported runtime behavior, harness infrastructure, and planned future work

### Requirement: Roadmap Must Be Staged Around Runnable Checkpoints
The project SHALL organize the next development roadmap into staged checkpoints that can each be verified by commands, tests, or harness scenarios.

#### Scenario: A roadmap item is accepted
- **WHEN** a roadmap item is marked ready for implementation
- **THEN** it includes the affected runtime surface, required tests or harness scenarios, and documentation updates

#### Scenario: A broad feature is proposed
- **WHEN** a feature spans multiple modules or public claims
- **THEN** it is split into smaller phases before implementation begins

### Requirement: Core Agent UX Must Precede Extension Ecosystem Work
The project SHALL prioritize prompt mode, model configuration, tool use, permissions, sessions, replay, doctor diagnostics, and harness coverage before expanding MCP, LSP, hooks, skills, or plugin ecosystem claims.

#### Scenario: Extension work competes with core reliability
- **WHEN** an extension-surface task lacks tested core prerequisites
- **THEN** the task is deferred or scoped to documentation/status cleanup rather than productized as complete

### Requirement: Existing Strengths Must Be Preserved
The project SHALL preserve and strengthen the existing provider abstraction, doctor diagnostics, permission policy, session replay, trace logging, Python harness, docs site, and OpenSpec workflow unless a replacement is explicitly justified.

#### Scenario: Refactor touches an existing strength
- **WHEN** implementation changes one of the existing core assets
- **THEN** the change keeps behavior covered by tests or documents why the old behavior was removed
