## ADDED Requirements

### Requirement: Harness enforces latency budgets
The harness SHALL allow scenarios to declare latency budgets and report budget violations separately from functional assertion failures.

#### Scenario: Scenario exceeds latency budget
- **WHEN** a scenario passes functional assertions but exceeds its latency budget
- **THEN** the result records a latency budget violation

### Requirement: Harness validates trace invariants
The harness MUST validate required trace invariants for scenarios that exercise tools, permissions, extensions, or recovery.

#### Scenario: Permission scenario lacks decision trace
- **WHEN** a permission scenario completes without a permission decision event in trace
- **THEN** the scenario fails its trace invariant
