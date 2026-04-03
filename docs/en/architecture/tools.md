---
title: Tool System
description: Deep dive into go-code's tool system — interface design, registry pattern, built-in tools, and MCP integration
---

# Tool System

go-code implements a flexible tool system that enables the agent to interact with the filesystem, execute commands, and integrate with external services via MCP.

## Tool Interface Design

All tools implement the `Tool` interface defined in `internal/tool/tool.go`:

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

### Result Type

```go
// Result represents the output of a tool execution.
type Result struct {
    Content string  // The output content
    IsError bool    // Whether this is an error result
}

// Helper constructors
func Success(content string) Result
func Error(msg string) Result
```

### Tool Definition

For API communication, tools also provide a definition:

```go
// ToolDefinition represents a tool's definition for API responses.
type ToolDefinition struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    InputSchema map[string]any `json:"input_schema"`
}
```

## Registry Pattern

The `Registry` in `internal/tool/registry.go` manages all available tools:

```go
type Registry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}
```

### Key Methods

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

### Thread Safety

The registry uses `sync.RWMutex` to allow concurrent reads while ensuring safe writes during registration.

### Registration

Tools are registered during startup in `internal/tool/init/register.go`:

```go
func RegisterBuiltinTools(r *Registry, workingDir string) error {
    // Register all built-in tools
    r.Register(builtin.NewReadTool())
    r.Register(builtin.NewWriteTool())
    r.Register(builtin.NewEditTool())
    r.Register(builtin.NewGlobTool())
    r.Register(builtin.NewGrepTool())
    r.Register(builtin.NewBashTool(workingDir))
    return nil
}
```

## Built-in Tools Overview

go-code provides six built-in tools that cover the essential operations for software development:

### 1. Read Tool

Reads file contents with optional offset and limit.

```go
type ReadTool struct{}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "file_path": {
      "type": "string",
      "description": "Path to the file to read"
    },
    "offset": {
      "type": "number",
      "description": "Line number to start reading from (0-based, default: 0)"
    },
    "limit": {
      "type": "number",
      "description": "Maximum number of lines to read (default: 2000)"
    }
  },
  "required": ["file_path"]
}
```

**Features:**
- Line-numbered output
- Maximum file size: 200KB
- No permission required

### 2. Write Tool

Creates or overwrites files.

```go
type WriteTool struct{}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "file_path": {
      "type": "string",
      "description": "Path to the file to write"
    },
    "content": {
      "type": "string",
      "description": "Content to write to the file"
    }
  },
  "required": ["file_path", "content"]
}
```

**Features:**
- Creates parent directories if needed
- Requires permission (can overwrite existing files)

### 3. Edit Tool

Makes targeted code edits using line-based replacement.

```go
type EditTool struct{}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "file_path": {
      "type": "string",
      "description": "Path to the file to edit"
    },
    "target_string": {
      "type": "string",
      "description": "The exact string to replace"
    },
    "replacement_string": {
      "type": "string",
      "description": "The replacement string"
    }
  },
  "required": ["file_path", "target_string", "replacement_string"]
}
```

**Features:**
- Exact string matching for safety
- Requires permission

### 4. Glob Tool

Finds files by pattern matching.

```go
type GlobTool struct{}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "pattern": {
      "type": "string",
      "description": "Glob pattern (e.g., *.go, **/*.ts)"
    }
  },
  "required": ["pattern"]
}
```

**Features:**
- Supports `*` and `**` patterns
- Recursive matching with `**`
- No permission required

### 5. Grep Tool

Searches file contents using regular expressions.

```go
type GrepTool struct{}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "pattern": {
      "type": "string",
      "description": "Regular expression pattern"
    },
    "path": {
      "type": "string",
      "description": "File path or directory to search"
    }
  },
  "required": ["pattern"]
}
```

**Features:**
- Full regex support
- Returns matching lines with context
- No permission required

### 6. Bash Tool

Executes shell commands.

```go
type BashTool struct {
    workingDir string
}
```

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "command": {
      "type": "string",
      "description": "The bash command to execute"
    },
    "timeout": {
      "type": "number",
      "description": "Timeout in seconds (default: 120)"
    }
  },
  "required": ["command"]
}
```

**Features:**
- Configurable timeout (default: 120s)
- Output truncation at 100KB
- Working directory set to project root
- Requires permission (can execute arbitrary commands)

## MCP Tool Adapter

The Model Context Protocol (MCP) integration allows go-code to use external tools and services.

### MCP Configuration

MCP servers are configured in `~/.config/go-code/mcp.json`:

```json
{
  "server-name": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "./workspace"],
    "env": {
      "HOME": "${HOME}"
    }
  }
}
```

### MCP Adapter

The adapter in `internal/tool/mcp/adapter.go` wraps MCP tools to match the go-code tool interface:

```go
type McpAdapter struct {
    client *McpClient
}

func (a *McpAdapter) Name() string { ... }
func (a *McpAdapter) Description() string { ... }
func (a *McpAdapter) InputSchema() map[string]any { ... }
func (a *McpAdapter) RequiresPermission() bool { ... }
func (a *McpAdapter) Execute(ctx context.Context, 
                              input map[string]any) tool.Result { ... }
```

### MCP Components

| File | Purpose |
|------|---------|
| `config.go` | Configuration loading and env variable interpolation |
| `client.go` | MCP protocol client implementation |
| `adapter.go` | Adapter to convert MCP tools to go-code interface |
| `transport.go` | Transport layer for MCP communication |
| `manager.go` | MCP server lifecycle management |

## Extending with Custom Tools

To add a custom tool:

1. **Create the tool implementation**:
```go
package mytool

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
    return false  // or true if needed
}

func (t *MyTool) Execute(ctx context.Context, 
                         input map[string]any) tool.Result {
    // Implementation
    return tool.Success("result")
}
```

2. **Register the tool** in `internal/tool/init/register.go`:
```go
r.Register(mytool.NewMyTool())
```

3. **Rebuild** — The tool will now be available to the agent.

## Permission Integration

Tools integrate with the permission system:

```go
func (a *Agent) checkPermission(toolName string, input map[string]any) bool {
    t := a.toolRegistry.GetTool(toolName)
    requiresPermission := t != nil && t.RequiresPermission()

    decision := a.permissionPolicy.Evaluate(toolName, input, requiresPermission)
    return decision == permission.Allow || decision == permission.Ask
}
```

If `RequiresPermission()` returns `true`, the permission system will prompt the user before execution.

## Related Documentation

- [Architecture Overview](overview.md) — System components
- [Agent Loop Deep Dive](agent-loop.md) — Tool execution in the loop
- [Configuration Guide](../guide/configuration.md) — MCP configuration