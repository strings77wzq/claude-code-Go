## ADDED Requirements

### Requirement: MCP transport process cleanup
The MCP transport SHALL properly handle process cleanup errors.

#### Scenario: Process termination
- **WHEN** an MCP transport is closed
- **THEN** Kill() and Wait() errors are logged, not ignored
