## ADDED Requirements

### Requirement: Provider profiles are transport-independent
Provider profiles SHALL describe provider identity, model defaults, compatibility flags, and capability metadata without embedding transport-specific request logic.

#### Scenario: Profile is loaded for an existing transport
- **WHEN** a provider profile is selected for a provider that already has a transport
- **THEN** the profile configures model defaults and capabilities without duplicating the transport implementation

### Requirement: Provider capabilities are diagnosable
The system SHALL expose provider profile and selected model capabilities to doctor, TUI status, and trace summaries.

#### Scenario: Model lacks a capability
- **WHEN** the selected model does not advertise a requested capability
- **THEN** the diagnostic identifies the provider, model, missing capability, and fallback behavior
