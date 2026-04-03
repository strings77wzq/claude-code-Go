---
title: MCP 集成
description: 深入解析模型上下文协议 — 传输层、JSON-RPC 格式、工具发现、适配器实现和配置
---

# MCP 集成

go-code 支持模型上下文协议 (MCP) 以集成外部工具和服务。本文提供 go-code 中 MCP 实现的全面概述。

## 什么是 MCP？

**模型上下文协议 (MCP)** 是一个开放标准，使 AI 模型能够通过标准化接口与外部工具和服务交互。它提供：

- **工具发现** — AI 模型可以动态发现可用工具
- **标准化通信** — 基于 JSON-RPC 的协议
- **传输灵活性** — 支持 stdio、HTTP 和其他传输方式
- **可扩展性** — 易于添加新工具和服务

### MCP 架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                       MCP 架构                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   go-code   │◄────────►│   MCP 服务器    │                     │
│   │   (客户端)  │  stdio   │  (子进程)       │                     │
│   └─────────────┘          └─────────────────┘                     │
│         │                            │                             │
│         │                            │                             │
│         ▼                            ▼                             │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │ 工具        │          │ 工具/服务       │                     │
│   │ 注册表      │          │ (例如：文件、   │                     │
│   │             │          │  数据库、Git)  │                    │
│   └─────────────┘          └─────────────────┘                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 扩展架构中的 MCP

MCP 是 go-code 中三种扩展机制之一，另外还有技能和钩子。各自服务不同目的：

| 扩展 | 用途 | 工作方式 |
|-----------|---------|----------|
| **技能** | 自定义 agent 行为 | 命名提示词注入系统提示词 |
| **钩子** | 监控工具执行 | 日志、审计的预/后回调 |
| **MCP** | 扩展工具能力 | 提供额外工具的外部服务器 |

MCP 在运行时与工具注册表集成，将外部 MCP 服务器的工具与内置工具（Read、Write、Edit、Glob、Grep、Bash）一起添加。

## Stdio 传输层

go-code 使用 stdio（标准输入/输出）作为 MCP 通信的传输层：

### 传输实现

```go
// StdioTransport implements stdio-based transport for MCP servers.
type StdioTransport struct {
    cmd     *exec.Cmd
    stdin   *os.File
    stdout  *bufio.Reader
    process *os.Process
    mu      sync.Mutex
    closed  bool
}
```

### 启动传输

```go
func NewStdioTransport(command string, args []string, env map[string]string) *StdioTransport {
    cmd := exec.Command(command, args...)
    if env != nil {
        for k, v := range env {
            cmd.Env = append(cmd.Env, k+"="+v)
        }
    }
    return &StdioTransport{cmd: cmd}
}

func (t *StdioTransport) Start() error {
    // Create pipes for stdin and stdout
    stdinPipe, err := t.cmd.StdinPipe()
    // ...
    stdoutPipe, err := t.cmd.StdoutPipe()
    // ...

    // Start the process
    if err := t.cmd.Start(); err != nil {
        return fmt.Errorf("failed to start process: %w", err)
    }

    t.process = t.cmd.Process
    return nil
}
```

传输层生成子进程并创建双向通信管道。

## JSON-RPC 格式

MCP 使用 JSON-RPC 2.0 进行请求/响应格式化：

### 请求格式

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

### 响应格式

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "tool-name",
        "description": "Tool description",
        "inputSchema": {...}
      }
    ]
  }
}
```

### 发送请求

```go
func (t *StdioTransport) SendRequest(method string, params map[string]any, id int) error {
    req := map[string]any{
        "jsonrpc": "2.0",
        "id":      id,
        "method":  method,
        "params":  params,
    }

    data, err := json.Marshal(req)
    // ...

    data = append(data, '\n')
    if _, err := t.stdin.Write(data); err != nil {
        return fmt.Errorf("failed to write request: %w", err)
    }

    if err := t.stdin.Sync(); err != nil {
        return fmt.Errorf("failed to sync stdin: %w", err)
    }

    return nil
}
```

### 读取响应

```go
func (t *StdioTransport) ReadResponse() (map[string]any, error) {
    line, err := t.stdout.ReadBytes('\n')
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Remove trailing newline
    if len(line) > 0 && line[len(line)-1] == '\n' {
        line = line[:len(line)-1]
    }

    var resp map[string]any
    if err := json.Unmarshal(line, &resp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return resp, nil
}
```

## 工具发现过程

当 MCP 服务器启动时，go-code 发现其可用工具：

### 发现流程

```
┌─────────────────────────────────────────────────────────────────────┐
│                    工具发现流程                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. MCP 服务器启动（通过 stdio 传输）                               │
│         │                                                           │
│         ▼                                                           │
│  2. 发送 "initialize" 请求及协议版本                               │
│         │                                                           │
│         ▼                                                           │
│  3. 接收服务器能力                                                  │
│         │                                                           │
│         ▼                                                           │
│  4. 发送 "tools/list" 请求                                         │
│         │                                                           │
│         ▼                                                           │
│  5. 接收工具定义                                                    │
│         │                                                           │
│         ▼                                                           │
│  6. 为每个工具创建 McpToolAdapter                                  │
│         │                                                           │
│         ▼                                                           │
│  7. 在工具注册表中注册适配器                                        │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 初始化请求

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "go-code",
      "version": "0.1.0"
    }
  }
}
```

### 工具列表响应

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "read_file",
        "description": "Read contents of a file",
        "inputSchema": {
          "type": "object",
          "properties": {
            "path": {"type": "string"}
          },
          "required": ["path"]
        }
      }
    ]
  }
}
```

## MCP 工具适配器实现

适配器包装 MCP 工具以匹配 go-code 的 `Tool` 接口：

### 适配器结构

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

### 接口实现

```go
// Name returns the unique name of the tool in format mcp__{serverName}__{toolName}.
func (a *McpToolAdapter) Name() string {
    return "mcp__" + a.serverName + "__" + a.toolName
}

// Description returns the human-readable description of the tool.
func (a *McpToolAdapter) Description() string {
    return a.description
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (a *McpToolAdapter) InputSchema() map[string]any {
    return a.inputSchema
}

// RequiresPermission always returns true since external MCP tools require approval.
func (a *McpToolAdapter) RequiresPermission() bool {
    return true
}

// Execute delegates to McpClient.CallTool and returns a tool.Result.
func (a *McpToolAdapter) Execute(ctx context.Context, input map[string]any) tool.Result {
    result, err := a.client.CallTool(a.toolName, input)
    if err != nil {
        return tool.Error(err.Error())
    }
    return tool.Success(result)
}
```

### 工具名称格式

MCP 工具使用 `mcp__` 前缀以避免与内置工具冲突：
- 格式：`mcp__{serverName}__{toolName}`
- 示例：`mcp__filesystem__read_file`

### 权限处理

MCP 工具始终需要权限，因为它们是可能执行任意操作的外部工具：

```go
func (a *McpToolAdapter) RequiresPermission() bool {
    return true
}
```

## 配置

MCP 服务器在 JSON 配置文件中配置。

### 配置文件位置

默认：`~/.config/go-code/mcp.json`

或通过配置自定义路径。

### 配置格式

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "./workspace"],
    "env": {
      "HOME": "${HOME}"
    }
  },
  "github": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-github"],
    "env": {
      "GITHUB_TOKEN": "${GITHUB_TOKEN}"
    }
  }
}
```

### 环境变量插值

配置支持 `${VAR}` 语法用于环境变量：

```go
var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

func InterpolateEnvVars(env map[string]string) map[string]string {
    result := make(map[string]string)
    for k, v := range env {
        result[k] = envVarPattern.ReplaceAllStringFunc(v, func(match string) string {
            varName := match[2 : len(match)-1]
            val := os.Getenv(varName)
            if val != "" {
                return val
            }
            return match
        })
    }
    return result
}
```

### 加载配置

```go
func LoadMcpConfigs(settingsPath string) (map[string]McpServerConfig, error) {
    data, err := os.ReadFile(settingsPath)
    // ...

    var configs map[string]McpServerConfig
    if err := json.Unmarshal(data, &configs); err != nil {
        return nil, fmt.Errorf("failed to parse settings file: %w", err)
    }

    for name, config := range configs {
        if config.Env != nil {
            interpolated := InterpolateEnvVars(config.Env)
            config.Env = interpolated
            configs[name] = config
        }
    }

    return configs, nil
}
```

## MCP 组件概览

| 文件 | 用途 |
|------|---------|
| `internal/tool/mcp/config.go` | 配置加载和环境变量插值 |
| `internal/tool/mcp/transport.go` | Stdio 传输层实现 |
| `internal/tool/mcp/client.go` | MCP 协议客户端（JSON-RPC 通信） |
| `internal/tool/mcp/adapter.go` | 将 MCP 工具转换为 go-code 接口的适配器 |
| `internal/tool/mcp/manager.go` | MCP 服务器生命周期管理 |

## 相关文档

- [技能系统](./skills.md) — 自定义 agent 行为的命名提示词
- [钩子系统](./hooks.md) — 预/后执行回调
- [工具系统概述](../tools/overview.md) — 工具接口和注册表
- [配置指南](../guide/configuration.md) — MCP 配置选项

---

<div class="nav-prev-next">

- [技能系统](./skills.md) ←
- → [钩子系统](./hooks.md)

</div>