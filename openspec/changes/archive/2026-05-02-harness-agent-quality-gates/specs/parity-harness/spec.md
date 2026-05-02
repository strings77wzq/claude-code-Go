## ADDED Requirements

### Requirement: Harness supports normalized comparison evidence
The harness SHALL support normalized evidence records that allow maintainers to compare this agent's behavior with manually supplied runs from other coding agents.

#### Scenario: External agent evidence is imported
- **WHEN** a maintainer supplies a normalized external run record
- **THEN** the harness report compares outcome, duration, tool count, and notes without requiring external credentials

### Requirement: Comparison reports avoid unsupported claims
The harness MUST distinguish measured local evidence from manual notes or inferred comparisons.

#### Scenario: Report includes manual competitor notes
- **WHEN** a comparison report includes manually entered notes about another agent
- **THEN** the notes are labeled as manual evidence rather than automated measurement
