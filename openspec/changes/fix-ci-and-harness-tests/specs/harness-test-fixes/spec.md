## ADDED Requirements

### Requirement: All harness tests pass without assertion errors
All harness tests SHALL pass when run locally and in CI.

#### Scenario: test_load_from_lines passes
- **WHEN** `TestSessionReplayer.test_load_from_lines` runs
- **THEN** the assertion for `tool_calls` matches the actual count
- **AND** no assertion error is raised

#### Scenario: test_timeline passes
- **WHEN** `TestTraceAnalyzer.test_timeline` runs
- **THEN** the assertion for timeline entry fields uses correct key names
- **AND** no KeyError is raised

#### Scenario: All tests pass
- **WHEN** `pytest harness/ -v` runs
- **THEN** all non-skipped tests pass
- **AND** no test fails due to assertion errors
