## ADDED Requirements

### Requirement: Architecture audit drives roadmap priorities
The project SHALL maintain an architecture audit that maps current implementation surfaces to prioritized risks, owner modules, and implementation milestones.

#### Scenario: Maintainer reviews next work
- **WHEN** a maintainer opens the roadmap hardening change
- **THEN** the design and tasks identify current architecture strengths, current risks, affected packages, priority order, and verification gates

#### Scenario: New issue is discovered during implementation
- **WHEN** implementation reveals a new architectural or product risk
- **THEN** the task list is updated with the risk, priority, acceptance criteria, and required verification before the issue is worked

### Requirement: OpenSpec artifacts remain implementation-ready
OpenSpec proposals, designs, specs, and tasks SHALL be specific enough for implementation without relying on unstated roadmap knowledge.

#### Scenario: Task is selected for implementation
- **WHEN** a task is started
- **THEN** it states the target files or surfaces, expected behavior, tests to add or update, and completion evidence

#### Scenario: Spec is archived into the baseline
- **WHEN** a change is archived
- **THEN** resulting baseline specs include a concrete Purpose section, normative requirements, scenarios, and no placeholder `TBD` purpose text

### Requirement: Agent usability is measured by runnable workflows
The project SHALL evaluate agent quality through runnable developer workflows rather than only feature inventory.

#### Scenario: Author compares the agent to other coding agents
- **WHEN** the author runs the project against a trial workflow
- **THEN** the workflow records task success, tool correctness, permission experience, error recovery, session replay usefulness, and known gaps

#### Scenario: Public docs claim support for a workflow
- **WHEN** docs or `PARITY.md` mark a workflow as verified
- **THEN** at least one Go test, harness scenario, release smoke check, or documented manual evidence item proves the claim

### Requirement: Stabilization precedes expansion
The project MUST complete core runtime stabilization gates before expanding extension features or marketing claims.

#### Scenario: Extension feature is proposed
- **WHEN** a new MCP, LSP, hooks, skills, IDE, or provider extension feature is proposed
- **THEN** the proposal verifies that prompt mode, TUI lifecycle, permission modes, session trace/replay, and harness gates remain passing

#### Scenario: Release candidate is prepared
- **WHEN** a release candidate is prepared
- **THEN** it runs Go tests, Go vet, harness tests, docs build, OpenSpec strict validation, doctor smoke checks, binary build, and any milestone-specific scenario tests

### Requirement: Roadmap tasks include stop conditions
The project SHALL define each roadmap task with clear acceptance criteria and stop conditions so execution can pause safely when evidence is missing.

#### Scenario: Task cannot be verified locally
- **WHEN** a task depends on credentials, external services, platform-specific behavior, or unavailable tooling
- **THEN** the task records the blocker, next-best local verification, and manual verification instructions instead of being marked complete silently

#### Scenario: Task scope expands
- **WHEN** a task uncovers work outside its acceptance criteria
- **THEN** the extra work is captured as a new task or follow-up issue before continuing
