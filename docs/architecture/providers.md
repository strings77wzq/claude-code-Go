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
    DefaultModel() string
    SendMessage(ctx context.Context, req *Request) (*Response, error)
    SendMessageStream(ctx context.Context, req *Request, onTextDelta func(text string)) (*Response, error)
}
```

### Core Types

```go
type Request struct {
    Model     string
    MaxTokens int
    System    string
    Stream    bool
    Messages  []Message
    Tools     []ToolDefinition
}

type Message struct {
    Role    string
    Content string
}

type Response struct {
    ID         string
    Content    string
    StopReason string
    Usage       Usage
}

type Usage struct {
    InputTokens  int
    OutputTokens int
}
```

## Available Providers

### Anthropic

The default provider using Anthropic's Messages API.

| Field | Value |
|-------|-------|
| Name | `anthropic` |
| Default Model | `claude-sonnet-4-20250514` |
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
| `LLM_PROVIDER` | Provider name | `anthropic` |
| `ANTHROPIC_API_KEY` | Anthropic API key | - |
| `OPENAI_API_KEY` | OpenAI API key | - |

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
    "github.com/user/go-code/internal/provider"
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
| `internal/provider/interface.go` | Provider interface and core types |
| `internal/provider/anthropic/provider.go` | Anthropic API adapter |
| `internal/provider/openai/provider.go` | OpenAI API adapter |

## Related Documentation

- [Agent Loop Implementation](../core-code/agent-loop-impl.md) — Provider usage in agent loop
- [Configuration Guide](../guide/configuration.md) — Full configuration reference
- [API Client](../core-code/entry-point.md) — Internal API details

---

<div class="nav-prev-next">

- [Tool System Overview](../tools/overview.md) ←
- → [Agent Loop Implementation](../core-code/agent-loop-impl.md)

</div>