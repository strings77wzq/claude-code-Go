---
title: 入口点详解
description: 深入解析 main.go —— 组件初始化顺序、信号处理和系统提示词
---

# 入口点详解

go-code 的主入口点位于 `cmd/go-code/main.go`。本文档详细说明了应用程序如何按正确顺序初始化所有组件。

## 概述

运行 `go-code` 时，会执行以下顺序：

```
┌─────────────────────────────────────────────────────────────────────┐
│                        main.go 执行流程                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. 日志设置            → 初始化结构化日志                            │
│  2. 信号处理            → 注册 SIGINT/SIGTERM 处理器                 │
│  3. 加载配置            → 从环境变量/配置文件加载                     │
│  4. 创建 API 客户端     → 初始化 Anthropic API 客户端               │
│  5. 创建工具注册表      → 创建空注册表                               │
│  6. 注册内置工具        → 注册 6 个内置工具                           │
│  7. 创建权限策略        → 设置三级权限策略                           │
│  8. 创建智能体          → 初始化所有依赖项                            │
│  9. 启动 REPL           → 开始交互式会话（阻塞）                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 源代码分析

### 导入

```go
import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/go-code/internal/agent"
	"github.com/user/go-code/internal/api"
	"github.com/user/go-code/internal/config"
	"github.com/user/go-code/internal/permission"
	"github.com/user/go-code/internal/tool"
	toolinit "github.com/user/go-code/internal/tool/init"
	"github.com/user/go-code/pkg/tty"
)
```

关键依赖：
- `log/slog` — 结构化日志（Go 1.21+）
- `os/signal` — 优雅关闭处理
- `internal/*` — 核心应用组件
- `pkg/tty` — REPL 实现

### 版本常量

```go
const version = "0.1.0"
```

### 系统提示词

```go
const systemPrompt = "You are an interactive agent that helps users with software engineering tasks. You have access to tools for reading files, editing files, executing shell commands, searching code, and more. Use your tools to complete tasks efficiently and accurately."
```

系统提示词定义了智能体的角色和能力。每次 API 请求都会发送此提示词来指导模型的行为。

## 初始化顺序

### 步骤 1：日志设置

```go
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)
```

应用程序使用 Go 内置的 `log/slog` 进行结构化日志记录。所有日志消息都包含时间戳和严重级别。

### 步骤 2：信号处理

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    sig := <-sigChan
    logger.Info("Received signal, shutting down", "signal", sig.String())
    logger.Info("Shutdown complete")
    os.Exit(0)
}()
```

应用程序为以下信号注册处理器：
- `SIGINT` — 中断（Ctrl+C）
- `SIGTERM` — 终止请求

收到信号时，应用程序记录事件并优雅退出。

### 步骤 3：加载配置

```go
logger.Info("Loading configuration")
cfg, err := config.Load(nil)
if err != nil {
    logger.Error("Failed to load configuration", "error", err)
    os.Exit(1)
}
logger.Info("Configuration loaded", "model", cfg.Model, "baseURL", cfg.BaseURL)
```

配置从多个来源加载（环境变量、配置文件）。详情请参阅[配置指南](../guide/configuration.md)。

### 步骤 4：创建 API 客户端

```go
logger.Info("Creating API client")
client := api.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
logger.Info("API client created")
```

API 客户端使用凭证和模型设置初始化。它负责与 Anthropic API 的 HTTP 通信。

### 步骤 5：创建工具注册表

```go
logger.Info("Creating tool registry")
registry := tool.NewRegistry()
logger.Info("Tool registry created")
```

注册表是所有可用工具的线程安全容器。它使用 `sync.RWMutex` 实现并发读取和安全写入。

### 步骤 6：注册内置工具

```go
logger.Info("Registering builtin tools")
wd := cfg.WorkingDir
if wd == "" {
    wd, _ = os.Getwd()
}
if err := toolinit.RegisterBuiltinTools(registry, wd); err != nil {
    logger.Error("Failed to register builtin tools", "error", err)
    os.Exit(1)
}
logger.Info("Builtin tools registered", "count", len(registry.GetAllDefinitions()))
```

注册六个内置工具：

| 工具 | 用途 |
|------|------|
| **Read** | 读取文件内容 |
| **Write** | 创建/覆盖文件 |
| **Edit** | 修改特定代码段 |
| **Glob** | 按模式查找文件 |
| **Grep** | 搜索文件内容 |
| **Bash** | 执行 Shell 命令 |

工作目录传递给 Bash 工具，以便命令在正确的上下文中执行。

### 步骤 7：创建权限策略

```go
logger.Info("Creating permission policy")
policy := permission.NewPolicy(permission.WorkspaceWrite)
logger.Info("Permission policy created")
```

权限策略确定哪些操作需要用户批准。详见[架构概述](overview.md)。

### 步骤 8：创建智能体

```go
logger.Info("Creating agent")
agentInstance := agent.NewAgent(client, registry, policy, systemPrompt, cfg.Model)
logger.Info("Agent started", "model", cfg.Model)
```

智能体使用所有依赖项创建：
- 用于模型通信的 API 客户端
- 用于执行的工具注册表
- 用于访问控制的权限策略
- 用于行为指导的系统提示词
- 用于 API 请求的模型标识符

### 步骤 9：启动 REPL

```go
logger.Info("Starting REPL")
repl := tty.NewREPL(agentInstance, version, cfg.Model)

// Run REPL - this blocks until exit
repl.Run()

logger.Info("REPL exited")
```

REPL（Read-Eval-Print Loop）是交互式界面。它：
1. 从终端读取用户输入
2. 将其传递给智能体处理
3. 显示响应
4. 重复直到用户退出（Ctrl+C 或 `exit` 命令）

## 关键初始化概念

### 依赖注入

应用程序始终使用依赖注入。每个组件通过构造函数接收其依赖项：

```go
agent.NewAgent(client, registry, policy, systemPrompt, model)
```

这使代码可测试且模块化。

### 顺序很重要

初始化顺序至关重要：

1. 日志必须首先就绪（其他组件在初始化期间记录日志）
2. 组件需要配置之前必须先加载配置
3. 工具注册表必须存在，然后才能注册工具
4. 创建智能体之前所有组件必须准备就绪

### 错误处理

每个初始化步骤检查错误并退出，并显示描述性消息：

```go
if err != nil {
    logger.Error("Failed to load configuration", "error", err)
    os.Exit(1)
}
```

## 相关文档

- [智能体循环实现](agent-loop-impl.md) — Run() 方法和执行流程
- [工具系统概述](../tools/overview.md) — 工具接口和注册表
- [MCP 集成](../architecture/mcp.md) — Model Context Protocol 支持
- [配置指南](../guide/configuration.md) — 配置选项

---

<div class="nav-prev-next">

- [架构概述](../architecture/overview.md) ←
- → [智能体循环实现](agent-loop-impl.md)

</div>