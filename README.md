# claude-code-Go — AI Coding Agent in Go

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/doc/install)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/strings77wzq/claude-code-Go/releases)
[![Stars](https://img.shields.io/github/stars/strings77wzq/claude-code-Go?style=social)](https://github.com/strings77wzq/claude-code-Go/stargazers)
[![CI](https://github.com/strings77wzq/claude-code-Go/actions/workflows/ci.yml/badge.svg)](https://github.com/strings77wzq/claude-code-Go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/strings77wzq/claude-code-Go)](https://goreportcard.com/report/github.com/strings77wzq/claude-code-Go)

> **Model provides intelligence, Harness provides reliability.**

A production-grade AI coding assistant with full agent loop, tool execution, permission management, SSE streaming, LSP integration, and auto-recovery — in a single Go binary.

## Installation

### Via go install (Recommended)

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

The binary installs to `$GOPATH/bin` (typically `~/go/bin`). Make sure it's in your `PATH`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

### Build from Source

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
make install
go-code
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/strings77wzq/claude-code-Go/releases):

```bash
# Example: Linux amd64
curl -fsSL https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64 -o go-code
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

## Quick Start

### 1. Set your API key

```bash
# Anthropic
export ANTHROPIC_API_KEY=sk-ant-...

# Or Tencent Cloud Coding Plan
export ANTHROPIC_API_KEY=sk-sp-...
export ANTHROPIC_BASE_URL=https://api.lkeap.cloud.tencent.com/coding/anthropic
export ANTHROPIC_MODEL=tc-code-latest
```

Or create `~/.go-code/settings.json`:
```json
{
  "apiKey": "sk-ant-...",
  "model": "claude-sonnet-4-6-20251001"
}
```

### 2. Run

```bash
# Interactive REPL
go-code

# Single prompt (non-interactive)
go-code -p "Explain the agent loop architecture"

# JSON output (for scripting)
go-code -p "List files in current directory" -f json

# Quiet mode (no spinner)
go-code -p "What is 2+2?" -q
```

## Features

- **Agent Loop**: Full "think → act → observe" cycle with stop_reason-driven state machine
- **10 Built-in Tools**: `Read`, `Write`, `Edit`, `Glob`, `Grep`, `Bash`, `Diff`, `Tree`, `WebFetch`, `TodoWrite`
- **Permission System**: 3-tier model (ReadOnly / WorkspaceWrite / DangerFullAccess) with glob rules and session memory
- **MCP Support**: Model Context Protocol with stdio transport
- **SSE Streaming**: Real-time token-by-token output with custom parser
- **Session Persistence**: Auto-save and resume conversations (JSONL format)
- **Skills System**: Custom commands and reusable workflows
- **Multi-Provider**: Anthropic, OpenAI, and any OpenAI-compatible API (DeepSeek, Qwen, GLM)
- **LSP Integration**: Language Server Protocol for code symbols, references, diagnostics
- **Auto-Recovery**: Automatic retry on API timeout, rate limit, and context full
- **Runtime Model Switching**: Change models mid-session with `/model <name>`
- **Auto-Update**: Check and download latest version with `/update`

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    claude-code-Go                            │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────┐  │
│  │  Bubbletea   │───▶│ Agent Loop   │───▶│ Tool Registry │  │
│  │     TUI      │    │   + Context  │    │               │  │
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

### Design Philosophy

**Model provides intelligence, Harness provides reliability.**

| Layer | Responsibility |
|-------|---------------|
| **Model (LLM)** | Intent understanding, tool selection, result interpretation, next-step planning |
| **Harness (Runtime)** | Permission control, timeout protection, output truncation, session persistence, error recovery |

### Go Runtime vs Python Harness

The main `go-code` binary is a pure Go implementation:
- Direct API communication with SSE streaming
- Local tool execution with safety guards
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
│   ├── config/           # Multi-source config loader
│   ├── cost/             # Cost tracking and estimation
│   ├── hooks/            # Pre/post execution hooks
│   ├── lsp/              # Language Server Protocol client
│   ├── permission/       # 3-tier permission system
│   ├── provider/         # Multi-provider abstraction
│   │   ├── anthropic/    # Anthropic API provider
│   │   ├── openai/       # OpenAI-compatible provider
│   │   └── registry/     # Provider auto-selection
│   ├── session/          # Session persistence + resume
│   ├── skills/           # Custom skills system
│   ├── tool/             # Tool interface + builtins
│   │   ├── builtin/      # 10 built-in tools
│   │   ├── mcp/          # MCP integration
│   │   └── init/         # Tool registration
│   └── update/           # Auto-update checker
├── pkg/
│   ├── tty/              # Legacy REPL (use --legacy-repl)
│   └── tui/              # Bubbletea TUI (default)
├── harness/              # Python test harness
├── docs/                 # VitePress documentation
├── scripts/              # Install and launch scripts
└── .github/workflows/    # CI/CD
```

## Development

### Build
```bash
make build
```

### Install
```bash
make install
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

### Documentation
```bash
make docs           # Serve docs locally
make docs-build     # Build for production
```

## Supported Providers

| Provider | Setup |
|----------|-------|
| **Anthropic** | Set `ANTHROPIC_API_KEY` |
| **OpenAI** | Set `ANTHROPIC_API_KEY` + `ANTHROPIC_BASE_URL=https://api.openai.com/v1` |
| **DeepSeek** | Set `ANTHROPIC_BASE_URL=https://api.deepseek.com` |
| **Qwen** | Set `ANTHROPIC_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1` |
| **GLM** | Set `ANTHROPIC_BASE_URL=https://open.bigmodel.cn/api/paas/v4` |
| **Tencent Cloud** | Set `ANTHROPIC_BASE_URL=https://api.lkeap.cloud.tencent.com/coding/anthropic` |

## Documentation

📖 Full documentation: [https://strings77wzq.github.io/claude-code-Go/](https://strings77wzq.github.io/claude-code-Go/)

---

⭐ If you find this project helpful, please give it a ⭐ Star!

## License

MIT License — see [LICENSE](LICENSE) for details.
