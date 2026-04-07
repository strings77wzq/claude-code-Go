## ADDED Requirements

### Requirement: Telemetry framework implemented
The project SHALL have optional anonymous usage tracking.

#### Scenario: Telemetry is opt-in
- **WHEN** a user first runs the tool
- **THEN** they are asked for consent to telemetry

#### Scenario: No PII collected
- **WHEN** telemetry is enabled
- **THEN** no personally identifiable information is sent

#### Scenario: Transparent data collection
- **WHEN** a user checks the code
- **THEN** they can see exactly what data is collected

#### Scenario: Easy opt-out
- **WHEN** a user wants to disable telemetry
- **THEN** they can do so via config or CLI flag
