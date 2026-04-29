## ADDED Requirements

### Requirement: Breakout Points Must Be Explicit
The project SHALL document its strongest open-source breakout points as evidence-backed differentiators rather than generic feature claims.

#### Scenario: User evaluates the project quickly
- **WHEN** a user opens the README or primary docs homepage
- **THEN** they see the core differentiators: Go single-binary runtime, Python harness reliability, local-agent safety, provider reach, and bilingual spec-driven engineering

#### Scenario: A differentiator lacks evidence
- **WHEN** a claimed differentiator is not backed by working code, tests, docs, or a parity row
- **THEN** it is labeled planned or experimental instead of presented as a current breakout point

### Requirement: Positioning Must Avoid Unsupported Clone Claims
The project SHALL describe itself as inspired by Claude Code-style workflows and Claw Code-style harness/product patterns, without claiming private Claude Code source compatibility unless legally usable source evidence exists.

#### Scenario: Public copy references Claude Code
- **WHEN** README, docs, or release notes reference Claude Code
- **THEN** the text frames the relationship as workflow inspiration or parity goals rather than proprietary-source equivalence

### Requirement: Public Narrative Must Be Narrow And Memorable
The project SHALL lead with a small set of memorable claims instead of an unranked list of every feature.

#### Scenario: README feature section is updated
- **WHEN** the README describes project value
- **THEN** it prioritizes the top breakout points before secondary roadmap features such as MCP, LSP, hooks, or ecosystem extensions

### Requirement: Showcase Must Demonstrate Real Agent Value
The project SHALL provide at least one showcase scenario that demonstrates a meaningful coding-agent workflow with visible inputs, tool activity, safety behavior, and final output.

#### Scenario: User opens showcase docs
- **WHEN** a user reads the showcase or watches the demo
- **THEN** they can understand what the agent did, which tools were used, what safety checks ran, and how to reproduce the scenario
