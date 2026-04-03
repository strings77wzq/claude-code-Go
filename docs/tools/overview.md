---
title: Tool System Overview
description: Deep dive into go-code's tool system — interface definition, registry pattern, built-in tools table, and extension guide
---

# Tool System Overview

The tool system is a core component of go-code that enables the agent to interact with the filesystem, execute shell commands, and integrate with external services via MCP. This document provides a comprehensive overview of the tool architecture.

## Tool Interface Definition

All tools must implement the `Tool` interface defined in `internal/tool/tool.go`:

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

### Interface Methods Explained

| Method | Purpose | Returns |
|--------|---------|---------|
| `Name()` | Unique identifier for the tool | `string` |
| `Description()` | Human-readable description shown to the model | `string` |
| `InputSchema()` | JSON Schema for input validation | `map[string]any` |
| `RequiresPermission()` | Whether execution requires user approval | `bool` |
| `Execute()` | Run the tool with provided input | `Result` |

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

### ToolDefinition

For API communication, tools provide a serializable definition:

```go
// ToolDefinition represents a tool's definition for API responses.
type ToolDefinition struct {
    Name        string         `json:"name"`
    Description string         `json:"description"`
    InputSchema map[string]any `json:"input_schema"`
}
```

## Registry Pattern

The `Registry` in `internal/tool/registry.go` manages all available tools using a thread-safe map:

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

The registry uses `sync.RWMutex` to allow:
- Concurrent reads (many readers can access simultaneously)
- Exclusive writes (only one writer at a time)

This ensures thread-safe access during agent execution.

## Built-in Tools Table

go-code provides nine built-in tools that cover essential software development operations:

| # | Tool Name | Purpose | Permission Required | Source File |
|---|-----------|---------|---------------------|--------------|
| 1 | **Read** | Read file contents with optional offset/limit | No | `internal/tool/builtin/read.go` |
| 2 | **Write** | Create or overwrite files | Yes | `internal/tool/builtin/write.go` |
| 3 | **Edit** | Make targeted code edits using exact string matching | Yes | `internal/tool/builtin/edit.go` |
| 4 | **Glob** | Find files by glob patterns (`*`, `**`) | No | `internal/tool/builtin/glob.go` |
| 5 | **Grep** | Search file contents using regular expressions | No | `internal/tool/builtin/grep.go` |
| 6 | **Bash** | Execute shell commands | Yes | `internal/tool/builtin/bash.go` |
| 7 | **Diff** | Compare two content strings and return unified diff output | No | `internal/tool/builtin/diff.go` |
| 8 | **Tree** | Display directory tree structure as text | No | `internal/tool/builtin/tree.go` |
| 9 | **WebFetch** | Fetch URL and return readable text (HTML stripped) | Yes | `internal/tool/builtin/webfetch.go` |

### Tool Details

#### Read Tool
```go
type ReadTool struct{}
```
- Reads file contents line by line
- Supports optional `offset` and `limit` parameters
- Maximum file size: 200KB
- Returns line-numbered output
- No permission required

#### Write Tool
```go
type WriteTool struct{}
```
- Creates new files or overwrites existing ones
- Auto-creates parent directories
- Requires permission (destructive operation)

#### Edit Tool
```go
type EditTool struct{}
```
- Replaces exact string sequences
- Safety feature: requires exact match to prevent unintended edits
- Requires permission

#### Glob Tool
```go
type GlobTool struct{}
```
- Pattern matching for file discovery
- Supports `*` (single level) and `**` (recursive)
- No permission required

#### Grep Tool
```go
type GrepTool struct{}
```
- Full regular expression search
- Returns matching lines with context
- No permission required

#### Bash Tool
```go
type BashTool struct {
    workingDir string
}
```
- Executes shell commands
- Configurable timeout (default: 120s)
- Output truncation at 100KB
- Working directory set to project root
- Requires permission

#### Diff Tool
```go
type DiffTool struct{}
```
- Compares two content strings
- Returns unified diff format
- Uses diff command if available, falls back to pure Go
- No permission required

#### Tree Tool
```go
type TreeTool struct{}
```
- Displays directory tree structure
- Configurable max depth (default: 3)
- Visual format with `├──` and `└──` connectors
- No permission required

#### WebFetch Tool
```go
type WebFetchTool struct{}
```
- Fetches URLs and returns readable text
- Strips HTML tags automatically
- Output limited to 50KB
- Requires permission (network access)

## How to Extend with Custom Tools

To add a custom tool:

### Step 1: Create the Tool Implementation

Create a new file in `internal/tool/builtin/` or a new package:

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
    return false  // Set to true if the tool needs user approval
}

func (t *MyTool) Execute(ctx context.Context, 
                         input map[string]any) tool.Result {
    // Implementation here
    param := input["param"].(string)
    // ... do something ...
    return tool.Success("result")
}
```

### Step 2: Register the Tool

In `internal/tool/init/register.go`:

```go
func RegisterBuiltinTools(r *Registry, workingDir string) error {
    // Existing tools
    r.Register(builtin.NewReadTool())
    r.Register(builtin.NewWriteTool())
    r.Register(builtin.NewEditTool())
    r.Register(builtin.NewGlobTool())
    r.Register(builtin.NewGrepTool())
    r.Register(builtin.NewBashTool(workingDir))

    // Add your custom tool
    r.Register(mytool.NewMyTool())
    
    return nil
}
```

### Step 3: Rebuild and Test

```bash
go build ./cmd/go-code
./go-code "Your test prompt"
```

## MCP Tool Adapter Overview

go-code integrates with the Model Context Protocol (MCP) to use external tools and services. The MCP adapter wraps external tools to match the go-code tool interface.

### MCP Adapter

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

The adapter:
- Converts MCP tool names to format `mcp__{serverName}__{toolName}`
- Delegates execution to the MCP client
- Always requires permission (external tools)

### MCP Components

| File | Purpose |
|------|---------|
| `internal/tool/mcp/config.go` | Configuration loading and environment variable interpolation |
| `internal/tool/mcp/client.go` | MCP protocol client implementation |
| `internal/tool/mcp/adapter.go` | Adapter to convert MCP tools to go-code interface |
| `internal/tool/mcp/transport.go` | Transport layer for stdio-based MCP communication |
| `internal/tool/mcp/manager.go` | MCP server lifecycle management |

See [MCP Integration](./mcp.md) for detailed documentation.

## Related Documentation

- [Agent Loop Implementation](../core-code/agent-loop-impl.md) — Tool execution in the loop
- [MCP Integration](./mcp.md) — Model Context Protocol support
- [Entry Point Walkthrough](../core-code/entry-point.md) — Tool registration during startup

---

<div class="nav-prev-next">

- [Agent Loop Implementation](../core-code/agent-loop-impl.md) ←
- → [MCP Integration](./mcp.md)

</div>