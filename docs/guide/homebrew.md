---
title: Homebrew Installation
description: Install go-code via Homebrew tap on macOS and Linux
---

# Homebrew Installation

You can install `go-code` using [Homebrew](https://brew.sh) on macOS and Linux.

## Prerequisites

- **Homebrew** — Install from [brew.sh](https://brew.sh)
- **API Key** — Anthropic, OpenAI, or any compatible provider

## Installation

### Add the Tap

```bash
brew tap strings77wzq/claude-code-go
```

### Install go-code

```bash
brew install claude-code-go
```

### Verify Installation

```bash
go-code --version
```

## Upgrading

To upgrade to the latest version:

```bash
brew update
brew upgrade claude-code-go
```

## Uninstalling

To remove go-code:

```bash
brew uninstall claude-code-go
brew untap strings77wzq/claude-code-go
```

## Formula Example

If you need to manually work with the formula, here's the structure:

```ruby
class ClaudeCodeGo < Formula
  desc "AI Coding Agent in Go"
  homepage "https://github.com/strings77wzq/claude-code-Go"
  url "https://github.com/strings77wzq/claude-code-Go/releases/download/v1.0.0/go-code_1.0.0_darwin_arm64.tar.gz"
  sha256 "..."
  license "MIT"
  version "1.0.0"

  arch arm64: "arm64", amd64: "amd64"

  def install
    bin.install "go-code"
  end

  test do
    system "#{bin}/go-code", "--version"
  end
end
```

## Next Steps

After installation, configure your API key:

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Or create `~/.go-code/settings.json`:

```json
{
  "apiKey": "sk-ant-...",
  "model": "claude-sonnet-4-20250514"
}
```

For more details, see the [Installation Guide](./installation.md) or [Quick Start](./quick-start.md).