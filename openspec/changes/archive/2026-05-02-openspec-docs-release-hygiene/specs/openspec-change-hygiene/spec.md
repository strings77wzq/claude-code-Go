## ADDED Requirements

### Requirement: Active changes are implementation-sized
OpenSpec implementation changes SHALL be scoped so their tasks can be implemented, verified, and archived as a coherent unit.

#### Scenario: Roadmap is too broad for direct apply
- **WHEN** a change contains broad roadmap work across unrelated runtime, docs, harness, and release surfaces
- **THEN** the project treats it as an umbrella and creates smaller implementation changes before apply

### Requirement: Specs have durable purpose and testable scenarios
Each active or release-relevant spec MUST include a concrete purpose and at least one testable scenario per requirement.

#### Scenario: Spec has placeholder purpose
- **WHEN** a release-relevant spec contains a placeholder purpose
- **THEN** the hygiene gate fails until the purpose is replaced with durable project context

### Requirement: Task completion records evidence
Completed OpenSpec tasks MUST be backed by validation evidence or an explicit not-tested note.

#### Scenario: Task is marked complete
- **WHEN** a task checkbox is marked complete
- **THEN** the change notes or task context identify the test, command, review evidence, or accepted validation gap

#### Scenario: Active change inventory is produced
- **WHEN** release hygiene checks are run
- **THEN** the project has an active-change inventory that identifies umbrella changes, implementation changes, and expected evidence

### Requirement: Strict validation gates archive readiness
The project MUST run strict OpenSpec validation before archiving a change.

#### Scenario: Maintainer prepares archive
- **WHEN** a change is ready to archive
- **THEN** `openspec validate <change> --strict` passes before archive proceeds
