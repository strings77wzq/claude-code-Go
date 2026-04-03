---
title: Architecture Overview
description: High-level architecture of go-code — components, data flow, and design philosophy
---

# Architecture Overview

go-code is a Go implementation of Anthropic's Claude Code agent system. It combines a large language model (LLM) with a robust execution harness to enable autonomous software engineering tasks.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                              go-code                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────────────┐ │
│  │   CLI/Repl   │───▶│  Agent Loop  │───▶│    Tool Registry     │ │
│  └──────────────┘    └──────────────┘    └───────────────────────┘ │
│                            │                         │              │
│                            ▼                         ▼              │
│                     ┌──────────────┐         ┌───────────────────┐  │
│                     │  API Client  │         │   Built-in Tools  │  │
│                     │   (Stream)   │         │  Bash, Read,      │  │
│                     └──────────────┘         │  Write, Edit,     │  │
│                            │                 │  Glob, Grep       │  │
│                            ▼                 └───────────────────┘  │
│                     ┌──────────────┐                                │
│                     │  Anthropic   │                                │
│                     │     API      │                                │
│                     └──────────────┘                                │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────────────┐ │
│  │    Config    │    │  Permission  │    │    MCP Integration   │ │
│  │    Loader   │    │    System    │    │       (Adapter)       │ │
│  └──────────────┘    └──────────────┘    └───────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

## Component Breakdown

### 1. CLI/RePL

The command-line interface provides the user-facing entry point:

- Parses command-line arguments
- Manages interactive REPL sessions
- Streams user input to the agent
- Displays output and tool execution results

**Location**: `pkg/tty/repl.go`

### 2. Agent Loop

The core orchestrator managing the think → act → observe execution cycle:

- Maintains conversation history
- Sends requests to the API with tools and context
- Processes model responses and stop reasons
- Executes tools and feeds results back to the model
- Manages the loop until the task completes

**Location**: `internal/agent/loop.go`

### 3. API Client

Handles communication with the Anthropic API:

- Sends messages with tool definitions
- Receives streaming responses
- Parses tool call requests from the model
- Manages authentication

**Location**: `internal/api/client.go`, `internal/api/stream.go`

### 4. Tool Registry

Central registry for all available tools:

- Manages built-in tool registration
- Handles MCP tool discovery
- Provides tool schemas to the model
- Routes tool calls to appropriate implementations

**Location**: `internal/tool/registry.go`

### 5. Built-in Tools

Six core tools that enable file system and command execution:

| Tool | Purpose | Permission Required |
|------|---------|---------------------|
| **Read** | Read file contents | No |
| **Write** | Create or overwrite files | Yes |
| **Edit** | Make targeted modifications | Yes |
| **Glob** | Find files by pattern | No |
| **Grep** | Search file contents | No |
| **Bash** | Execute shell commands | Yes |

**Location**: `internal/tool/builtin/`

### 6. Permission System

Ensures user control over dangerous operations:

- Intercepts potentially harmful tool calls
- Prompts for user approval
- Manages auto-approval policies
- Logs permission decisions

**Location**: `internal/permission/`

### 7. Config Loader

Loads and manages configuration:

- Reads JSON config files
- Processes environment variables
- Applies config precedence rules

**Location**: `internal/config/loader.go`

### 8. MCP Integration

Provides Model Context Protocol support:

- MCP server discovery
- Tool adapter for MCP tools
- Transport layer management

**Location**: `internal/tool/mcp/`

## Data Flow

```
User Input
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Agent Loop                               │
│  1. Add user message to history                                │
│  2. Compact history if needed                                  │
│  3. Build API request (tools + messages)                       │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Client                                 │
│  4. Send request to Anthropic API (streaming)                  │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ Anthropic    │
                    │     API      │
                    └──────┬───────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Agent Loop                               │
│  5. Process response based on stop_reason:                     │
│     - end_turn: Return final response                          │
│     - tool_use: Execute tools and continue                     │
│     - max_tokens: Return truncated response                    │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                    ┌──────▼───────┐
                    │   Execute    │
                    │    Tools     │
                    └──────┬───────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Permission System                             │
│  6. Check if tool requires permission                          │
│  7. Prompt user if needed                                      │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Tool Registry                                │
│  8. Route to appropriate tool (built-in or MCP)               │
│  9. Execute tool and capture result                            │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Agent Loop                               │
│  10. Add tool results to history                               │
│  11. Continue loop (back to step 2)                            │
└─────────────────────────────────────────────────────────────────┘
```

## Design Philosophy

go-code follows a clear separation of concerns:

### Model Provides Intelligence

The LLM (Claude) is responsible for:
- Understanding user intent
- Deciding which tools to use
- Interpreting tool results
- Generating natural language responses

### Harness Provides Reliability

The Go runtime provides:
- Tool execution and safety
- Conversation history management
- Permission enforcement
- Configuration management
- Error handling and recovery

This division allows the model to focus on the cognitive task while the harness ensures reliable, safe execution.

## Directory Structure

```
internal/
├── agent/           # Agent loop implementation
│   ├── loop.go     # Main loop logic
│   └── history.go  # Message history management
├── api/            # Anthropic API client
│   ├── client.go   # HTTP client
│   ├── stream.go   # Streaming handling
│   └── types.go    # API types
├── config/         # Configuration
│   ├── loader.go   # Config file loading
│   └── types.go    # Config types
├── permission/    # Permission system
│   ├── policy.go   # Permission policies
│   ├── rules.go    # Permission rules
│   └── prompter.go # User prompts
└── tool/           # Tool system
    ├── registry.go    # Tool registry
    ├── tool.go        # Tool interface
    ├── init/          # Tool initialization
    ├── builtin/       # Built-in tools
    │   ├── read.go    # File reading
    │   ├── write.go   # File writing
    │   ├── edit.go    # File editing
    │   ├── glob.go    # File globbing
    │   ├── grep.go    # Content search
    │   └── bash.go    # Shell execution
    └── mcp/           # MCP integration
        ├── adapter.go # MCP tool adapter
        ├── client.go  # MCP client
        ├── config.go  # MCP configuration
        └── transport.go# MCP transport
```

## Related Documentation

- [Agent Loop Deep Dive](agent-loop.md) — Detailed execution cycle explanation
- [Tool System](tools.md) — Tool implementations and extension
- [Quick Start Guide](../guide/quick-start.md) — Running go-code