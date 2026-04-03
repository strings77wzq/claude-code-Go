# go-code — Claude Code in Go

[![Go Version](https://img.shields.io/badge/Go-1.23-blue)](https://go.dev/doc/install)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A Go implementation of Anthropic's Claude Code (Cline) agent system. Implements the full agent loop with built-in tools, permission system, and MCP support.

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/user/go-code.git
cd go-code

# Build the binary
make build

# Or install to your PATH
go install ./cmd/go-code
```

### Configure API Key

Set the `ANTHROPIC_API_KEY` environment variable:

```bash
export ANTHROPIC_API_KEY=your-api-key-here
```

Or create a config file at `~/.config/go-code/config.yaml`:

```yaml
api_key: "your-api-key-here"
```

### Run

```bash
# Run in interactive REPL mode
./bin/go-code

# Or run with a prompt directly
./bin/go-code "Write a hello world program in Go"
```

## Features

- **Agent Loop**: Full implementation of the Claude Code agent loop with message history and state management
- **6 Built-in Tools**:
  - `Read` - Read files from the filesystem
  - `Write` - Write/create files
  - `Edit` - Make targeted code edits
  - `Glob` - Find files by pattern
  - `Grep` - Search file contents
  - `Bash` - Execute shell commands
- **Permission System**: User approval for dangerous operations (file deletions, network requests)
- **MCP Support**: Model Context Protocol integration for extending capabilities
- **Streaming**: Real-time token streaming from the API

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        go-code                               │
├─────────────────────────────────────────────────────────────┤
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
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     Python Harness (Optional)               │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────┐  │
│  │ Mock Server  │    │  Replay      │    │  Test Suite   │  │
│  │              │    │  Tests       │    │               │  │
│  └──────────────┘    └──────────────┘    └───────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
go-code/
├── cmd/go-code/          # Main entry point
├── internal/
│   ├── agent/            # Agent loop implementation
│   ├── api/              # Anthropic API client
│   ├── config/           # Configuration loading
│   ├── permission/       # Permission system
│   └── tool/             # Tool implementations
│       └── builtin/      # Built-in tools (Read, Write, Edit, Glob, Grep, Bash)
├── harness/             # Python test harness (optional)
│   └── mock_server/      # Mock server for testing
├── docs/                 # MkDocs documentation
├── Makefile              # Build targets
└── README.md             # This file
```

### Go Runtime vs Python Harness

The main `go-code` binary is a pure Go implementation that:
- Communicates directly with the Anthropic API
- Executes tools locally
- Manages agent state

The Python harness (`harness/`) is optional and provides:
- Mock API server for testing
- Replay functionality for debugging
- Integration test suite

## Development

### Build

```bash
make build
```

### Test

```bash
# Run all tests (Go + Python if harness exists)
make test

# Run only Go tests
go test -v ./...

# Run only Python tests (if harness exists)
cd harness && python -m pytest -v
```

### Static Analysis

```bash
make vet
```

### Documentation

```bash
# Serve documentation locally
make docs

# Build for production
cd docs && mkdocs build
```

### Cross-Platform Build

```bash
# Build for Linux amd64, macOS amd64, macOS arm64
make build-all
```

Output binaries are placed in `bin/`.

### Clean

```bash
make clean
```

## License

MIT License - see [LICENSE](LICENSE) for details.