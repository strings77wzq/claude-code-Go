## Purpose
Define the documentation experience required for a trustworthy product surface: users can reach a verified first success path, understand feature maturity, and trace public claims back to implemented behavior or reproducible evidence.

## Requirements

### Requirement: Documentation starts with verified quick success
The documentation SHALL guide new users from installation to a successful `doctor` check and first prompt using commands that are covered by tests or documented smoke checks.

#### Scenario: Chinese quick start
- **WHEN** a Chinese-speaking user opens the Chinese quick start
- **THEN** the page shows install, configure, doctor, first prompt, and troubleshooting steps in Chinese

#### Scenario: Failed doctor link
- **WHEN** doctor reports a failed check
- **THEN** the referenced documentation page explains the failure and remediation path

### Requirement: Docs distinguish implemented, planned, and experimental features
The documentation MUST clearly mark whether features are stable, experimental, planned, or unsupported.

#### Scenario: Planned IDE integration
- **WHEN** docs mention IDE integration
- **THEN** the page marks it as planned unless the implementation exists and is tested

### Requirement: Architecture docs teach the actual code
The documentation SHALL include architecture pages that map concepts to real packages, interfaces, and runtime flows.

#### Scenario: Agent loop doc
- **WHEN** a reader opens the agent loop architecture page
- **THEN** the page explains request construction, stop reasons, tool execution, history updates, compaction, recovery, and session save points

### Requirement: Public claims are evidence-backed
The documentation MUST remove or revise placeholder testimonials, stale metrics, unsupported benchmark comparisons, and claims that are not backed by reproducible evidence.

#### Scenario: Benchmark page
- **WHEN** the benchmark page lists performance numbers
- **THEN** it includes methodology, command, environment, date, and reproduction instructions
