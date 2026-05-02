---
layout: home
title: claude-code-Go
titleTemplate: AI-powered coding assistant built in pure Go

hero:
  name: claude-code-Go
  text: Intelligence meets reliability.
  tagline: A production-grade AI coding assistant with full agent loop, permission system, and multi-provider support — built in pure Go. Single binary. Zero runtime dependencies.
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
      <div class="metric-value">26</div>
      <div class="metric-label">Go Packages</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">11</div>
      <div class="metric-label">Built-in Tools</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">3</div>
      <div class="metric-label">Permission Tiers</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">100%</div>
      <div class="metric-label">Test Coverage</div>
    </div>
  </div>
</div>

<div class="code-preview-section fade-in-section">
  <CodePreview />
</div>

<div class="features-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">&#x26A1;</div>
    <div class="feature-title">Single Binary Deployment</div>
    <div class="feature-desc">One `go install` and you're done. No runtime, no virtual environment, no node_modules. Runs everywhere Go compiles — Linux, macOS, Windows.</div>
    <div class="feature-tags">
      <span class="tag">Linux</span>
      <span class="tag">macOS</span>
      <span class="tag">Windows</span>
      <span class="tag">go install</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F512;</div>
    <div class="feature-title">Harness-First Reliability</div>
    <div class="feature-desc">Three-tier permission system with rule-based matching and session memory. Every tool call is validated before execution. The harness ensures safety so the model can focus on intelligence.</div>
    <div class="feature-tags">
      <span class="tag">3-Tier Permissions</span>
      <span class="tag">Glob Rules</span>
      <span class="tag">Session Memory</span>
      <span class="tag">Path Validation</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F50C;</div>
    <div class="feature-title">Multi-Provider Ecosystem</div>
    <div class="feature-desc">Native support for Anthropic, OpenAI, and OpenAI-compatible providers (DeepSeek, Qwen, GLM). Extensible via MCP protocol, Hooks, and Skills — integrate with your workflow.</div>
    <div class="feature-tags">
      <span class="tag">Anthropic</span>
      <span class="tag">OpenAI</span>
      <span class="tag">DeepSeek</span>
      <span class="tag">MCP</span>
      <span class="tag">Skills</span>
    </div>
  </div>
</div>

<div class="architecture-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">&#x1F9E0;</div>
    <div class="feature-title">Model Provides Intelligence</div>
    <div class="feature-desc">The LLM handles: understanding intent, deciding which tool to use, interpreting results, and planning next steps. It's the brain of the system — reasoning, adapting, creating.</div>
    <div class="feature-tags">
      <span class="tag">Intent Understanding</span>
      <span class="tag">Tool Selection</span>
      <span class="tag">Result Interpretation</span>
      <span class="tag">Next-Step Planning</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F6E1;</div>
    <div class="feature-title">Harness Provides Reliability</div>
    <div class="feature-desc">The runtime handles: permission control, timeout protection, output truncation, session persistence, and error recovery. It's the safety net that makes the system production-ready.</div>
    <div class="feature-tags">
      <span class="tag">Permission Control</span>
      <span class="tag">Timeout Protection</span>
      <span class="tag">Session Persistence</span>
      <span class="tag">Error Recovery</span>
    </div>
  </div>
</div>

<div class="terminal-section fade-in-section">
  <TerminalTypewriter />
</div>

## Feature Highlights

| Feature | Description |
|---------|-------------|
| &#x1F504; Agent Loop | Autonomous "think &rarr; act &rarr; observe" cycle driven by `stop_reason` dispatch |
| &#x1F6E0; 11 Built-in Tools | Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch, WebSearch, TodoWrite |
| &#x1F512; Permission System | Three-tier model with rule-based matching, glob patterns, and session memory |
| &#x1F50C; MCP Integration | Model Context Protocol with stdio transport, JSON-RPC, and permission gates |
| &#x1F30A; SSE Streaming | Real-time token-by-token output with custom parser, zero external dependencies |
| &#x1F9E0; Context Management | Intelligent token estimation and automatic conversation compaction |

<div class="playground-section fade-in-section">
  <Playground />
</div>

<div class="quick-start-section fade-in-section">
  <div class="quick-start-header">
    <h2 class="quick-start-title">5 Minutes to Get Started</h2>
    <p class="quick-start-subtitle">Three simple steps to start coding with AI assistance</p>
  </div>
  <div class="quick-start-steps">
    <div class="quick-start-step">
      <div class="step-number">1</div>
      <div class="step-title">Install</div>
      <div class="step-desc">One command: <code>go install</code>. Zero dependencies, zero configuration needed.</div>
    </div>
    <div class="quick-start-step">
      <div class="step-number">2</div>
      <div class="step-title">Configure</div>
      <div class="step-desc">Set your API key. Anthropic, OpenAI, DeepSeek, and more — choose your provider.</div>
    </div>
    <div class="quick-start-step">
      <div class="step-number">3</div>
      <div class="step-title">Start Building</div>
      <div class="step-desc">Launch <code>go-code</code> and start building. The agent understands your codebase and gets to work.</div>
    </div>
  </div>
</div>

## Choose Your Path

<div class="role-cards">
  <a href="/guide/quick-start" class="role-card">
    <div class="role-icon">&#x1F468;&#x200D;&#x1F4BB;</div>
    <div class="role-title">Full-stack Developer</div>
    <div class="role-desc">Build real applications with hands-on guides. From your first tool call to production-ready code, learn by doing.</div>
    <div class="role-link">Quick start guide &rarr;</div>
  </a>

  <a href="/architecture/overview" class="role-card">
    <div class="role-icon">&#x1F3D7;</div>
    <div class="role-title">Architect</div>
    <div class="role-desc">Deep dive into the agent loop, tool registry, and permission system. Design extensible, safe AI systems.</div>
    <div class="role-link">Architecture deep dive &rarr;</div>
  </a>

  <a href="/guide/introduction" class="role-card">
    <div class="role-icon">&#x1F393;</div>
    <div class="role-title">Learner</div>
    <div class="role-desc">Understand the core principles: how AI agents reason, plan, and execute tools safely. Start from zero.</div>
    <div class="role-link">Learn the principles &rarr;</div>
  </a>
</div>

## Why Go?

<div class="comparison-table">

| Feature | Go | Python | Rust | TypeScript |
|---------|-----|--------|------|------------|
| **Single binary** | &#x2705; | &#x274C; | &#x2705; | &#x274C; |
| **Zero runtime deps** | &#x2705; | &#x274C; | &#x2705; | &#x274C; |
| **Built-in concurrency** | &#x2705; Goroutines | &#x274C; asyncio | &#x2705; async/await | &#x2705; event loop |
| **Deployment** | `go install` | `pip install` | `cargo build` | `npm install` |
| **Learning curve** | Moderate | Easy | Steep | Moderate |

</div>

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

Set your API key and start coding:
```bash
export GO_CODE_API_KEY=sk-ant-...
./go-code
```
