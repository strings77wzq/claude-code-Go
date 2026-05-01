## ADDED Requirements

### Requirement: Replay includes extension and permission events
Replay output SHALL include enough structured information to debug MCP, LSP, hook, skill, and permission behavior without calling a real provider.

#### Scenario: MCP tool is denied
- **WHEN** a saved session includes an MCP tool call denied by policy
- **THEN** replay shows the tool name, permission decision, denial reason, and final agent-visible result

#### Scenario: LSP capability is unavailable
- **WHEN** a saved session attempts an unavailable LSP operation
- **THEN** replay shows the unavailable capability, configuration state, and non-fatal result

### Requirement: Trace output redacts sensitive values
Trace and replay output MUST redact API keys, provider tokens, authorization headers, and configured secrets.

#### Scenario: Provider configuration contains a secret
- **WHEN** replay or trace output includes provider configuration
- **THEN** secret values are replaced with a redacted marker

### Requirement: Replay has a release-evidence mode
The replay command SHALL provide a concise mode suitable for attaching to release evidence or issue reports.

#### Scenario: Maintainer collects replay evidence
- **WHEN** the maintainer runs replay in evidence mode
- **THEN** the output summarizes prompts, tool calls, permission decisions, extension availability, errors, and final status
- **AND** it omits raw provider secrets and excessive token payloads
