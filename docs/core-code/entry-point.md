---
title: Entry Point Architecture
description: Understanding the initialization sequence as a dependency graph — why order matters, dependency injection patterns, and design decisions
---

# Entry Point Architecture

The entry point is where the application's initialization contract is established. Rather than viewing `main.go` as a sequence of function calls, consider it as **dependency graph assembly** — each component is wired into the system in the correct order according to its dependencies, ensuring the system is fully operational before receiving user input.

## Initialization Dependency Graph

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Initialization Dependency Graph                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐                                                  │
│  │   Logger    │  (Foundation — all components log during init)  │
│  └──────────────┘                                                  │
│         │                                                          │
│         ▼                                                          │
│  ┌──────────────┐                                                  │
│  │   Signal    │  (Graceful shutdown contract)                   │
│  │   Handling  │                                                  │
│  └──────────────┘                                                  │
│         │                                                          │
│         ▼                                                          │
│  ┌──────────────┐    ┌──────────────┐                             │
│  │    Config    │───▶│  API Client  │  (Requires credentials)      │
│  └──────────────┘    └──────────────┘                             │
│         │                  │                                       │
│         │                  ▼                                       │
│         │          ┌──────────────┐                               │
│         └─────────▶│ Tool Registry │  (Requires working dir)     │
│                     └──────────────┘                               │
│                            │                                        │
│                            ▼                                       │
│                     ┌──────────────┐    ┌──────────────┐          │
│                     │ Permission   │◀───│    Agent     │          │
│                     │   Policy     │    │              │          │
│                     └──────────────┘    └──────────────┘          │
│                            │                     │                 │
│                            └─────────────────────┘                 │
│                                              ▼                     │
│                                    ┌──────────────┐                │
│                                    │    REPL      │  (Blocking)    │
│                                    └──────────────┘                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Why This Order Matters

The initialization sequence is not arbitrary — each step establishes preconditions for subsequent steps:

### 1. Logger First — Foundation Layer

The logger must be initialized before any other component because **every subsequent component logs during its initialization**. Without a log destination, errors during config loading, client creation, or tool registration would have nowhere to go.

This follows the **Inversion of Control** principle: the application provides a global logger rather than each component creating its own. This enables:
- Consistent log format across all components
- Centralized log level control
- Easy output redirection (file, stderr, syslog)

### 2. Signal — Shutdown Contract

Signal handlers are registered early to establish a **graceful shutdown contract**. The application commits that any signal (SIGINT, SIGTERM) will result in clean shutdown with proper logging — no orphan processes, no lost state.

The signal handler requires no other components to be ready; it only needs the logger. This minimal dependency allows shutdown to work even if initialization fails midway.

### 3. Config → Client → Tools — Dependency Chain

- Config must be loaded before API client creation (client requires credentials)
- API client must exist before agent instantiation
- But there's a subtle ordering between tools and config:
  - **Tool registry** must exist before tool registration (it serves as the container)
  - **Working directory** from config must be passed to tools (especially Bash tool needs to know where to execute commands)

This creates a narrow window: config load → create registry → register tools with working directory. Skipping this order causes tools to execute in the wrong directory or fail to find their dependencies.

### 4. Permission Policy — Security Boundary

Permission policy is created after all other components because it needs to know **which tools exist** to make access decisions. The policy queries the tool registry to determine which tools require permission and what the user has approved.

### 5. Agent — Composition Root

The agent is the **composition root** in DDD terminology — it's where all dependencies are assembled into the single object that drives the application. At this point:
- API client is ready to make requests
- Tool registry contains all available tools
- Permission policy knows what restrictions apply
- System prompt defines the agent's role

The agent doesn't create its dependencies; it receives them. This is **constructor injection**, the purest form of dependency injection.

### 6. REPL — Endpoint

The REPL blocks, waiting for user input. It receives the fully-configured agent as a dependency. This follows the **application controller** pattern — the REPL translates user commands into method calls on the agent.

## Dependency Injection Benefits

The initialization sequence demonstrates three key benefits of dependency injection:

### 1. Testability

Each component can be isolated and tested by providing mock dependencies:

```go
// Production: real client
agent := agent.NewAgent(realClient, registry, policy, prompt, model)

// Testing: mock client that returns predictable responses
agent := agent.NewAgent(mockClient, registry, policy, prompt, model)
```

The agent doesn't know the difference — it depends on the `APIClient` interface, not a concrete implementation.

### 2. Modularity

Components can be swapped without changing other parts of the system. Permission policy can be replaced (e.g., for different security levels), API client can be swapped (e.g., for testing), tools can be added or removed — all without modifying the agent.

### 3. Explicit Dependencies

The constructor signature explicitly documents what the agent requires:

```go
func NewAgent(
    client APIClient,
    registry *Registry,
    policy *Policy,
    systemPrompt string,
    model string,
) *Agent
```

No hidden state, no global variables, no accidental initialization order. **The compiler enforces the dependency graph.**

## Error Handling Strategy — Fail Fast

The initialization sequence follows the **fail-fast** principle:

```go
if err != nil {
    logger.Error("Failed to load configuration", "error", err)
    os.Exit(1)
}
```

Why fail fast?
- **Diagnostics**: Early errors are easier to diagnose — stack trace points directly at the problem
- **Recovery**: Immediate exit is better than continuing in a partially-initialized state
- **User Experience**: Descriptive error messages tell the user what went wrong and how to fix it

Each initialization step validates its preconditions and exits with a clear error message if something goes wrong. No retry logic, no fallback behavior — the system either starts correctly or doesn't.

## Signal Handling — Graceful Degradation

The signal handler is deliberately minimal:

```go
go func() {
    sig := <-sigChan
    logger.Info("Received signal, shutting down", "signal", sig.String())
    logger.Info("Shutdown complete")
    os.Exit(0)
}()
```

Why not more complex?
- The REPL handles user-level cancellation (Ctrl+C cancels the current turn)
- Signal handler handles system-level termination
- No complex state needs to be saved — sessions are saved after each turn
- The application is short-lived (run once, execute a task, exit)

This is **intentional simplicity** — adding more complex shutdown logic would introduce complexity without benefit for this use case.

## Design Patterns Used

### Registry Pattern — Tool Registry

The tool registry uses the **registry pattern** to provide a centralized container for all available tools:

```go
registry := tool.NewRegistry()
```

Key characteristics:
- **Lazy lookup**: Tools are retrieved by name at execution time, not registration time
- **Thread-safe**: Uses `sync.RWMutex` for concurrent reads and safe writes
- **Extensible**: MCP tools can be added at runtime without modifying existing code

The registry decouples **tool definition** (what the API sees) from **tool execution** (what happens when called). The agent doesn't need to know whether a tool is built-in or from an MCP server.

### Strategy Pattern — Permission System

The permission system uses the **strategy pattern** for access decisions:

```go
policy := permission.NewPolicy(permission.WorkspaceWrite)
```

Key characteristics:
- **Declarative**: Policy expresses rules, not implementation
- **Composable**: Multiple policies can be combined (though only one is currently used)
- **Runtime evaluation**: Policy checks specific tool and input, not just tool name

Permission policy is a form of **strategy pattern** — behavior (allow/deny/ask) can vary based on configuration without changing the code that uses it.

### Factory Pattern — Agent Creation

The agent constructor acts as a **factory**, assembling the agent with all dependencies:

```go
agentInstance := agent.NewAgent(client, registry, policy, systemPrompt, model)
```

This is also an example of **builder pattern** semantics — each dependency is explicitly required, making construction self-documenting and impossible to get wrong (compile errors rather than runtime errors).

## Architectural Summary

| Principle | Implementation |
|-----------|----------------|
| **Dependency Graph** | Initialization order respects component dependencies |
| **Composition Root** | Agent is the single location where all dependencies converge |
| **Fail Fast** | Errors are caught early with descriptive messages |
| **Explicit Dependencies** | Constructor injection makes dependencies visible |
| **Registry Pattern** | Centralized tool container with runtime lookup |
| **Strategy Pattern** | Declarative permission rules evaluated at execution time |

Understanding this initialization sequence is fundamental to understanding the entire system — it's where the dependency contract is established, upon which all subsequent operations are built.

---

<div class="nav-prev-next">

- [Architecture Overview](../architecture/overview.md) ←
- → [Agent Loop Implementation](agent-loop-impl.md)

</div>