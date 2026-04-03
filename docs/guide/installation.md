---
title: Installation Guide
description: Learn how to install go-code — one-command install, go install, source build, or pre-built binaries
---

# Installation Guide

This guide covers all ways to install go-code on Linux, macOS, and Windows.

## Prerequisites

- **Go 1.23+** — Only needed for building from source or `go install`
- **An API key** — Anthropic, OpenAI, or any compatible provider
- **Python 3.x** (optional) — Only for the test harness

## One-Command Install

**Linux / macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.ps1 | iex
```

## Via go install

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

This installs the binary to `$GOPATH/bin` (typically `~/go/bin`). Ensure this directory is in your `PATH`.

## Build from Source

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
make build
./bin/go-code
```

Or build manually:
```bash
go build -o bin/go-code ./cmd/go-code
```

## Pre-built Binaries

Download from [GitHub Releases](https://github.com/strings77wzq/claude-code-Go/releases):

| Platform  | Architecture    | Binary                    |
|-----------|-----------------|---------------------------|
| Linux     | amd64           | `go-code-linux-amd64`     |
| macOS     | amd64           | `go-code-darwin-amd64`    |
| macOS     | arm64 (M1/M2)   | `go-code-darwin-arm64`    |
| Windows   | amd64           | `go-code-windows-amd64.exe` |

### Linux / macOS
```bash
curl -fsSL https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64 -o go-code
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

### Windows
Download `go-code-windows-amd64.exe` from [Releases](https://github.com/strings77wzq/claude-code-Go/releases), rename to `go-code.exe`, and add to your `PATH`.

## Verify Installation

```bash
go-code --help
```

## Supported Platforms

| OS      | Architecture | Build from Source | Pre-built Binary |
|---------|-------------|-------------------|------------------|
| Linux   | amd64       | ✅                | ✅               |
| macOS   | amd64       | ✅                | ✅               |
| macOS   | arm64       | ✅                | ✅               |
| Windows | amd64       | ✅                | ✅               |

## Next Steps

- [Quick Start Guide](quick-start.md) — Run your first command
- [Configuration Guide](configuration.md) — Set up API keys and preferences
