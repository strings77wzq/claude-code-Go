## ADDED Requirements

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
