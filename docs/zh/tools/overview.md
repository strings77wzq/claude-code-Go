---
title: 工具系统概述
description: 深入解析 go-code 的工具系统 — 接口定义、注册表模式、内置工具表和扩展指南
---

# 工具系统概述

工具系统是 go-code 的核心组件，使智能体能够与文件系统交互、执行 Shell 命令，并通过 MCP 与外部服务集成。本文档全面概述工具架构。

## 工具接口定义

所有工具必须实现 `internal/tool/tool.go` 中定义的 `Tool` 接口：

```go
// Tool represents an executable tool that can be called by the agent.
type Tool interface {
    // Name returns the unique name of the tool.
    Name() string

    // Description returns a human-readable description of what the tool does.
    Description() string

    // InputSchema returns the JSON schema for the tool's input parameters.
    InputSchema() map[string]any

    // RequiresPermission returns true if the tool requires special permissions.
    RequiresPermission() bool

    // Execute runs the tool with the given input and returns a result.
    Execute(ctx context.Context, input map[string]any) Result
}
```

### 接口方法说明

| 方法 | 用途 | 返回值 |
|------|------|--------|
| `Name()` | 工具的唯一标识符 | `string` |
| `Description()` | 供模型查看的描述 | `string` |
| `InputSchema()` | 输入验证的 JSON Schema | `map[string]any` |
| `RequiresPermission()` | 执行是否需要用户批准 | `bool` |
| `Execute()` | 使用给定输入运行工具 | `Result` |

### Result 类型

```go
// Result represents the output of a tool execution.
type Result struct {
    Content string  // 输出内容
    IsError bool    // 是否为错误结果
}

// 辅助构造函数
func Success(content string) Result
func Error(msg string) Result
```

### ToolDefinition

对于 API 通信，工具提供可序列化的定义：

```go
// ToolDefinition represents a tool's definition for API responses.
type ToolDefinition struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    InputSchema map[string]any `json:"input_schema"`
}
```

## 注册表模式

`internal/tool/registry.go` 中的 `Registry` 使用线程安全的 map 管理所有可用工具：

```go
type Registry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}
```

### 关键方法

```go
// Register adds a new tool to the registry
func (r *Registry) Register(tool Tool) error

// GetTool retrieves a tool by name
func (r *Registry) GetTool(name string) Tool

// Execute runs a tool by name with given input
func (r *Registry) Execute(ctx context.Context, name string, 
                           input map[string]any) Result

// GetAllDefinitions returns all tool definitions for API requests
func (r *Registry) GetAllDefinitions() []ToolDefinition
```

### 线程安全

注册表使用 `sync.RWMutex` 允许：
- 并发读取（多个读取器可以同时访问）
- 独占写入（一次只能有一个写入器）

这确保了智能体执行期间的线程安全访问。

## 内置工具表

go-code 提供六个内置工具，涵盖基本的软件开发操作：

| # | 工具名称 | 用途 | 需要权限 | 源文件 |
|---|----------|------|----------|--------|
| 1 | **Read** | 读取文件内容，支持可选的 offset/limit | 否 | `internal/tool/builtin/read.go` |
| 2 | **Write** | 创建或覆盖文件 | 是 | `internal/tool/builtin/write.go` |
| 3 | **Edit** | 使用精确字符串匹配进行针对性代码编辑 | 是 | `internal/tool/builtin/edit.go` |
| 4 | **Glob** | 按 glob 模式查找文件（`*`, `**`） | 否 | `internal/tool/builtin/glob.go` |
| 5 | **Grep** | 使用正则表达式搜索文件内容 | 否 | `internal/tool/builtin/grep.go` |
| 6 | **Bash** | 执行 Shell 命令 | 是 | `internal/tool/builtin/bash.go` |

### 工具详情

#### Read 工具
```go
type ReadTool struct{}
```
- 逐行读取文件内容
- 支持可选的 `offset` 和 `limit` 参数
- 最大文件大小：200KB
- 返回带行号的输出
- 不需要权限

#### Write 工具
```go
type WriteTool struct{}
```
- 创建新文件或覆盖现有文件
- 自动创建父目录
- 需要权限（破坏性操作）

#### Edit 工具
```go
type EditTool struct{}
```
- 替换精确的字符串序列
- 安全特性：需要精确匹配以防止意外编辑
- 需要权限

#### Glob 工具
```go
type GlobTool struct{}
```
- 文件发现的模式匹配
- 支持 `*`（单层）和 `**`（递归）
- 不需要权限

#### Grep 工具
```go
type GrepTool struct{}
```
- 完整的正则表达式搜索
- 返回带上下文的匹配行
- 不需要权限

#### Bash 工具
```go
type BashTool struct {
    workingDir string
}
```
- 执行 Shell 命令
- 可配置超时（默认：120s）
- 输出在 100KB 处截断
- 工作目录设置为项目根目录
- 需要权限

## 如何扩展自定义工具

添加自定义工具：

### 步骤 1：创建工具实现

在 `internal/tool/builtin/` 或新包中创建新文件：

```go
package mytool

import (
    "context"
    "github.com/user/go-code/internal/tool"
)

type MyTool struct{}

func NewMyTool() tool.Tool {
    return &MyTool{}
}

func (t *MyTool) Name() string {
    return "MyTool"
}

func (t *MyTool) Description() string {
    return "Description of what my tool does"
}

func (t *MyTool) InputSchema() map[string]any {
    return map[string]any{
        "type": "object",
        "properties": map[string]any{
            "param": map[string]any{
                "type":        "string",
                "description": "Parameter description",
            },
        },
        "required": []string{"param"},
    }
}

func (t *MyTool) RequiresPermission() bool {
    return false  // 如果需要用户批准则设置为 true
}

func (t *MyTool) Execute(ctx context.Context, 
                         input map[string]any) tool.Result {
    // 在此实现
    param := input["param"].(string)
    // ... 做某事 ...
    return tool.Success("result")
}
```

### 步骤 2：注册工具

在 `internal/tool/init/register.go` 中：

```go
func RegisterBuiltinTools(r *Registry, workingDir string) error {
    // 现有工具
    r.Register(builtin.NewReadTool())
    r.Register(builtin.NewWriteTool())
    r.Register(builtin.NewEditTool())
    r.Register(builtin.NewGlobTool())
    r.Register(builtin.NewGrepTool())
    r.Register(builtin.NewBashTool(workingDir))

    // 添加你的自定义工具
    r.Register(mytool.NewMyTool())
    
    return nil
}
```

### 步骤 3：构建和测试

```bash
go build ./cmd/go-code
./go-code "Your test prompt"
```

## MCP 工具适配器概述

go-code 与 Model Context Protocol (MCP) 集成，以使用外部工具和服务。MCP 适配器将外部工具包装以匹配 go-code 工具接口。

### MCP 适配器

```go
// McpToolAdapter wraps an MCP tool and implements the tool.Tool interface.
type McpToolAdapter struct {
    serverName  string
    toolName    string
    description string
    inputSchema map[string]any
    client      *McpClient
}
```

适配器：
- 将 MCP 工具名称转换为格式 `mcp__{serverName}__{toolName}`
- 将执行委托给 MCP 客户端
- 始终需要权限（外部工具）

### MCP 组件

| 文件 | 用途 |
|------|------|
| `internal/tool/mcp/config.go` | 配置加载和环境变量插值 |
| `internal/tool/mcp/client.go` | MCP 协议客户端实现 |
| `internal/tool/mcp/adapter.go` | 将 MCP 工具转换为 go-code 接口的适配器 |
| `internal/tool/mcp/transport.go` | 基于 stdio 的 MCP 通信传输层 |
| `internal/tool/mcp/manager.go` | MCP 服务器生命周期管理 |

详见 [MCP 集成](./mcp.md)。

## 相关文档

- [智能体循环实现](../core-code/agent-loop-impl.md) — 循环中的工具执行
- [MCP 集成](./mcp.md) — Model Context Protocol 支持
- [入口点详解](../core-code/entry-point.md) — 启动时的工具注册

---

<div class="nav-prev-next">

- [智能体循环实现](../core-code/agent-loop-impl.md) ←
- → [MCP 集成](./mcp.md)

</div>