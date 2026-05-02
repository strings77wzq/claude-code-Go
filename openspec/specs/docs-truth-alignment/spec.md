# docs-truth-alignment Specification

## Purpose
TBD - created by archiving change v02-consolidation-release. Update Purpose after archive.
## Requirements
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

### Requirement: Generated documentation has declared source
Generated documentation MUST identify its source files or generation command so reviewers can trace changes to source-of-truth inputs.

#### Scenario: Generated docs change in a PR
- **WHEN** generated documentation changes
- **THEN** the corresponding source change or generation command is identified in review evidence

#### Scenario: Documentation inventory is inspected
- **WHEN** a maintainer reviews release readiness
- **THEN** the documentation inventory distinguishes source-of-truth docs from generated outputs

### Requirement: Documentation drift is checked before release
The project SHALL run a docs drift check before publishing a release.

#### Scenario: Generated docs are stale
- **WHEN** generated docs do not match their source inputs
- **THEN** the release readiness check fails with the stale output identified

### Requirement: Claims distinguish implemented and planned behavior
Documentation MUST distinguish shipped behavior from planned roadmap behavior.

#### Scenario: README mentions a planned feature
- **WHEN** documentation describes a feature that is not implemented in the current release
- **THEN** the documentation labels it as planned or experimental
