---
title: Hooks System
description: Technical deep dive into the Hooks system — Hook interface, built-in implementations, custom hook creation, registration, and execution flow
---

# Hooks System

go-code implements a Hooks system that provides pre/post execution callbacks for monitoring and extending tool behavior. This document provides a comprehensive technical overview.

## What Are Hooks?

**Hooks** are callbacks that get triggered before and after tool execution. They provide a mechanism for:

- **Monitoring** — Track tool invocations and results
- **Logging** — Record detailed execution information
- **Auditing** — Maintain compliance-ready audit trails
- **Metrics** — Collect performance and usage statistics
- **Notifications** — Alert on specific tool operations
- **Validation** — Pre-check inputs or validate outputs

## Hook Interface

The Hook interface defines two methods:

```go
type Hook interface {
    // Name returns the unique identifier for this hook.
    Name() string

    // PreExecute is called before a tool is executed.
    // It receives the tool name and input parameters.
    // Returning an error will prevent tool execution.
    PreExecute(toolName string, input map[string]any) error

    // PostExecute is called after a tool has been executed.
    // It receives the tool name, input parameters, execution result, and whether the result is an error.
    // Errors from PostExecute are logged but do not affect the tool result.
    PostExecute(toolName string, input map[string]any, result string, isError bool) error
}
```

### PreExecute

The `PreExecute` method is called before a tool runs. Use cases include:
- Validating input parameters
- Checking permissions
- Starting performance timers
- Logging the imminent execution
- Preventing execution based on conditions

If `PreExecute` returns an error, the tool execution is **blocked** and the error is propagated to the caller.

### PostExecute

The `PostExecute` method is called after a tool completes. Use cases include:
- Logging execution results
- Recording metrics and durations
- Writing audit records
- Processing results
- Sending notifications

Errors from `PostExecute` are **logged but ignored** — they do not affect the tool's result or the agent's execution flow.

## Hook Registry

The `Registry` manages hook registration and execution:

```go
type Registry struct {
    mu        sync.RWMutex
    preHooks  []Hook
    postHooks []Hook
}
```

### Key Methods

```go
// Register adds a hook to the registry.
func (r *Registry) Register(hook Hook) error

// RunPreHooks executes all pre-execute hooks in order.
// Returns the first error encountered (stops execution).
func (r *Registry) RunPreHooks(toolName string, input map[string]any) error

// RunPostHooks executes all post-execute hooks in order.
// Errors are logged but do not stop execution of subsequent hooks.
func (r *Registry) RunPostHooks(toolName string, input map[string]any, result string, isError bool)
```

### Registration Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Hook Registration Flow                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   1. Create hook instance (LoggingHook, AuditHook, custom)        │
│         │                                                           │
│         ▼                                                           │
│   2. Call registry.Register(hook)                                 │
│         │                                                           │
│         ▼                                                           │
│   3. Validate no duplicate names                                   │
│         │                                                           │
│         ▼                                                           │
│   4. Add hook to preHooks and postHooks slices                    │
│         │                                                           │
│         ▼                                                           │
│   5. Hook ready for execution on all tool calls                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Execution Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Hook Execution Flow                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Tool.Execute(input) called                                       │
│         │                                                           │
│         ▼                                                           │
│   ┌─────────────────────────────────────────────────────────────┐ │
│   │ RunPreHooks(toolName, input)                                │ │
│   │   - Hook1.PreExecute() → OK                                 │ │
│   │   - Hook2.PreExecute() → OK                                 │ │
│   │   - Hook3.PreExecute() → ERROR → STOP                      │ │
│   └─────────────────────────────────────────────────────────────┘ │
│         │                                                           │
│         ▼ (if all pre-hooks pass)                                  │
│   Tool executes (Read, Write, Edit, Glob, Grep, Bash, MCP)         │
│         │                                                           │
│         ▼                                                           │
│   ┌─────────────────────────────────────────────────────────────┐ │
│   │ RunPostHooks(toolName, input, result, isError)              │ │
│   │   - Hook1.PostExecute() → OK                                │ │
│   │   - Hook2.PostExecute() → OK                                │ │
│   │   - Hook3.PostExecute() → ERROR (logged, ignored)           │ │
│   └─────────────────────────────────────────────────────────────┘ │
│         │                                                           │
│         ▼                                                           │
│   Return tool result to agent                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Built-in Hooks

go-code provides two built-in hooks:

### LoggingHook

The `LoggingHook` logs all tool executions using Go's `slog` package:

```go
type LoggingHook struct {
    logger *slog.Logger
}
```

**Features:**
- Logs tool name, input (truncated), and result (truncated)
- Uses DEBUG level for pre-execute, INFO for success, ERROR for failures
- Configurable truncation length (default: 500 characters)

**Usage:**
```go
hook := hooks.NewLoggingHook()              // Uses default logger
hook := hooks.NewLoggingHookWithLogger(l)   // Uses custom logger
registry.Register(hook)
```

### AuditHook

The `AuditHook` records tool executions to a JSONL audit log file:

```go
type AuditHook struct {
    filePath string
    file     *os.File
}

type AuditRecord struct {
    Timestamp  string         `json:"timestamp"`
    ToolName   string         `json:"tool_name"`
    Input      map[string]any `json:"input"`
    Result     string         `json:"result"`
    IsError    bool           `json:"is_error"`
    DurationMS int64          `json:"duration_ms"`
    PreHookErr string         `json:"pre_hook_error,omitempty"`
}
```

**Features:**
- Writes one JSON record per line (JSONL format)
- Records timestamp, tool name, input, result, error status
- Automatically creates directories and files
- Thread-safe with proper file syncing

**Usage:**
```go
hook, err := hooks.NewAuditHook("/var/log/go-code/audit.jsonl")
if err != nil {
    log.Fatal(err)
}
defer hook.Close()
registry.Register(hook)
```

## Creating Custom Hooks

To create a custom hook, implement the `Hook` interface:

```go
import "log/slog"

type MetricsHook struct {
    counters map[string]int64
    mu       sync.Mutex
}

func (h *MetricsHook) Name() string {
    return "metrics"
}

func (h *MetricsHook) PreExecute(toolName string, input map[string]any) error {
    h.mu.Lock()
    h.counters[toolName]++
    h.mu.Unlock()
    return nil
}

func (h *MetricsHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    // Could record timing, success rate, etc.
    return nil
}
```

### Example: Notification Hook

```go
type NotificationHook struct {
    channels []chan string
}

func (h *NotificationHook) Name() string {
    return "notification"
}

func (h *NotificationHook) PreExecute(toolName string, input map[string]any) error {
    if toolName == "Bash" {
        for _, ch := range h.channels {
            ch <- fmt.Sprintf("Executing dangerous command: %v", input["command"])
        }
    }
    return nil
}

func (h *NotificationHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    return nil
}
```

### Example: Validation Hook

```go
type ValidationHook struct {
    allowedPaths []string
}

func (h *ValidationHook) Name() string {
    return "validation"
}

func (h *ValidationHook) PreExecute(toolName string, input map[string]any) error {
    if toolName == "Write" || toolName == "Edit" {
        if path, ok := input["filePath"].(string); ok {
            if !h.isPathAllowed(path) {
                return fmt.Errorf("path not allowed: %s", path)
            }
        }
    }
    return nil
}

func (h *ValidationHook) isPathAllowed(path string) bool {
    for _, allowed := range h.allowedPaths {
        if strings.HasPrefix(path, allowed) {
            return true
        }
    }
    return false
}

func (h *ValidationHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    return nil
}
```

## Use Cases

### Logging and Monitoring

```go
// Track all tool usage
registry.Register(hooks.NewLoggingHook())

// Query logs with:
grep "tool execution" /var/log/app.log
```

### Compliance Auditing

```go
// Record all operations for compliance
hook, _ := hooks.NewAuditHook("/audit/tool-access.jsonl")
registry.Register(hook)

// Review audit log:
cat /audit/tool-access.jsonl | jq '.'
```

### Performance Metrics

```go
// Track tool usage and performance
type PerfHook struct{}

func (h *PerfHook) Name() string { return "perf" }
func (h *PerfHook) PreExecute(t string, i map[string]any) error { /* start timer */ return nil }
func (h *PerfHook) PostExecute(t string, i map[string]any, r string, e bool) error {
    // Record duration
    return nil
}
registry.Register(&PerfHook{})
```

### Security Validation

```go
// Block dangerous operations
type SecurityHook struct{}

func (h *SecurityHook) Name() string { return "security" }
func (h *SecurityHook) PreExecute(t string, i map[string]any) error {
    if t == "Bash" && strings.Contains(i["command"].(string), "rm -rf") {
        return errors.New("blocked dangerous command")
    }
    return nil
}
func (h *SecurityHook) PostExecute(t string, i map[string]any, r string, e bool) error {
    return nil
}
registry.Register(&SecurityHook{})
```

### Notifications

```go
// Alert on specific events
notifHook := &NotificationHook{channels: []chan string{alertChannel}}
registry.Register(notifHook)
```

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Hooks Architecture                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    Tool Execution                             │   │
│   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐       │   │
│   │  │  PreExecute │───▶│    Tool    │───▶│ PostExecute │       │   │
│   │  │   Hooks     │    │   (Read,   │    │    Hooks    │       │   │
│   │  │             │    │  Write,    │    │             │       │   │
│   │  │ - Logging   │    │  Edit,     │    │ - Logging   │       │   │
│   │  │ - Audit     │    │  Glob,     │    │ - Audit     │       │   │
│   │  │ - Custom    │    │  Grep,     │    │ - Custom    │       │   │
│   │  └─────────────┘    │  Bash,     │    └─────────────┘       │   │
│   │                      │  MCP)      │                          │   │
│   │                      └─────────────┘                          │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│                              ▼                                       │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    Hook Registry                            │   │
│   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │   │
│   │  │  Register  │    │  RunPreHooks│    │ RunPostHooks│      │   │
│   │  └─────────────┘    └─────────────┘    └─────────────┘      │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Related Documentation

- [Skills System](./skills.md) — Named prompts for customizing agent behavior
- [MCP Integration](./mcp.md) — External tool integration
- [Tool System Overview](../tools/overview.md) — Tool interface and registry
- [Agent Loop Implementation](../core-code/agent-loop-impl.md) — Tool execution flow

---

<div class="nav-prev-next">

- [MCP Integration](./mcp.md) ←
- → [Extension Overview](./overview.md)

</div>