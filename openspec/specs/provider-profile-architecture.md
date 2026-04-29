# Provider Profile Architecture

Extracted from `recenter-claudecodego-agent-roadmap` (archived).

## Architecture Decision

Go provides the shipped CLI/TUI runtime (agent loop, providers, tools, permissions, sessions, traces). Python provides deterministic evaluation (mock providers, replay analysis, CI harness scenarios).

## Provider Profile System

Providers are represented as explicit named compatibility profiles, not merely transports:

| Provider | Current Models | Transport | Status |
|----------|---------------|-----------|--------|
| Anthropic | claude-opus-4-6, claude-sonnet-4-6, claude-haiku-4 | Native | Verified |
| OpenAI | gpt-4o, gpt-4o-mini, o1, o3 | Native | Verified |
| DeepSeek | deepseek-v4-pro, deepseek-v4-flash | OpenAI-compatible | Verified |
| MiMo | mimo-v2.5-pro | OpenAI-compatible | Planned |
| Qwen | qwen-max, qwen-plus, qwen-turbo | OpenAI-compatible | Partial |
| GLM | glm-4-plus, glm-4, glm-4-flash | OpenAI-compatible | Partial |

## DeepSeek Model Names

- **Current**: `deepseek-v4-pro`, `deepseek-v4-flash`
- **Legacy (deprecated)**: `deepseek-chat`, `deepseek-reasoner` — keep as aliases with deprecation warnings

## MiMo-V2.5 Support

- Model: `mimo-v2.5-pro`
- Transport: OpenAI-compatible (verify against public docs)
- Positioned for long-context agentic coding harnesses

## Harness Engineering Gate

Every user-facing "supported" claim must be backed by:
1. Go unit test, OR
2. Python harness scenario, OR
3. Docs build check, OR
4. Manual verification note in PARITY.md

Release gate: `go test ./...` passes + `./scripts/run-harness.sh` passes.
