# claude-code-Go — AI Coding Agent in Go

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/doc/install)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/strings77wzq/claude-code-Go/releases)
[![CI](https://github.com/strings77wzq/claude-code-Go/actions/workflows/ci.yml/badge.svg)](https://github.com/strings77wzq/claude-code-Go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/strings77wzq/claude-code-Go)](https://goreportcard.com/report/github.com/strings77wzq/claude-code-Go)

A Go-native AI coding agent with full agent loop, tool execution, permission management, SSE streaming, and auto-recovery — in a single binary.

> **Status: v0.3 verified release** — Core agent runtime, permission system, session persistence, multi-provider support, doctor checks, replay evidence, MCP/LSP extension diagnostics, and manifest-driven harness quality gates are implemented and tested. Real external MCP/LSP server smoke checks and competitor-agent comparisons remain manual. See [PARITY.md](PARITY.md) for detailed feature status.

> **Disclaimer:** This is an independent open-source project. It is not affiliated with, endorsed by, or connected to Anthropic PBC. "Claude" and "Claude Code" are trademarks of Anthropic PBC.

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

For script-based install, see [scripts/](scripts/) for `install.sh` and `install.ps1`.

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

### 2. Verify your setup

Run the health check before starting a real session:

```bash
go-code doctor
```

For offline environments:

```bash
go-code doctor --offline
```

### 3. Run

```bash
# Interactive TUI
go-code

# Legacy REPL
go-code --legacy-repl

# Single prompt (non-interactive)
go-code -p "Explain the agent loop architecture"

# JSON output (for scripting)
go-code -p "List files in current directory" -f json

# Replay a saved session
go-code replay latest

# Collect concise release/issue evidence
go-code replay --evidence latest
```

## Verified Features (v0.3)

These features are tested and covered by the parity harness. See [PARITY.md](PARITY.md) for evidence links.

| Feature | Status | Tests |
|---------|--------|-------|
| Agent Loop (think → act → observe) | Verified | Go unit + harness |
| 11 Built-in Tools | Verified | Go tests |
| Permission System (3-tier) | Verified | Go + harness |
| Doctor Health Check | Verified | Go tests |
| Multi-Provider (Anthropic, OpenAI-compatible) | Verified | Go tests |
| Session Persistence + Resume | Verified | Go tests |
| Session Replay + Evidence Mode | Partial v0.3 | Go tests + harness |
| Slash Commands (/help, /model, /sessions, etc.) | Verified | Go tests |
| Skills System | Verified | Go tests |
| Hooks System | Verified | Go tests |
| MCP Integration | Partial v0.3 | Go tests + harness |
| LSP Integration | Partial v0.3 | Go tests + harness |

### Built-in Tools

`Read`, `Write`, `Edit`, `Glob`, `Grep`, `Bash`, `Diff`, `Tree`, `WebFetch`, `TodoWrite`, `NotebookEdit`

## Supported Providers

| Provider | Setup |
|----------|-------|
| **Anthropic** | `ANTHROPIC_API_KEY=...`, optional `LLM_PROVIDER=anthropic` |
| **OpenAI** | `LLM_PROVIDER=openai`, `ANTHROPIC_API_KEY=...`, `ANTHROPIC_BASE_URL=https://api.openai.com` |
| **DeepSeek** | `LLM_PROVIDER=openai`, `ANTHROPIC_BASE_URL=https://api.deepseek.com`, model `deepseek-chat` or `deepseek-reasoner` |
| **Qwen** | `LLM_PROVIDER=openai`, `ANTHROPIC_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode` |
| **GLM** | `LLM_PROVIDER=openai`, `ANTHROPIC_BASE_URL=https://open.bigmodel.cn/api/paas` |
| **Tencent Cloud** | `LLM_PROVIDER=anthropic`, `ANTHROPIC_BASE_URL=https://api.lkeap.cloud.tencent.com/coding/anthropic`, model `tc-code-latest` |

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
└─────────────────────────────────────────────────────────────┘
```

For a detailed architecture overview, see [docs/architecture/](docs/architecture/).

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
│   │   ├── builtin/      # 11 built-in tools
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

```bash
make build          # Build binary
make test           # Go + Python harness tests
make vet            # Static analysis
make build-all      # Cross-compile all platforms
make docs           # Serve docs locally
make docs-build     # Build docs for production
```

## Documentation

Full documentation: [https://strings77wzq.github.io/claude-code-Go/](https://strings77wzq.github.io/claude-code-Go/)

Key pages:
- [Quick Start](docs/guide/quick-start.md)
- [Architecture Overview](docs/architecture/overview.md)
- [Roadmap](docs/roadmap.md)
- [Troubleshooting](docs/troubleshooting/)
- [MCP Integration](docs/extension/mcp.md)
- [LSP Integration](docs/extension/lsp.md)
- [Contributing](CONTRIBUTING.md)
- [Parity Status](PARITY.md)

## License

MIT License — see [LICENSE](LICENSE) for details.
