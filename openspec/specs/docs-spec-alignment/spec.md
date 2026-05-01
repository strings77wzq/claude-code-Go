# docs-spec-alignment Specification

## Purpose
TBD - created by archiving change recenter-claudecodego-agent-roadmap. Update Purpose after archive.
## Requirements
### Requirement: Public Claims Must Match Verified Implementation
Docs, README, parity tables, and OpenSpec task status SHALL distinguish supported, partial, planned, and removed behavior based on implementation and verification evidence.

#### Scenario: README describes a capability
- **WHEN** README or docs describe a capability as available
- **THEN** the capability has working code plus verification evidence, or the text labels it as planned/experimental

#### Scenario: Spec task is marked complete
- **WHEN** an OpenSpec task is checked complete
- **THEN** the task references passing verification or a deliberate not-tested note

### Requirement: Bilingual Docs Must Share The Same Product Truth
Chinese and English docs SHALL present the same supported provider list, quick-start path, architecture boundaries, parity status, and known gaps.

#### Scenario: DeepSeek or MiMo docs are updated
- **WHEN** provider support docs change in English or Chinese
- **THEN** the corresponding language page is updated or marked as pending translation with the same support status

### Requirement: Parity Must Be A Matrix, Not A Marketing Claim
The project SHALL maintain a parity/status matrix for Claude Code-style workflows that records status, evidence, known gaps, and next task links.

#### Scenario: User evaluates Claude Code-style compatibility
- **WHEN** a user opens `PARITY.md` or the docs parity page
- **THEN** they see workflow-level status for prompt mode, streaming, tools, edits, bash, permissions, sessions, replay, compacting, provider switching, MCP, recovery, and docs

### Requirement: Architecture Docs Must Explain The Go/Python Boundary
Architecture docs SHALL explain which responsibilities belong to the Go runtime and which belong to the Python harness, including why this split exists.

#### Scenario: Contributor chooses where to implement a change
- **WHEN** a contributor reads the architecture docs
- **THEN** they can decide whether a change belongs in Go runtime code, Python harness code, or documentation/specs

### Requirement: Public docs distinguish verified, partial, experimental, and planned support
Public documentation SHALL label each user-facing feature according to the evidence available in PARITY.md and OpenSpec.

#### Scenario: MCP is described before full productization
- **WHEN** docs mention MCP support before v0.3 extension gates pass
- **THEN** the docs label MCP as partial or experimental
- **AND** they link to configuration limits and verification status

#### Scenario: LSP is described before command exposure is complete
- **WHEN** docs mention LSP capabilities before they are user-facing
- **THEN** the docs label LSP as planned or experimental rather than verified

### Requirement: English and Chinese docs stay status-aligned
The English and Chinese documentation SHALL present the same feature status, provider model names, commands, and known limitations.

#### Scenario: English roadmap changes
- **WHEN** `docs/roadmap.md` changes a feature status
- **THEN** the corresponding Chinese roadmap or guide page is updated in the same change

### Requirement: Marketing claims require evidence
The docs SHALL avoid testimonials, benchmark numbers, parity claims, or competitive superiority statements unless they reference reproducible evidence.

#### Scenario: Benchmark claim is added
- **WHEN** a page claims speed, lower latency, or better performance
- **THEN** it includes the benchmark command, date, environment, and raw result location

### Requirement: Parked business proposals are not presented as roadmap commitments
Enterprise and content-marketing proposals SHALL not appear as committed roadmap items until they have approved tasks and implementation scope.

#### Scenario: Roadmap is updated
- **WHEN** enterprise or content-marketing ideas are mentioned
- **THEN** they are labeled as parked or future concepts unless an active approved OpenSpec change exists

