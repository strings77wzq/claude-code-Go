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

<div class="metrics-section fade-in-section">
  <div class="metrics-grid">
    <div class="metric-item">
      <div class="metric-value">50+</div>
      <div class="metric-label">Source Files</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">8</div>
      <div class="metric-label">Modules</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">9</div>
      <div class="metric-label">Built-in Tools</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">7,000+</div>
      <div class="metric-label">Lines of Code</div>
    </div>
  </div>
</div>

<div class="features-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">⚡</div>
    <div class="feature-title">Single Binary Deployment</div>
    <div class="feature-desc">Zero dependencies. One file runs everywhere — Linux, macOS, Windows. No runtime, no virtualenv, no node_modules.</div>
    <div class="feature-tags">
      <span class="tag">Linux</span>
      <span class="tag">macOS</span>
      <span class="tag">Windows</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🔒</div>
    <div class="feature-title">Harness-First Reliability</div>
    <div class="feature-desc">Permission control, timeout protection, session persistence. The harness ensures safety so the model can focus on intelligence.</div>
    <div class="feature-tags">
      <span class="tag">3-Tier Permissions</span>
      <span class="tag">Glob Rules</span>
      <span class="tag">Session Memory</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🔌</div>
    <div class="feature-title">Extensible Ecosystem</div>
    <div class="feature-desc">MCP protocol, Hooks, and Skills — extend with custom tools and integrate with your workflow.</div>
    <div class="feature-tags">
      <span class="tag">MCP</span>
      <span class="tag">Hooks</span>
      <span class="tag">Skills</span>
    </div>
  </div>
</div>

<div class="architecture-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">🧠</div>
    <div class="feature-title">Model Provides Intelligence</div>
    <div class="feature-desc">The LLM handles: understanding intent, deciding which tool to use, interpreting results, and planning next steps. It's the brain of the system.</div>
    <div class="feature-tags">
      <span class="tag">Intent Understanding</span>
      <span class="tag">Tool Selection</span>
      <span class="tag">Result Interpretation</span>
      <span class="tag">Next-Step Planning</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🛡️</div>
    <div class="feature-title">Harness Provides Reliability</div>
    <div class="feature-desc">The runtime handles: permission control, timeout protection, output truncation, session persistence, and error recovery. It's the safety net that makes the system production-ready.</div>
    <div class="feature-tags">
      <span class="tag">Permission Control</span>
      <span class="tag">Timeout Protection</span>
      <span class="tag">Output Truncation</span>
      <span class="tag">Session Persistence</span>
    </div>
  </div>
</div>

<div class="terminal-section fade-in-section">
  <div class="terminal-window">
    <div class="terminal-header">
      <span class="terminal-dot red"></span>
      <span class="terminal-dot yellow"></span>
      <span class="terminal-dot green"></span>
      <span class="terminal-title">claude-code-Go</span>
    </div>
    <div class="terminal-body">
      <div class="terminal-line"><span class="terminal-prompt">$</span> <span class="terminal-cmd">go-code</span></div>
      <div class="terminal-line terminal-output">claude-code-Go v0.1.0</div>
      <div class="terminal-line terminal-output">Type /help for commands, /exit to quit.</div>
      <div class="terminal-line">&nbsp;</div>
      <div class="terminal-line"><span class="terminal-prompt">go-code></span> <span class="terminal-cmd">write an HTTP server on port 8080</span></div>
      <div class="terminal-line terminal-output">🔄 Agent thinking...</div>
      <div class="terminal-line terminal-output">🛠️ Tool call: Write → main.go</div>
      <div class="terminal-line terminal-output terminal-success">✓ File written</div>
      <div class="terminal-line terminal-output">🔄 Agent continuing...</div>
      <div class="terminal-line terminal-output">🛠️ Tool call: Bash → go run main.go</div>
      <div class="terminal-line terminal-output terminal-success">✓ Server started on port 8080</div>
      <div class="terminal-line terminal-output terminal-success">✓ Done! HTTP server created and running.</div>
      <div class="terminal-line">&nbsp;</div>
      <div class="terminal-line"><span class="terminal-prompt">go-code></span> <span class="terminal-cursor"></span></div>
    </div>
  </div>
</div>

## Feature Highlights

| Feature | Description |
|---------|-------------|
| 🔄 Agent Loop | Autonomous "think → act → observe" cycle driven by stop_reason dispatch |
| 🛠️ 9 Built-in Tools | Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch |
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

<div class="cta-section fade-in-section">
  <h2 class="cta-title">Ready to start?</h2>
  <p class="cta-desc">Get started with claude-code-Go in seconds.</p>
  <div class="cta-actions">
    <a href="/guide/introduction" class="cta-button primary">Get Started</a>
    <a href="/guide/quick-start" class="cta-button secondary">Quick Start</a>
    <a href="https://github.com/strings77wzq/claude-code-Go" class="cta-button secondary">GitHub</a>
  </div>
</div>
