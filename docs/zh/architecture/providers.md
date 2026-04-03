---
title: Provider 系统
description: LLM 后端多 Provider 支持 — 接口设计、Anthropic 和 OpenAI 适配器、配置
---

# Provider 系统

go-code 通过统一的 Provider 接口支持多个 LLM Provider。本文档描述了 Provider 架构以及如何配置不同的后端。

## 架构概述

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Provider 架构                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌──────────────┐                       │
│   │  Agent 循环  │─────────▶│ Provider      │                       │
│   │             │          │ 接口          │                       │
│   └─────────────┘          └──────┬───────┘                       │
│                                   │                                │
│                    ┌──────────────┼──────────────┐                 │
│                    ▼              ▼              ▼                 │
│            ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│            │  Anthropic │  │  OpenAI   │  │ 自定义     │             │
│            │ Provider   │  │ Provider  │  │ Provider  │             │
│            └───────────┘  └───────────┘  └───────────┘             │
│                    │              │              │                 │
│                    ▼              ▼              ▼                 │
│            ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│            │  Anthropic │  │  OpenAI    │  │ 自定义     │             │
│            │   API      │  │   API      │  │   API      │             │
│            └───────────┘  └───────────┘  └───────────┘             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Provider 接口

所有 Provider 都实现了 `internal/provider/interface.go` 中定义的 `Provider` 接口：

```go
type Provider interface {
    Name() string
    DefaultModel() string
    SendMessage(ctx context.Context, req *Request) (*Response, error)
    SendMessageStream(ctx context.Context, req *Request, onTextDelta func(text string)) (*Response, error)
}
```

### 核心类型

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

## 可用 Provider

### Anthropic

使用 Anthropic Messages API 的默认 Provider。

| 字段 | 值 |
|-------|-------|
| 名称 | `anthropic` |
| 默认模型 | `claude-sonnet-4-20250514` |
| API 基础 URL | `https://api.anthropic.com` |

```go
provider := anthropic.NewProvider(apiKey, baseURL, model)
```

### OpenAI

使用 OpenAI Chat Completions API 的 Provider。

| 字段 | 值 |
|-------|-------|
| 名称 | `openai` |
| 默认模型 | `gpt-4o` |
| API 基础 URL | `https://api.openai.com` |

```go
provider := openai.NewProvider(apiKey, baseURL, model)
```

## 配置

### 环境变量

| 变量 | 描述 | 默认值 |
|----------|-------------|---------|
| `LLM_PROVIDER` | Provider 名称 | `anthropic` |
| `ANTHROPIC_API_KEY` | Anthropic API 密钥 | - |
| `OPENAI_API_KEY` | OpenAI API 密钥 | - |

### 配置文件

创建 `~/.go-code/settings.json` 文件：

```json
{
  "apiKey": "your-api-key",
  "provider": "anthropic",
  "model": "claude-sonnet-4-20250514",
  "baseUrl": "https://api.anthropic.com"
}
```

OpenAI 配置：

```json
{
  "apiKey": "sk-...",
  "provider": "openai",
  "model": "gpt-4o",
  "baseUrl": "https://api.openai.com"
}
```

### CLI 覆盖

```go
overrides := &config.CLIOverrides{
    APIKey:   "your-api-key",
    Provider: "openai",
    Model:    "gpt-4o",
}
```

## 配置优先级

配置按以下顺序加载（后面的来源会覆盖前面的）：

1. `DefaultConfig()` 中的默认值
2. 用户配置文件 (`~/.go-code/settings.json`)
3. 项目配置文件 (`./.go-code/settings.json`)
4. 环境变量
5. CLI 覆盖

## 添加自定义 Provider

要添加新的 Provider，实现 `Provider` 接口：

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
    // 实现非流式请求
}

func (p *MyProvider) SendMessageStream(ctx context.Context, req *provider.Request, onTextDelta func(text string)) (*provider.Response, error) {
    // 实现流式请求
}
```

## Provider 组件

| 文件 | 用途 |
|------|---------|
| `internal/provider/interface.go` | Provider 接口和核心类型 |
| `internal/provider/anthropic/provider.go` | Anthropic API 适配器 |
| `internal/provider/openai/provider.go` | OpenAI API 适配器 |

## 相关文档

- [Agent 循环实现](../core-code/agent-loop-impl.md) — Provider 在 Agent 循环中的使用
- [配置指南](../guide/configuration.md) — 完整配置参考
- [入口点](../core-code/entry-point.md) — 内部 API 详情

---

<div class="nav-prev-next">

- [工具系统概述](../tools/overview.md) ←
- → [Agent 循环实现](../core-code/agent-loop-impl.md)

</div>