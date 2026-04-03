## ADDED Requirements

### Requirement: Graceful shutdown
The system SHALL clean up resources before exiting on signal.

#### Scenario: Signal received
- **WHEN** SIGINT or SIGTERM is received
- **THEN** the system cancels context, saves session, closes MCP connections, and exits
