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
  <TerminalTypewriter />
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

## Choose Your Role

<div class="role-cards">
  <a href="/guide/quick-start" class="role-card">
    <div class="role-icon">👨‍💻</div>
    <div class="role-title">Full-stack Developer</div>
    <div class="role-desc">Build real applications with hands-on guides. From REPL to production, learn by doing.</div>
    <div class="role-link">Quick start guide →</div>
  </a>

  <a href="/architecture/overview" class="role-card">
    <div class="role-icon">🏗️</div>
    <div class="role-title">Architect</div>
    <div class="role-desc">Deep dive into the agent loop, tool registry, and permission system. Design extensible systems.</div>
    <div class="role-link">Architecture deep dive →</div>
  </a>

  <a href="/guide/introduction" class="role-card">
    <div class="role-icon">🎓</div>
    <div class="role-title">Student</div>
    <div class="role-desc">Learn the core principles: how AI agents reason, plan, and execute tools safely.</div>
    <div class="role-link">Learn the principles →</div>
  </a>
</div>

<style>
.role-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin: 2rem 0;
}

.role-card {
  display: block;
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  background: var(--vp-c-bg-soft);
  text-decoration: none;
  transition: all 0.2s ease;
}

.role-card:hover {
  border-color: var(--vp-c-brand);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.role-icon {
  font-size: 2.5rem;
  margin-bottom: 0.75rem;
}

.role-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  margin-bottom: 0.5rem;
}

.role-desc {
  font-size: 0.9rem;
  color: var(--vp-c-text-2);
  line-height: 1.5;
  margin-bottom: 1rem;
}

.role-link {
  font-size: 0.85rem;
  color: var(--vp-c-brand);
  font-weight: 500;
}
</style>

## Why Go?

<div class="comparison-table">

| Feature | Go | Python | Rust | TypeScript |
|---------|-----|--------|------|------------|
| **Single binary** | ✅ | ❌ | ✅ | ❌ |
| **Zero runtime deps** | ✅ | ❌ | ✅ | ❌ |
| **Concurrency** | ✅ Goroutines | ❌ asyncio | ✅ async/await | ✅ event loop |
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

Then set your API key and start:
```bash
export ANTHROPIC_API_KEY=sk-ant-...
./go-code
```
