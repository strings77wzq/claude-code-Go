# Quick Start

This guide walks through running go-code for the first time.

## Prerequisites

Before running go-code, you need:

1. An Anthropic API key
2. The binary built (see [Installation](install.md))

## Configure Your API Key

### Option 1: Environment Variable

```bash
export ANTHROPIC_API_KEY=sk-ant-your-api-key-here
```

Add this to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) for persistence.

### Option 2: Config File

Create `~/.config/go-code/config.yaml`:

```yaml
api_key: "sk-ant-your-api-key-here"
```

## Run in REPL Mode

Start an interactive session:

```bash
./bin/go-code
```

You should see a prompt where you can type requests:

```
go-code> What is the current directory?
```

The agent will use tools as needed to answer your question.

## Run with a Single Prompt

For one-off commands:

```bash
./bin/go-code "Create a hello world program in Go"
```

This will execute the command and exit.

## Example Session

```
$ ./bin/go-code
go-code> List the files in the current directory

[Agent thinking...]
[Using tool: Glob with pattern: *]
[Tool result: Found 3 files: main.go, Makefile, README.md]

The current directory contains 3 files:
- main.go
- Makefile
- README.md
```

## Permission Prompts

When a tool needs to perform a potentially dangerous operation, go-code will prompt for permission:

```
go-code> Delete all files in the current directory

⚠️  This will delete 3 files. Approve? (yes/no): no
```

Type `yes` to approve or `no` to deny.

## Next Steps

- [Configuration](config.md) - Customize behavior
- [Architecture Overview](../architecture/overview.md) - Understand how it works