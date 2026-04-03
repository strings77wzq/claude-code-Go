---
title: Installation Guide
description: Learn how to install go-code from source or pre-built binaries
---

# Installation Guide

This guide covers how to install go-code, including building from source and using pre-built binaries.

## Prerequisites

Before installing go-code, ensure you have:

- **Go 1.23 or later** — Required for building from source
- **An Anthropic API key** — Required for running the application
- **Python 3.x** (optional) — Only needed if using the test harness

## Install via go install

The fastest way to install go-code is using Go's built-in install command:

```bash
go install github.com/strings77wzq/claude-code-go/cmd/go-code@latest
```

This installs the binary to `$GOPATH/bin` (typically `~/go/bin`). Ensure this directory is in your PATH.

## Build from Source

### Clone the Repository

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
```

### Build the Binary

Using the Makefile:

```bash
make build
```

This creates the binary at `bin/go-code`.

Alternatively, build manually:

```bash
go build -o bin/go-code ./cmd/go-code
```

### Install to System PATH

If you prefer to install globally:

```bash
go install ./cmd/go-code
```

This installs the binary to `$GOPATH/bin`.

## Pre-built Binaries

For platforms without Go installed, pre-built binaries are available on the GitHub releases page:

| Platform | Architecture | Filename |
|----------|--------------|----------|
| Linux    | amd64        | go-code-linux-amd64 |
| macOS    | amd64        | go-code-darwin-amd64 |
| macOS    | arm64        | go-code-darwin-arm64 |

Download the appropriate binary and make it executable:

```bash
# Example: Download and install Linux amd64
curl -L -o go-code https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

## Verify Installation

After installation, verify the binary works:

```bash
# Check the binary exists
which go-code

# Or if using bin/ directory
ls -la bin/go-code

# Run the help command (if supported)
go-code --help
```

You should see output indicating the version and available options.

## Next Steps

- [Quick Start Guide](quick-start.md) — Run your first command
- [Configuration Guide](configuration.md) — Set up API keys and preferences