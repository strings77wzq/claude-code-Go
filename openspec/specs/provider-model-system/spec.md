## Purpose
Define provider and model selection behavior, including explicit validation, safe runtime switching, adapter routing, and normalized provider errors.
## Requirements
### Requirement: Provider configuration is explicit and validated
The system SHALL validate provider, base URL, API key source, and model before starting an agent request.

#### Scenario: Unknown provider
- **WHEN** configuration names an unsupported provider
- **THEN** the system fails before making a network request and lists supported providers

#### Scenario: OpenAI-compatible provider
- **WHEN** the user configures an OpenAI-compatible base URL and model
- **THEN** the system routes requests through the OpenAI-compatible adapter and documents compatibility limits

### Requirement: Runtime model switching is safe
The system MUST support runtime model switching only when the selected model can be mapped to a provider adapter.

#### Scenario: Supported model switch
- **WHEN** the user enters `/model <supported-model>`
- **THEN** subsequent requests use that model and the UI confirms the active provider/model

#### Scenario: Unsupported model switch
- **WHEN** the user enters `/model <unknown-model>`
- **THEN** the system rejects the switch and keeps the previous model active

### Requirement: Provider errors are normalized
The system SHALL classify provider errors into auth, rate limit, timeout, server, network, invalid request, and unexpected categories.

#### Scenario: Rate limit error
- **WHEN** the provider returns a rate limit response
- **THEN** the system classifies it as rate limit and applies the configured retry policy

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
