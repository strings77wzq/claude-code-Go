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
    SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error)
    SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}
```

首次请求前会通过 `internal/provider/registry` 校验 Provider 名称、base URL、模型和 API key。未知 Provider 会在发起网络请求前失败，并提示支持的 Provider。

## 可用 Provider

### Anthropic

使用 Anthropic Messages API 的默认 Provider。

| 字段 | 值 |
|-------|-------|
| 名称 | `anthropic` |
| 默认模型 | `claude-sonnet-4-6-20251001` |
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
| `LLM_PROVIDER` | Provider 名称（`anthropic` 或 `openai`） | 按模型自动判断 |
| `ANTHROPIC_API_KEY` | 当前 Provider 使用的 API key | - |
| `ANTHROPIC_BASE_URL` | Provider base URL 覆盖 | Provider 默认值 |
| `ANTHROPIC_MODEL` | 模型覆盖 | `claude-sonnet-4-6-20251001` |

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

OpenAI-compatible Provider 需要显式设置 `provider` 为 `openai`，并填写对应厂商的 base URL：

| 厂商 | Provider | Base URL | 说明 |
| --- | --- | --- | --- |
| Anthropic | `anthropic` | `https://api.anthropic.com` | 原生 Messages API 路径。 |
| OpenAI | `openai` | `https://api.openai.com` | Chat Completions 兼容路径。 |
| DeepSeek | `openai` | `https://api.deepseek.com` | 主力模型：`deepseek-v4-pro`、`deepseek-v4-flash`。旧名称（`deepseek-chat`、`deepseek-reasoner`）仍可使用但会显示弃用警告。 |
| Qwen | `openai` | `https://dashscope.aliyuncs.com/compatible-mode` | 使用 OpenAI 兼容模式，具体模型以厂商文档为准。 |
| GLM | `openai` | `https://open.bigmodel.cn/api/paas` | 使用 OpenAI 兼容模式，具体模型以厂商文档为准。 |
| MiMo | `openai` | `https://api.mimo.com` | 使用 OpenAI 兼容模式，推荐模型 `mimo-v2.5-pro`。 |
| Tencent Cloud Coding Plan | `anthropic` | `https://api.lkeap.cloud.tencent.com/coding/anthropic` | Anthropic 兼容路径，常配合 `tc-code-latest`。 |

运行时 `/model <name>` 接受 registry 已知模型以及未知模型的透传（根据模型名称前缀推断 Provider）。未知模型会显示警告并使用推断的 Provider 继续运行。

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
| `internal/provider/interface.go` | Provider 接口 |
| `internal/provider/registry/registry.go` | Provider/model 校验与路由 |
| `internal/provider/errors.go` | 统一 Provider 错误分类 |
| `internal/provider/anthropic/provider.go` | Anthropic API 适配器 |
| `internal/provider/openai/provider.go` | OpenAI API 适配器 |

## 相关文档

- [Agent 循环实现](../core-code/agent-loop-impl) — Provider 在 Agent 循环中的使用
- [配置指南](../guide/configuration) — 完整配置参考
- [入口点](../core-code/entry-point) — 内部 API 详情

---

<div class="nav-prev-next">

- [工具系统概述](../tools/overview) ←
- → [Agent 循环实现](../core-code/agent-loop-impl)

</div>
