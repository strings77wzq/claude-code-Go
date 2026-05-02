---
title: Provider System
description: Multi-provider support for LLM backends — interface design, Anthropic and OpenAI adapters, configuration
---

# Provider System

go-code supports multiple LLM providers through a unified provider interface. This document describes the provider architecture and how to configure different backends.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Provider Architecture                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌──────────────┐                       │
│   │ Agent Loop  │─────────▶│ Provider     │                       │
│   │             │          │ Interface    │                       │
│   └─────────────┘          └──────┬───────┘                       │
│                                   │                                │
│                    ┌──────────────┼──────────────┐                 │
│                    ▼              ▼              ▼                 │
│            ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│            │ Anthropic │  │  OpenAI   │  │ Custom    │             │
│            │ Provider  │  │ Provider  │  │ Provider  │             │
│            └───────────┘  └───────────┘  └───────────┘             │
│                    │              │              │                 │
│                    ▼              ▼              ▼                 │
│            ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│            │ Anthropic │  │ OpenAI    │  │ Custom    │             │
│            │   API     │  │   API     │  │   API     │             │
│            └───────────┘  └───────────┘  └───────────┘             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Provider Interface

All providers implement the `Provider` interface defined in `internal/provider/interface.go`:

```go
type Provider interface {
    Name() string
    SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error)
    SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}
```

Provider selection is validated before the first request by `internal/provider/registry`.
The registry checks provider name, base URL, model, and API key presence. Unknown
providers fail before any network request.

## Available Providers

### Anthropic

The default provider using Anthropic's Messages API.

| Field | Value |
|-------|-------|
| Name | `anthropic` |
| Default Model | `claude-sonnet-4-6-20251001` |
| API Base URL | `https://api.anthropic.com` |

```go
provider := anthropic.NewProvider(apiKey, baseURL, model)
```

### OpenAI

Provider using OpenAI's Chat Completions API.

| Field | Value |
|-------|-------|
| Name | `openai` |
| Default Model | `gpt-4o` |
| API Base URL | `https://api.openai.com` |

```go
provider := openai.NewProvider(apiKey, baseURL, model)
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_PROVIDER` | Provider name (`anthropic` or `openai`) | auto-detected from model |
| `ANTHROPIC_API_KEY` | API key used by the active provider | - |
| `ANTHROPIC_BASE_URL` | Provider base URL override | provider default |
| `ANTHROPIC_MODEL` | Model override | `claude-sonnet-4-6-20251001` |

### Configuration File

Create a `~/.go-code/settings.json` file:

```json
{
  "apiKey": "your-api-key",
  "provider": "anthropic",
  "model": "claude-sonnet-4-20250514",
  "baseUrl": "https://api.anthropic.com"
}
```

For OpenAI:

```json
{
  "apiKey": "sk-...",
  "provider": "openai",
  "model": "gpt-4o",
  "baseUrl": "https://api.openai.com"
}
```

For OpenAI-compatible providers, set `provider` to `openai` and use the vendor
base URL:

| Vendor | Provider | Base URL | Notes |
| --- | --- | --- | --- |
| Anthropic | `anthropic` | `https://api.anthropic.com` | Native Messages API path. |
| OpenAI | `openai` | `https://api.openai.com` | Chat Completions compatibility. |
| DeepSeek | `openai` | `https://api.deepseek.com` | Primary models: `deepseek-v4-pro`, `deepseek-v4-flash`. Legacy aliases (`deepseek-chat`, `deepseek-reasoner`) still work with deprecation warning. |
| Qwen | `openai` | `https://dashscope.aliyuncs.com/compatible-mode` | OpenAI-compatible mode; vendor-specific model availability applies. |
| GLM | `openai` | `https://open.bigmodel.cn/api/paas` | OpenAI-compatible mode; verify model names with GLM docs. |
| MiMo | `openai` | `https://api.mimo.com` | OpenAI-compatible mode; use `mimo-v2.5-pro` model. |
| Tencent Cloud Coding Plan | `anthropic` | `https://api.lkeap.cloud.tencent.com/coding/anthropic` | Anthropic-compatible path, commonly with `tc-code-latest`. |

Runtime `/model <name>` accepts known registry models and unknown models via
passthrough (provider inferred from model name prefix). Unknown models produce
a warning and proceed with the inferred provider.

### CLI Overrides

```go
overrides := &config.CLIOverrides{
    APIKey:   "your-api-key",
    Provider: "openai",
    Model:    "gpt-4o",
}
```

## Configuration Priority

Configuration is loaded in the following order (later sources override earlier ones):

1. Default values in `DefaultConfig()`
2. User config file (`~/.go-code/settings.json`)
3. Project config file (`./.go-code/settings.json`)
4. Environment variables
5. CLI overrides

## Adding Custom Providers

To add a new provider, implement the `Provider` interface:

```go
package myprovider

import (
    "context"
    "github.com/strings77wzq/claude-code-Go/internal/provider"
)

type MyProvider struct {
    apiKey string
    model  string
}

func NewProvider(apiKey, model string) *MyProvider {
    return &MyProvider{
        apiKey: apiKey,
        model:  model,
    }
}

func (p *MyProvider) Name() string {
    return "myprovider"
}

func (p *MyProvider) DefaultModel() string {
    return "my-model"
}

func (p *MyProvider) SendMessage(ctx context.Context, req *provider.Request) (*provider.Response, error) {
    // Implement non-streaming request
}

func (p *MyProvider) SendMessageStream(ctx context.Context, req *provider.Request, onTextDelta func(text string)) (*provider.Response, error) {
    // Implement streaming request
}
```

## Provider Components

| File | Purpose |
|------|---------|
| `internal/provider/interface.go` | Provider interface |
| `internal/provider/registry/registry.go` | Provider/model validation and routing |
| `internal/provider/errors.go` | Normalized provider error categories |
| `internal/provider/anthropic/provider.go` | Anthropic API adapter |
| `internal/provider/openai/provider.go` | OpenAI API adapter |

## Related Documentation

- [Agent Loop Implementation](./core-code-agent-loop) — Provider usage in agent loop
- [Configuration Guide](../api/config) — Full configuration reference
- [API Client](./core-code-entry) — Internal API details

---

<div class="nav-prev-next">

- [Tool System Overview](../tools/overview) ←
- → [Agent Loop Implementation](./core-code-agent-loop)

</div>
