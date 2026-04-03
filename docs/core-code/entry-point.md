---
title: Entry Point Walkthrough
description: Deep dive into main.go — component initialization sequence, signal handling, and system prompt
---

# Entry Point Walkthrough

The main entry point of go-code is located at `cmd/go-code/main.go`. This document provides a detailed walkthrough of how the application initializes all its components in the correct order.

## Overview

When you run `go-code`, the following sequence occurs:

```
┌─────────────────────────────────────────────────────────────────────┐
│                        main.go Execution                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Logger Setup          → Initialize structured logging           │
│  2. Signal Handler       → Register SIGINT/SIGTERM handlers        │
│  3. Load Configuration    → Load from env vars / config file       │
│  4. Create API Client     → Initialize Anthropic API client        │
│  5. Create Tool Registry  → Create empty registry for tools        │
│  6. Register Built-in     → Register 6 built-in tools              │
│  7. Create Permission     → Set up 3-tier permission policy        │
│  8. Create Agent          → Initialize agent with all dependencies  │
│  9. Start REPL            → Begin interactive session (blocks)       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Source Code Analysis

### Imports

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

Key dependencies:
- `log/slog` — Structured logging (Go 1.21+)
- `os/signal` — Graceful shutdown handling
- `internal/*` — Core application components
- `pkg/tty` — REPL implementation

### Version Constant

```go
const version = "0.1.0"
```

### System Prompt

```go
const systemPrompt = "You are an interactive agent that helps users with software engineering tasks. You have access to tools for reading files, editing files, executing shell commands, searching code, and more. Use your tools to complete tasks efficiently and accurately."
```

The system prompt defines the agent's role and capabilities. This is sent with every API request to guide the model's behavior.

## Initialization Sequence

### Step 1: Logger Setup

```go
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)
```

The application uses Go's built-in `log/slog` for structured logging. All log messages include timestamps and severity levels.

### Step 2: Signal Handling

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

The application registers handlers for:
- `SIGINT` — Interrupt (Ctrl+C)
- `SIGTERM` — Termination request

When received, the application logs the event and exits gracefully.

### Step 3: Load Configuration

```go
logger.Info("Loading configuration")
cfg, err := config.Load(nil)
if err != nil {
    logger.Error("Failed to load configuration", "error", err)
    os.Exit(1)
}
logger.Info("Configuration loaded", "model", cfg.Model, "baseURL", cfg.BaseURL)
```

Configuration is loaded from multiple sources (environment variables, config file). See [Configuration Guide](../guide/configuration.md) for details.

### Step 4: Create API Client

```go
logger.Info("Creating API client")
client := api.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
logger.Info("API client created")
```

The API client is initialized with credentials and model settings. It handles HTTP communication with the Anthropic API.

### Step 5: Create Tool Registry

```go
logger.Info("Creating tool registry")
registry := tool.NewRegistry()
logger.Info("Tool registry created")
```

The registry is a thread-safe container for all available tools. It uses `sync.RWMutex` for concurrent reads and safe writes.

### Step 6: Register Built-in Tools

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

Six built-in tools are registered:

| Tool | Purpose |
|------|---------|
| **Read** | Read file contents |
| **Write** | Create/overwrite files |
| **Edit** | Modify specific code sections |
| **Glob** | Find files by pattern |
| **Grep** | Search file contents |
| **Bash** | Execute shell commands |

The working directory is passed to the Bash tool so commands execute in the correct context.

### Step 7: Create Permission Policy

```go
logger.Info("Creating permission policy")
policy := permission.NewPolicy(permission.WorkspaceWrite)
logger.Info("Permission policy created")
```

The permission policy determines which operations require user approval. See [Architecture Overview](overview.md) for details.

### Step 8: Create Agent

```go
logger.Info("Creating agent")
agentInstance := agent.NewAgent(client, registry, policy, systemPrompt, cfg.Model)
logger.Info("Agent started", "model", cfg.Model)
```

The agent is created with all its dependencies:
- API client for model communication
- Tool registry for execution
- Permission policy for access control
- System prompt for behavior guidance
- Model identifier for API requests

### Step 9: Start REPL

```go
logger.Info("Starting REPL")
repl := tty.NewREPL(agentInstance, version, cfg.Model)

// Run REPL - this blocks until exit
repl.Run()

logger.Info("REPL exited")
```

The REPL (Read-Eval-Print Loop) is the interactive interface. It:
1. Reads user input from the terminal
2. Passes it to the agent for processing
3. Displays the response
4. Repeats until the user exits (Ctrl+C or `exit` command)

## Key Initialization Concepts

### Dependency Injection

The application uses dependency injection throughout. Each component receives its dependencies through constructor functions:

```go
agent.NewAgent(client, registry, policy, systemPrompt, model)
```

This makes the code testable and modular.

### Order Matters

The initialization order is critical:

1. Logger must be ready first (other components log during init)
2. Config must load before components that need it
3. Registry must exist before tools can be registered
4. All components must be ready before the Agent is created

### Error Handling

Each initialization step checks for errors and exits with a descriptive message:

```go
if err != nil {
    logger.Error("Failed to load configuration", "error", err)
    os.Exit(1)
}
```

## Related Documentation

- [Agent Loop Implementation](agent-loop-impl.md) — The Run() method and execution flow
- [Tool System Overview](../tools/overview.md) — Tool interface and registry
- [MCP Integration](../architecture/mcp.md) — Model Context Protocol support
- [Configuration Guide](../guide/configuration.md) — Configuration options

---

<div class="nav-prev-next">

- [Architecture Overview](../architecture/overview.md) ←
- → [Agent Loop Implementation](agent-loop-impl.md)

</div>