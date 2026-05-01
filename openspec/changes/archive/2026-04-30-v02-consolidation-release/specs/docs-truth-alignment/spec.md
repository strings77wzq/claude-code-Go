## ADDED Requirements

### Requirement: README claims match verified implementation
The README SHALL only claim features that are verified by automated tests or documented manual smoke checks. Any feature listed in the README that PARITY.md marks as `partial` or `planned` SHALL be explicitly labeled with its status (e.g., "Experimental", "Planned (v0.3)").

#### Scenario: README feature list matches PARITY.md
- **WHEN** the README is compared against PARITY.md
- **THEN** every claimed feature in README has status `verified` in PARITY.md, or is labeled with its actual status
- **AND** no feature marked `unsupported` in PARITY.md appears as a supported feature in README

#### Scenario: Placeholder content is removed
- **WHEN** the README and docs are inspected
- **THEN** the demo GIF placeholder note is removed (replaced with a text description until a real recording exists)
- **AND** placeholder testimonials are removed
- **AND** benchmark numbers are replaced with "Methodology defined, results pending" or removed

### Requirement: Chinese docs sync with English docs
The Chinese documentation (`docs/zh/`) SHALL reflect the same feature status, model names, and installation instructions as the English docs. Any translation gap SHALL be explicitly marked.

#### Scenario: Chinese provider docs match English
- **WHEN** the Chinese provider docs are compared against the English provider docs
- **THEN** both list the same model names and compatibility status
- **AND** both link to the same verification evidence

### Requirement: PARITY.md reflects verified-only status
PARITY.md SHALL only mark a workflow as `verified` when it has both an automated test AND a documented manual smoke check path. Workflows currently marked `verified` without evidence SHALL be downgraded to `partial`.

#### Scenario: Each verified workflow has linked evidence
- **WHEN** a PARITY.md row is inspected
- **THEN** any `verified` status includes a link to the test file or smoke check doc that proves it
