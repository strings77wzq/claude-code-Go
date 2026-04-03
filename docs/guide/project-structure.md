---
title: Project Structure
description: Understanding the directory layout and module organization of claude-code-Go
---

# Project Structure

This document provides a complete overview of the claude-code-Go project directory structure and explains the responsibility of each module.

## Directory Tree

```
claude-code-Go/
├── cmd/go-code/              # 🚀 Main entry point
│   └── main.go               # Application bootstrap + signal handling
│
├── internal/                 # 🏗️ Core modules (private)
│   ├── agent/                # 🧠 Agent loop + context management
│   │   ├── loop.go          # Core agent execution cycle
│   │   ├── history.go       # Message history tracking
│   │   └── compact.go       # Context compaction logic
│   │
│   ├── api/                  # 🌐 Anthropic API client
│   │   ├── client.go        # HTTP client for Messages API
│   │   ├── stream.go        # SSE token streaming handler
│   │   └── types.go         # API request/response types
│   │
│   ├── config/              # ⚙️ Configuration management
│   │   ├── loader.go        # Multi-source config loader
│   │   ├── loader_test.go   # Config loading tests
│   │   └── types.go         # Configuration structures
│   │
│   ├── permission/          # 🛡️ Permission system
│   │   ├── policy.go        # Permission policy engine
│   │   ├── rules.go         # Rule definitions
│   │   ├── prompter.go      # User permission prompts
│   │   └── rules_test.go    # Permission rule tests
│   │
│   ├── session/             # 💾 Session persistence
│   │   ├── session.go       # Session state management
│   │   └── session_test.go  # Session tests
│   │
│   ├── hooks/               # 🔌 Pre/post execution hooks
│   │   ├── hooks.go         # Hook interface + registry
│   │   └── builtin.go       # Built-in hook implementations
│   │
│   └── tool/                # 🔧 Tool system
│       ├── tool.go          # Tool interface definition
│       ├── registry.go      # Tool registration + lookup
│       ├── builtin/         # Built-in tools
│       │   ├── read.go      # File reading tool
│       │   ├── write.go     # File writing tool
│       │   ├── edit.go      # Code editing tool
│       │   ├── glob.go      # File pattern matching
│       │   ├── grep.go      # Content search
│       │   └── bash.go      # Shell command execution
│       ├── mcp/             # MCP integration
│       │   ├── manager.go   # MCP server lifecycle
│       │   ├── client.go    # MCP protocol client
│       │   ├── adapter.go   # MCP tool adapter
│       │   ├── config.go    # MCP configuration
│       │   └── transport.go # Transport layer
│       └── init/            # Tool registration
│           └── register.go  # Built-in tool registration
│
├── pkg/tty/                 # 🎨 Terminal UI
│   ├── repl.go             # REPL main loop
│   ├── renderer.go         # Terminal output rendering
│   └── repl_test.go        # REPL tests
│
├── harness/                 # 🧪 Python test harness (optional)
│   ├── mock_server/        # Mock Anthropic API
│   │   └── server.py       # Mock API server
│   ├── evaluators/         # Quality evaluation
│   │   └── evaluator.py    # Response quality checks
│   └── replay/             # Session replay + trace
│       └── replay.py       # Debug replay tool
│
├── docs/                    # 📚 VitePress documentation
│   ├── en/                  # English docs
│   │   ├── guide/          # User guides
│   │   └── architecture/   # Architecture docs
│   └── zh/                  # Chinese docs
│       ├── guide/
│       └── architecture/
│
├── bin/                     # 📦 Built binaries (generated)
│
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies
├── Makefile                # Build automation
└── README.md               # Project readme
```

## Module Responsibilities

| Module | Responsibility | Key Files |
|--------|---------------|------------|
| **cmd/go-code** | Application entry point, signal handling, component initialization | `main.go` |
| **internal/agent** | Core agent loop execution, message history, context compaction | `loop.go`, `history.go`, `compact.go` |
| **internal/api** | Anthropic Messages API communication, SSE streaming | `client.go`, `stream.go` |
| **internal/config** | Multi-source configuration loading (env, file, defaults) | `loader.go`, `types.go` |
| **internal/permission** | Three-tier permission system, user approval prompts | `policy.go`, `rules.go`, `prompter.go` |
| **internal/session** | Session state persistence and management | `session.go` |
| **internal/hooks** | Pre/post execution hooks for extensibility | `hooks.go`, `builtin.go` |
| **internal/tool** | Tool registry, built-in tools, MCP integration | `registry.go`, `builtin/*`, `mcp/*` |
| **pkg/tty** | Terminal REPL, input handling, output rendering | `repl.go`, `renderer.go` |

## Dependency Relationships

The module dependencies follow a unidirectional flow from top to bottom:

```
                    ┌──────────────┐
                    │ cmd/go-code  │  (Entry point)
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │  config  │ │   agent   │ │   tool   │
        └────┬─────┘ └─────┬─────┘ └────┬─────┘
             │             │            │
             ▼             ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │   api    │ │ session  │ │permission│
        └──────────┘ └────┬─────┘ └────┬─────┘
                           │            │
                           ▼            ▼
                     ┌──────────┐ ┌──────────┐
                     │  hooks   │ │  (via    │
                     └──────────┘ │ agent)   │
                                   └──────────┘

              ┌─────────────────────────────────┐
              │           pkg/tty               │  (Uses agent)
              └─────────────────────────────────┘
```

**Key principles:**
- Dependencies are **unidirectional** — no circular dependencies
- `internal/*` modules are private and form the core
- `pkg/tty` depends on `internal/agent` for agent functionality
- All modules ultimately flow through `cmd/go-code` as the composition root

## Design Note: Facade Pattern

The `AgentLoop` in `internal/agent/loop.go` serves as a **facade** that abstracts away the complexity of:

- API communication
- Tool execution
- Permission checking
- Session management
- Hook invocation

The REPL in `pkg/tty` only needs to interact with this single interface, keeping the presentation layer simple and decoupled from the core logic.

## Next Steps

- [Architecture Overview](../architecture/overview.md) — Deep dive into the system design
- [Agent Loop Deep Dive](../architecture/agent-loop.md) — Understand the core execution cycle
- [Tools Architecture](../architecture/tools.md) — Learn about the tool system and MCP support