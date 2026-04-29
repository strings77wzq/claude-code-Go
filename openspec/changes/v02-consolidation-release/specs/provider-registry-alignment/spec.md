## ADDED Requirements

### Requirement: Model registry contains current model names
The model registry in `internal/provider/registry/registry.go` SHALL list the current preferred model names for all supported providers: Anthropic (claude-opus-4-6, claude-sonnet-4-6, claude-haiku-4), OpenAI (gpt-4o, gpt-4o-mini, o1, o3), DeepSeek (deepseek-v4-pro, deepseek-v4-flash), Qwen (qwen-max, qwen-plus, qwen-turbo), and GLM (glm-4-plus, glm-4, glm-4-flash).

#### Scenario: DeepSeek models use current names
- **WHEN** a user configures model `deepseek-v4-pro` or `deepseek-v4-flash`
- **THEN** the model is found in the registry with provider `openai` (OpenAI-compatible transport)
- **AND** the system constructs a valid API request

#### Scenario: Legacy DeepSeek model names are rejected with guidance
- **WHEN** a user configures model `deepseek-chat` or `deepseek-reasoner`
- **THEN** the system logs a deprecation warning
- **AND** suggests migrating to `deepseek-v4-pro` or `deepseek-v4-flash`

### Requirement: Unknown model passthrough
When a model name is not found in the registry, the system SHALL not reject it outright. Instead, it SHALL construct a best-effort provider configuration by inferring the provider from name heuristics (e.g., `gpt-*` → openai, `claude-*` → anthropic) or defaulting to the configured base URL's transport, log a warning, and proceed with the request.

#### Scenario: Unknown model with recognizable prefix
- **WHEN** a user configures model `gpt-5-turbo` (not in registry)
- **THEN** the system infers provider `openai` from the `gpt-` prefix
- **AND** logs a warning: "model gpt-5-turbo not in verified registry, proceeding with inferred provider openai"
- **AND** constructs a valid API request

#### Scenario: Unknown model with no recognizable prefix
- **WHEN** a user configures model `my-custom-model` (not in registry, no recognizable prefix)
- **THEN** the system defaults to the configured base URL's transport or anthropic transport
- **AND** logs a warning with a link to the verified models documentation

### Requirement: MiMo-V2.5 provider profile
The provider system SHALL include a MiMo provider profile with `mimo-v2.5-pro` as the default model. If the official MiMo API is OpenAI-compatible, the profile SHALL use the shared OpenAI-compatible transport. If it requires a different API shape, a narrow adapter SHALL be added.

#### Scenario: MiMo profile resolves with OpenAI-compatible transport
- **WHEN** a user configures provider `mimo` and model `mimo-v2.5-pro`
- **THEN** the system selects the MiMo provider profile
- **AND** routes requests through the OpenAI-compatible transport (if verified compatible)
- **AND** documents the compatibility status in provider docs
