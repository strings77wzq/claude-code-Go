# go-code — Go implementation of Claude Code

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue)](https://go.dev/doc/install)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/strings77wzq/claude-code-Go/releases)

A Go implementation of Anthropic's Claude Code agent system. Full agent loop with built-in tools, permission system, MCP support, and SSE streaming — in a single binary.

## Installation

### One-Command Install

**Linux / macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.ps1 | iex
```

### Via go install

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

The binary installs to `$GOPATH/bin` (typically `~/go/bin`). Make sure it's in your `PATH`.

### Build from Source

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
make build
./bin/go-code
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/strings77wzq/claude-code-Go/releases):

| Platform  | Architecture    | Binary                    |
|-----------|-----------------|---------------------------|
| Linux     | amd64           | `go-code-linux-amd64`     |
| macOS     | amd64           | `go-code-darwin-amd64`    |
| macOS     | arm64 (M1/M2)   | `go-code-darwin-arm64`    |
| Windows   | amd64           | `go-code-windows-amd64.exe` |

```bash
# Example: Linux amd64
curl -fsSL https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64 -o go-code
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

## Quick Start

### For AI Agents

If you're using an AI coding assistant, paste this URL into its session:
[Installation Guide for Agents](https://github.com/strings77wzq/claude-code-Go/blob/main/docs/guide/installation-for-agents.md)

### 1. Set your API key

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Or create `~/.go-code/settings.json`:
```json
{
  "apiKey": "sk-ant-..."
}
```

### 2. Run

```bash
# Interactive REPL
go-code

# Or with a direct prompt
go-code "Explain the agent loop architecture"
```

## Features

- **Agent Loop**: Full "think → act → observe" cycle with message history and context compaction
- **6 Built-in Tools**: `Read`, `Write`, `Edit`, `Glob`, `Grep`, `Bash`
- **Permission System**: 3-tier model (ReadOnly / WorkspaceWrite / DangerFullAccess)
- **MCP Support**: Model Context Protocol with stdio transport
- **SSE Streaming**: Real-time token-by-token output
- **Session Persistence**: Auto-save and resume conversations
- **Skills System**: Custom commands and reusable workflows
- **Multi-Provider**: Anthropic, OpenAI, and any OpenAI-compatible API

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        go-code                              │
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
```

### Go Runtime vs Python Harness

The main `go-code` binary is a pure Go implementation:
- Direct Anthropic API communication
- Local tool execution
- Agent state and context management

The Python harness (`harness/`) is optional:
- Mock API server for testing without API costs
- Session replay for debugging
- Quality evaluators for integration testing

## Project Structure

```
claude-code-Go/
├── cmd/go-code/          # Main entry point
├── internal/
│   ├── agent/            # Agent loop + context management
│   ├── api/              # Anthropic API client + SSE
│   ├── config/           # Multi-source config loader
│   ├── permission/       # 3-tier permission system
│   ├── tool/             # Tool interface + builtins
│   │   ├── builtin/      # Bash, Read, Write, Edit, Glob, Grep
│   │   ├── mcp/          # MCP integration
│   │   └── init/         # Tool registration
│   └── hooks/            # Pre/post execution hooks
├── pkg/tty/              # REPL + terminal rendering
├── harness/              # Python test harness (optional)
├── docs/                 # Documentation
│   └── guide/            # User guides
└── .github/workflows/   # CI/CD
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test          # Go + Python tests
go test -v ./...   # Go tests only
```

### Cross-Platform Build
```bash
make build-all     # Linux amd64, macOS amd64/arm64, Windows amd64
```

Output binaries are placed in `bin/`.

### Clean
```bash
make clean
```

### Documentation
```bash
make docs           # Serve docs locally
make docs-build     # Build for production
```

---

⭐ If you find this project helpful, please give it a ⭐ Star!

## License

MIT License — see [LICENSE](LICENSE) for details.
