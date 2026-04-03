# Architecture Overview

go-code is organized into several key components that work together to implement the Claude Code agent system.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        go-code                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────┐  │
│  │   CLI/Repl   │───▶│ Agent Loop   │───▶│ Tool Registry │  │
│  └──────────────┘    └──────────────┘    └───────────────┘  │
│                            │                    │            │
│                            ▼                    ▼            │
│                     ┌──────────────┐    ┌───────────────┐  │
│                     │ API Client   │    │   Built-in    │  │
│                     │   (Stream)   │    │    Tools      │  │
│                     └──────────────┘    └───────────────┘  │
│                            │                               │
│                            ▼                               │
│                     ┌──────────────┐                       │
│                     │ Anthropic    │                       │
│                     │   API        │                       │
│                     └──────────────┘                       │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐                      │
│  │   Config    │    │  Permission  │                      │
│  │   Loader    │    │   System     │                      │
│  └──────────────┘    └──────────────┘                      │
└─────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### CLI/Repl

The command-line interface provides the user-facing entry point. It handles:
- Parsing command-line arguments
- Managing the interactive REPL session
- Streaming user input to the agent

### Agent Loop

The core orchestrator that manages the agent's execution cycle. It:
- Maintains conversation history
- Decides which tools to call based on model responses
- Handles tool execution and result processing
- Manages the loop until the task is complete

### API Client

Handles communication with the Anthropic API:
- Sends messages with tool definitions
- Receives streaming responses
- Parses tool call requests from the model
- Manages authentication and rate limiting

### Tool Registry

Central registry for all available tools:
- Manages built-in tool registration
- Handles MCP tool discovery
- Provides tool schema to the model
- Routes tool calls to appropriate implementations

### Built-in Tools

Six core tools that enable file system and command execution:
- **Read**: Read file contents
- **Write**: Create or overwrite files
- **Edit**: Make targeted modifications
- **Glob**: Find files by pattern
- **Grep**: Search file contents
- **Bash**: Execute shell commands

### Permission System

Ensures user control over dangerous operations:
- Intercepts potentially harmful tool calls
- Prompts for user approval
- Manages auto-approval rules
- Logs permission decisions

### Config Loader

Loads and manages configuration:
- Reads config files (YAML)
- Processes environment variables
- Provides configuration to all components
- Handles default values

## Directory Structure

```
internal/
├── agent/          # Agent loop implementation
│   ├── loop.go     # Main loop logic
│   └── history.go  # Message history management
├── api/            # Anthropic API client
│   ├── client.go   # HTTP client
│   ├── stream.go   # Streaming handling
│   └── types.go    # API types
├── config/         # Configuration
│   ├── loader.go   # Config file loading
│   └── types.go    # Config types
├── permission/     # Permission system
│   ├── policy.go   # Permission policies
│   ├── rules.go    # Permission rules
│   └── prompter.go # User prompts
└── tool/           # Tool system
    ├── registry.go # Tool registry
    ├── tool.go     # Tool interface
    ├── builtin/    # Built-in tools
    │   ├── read.go
    │   ├── write.go
    │   ├── edit.go
    │   ├── glob.go
    │   ├── grep.go
    │   └── bash.go
    └── mcp/        # MCP integration
        ├── client.go
        └── transport.go
```

## Data Flow

1. User enters a prompt in the CLI/Repl
2. Agent Loop receives the prompt and adds it to history
3. Agent Loop sends the full message history to the API Client
4. API Client sends request to Anthropic API and receives streaming response
5. If the response contains tool calls, the Agent Loop:
   - Routes each tool call to the Tool Registry
   - Checks permissions via the Permission System
   - Executes the tool and captures the result
   - Adds the tool result to the message history
6. Loop continues until the model provides a final response
7. Final response is displayed to the user

## Related Documentation

- [Agent Loop](agent-loop.md) - Detailed agent loop explanation
- [Built-in Tools](tools.md) - Tool implementations
- [Python Harness](../harness/overview.md) - Test infrastructure