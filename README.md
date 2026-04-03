# claude-code-Go — Claude Code in Go

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue)](https://go.dev/doc/install)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A Go implementation of Anthropic's Claude Code (Cline) agent system. Implements the full agent loop with built-in tools, permission system, MCP support, and SSE streaming.

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go

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

Or create a config file at `~/.config/claude-code-go/config.yaml`:

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

## Use Cases

What you can do with claude-code-Go:

- **Automated Code Editing** — Edit multiple files across your project with AI-assisted precision
- **Codebase Exploration** — Search and navigate large codebases with Glob and Grep tools
- **Shell Command Automation** — Execute complex build scripts and CLI operations
- **Interactive Development** — Work with an AI partner in a terminal REPL
- **MCP Integration** — Extend capabilities with Model Context Protocol servers

## Features

- **Agent Loop**: Full implementation of the Claude Code agent loop with message history and state management
- **6 Built-in Tools**:
  - `Read` - Read files from the filesystem
  - `Write` - Write/create files
  - `Edit` - Make targeted code edits
  - `Glob` - Find files by pattern
  - `Grep` - Search file contents
  - `Bash` - Execute shell commands
- **Permission System**: 3-tier permission system with user approval for dangerous operations (file deletions, network requests)
- **MCP Support**: Model Context Protocol integration for extending capabilities
- **SSE Streaming**: Real-time Server-Sent Events token streaming from the API
- **Context Management**: Automatic conversation compaction and context tracking

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      claude-code-Go                         │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────┐  │
│  │   CLI/RePL   │───▶│ Agent Loop   │───▶│ Tool Registry │  │
│  │   (pkg/tty)  │    │   + Context  │    │               │  │
│  └──────────────┘    └──────────────┘    └───────────────┘  │
│                            │                    │            │
│                            ▼                    ▼            │
│                     ┌──────────────┐    ┌───────────────┐  │
│                     │   API Client │    │   Built-in    │  │
│                     │ (SSE Stream) │    │    Tools      │  │
│                     └──────────────┘    └───────────────┘  │
│                            │                               │
│                            ▼                               │
│                     ┌──────────────┐                       │
│                     │  Anthropic   │                       │
│                     │     API      │                       │
│                     └──────────────┘                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   Go Runtime vs Python Harness               │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐    ┌────────────────────────────┐  │
│  │   Go Implementation  │    │   Python Harness (Optional)│  │
│  │ ───────────────────  │    │ ────────────────────────── │  │
│  │ • Direct Anthropic   │    │ • Mock API Server          │  │
│  │   API calls          │    │ • Replay + Debugging       │  │
│  │ • Local tool exec    │    │ • Quality Evaluators       │  │
│  │ • Agent state mgmt   │    │ • Test Suite               │  │
│  └─────────────────────┘    └────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
claude-code-Go/
├── cmd/go-code/          # Main entry point
├── internal/
│   ├── agent/            # Agent loop + context management
│   ├── api/              # Anthropic API client + SSE
│   ├── config/           # Multi-source config loader
│   ├── permission/       # 3-tier permission system
│   ├── tool/             # Tool interface + 6 builtins
│   │   ├── builtin/      # Bash, Read, Write, Edit, Glob, Grep
│   │   ├── mcp/          # MCP integration
│   │   └── init/         # Tool registration
│   └── hooks/            # Pre/post execution hooks
├── pkg/tty/              # REPL + terminal rendering
├── harness/              # Python test harness
│   ├── mock_server/      # Mock Anthropic API
│   ├── evaluators/      # Quality evaluation
│   └── replay/          # Session replay + trace analysis
├── docs/                 # VitePress documentation
│   ├── en/               # English docs
│   └── zh/               # Chinese docs
└── .github/workflows/   # CI/CD
```

### Go Runtime vs Python Harness

The main `go-code` binary is a pure Go implementation that:
- Communicates directly with the Anthropic API
- Executes tools locally
- Manages agent state and context

The Python harness (`harness/`) is optional and provides:
- Mock API server for testing without API costs
- Replay functionality for debugging sessions
- Integration test suite with quality evaluators

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
# Serve VitePress docs locally
make docs

# Build for production
make docs-build
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

---

如果这个项目对你有帮助，请给个 ⭐ Star！

## License

MIT License - see [LICENSE](LICENSE) for details.