# Quick Start Guide

Get started with claude-code-Go in 5 minutes.

## Prerequisites

- Go 1.24 or later installed
- An API key from Anthropic, OpenAI, or compatible provider

## Installation

### Option 1: Using go install (Recommended)

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

### Option 2: Build from source

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
go build -o bin/go-code ./cmd/go-code
export PATH="$PWD/bin:$PATH"
```

## Configuration

### Set your API key

**Option A: Environment variable (recommended for quick start)**

```bash
export ANTHROPIC_API_KEY="sk-ant-api03-your-key-here"
```

**Option B: Settings file (persistent)**

Create `~/.go-code/settings.json`:

```json
{
  "apiKey": "sk-ant-api03-your-key-here",
  "model": "claude-sonnet-4-20250514"
}
```

## First Run

Start the interactive REPL:

```bash
go-code
```

You'll see a welcome message and prompt:

```
🤖 claude-code-Go v0.1.0
Model: claude-sonnet-4-20250514

Type /help for available commands

> 
```

## Your First Interaction

Try asking a simple question:

```
> What can you help me with?
```

The AI will respond with an overview of its capabilities.

## Try a Tool

Let's read a file:

```
> Read the README.md file
```

You'll see:
1. The AI decides to use the Read tool
2. The file content is displayed
3. The AI summarizes what it found

## Using Commands

Available slash commands:

- `/help` - Show all commands
- `/quit` or `/q` - Exit the REPL
- `/model <name>` - Switch AI model
- `/compact` - Compact conversation context
- `/clear` - Clear conversation history

## Next Steps

- [Tutorial 2: Your First Tool Call](02-first-tool-call.md) - Learn how tools work
- [Tutorial 3: Understanding the Agent Loop](03-agent-loop.md) - Deep dive into how the AI thinks
- Check out the [Architecture Overview](../architecture/overview.md) for technical details

## Troubleshooting

### "ANTHROPIC_API_KEY not set"

Make sure you've set your API key:

```bash
export ANTHROPIC_API_KEY="your-key-here"
```

Or create the settings file as shown above.

### "command not found: go-code"

Make sure the binary is in your PATH:

```bash
which go-code
# If nothing is returned, add to PATH:
export PATH="$HOME/go/bin:$PATH"
```

### Connection errors

Check your internet connection and API key validity. Try:

```bash
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/models
```

## Getting Help

- [Troubleshooting Guide](../troubleshooting/common-issues.md)
- [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)
- [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
