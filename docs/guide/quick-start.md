---
title: Quick Start Guide
description: Get started with go-code in minutes — configure your API key and run your first command
---

# Quick Start Guide

This guide walks you through running go-code for the first time.

## Prerequisites

Before running go-code, ensure you have:

1. A built binary (see [Installation Guide](installation.md))
2. An API key (Anthropic, OpenAI, or compatible)

## Configure Your API Key

go-code requires an API key to communicate with the LLM. You can set this in two ways:

### Option 1: Environment Variable

```bash
export ANTHROPIC_API_KEY=sk-ant-your-api-key-here
```

Add this to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) for persistence.

### Option 2: Config File

Create a configuration file at `~/.go-code/settings.json`:

```json
{
  "apiKey": "sk-ant-your-api-key-here"
}
```

The config loader searches in this order (later sources override earlier ones):
1. CLI arguments (highest priority)
2. Environment variables: `ANTHROPIC_API_KEY`
3. Project config: `./.go-code/settings.json`
4. User config: `~/.go-code/settings.json`
5. Built-in defaults (lowest priority)

## Running go-code

### Interactive REPL Mode

Start an interactive session:

```bash
# If installed via go install
go-code

# Or if using the built binary
./bin/go-code
```

You should see the welcome screen:

```
  ____   _    ____ ___ 
 |  _ \ / \  / ___|_ _|
 | |_) / _ \ \___ \| | 
 |  __/ ___ \ ___) | | 
 | |_| /_/   \_\____/___|

Welcome to go-code 0.1.0
Type /help for available commands

go-code> 
```

Try asking a question:

```
go-code> What files are in the current directory?
```

The agent will use its tools to explore the filesystem and provide an answer.

### Single Command Mode

For one-off commands, use the `-p` flag:

```bash
go-code -p "Create a hello world program in Go"
```

This executes the command and exits when complete.

#### Output Formats

```bash
# Plain text (default)
go-code -p "Explain the agent loop"

# JSON output (for scripting)
go-code -p "Explain the agent loop" -f json

# Quiet mode (no spinner, for scripts)
go-code -p "What is 2+2?" -q
```

## Available Commands

In interactive mode, you can use these special commands:

| Command | Description |
|---------|-------------|
| `/help` | Show available commands |
| `/clear` | Clear conversation history |
| `/exit` | Exit the program |
| `/quit` | Exit the program (same as /exit) |
| `/model` | Show current model |
| `/model <name>` | Switch to a different model |
| `/models` | List all available models |
| `/sessions` | List saved sessions |
| `/resume <id>` | Resume a saved session |
| `/compact` | Compact conversation context |
| `/update` | Check for updates |

## Startup Parameters

go-code supports the following startup options:

### Positional Arguments

```bash
go-code [prompt]
```

- `prompt` (optional): If provided, go-code executes this single prompt and exits

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ANTHROPIC_API_KEY` | Your API key | Yes |
| `ANTHROPIC_BASE_URL` | Override default API endpoint (optional) | No |
| `ANTHROPIC_MODEL` | Specify model (default: claude-sonnet-4-20250514) | No |

### Startup Parameters

| Flag | Description |
|------|-------------|
| `-p "prompt"` | Run a single prompt and exit |
| `-f json` | Output in JSON format (use with `-p`) |
| `-q` | Quiet mode, no spinner (use with `-p`) |
| `--legacy-repl` | Use the old REPL instead of bubbletea TUI |

### Configuration File

Create `~/.go-code/settings.json` for persistent settings:

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "model": "claude-sonnet-4-20250514",
  "baseUrl": "https://api.anthropic.com"
}
```

## Basic Usage Examples

### Explore a Project

```
go-code> Find all Go files in this project and list their names
```

### Read and Understand Code

```
go-code> Read the main.go file and explain what it does
```

### Create New Files

```
go-code> Write a simple HTTP server in Go that listens on port 8080
```

### Execute Commands

```
go-code> Run the tests in this project
```

## Permission System

When a tool needs to perform a potentially dangerous operation, go-code prompts for permission:

```
go-code> Delete all files in the current directory

⚠️  This will delete multiple files. Approve? (yes/no): no
```

Type `yes` to approve or `no` to deny.

The permission system controls:
- File deletion and overwriting
- Shell command execution
- Network requests
- Other potentially harmful operations

## Next Steps

- [Configuration Guide](configuration.md) — Customize behavior with advanced settings
- [Architecture Overview](../architecture/overview.md) — Understand how go-code works
- [Agent Loop Deep Dive](../architecture/agent-loop.md) — Learn about the core execution cycle
