---
title: MCP Integration
description: Deep dive into Model Context Protocol — transport layer, JSON-RPC format, tool discovery, adapter implementation, and configuration
---

# MCP Integration

go-code supports the Model Context Protocol (MCP) for integrating external tools and services. This document provides a comprehensive overview of MCP implementation in go-code.

## What is MCP?

The **Model Context Protocol (MCP)** is an open standard that enables AI models to interact with external tools and services through a standardized interface. It provides:

- **Tool Discovery** — AI models can discover available tools dynamically
- **Standardized Communication** — JSON-RPC based protocol
- **Transport Flexibility** — Supports stdio, HTTP, and other transports
- **Extensibility** — Easy to add new tools and services

### MCP Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                       MCP Architecture                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   go-code   │◄────────►│   MCP Server    │                     │
│   │   (Client)  │  stdio   │  (Subprocess)   │                     │
│   └─────────────┘          └─────────────────┘                     │
│         │                            │                             │
│         │                            │                             │
│         ▼                            ▼                             │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │ Tool        │          │ Tools/Services  │                     │
│   │ Registry    │          │ (e.g., Files,   │                     │
│   │             │          │  DB, Git, etc.) │                    │
│   └─────────────┘          └─────────────────┘                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## MCP in the Extension Architecture

MCP is one of three extension mechanisms in go-code, alongside Skills and Hooks. Each serves a different purpose:

| Extension | Purpose | How it Works |
|-----------|---------|--------------|
| **Skills** | Customize agent behavior | Named prompts injected into system prompt |
| **Hooks** | Monitor tool execution | Pre/post callbacks for logging, auditing |
| **MCP** | Extend tool capabilities | External servers providing additional tools |

MCP integrates with the tool registry at runtime, adding tools from external MCP servers alongside built-in tools (Read, Write, Edit, Glob, Grep, Bash).

## Stdio Transport Layer

go-code uses stdio (standard input/output) as the transport layer for MCP communication:

### Transport Implementation

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

### Starting the Transport

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

The transport spawns a subprocess and creates pipes for bidirectional communication.

## JSON-RPC Format

MCP uses JSON-RPC 2.0 for request/response formatting:

### Request Format

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

### Response Format

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

### Sending Requests

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

### Reading Responses

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

## Tool Discovery Process

When an MCP server starts, go-code discovers its available tools:

### Discovery Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Tool Discovery Flow                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. MCP Server starts (via stdio transport)                        │
│         │                                                           │
│         ▼                                                           │
│  2. Send "initialize" request with protocol version                 │
│         │                                                           │
│         ▼                                                           │
│  3. Receive server capabilities                                    │
│         │                                                           │
│         ▼                                                           │
│  4. Send "tools/list" request                                      │
│         │                                                           │
│         ▼                                                           │
│  5. Receive tool definitions                                       │
│         │                                                           │
│         ▼                                                           │
│  6. Create McpToolAdapter for each tool                            │
│         │                                                           │
│         ▼                                                           │
│  7. Register adapters in tool registry                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Initialization Request

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

### Tool List Response

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

## MCP Tool Adapter Implementation

The adapter wraps MCP tools to match the go-code `Tool` interface:

### Adapter Structure

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

### Interface Implementation

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

### Tool Name Format

MCP tools are prefixed with `mcp__` to avoid conflicts with built-in tools:
- Format: `mcp__{serverName}__{toolName}`
- Example: `mcp__filesystem__read_file`

### Permission Handling

MCP tools always require permission because they are external tools that could perform arbitrary operations:

```go
func (a *McpToolAdapter) RequiresPermission() bool {
    return true
}
```

## Configuration

MCP servers are configured in a JSON configuration file.

### Configuration File Location

Default: `~/.go-code/mcp.json`

Or custom path via configuration.

### Configuration Format

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

### Environment Variable Interpolation

The configuration supports `${VAR}` syntax for environment variables:

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

### Loading Configuration

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

## MCP Components Summary

| File | Purpose |
|------|---------|
| `internal/tool/mcp/config.go` | Configuration loading and env variable interpolation |
| `internal/tool/mcp/transport.go` | Stdio transport layer implementation |
| `internal/tool/mcp/client.go` | MCP protocol client (JSON-RPC communication) |
| `internal/tool/mcp/adapter.go` | Adapter converting MCP tools to go-code interface |
| `internal/tool/mcp/manager.go` | MCP server lifecycle management |

## Related Documentation

- [Skills System](./skills.md) — Named prompts for customizing agent behavior
- [Hooks System](./hooks.md) — Pre/post execution callbacks
- [Tool System Overview](../tools/overview.md) — Tool interface and registry
- [Configuration Guide](../guide/configuration.md) — MCP configuration options

---

<div class="nav-prev-next">

- [Skills System](./skills.md) ←
- → [Hooks System](./hooks.md)

</div>