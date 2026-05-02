---
title: LSP 集成
description: LSP 健康检查、能力门控、诊断和降级行为
---

# LSP 集成

> **状态：实验性功能（v0.3 产品化中）** — LSP 客户端和健康检查门控已经实现，但只有配置的服务器通过健康检查后，才会暴露 LSP 代码智能能力。`go-code doctor --offline` 只报告本地配置状态，不发起网络探测。

go-code 将 Language Server Protocol (LSP) 作为可选扩展能力，用于代码智能。未配置 LSP 服务器时，核心提示词流程、TUI、内置工具、MCP 工具、Hooks 和 Skills 仍会继续工作。

## 配置

LSP 健康检查通过 `GO_CODE_LSP_URL` 显式启用：

```bash
export GO_CODE_LSP_URL="http://127.0.0.1:8080/lsp"
```

未设置该变量时，LSP 会被视为 unavailable，而不是 failed。

## 健康检查

LSP gate 会先发送 `initialize` 请求，再决定是否暴露代码智能操作。只有初始化成功的服务器才会被视为健康。

```bash
go-code doctor
```

离线诊断不会连接服务器：

```bash
go-code doctor --offline
```

常见 doctor 状态：

| 状态 | 含义 |
| --- | --- |
| `[SKIP] lsp: not configured` | 未设置 `GO_CODE_LSP_URL`。 |
| `[SKIP] lsp: configured ... skipped by --offline` | 已配置 LSP URL，但 `--offline` 跳过了网络健康检查。 |
| `[PASS] lsp` | 初始化成功。 |
| `[FAIL] lsp` | 已配置服务器初始化失败或返回无效响应。 |

## 能力门控

go-code 只会在健康检查通过后暴露 LSP 操作。当前 gate 识别：

| 操作 | LSP capability |
| --- | --- |
| Diagnostics | `publishDiagnosticsProvider` |
| Symbols | `workspaceSymbolProvider` 或 `documentSymbolProvider` |
| Definitions | `definitionProvider` |
| References | `referencesProvider` |
| Hover | `hoverProvider` |

如果服务器未配置、初始化失败，或没有声明某项能力，对应操作不会被暴露。

## Trace 证据

LSP 健康检查结果会作为非致命 extension 事件写入 session trace：

```json
{"type":"extension","name":"lsp","event":"health_check","status":"unavailable"}
```

成功检查会包含可暴露操作，以及服务器返回的身份信息。失败检查会记录初始化错误，但不会阻断无关工作流。

## 降级行为

- 未配置 LSP：继续运行，但不暴露 LSP 操作。
- 离线 doctor：只报告配置状态，不访问网络。
- 初始化失败：核心 prompt/tool 工作流保持可用，并通过 diagnostics 和 trace evidence 报告 LSP 失败。
- 缺少服务器能力：不暴露对应操作。

## 验证

当前覆盖位于 `internal/lsp/gate_test.go`，证明 unavailable、初始化成功、初始化失败、trace 记录、以及 operation gating 行为。

## 相关

- [MCP 集成](./mcp)
- [Hooks 系统](./hooks)
- [Skills 系统](./skills)
