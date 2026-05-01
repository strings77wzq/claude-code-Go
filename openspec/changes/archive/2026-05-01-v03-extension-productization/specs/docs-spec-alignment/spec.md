## ADDED Requirements

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
