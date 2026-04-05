---
title: Homebrew 安装
description: 通过 Homebrew tap 在 macOS 和 Linux 上安装 go-code
---

# Homebrew 安装

您可以使用 [Homebrew](https://brew.sh) 在 macOS 和 Linux 上安装 `go-code`。

## 前置条件

- **Homebrew** — 从 [brew.sh](https://brew.sh) 安装
- **API 密钥** — Anthropic、OpenAI 或任何兼容的提供商

## 安装

### 添加 Tap

```bash
brew tap strings77wzq/claude-code-go
```

### 安装 go-code

```bash
brew install claude-code-go
```

### 验证安装

```bash
go-code --version
```

## 升级

要升级到最新版本：

```bash
brew update
brew upgrade claude-code-go
```

## 卸载

要删除 go-code：

```bash
brew uninstall claude-code-go
brew untap strings77wzq/claude-code-go
```

## Formula 示例

如果您需要手动处理 formula，以下是结构：

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

## 下一步

安装后，配置您的 API 密钥：

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

或创建 `~/.go-code/settings.json`：

```json
{
  "apiKey": "sk-ant-...",
  "model": "claude-sonnet-4-20250514"
}
```

更多详情，请参阅[安装指南](./installation.md)或[快速开始](./quick-start.md)。