---
layout: home
title: claude-code-Go
titleTemplate: AI-powered coding assistant built in Go

hero:
  name: claude-code-Go
  text: Model provides intelligence, harness provides reliability.
  tagline: A production-grade AI coding assistant with full agent loop, tool execution, and permission management — built in pure Go.
  image:
    src: /logo.svg
    alt: claude-code-Go Logo
  actions:
    - theme: brand
      text: Get Started
      link: /guide/introduction
    - theme: alt
      text: View on GitHub
      link: https://github.com/strings77wzq/claude-code-Go

---

<div class="terminal-window">
  <div class="terminal-header">
    <span class="terminal-dot red"></span>
    <span class="terminal-dot yellow"></span>
    <span class="terminal-dot green"></span>
  </div>
  <div class="terminal-body">
    <div><span class="terminal-prompt">$</span> go-code "Hello, build me a REST API"</div>
    <div style="margin-top: 8px;">Initializing agent loop...</div>
    <div>Thinking... <span class="terminal-cursor"></span></div>
  </div>
</div>

features:
  - icon: ⚡
    title: Single Binary
    details: Zero dependencies. One file runs everywhere — Linux, macOS, Windows. No runtime, no virtualenv, no node_modules.
  - icon: 🔒
    title: Reliability First
    details: Permission system, timeout protection, session persistence. The harness ensures safety so the model can focus on intelligence.
  - icon: 🔌
    title: Extensible Ecosystem
    details: MCP protocol, Hooks, and Skills — extend with custom tools and integrate with your workflow.
  - icon: 🎯
    title: Skills System
    details: Custom commands and workflows. Define your own /review-pr, /deploy, or any reusable workflow.
  - icon: 🔄
    title: Multi-Provider
    details: Support for Anthropic, OpenAI, and any OpenAI-compatible API. Switch models without changing your workflow.
  - icon: 💾
    title: Session Resume
    details: Pick up where you left off. Load previous conversations and continue seamlessly.

stats:
  - label: Source Files
    value: 50+
  - label: Modules
    value: 8
  - label: Built-in Tools
    value: 6
  - label: Lines of Code
    value: 7,000+
---

## Architecture Philosophy

::: details 🧠 Model Provides Intelligence
The LLM handles: understanding intent, deciding which tool to use, interpreting results, and planning next steps. It's the brain of the system.
:::

::: details 🛡️ Harness Provides Reliability
The runtime handles: permission control, timeout protection, output truncation, session persistence, and error recovery. It's the safety net that makes the system production-ready.
:::

::: details 🔌 Extensible Ecosystem
MCP protocol for external tool discovery, Hooks for pre/post execution callbacks, and a clean Tool interface — just implement one interface to add new capabilities.
:::

## Feature Highlights

| Feature | Description |
|---------|-------------|
| 🔄 Agent Loop | Autonomous "think → act → observe" cycle driven by stop_reason dispatch |
| 🛠️ 6 Built-in Tools | Read, Write, Edit, Glob, Grep, Bash — complete toolset out of the box |
| 🔒 Permission System | Three-tier model with rule-based matching and session memory |
| 🔌 MCP Integration | Model Context Protocol with stdio transport and JSON-RPC |
| 🌊 SSE Streaming | Real-time token-by-token with custom parser, zero dependencies |
| 🧠 Context Management | Intelligent token estimation and automatic conversation compaction |

## Why Go?

| | Go | Python | Rust |
|---|---|---|---|
| **Deployment** | Single binary | Requires runtime | Single binary |
| **Dependencies** | Zero | pip install | Cargo build |
| **Cross-compile** | `GOOS=linux go build` | Platform-specific | Cross-toolchain needed |
| **Concurrency** | Goroutines (built-in) | asyncio (library) | async/await |
| **Learning curve** | Moderate | Easy | Steep |

Go gives you the best of both worlds: **single-binary deployment like Rust, with development speed closer to Python.**

## Quick Start

::: code-group
```bash [go install]
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

```bash [Build from source]
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
go build -o bin/go-code ./cmd/go-code
```

```bash [Pre-built binary]
curl -fsSL https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64 -o go-code
chmod +x go-code
```
:::

Then set your API key and start:
```bash
export ANTHROPIC_API_KEY=sk-ant-...
./go-code
```
