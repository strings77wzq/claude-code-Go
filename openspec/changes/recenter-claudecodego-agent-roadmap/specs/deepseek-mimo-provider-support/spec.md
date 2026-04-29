## ADDED Requirements

### Requirement: DeepSeek Must Be A First-Class Provider Profile
The system SHALL support DeepSeek as an explicit provider profile with documented base URLs, model IDs, environment variables, compatibility notes, and tests.

#### Scenario: User configures DeepSeek
- **WHEN** a user selects the DeepSeek profile with an API key and `deepseek-v4-pro` or `deepseek-v4-flash`
- **THEN** the system resolves the correct base URL, provider transport, model metadata, and runtime configuration without requiring a generic custom-provider workaround

#### Scenario: User uses legacy DeepSeek aliases
- **WHEN** a user configures `deepseek-chat` or `deepseek-reasoner`
- **THEN** the system either supports the alias with a deprecation warning and migration guidance or rejects it with a clear error that names the replacement models

### Requirement: MiMo-V2.5 Must Be A First-Class Provider Profile
The system SHALL support the MiMo-V2.5 series as an explicit provider profile beginning with `mimo-v2.5-pro`, including config examples, model metadata, compatibility notes, and tests.

#### Scenario: User configures MiMo-V2.5-Pro
- **WHEN** a user selects the MiMo profile with `mimo-v2.5-pro`
- **THEN** the system resolves a MiMo-specific profile and does not require the user to guess a generic OpenAI-compatible model prefix

#### Scenario: MiMo endpoint behavior is not fully verified
- **WHEN** official MiMo API details are unavailable or incomplete
- **THEN** docs and runtime errors mark the missing endpoint assumptions explicitly instead of claiming full support

### Requirement: Provider Profiles Must Preserve Transport Reuse
The system SHALL allow provider profiles to reuse OpenAI-compatible or Anthropic-compatible transports internally while exposing user-facing provider names and defaults that match the model family.

#### Scenario: DeepSeek uses OpenAI-compatible transport
- **WHEN** DeepSeek is configured through the profile
- **THEN** the request path uses the compatible transport while logs, errors, and docs identify the provider as DeepSeek

#### Scenario: Provider-specific behavior differs
- **WHEN** a provider requires non-standard request fields, streaming behavior, or error handling
- **THEN** the provider profile captures that behavior with tests rather than hiding it in generic transport code

### Requirement: Model Metadata Must Be Verifiable
The system SHALL store model metadata for supported first-class models, including provider family, model ID, context capability when known, default reasoning/thinking options when known, and documentation source.

#### Scenario: User lists models
- **WHEN** a user runs the model listing command or reads provider docs
- **THEN** DeepSeek and MiMo model entries show current preferred IDs and compatibility status

#### Scenario: Upstream docs change
- **WHEN** official model IDs, pricing, deprecation, or base URLs change
- **THEN** updating the metadata and docs is treated as a required maintenance task before release claims are refreshed
