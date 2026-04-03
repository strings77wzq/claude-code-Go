---
layout: home
title: claude-code-Go
titleTemplate: Go 语言实现的 AI 编程助手

hero:
  name: claude-code-Go
  text: 模型提供智能，Harness 提供可靠性
  tagline: 完整的 Agent Loop、工具执行、权限管理——纯 Go 打造的生产级 AI 编程助手。
  image:
    src: /logo.svg
    alt: claude-code-Go 标志
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/introduction
    - theme: alt
      text: 查看源码
      link: https://github.com/strings77wzq/claude-code-Go
---

<div class="metrics-section fade-in-section">
  <div class="metrics-grid">
    <div class="metric-item">
      <div class="metric-value">50+</div>
      <div class="metric-label">源代码文件</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">8</div>
      <div class="metric-label">模块数量</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">9</div>
      <div class="metric-label">内置工具</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">7,000+</div>
      <div class="metric-label">代码行数</div>
    </div>
  </div>
</div>

<div class="features-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">⚡</div>
    <div class="feature-title">单二进制部署</div>
    <div class="feature-desc">零依赖，一个文件跑遍全平台——Linux、macOS、Windows。无需运行时、无需虚拟环境、无需 node_modules。</div>
    <div class="feature-tags">
      <span class="tag">Linux</span>
      <span class="tag">macOS</span>
      <span class="tag">Windows</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🔒</div>
    <div class="feature-title">Harness-First 可靠性</div>
    <div class="feature-desc">权限控制、超时保护、会话持久化。Harness 保障安全，让模型专注于智能。</div>
    <div class="feature-tags">
      <span class="tag">三级权限</span>
      <span class="tag">glob 规则</span>
      <span class="tag">会话记忆</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🔌</div>
    <div class="feature-title">可扩展生态</div>
    <div class="feature-desc">MCP 协议、Hooks、Skills——通过自定义工具扩展能力，无缝集成你的工作流。</div>
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
    <div class="feature-title">模型提供智能</div>
    <div class="feature-desc">LLM 负责：理解意图、决策使用哪个工具、解释结果、规划下一步。它是系统的大脑。</div>
    <div class="feature-tags">
      <span class="tag">意图理解</span>
      <span class="tag">工具选择</span>
      <span class="tag">结果解释</span>
      <span class="tag">下一步规划</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">🛡️</div>
    <div class="feature-title">Harness 提供可靠性</div>
    <div class="feature-desc">运行时负责：权限控制、超时保护、输出截断、会话持久化、错误恢复。它是让系统达到生产级别的安全网。</div>
    <div class="feature-tags">
      <span class="tag">权限控制</span>
      <span class="tag">超时保护</span>
      <span class="tag">输出截断</span>
      <span class="tag">会话持久化</span>
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
      <div class="terminal-line terminal-output">输入 /help 查看命令，/exit 退出。</div>
      <div class="terminal-line">&nbsp;</div>
      <div class="terminal-line"><span class="terminal-prompt">go-code></span> <span class="terminal-cmd">写一个 8080 端口的 HTTP 服务器</span></div>
      <div class="terminal-line terminal-output">🔄 Agent 思考中...</div>
      <div class="terminal-line terminal-output">🛠️ 调用工具: Write → main.go</div>
      <div class="terminal-line terminal-output terminal-success">✓ 文件已写入</div>
      <div class="terminal-line terminal-output">🔄 Agent 继续思考...</div>
      <div class="terminal-line terminal-output">🛠️ 调用工具: Bash → go run main.go</div>
      <div class="terminal-line terminal-output terminal-success">✓ 服务器启动在 8080 端口</div>
      <div class="terminal-line terminal-output terminal-success">✓ 完成！HTTP 服务器已创建并运行。</div>
      <div class="terminal-line">&nbsp;</div>
      <div class="terminal-line"><span class="terminal-prompt">go-code></span> <span class="terminal-cursor"></span></div>
    </div>
  </div>
</div>

## 功能特性

| 功能 | 说明 |
|------|------|
| 🔄 Agent Loop | 基于 stop_reason 驱动的「思考→行动→观察」自主循环 |
| 🛠️ 9 大内置工具 | Read、Write、Edit、Glob、Grep、Bash、Diff、Tree、WebFetch |
| 🔒 权限系统 | 三级权限模型，支持 glob 规则匹配与会话记忆 |
| 🔌 MCP 集成 | Model Context Protocol，stdio 传输、JSON-RPC 客户端 |
| 🌊 SSE 流式 | 逐 token 实时流式响应，自研解析器，零外部依赖 |
| 🧠 上下文管理 | 智能 token 估算与自动对话压缩 |

## 为什么选择 Go？

| | Go | Python | Rust |
|---|---|---|---|
| **部署方式** | 单二进制文件 | 需要运行时环境 | 单二进制文件 |
| **外部依赖** | 零依赖 | pip install | Cargo build |
| **交叉编译** | `GOOS=linux go build` | 平台相关 | 需要交叉工具链 |
| **并发模型** | Goroutine（内置） | asyncio（库） | async/await |
| **学习曲线** | 适中 | 简单 | 陡峭 |

Go 兼具两者的优势：**像 Rust 一样的单二进制部署，接近 Python 的开发效率。**

## 快速开始

::: code-group
```bash [go install]
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

```bash [源码编译]
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
go build -o bin/go-code ./cmd/go-code
```

```bash [预编译二进制]
curl -fsSL https://github.com/strings77wzq/claude-code-Go/releases/latest/download/go-code-linux-amd64 -o go-code
chmod +x go-code
```
:::

设置 API Key 后即可使用：
```bash
export ANTHROPIC_API_KEY=sk-ant-...
./go-code
```

<div class="cta-section fade-in-section">
  <h2 class="cta-title">准备开始？</h2>
  <p class="cta-desc">几秒钟即可开始使用 claude-code-Go。</p>
  <div class="cta-actions">
    <a href="/zh/guide/introduction" class="cta-button primary">快速开始</a>
    <a href="/zh/guide/quick-start" class="cta-button secondary">安装指南</a>
    <a href="https://github.com/strings77wzq/claude-code-Go" class="cta-button secondary">GitHub</a>
  </div>
</div>