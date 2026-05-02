---
layout: home
title: claude-code-Go
titleTemplate: 纯 Go 打造的 AI 编程助手

hero:
  name: claude-code-Go
  text: 智能与可靠的结合
  tagline: 完整的 Agent Loop、权限系统、多供应商支持 — 纯 Go 打造的生产级 AI 编程助手。单二进制文件，零运行时依赖。
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
      <div class="metric-value">26</div>
      <div class="metric-label">Go 包</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">11</div>
      <div class="metric-label">内置工具</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">3</div>
      <div class="metric-label">权限层级</div>
    </div>
    <div class="metric-item">
      <div class="metric-value">100%</div>
      <div class="metric-label">测试覆盖</div>
    </div>
  </div>
</div>

<div class="code-preview-section fade-in-section">
  <CodePreview />
</div>

<div class="features-section fade-in-section">
  <div class="feature-block">
    <div class="feature-icon">&#x26A1;</div>
    <div class="feature-title">单二进制部署</div>
    <div class="feature-desc">一条 `go install` 命令即可完成。无需运行时、无需虚拟环境、无需 node_modules。Go 编译的所有平台均可运行 — Linux、macOS、Windows。</div>
    <div class="feature-tags">
      <span class="tag">Linux</span>
      <span class="tag">macOS</span>
      <span class="tag">Windows</span>
      <span class="tag">go install</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F512;</div>
    <div class="feature-title">Harness-First 可靠性</div>
    <div class="feature-desc">三级权限系统，支持基于规则的匹配和会话记忆。每个工具调用在执行前都经过验证。Harness 保障安全，让模型专注于智能处理。</div>
    <div class="feature-tags">
      <span class="tag">三级权限</span>
      <span class="tag">Glob 规则</span>
      <span class="tag">会话记忆</span>
      <span class="tag">路径验证</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F50C;</div>
    <div class="feature-title">多供应商生态</div>
    <div class="feature-desc">原生支持 Anthropic、OpenAI 和 OpenAI 兼容供应商（DeepSeek、Qwen、GLM）。通过 MCP 协议、Hooks、Skills 扩展 — 与你的工作流无缝集成。</div>
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
    <div class="feature-title">模型提供智能</div>
    <div class="feature-desc">LLM 负责：理解意图、选择工具、解释结果、规划下一步。它是系统的大脑 — 推理、适应、创造。</div>
    <div class="feature-tags">
      <span class="tag">意图理解</span>
      <span class="tag">工具选择</span>
      <span class="tag">结果解释</span>
      <span class="tag">下一步规划</span>
    </div>
  </div>

  <div class="feature-block">
    <div class="feature-icon">&#x1F6E1;</div>
    <div class="feature-title">Harness 提供可靠性</div>
    <div class="feature-desc">运行时负责：权限控制、超时保护、输出截断、会话持久化、错误恢复。它是让系统达到生产级别的安全网。</div>
    <div class="feature-tags">
      <span class="tag">权限控制</span>
      <span class="tag">超时保护</span>
      <span class="tag">会话持久化</span>
      <span class="tag">错误恢复</span>
    </div>
  </div>
</div>

<div class="terminal-section fade-in-section">
  <TerminalTypewriter />
</div>

## 功能特性

| 功能 | 说明 |
|------|------|
| &#x1F504; Agent Loop | 由 `stop_reason` 驱动的「思考 &rarr; 行动 &rarr; 观察」自主循环 |
| &#x1F6E0; 11 大内置工具 | Read、Write、Edit、Glob、Grep、Bash、Diff、Tree、WebFetch、WebSearch、TodoWrite |
| &#x1F512; 权限系统 | 三级权限模型，支持 glob 规则匹配与会话记忆 |
| &#x1F50C; MCP 集成 | Model Context Protocol，stdio 传输、JSON-RPC、权限门控 |
| &#x1F30A; SSE 流式 | 逐 token 实时流式输出，自研解析器，零外部依赖 |
| &#x1F9E0; 上下文管理 | 智能 token 用量估算与自动对话压缩 |

<div class="playground-section fade-in-section">
  <Playground />
</div>

<div class="quick-start-section fade-in-section">
  <div class="quick-start-header">
    <h2 class="quick-start-title">5 分钟快速上手</h2>
    <p class="quick-start-subtitle">三步开始 AI 辅助编程</p>
  </div>
  <div class="quick-start-steps">
    <div class="quick-start-step">
      <div class="step-number">1</div>
      <div class="step-title">安装</div>
      <div class="step-desc">一条命令：<code>go install</code>。零依赖，无需任何配置。</div>
    </div>
    <div class="quick-start-step">
      <div class="step-number">2</div>
      <div class="step-title">配置</div>
      <div class="step-desc">设置 API Key。支持 Anthropic、OpenAI、DeepSeek 等多种供应商。</div>
    </div>
    <div class="quick-start-step">
      <div class="step-number">3</div>
      <div class="step-title">开始构建</div>
      <div class="step-desc">启动 <code>go-code</code> 开始构建。Agent 理解你的代码库并开始工作。</div>
    </div>
  </div>
</div>

## 选择你的路径

<div class="role-cards">
  <a href="/zh/guide/quick-start" class="role-card">
    <div class="role-icon">&#x1F468;&#x200D;&#x1F4BB;</div>
    <div class="role-title">全栈开发者</div>
    <div class="role-desc">通过动手实践构建真实应用。从第一个工具调用到生产就绪代码，边做边学。</div>
    <div class="role-link">快速开始指南 &rarr;</div>
  </a>

  <a href="/zh/architecture/overview" class="role-card">
    <div class="role-icon">&#x1F3D7;</div>
    <div class="role-title">架构师</div>
    <div class="role-desc">深入 Agent Loop、工具注册表、权限系统。设计可扩展、安全的 AI 系统。</div>
    <div class="role-link">架构深度探索 &rarr;</div>
  </a>

  <a href="/zh/guide/introduction" class="role-card">
    <div class="role-icon">&#x1F393;</div>
    <div class="role-title">学习者</div>
    <div class="role-desc">从零开始理解核心原理：AI Agent 如何推理、规划、安全执行工具。</div>
    <div class="role-link">学习核心原理 &rarr;</div>
  </a>
</div>

## 为什么选择 Go？

<div class="comparison-table">

| 特性 | Go | Python | Rust | TypeScript |
|------|-----|--------|------|------------|
| **单二进制部署** | &#x2705; | &#x274C; | &#x2705; | &#x274C; |
| **零运行时依赖** | &#x2705; | &#x274C; | &#x2705; | &#x274C; |
| **内置并发** | &#x2705; Goroutines | &#x274C; asyncio | &#x2705; async/await | &#x2705; event loop |
| **安装方式** | `go install` | `pip install` | `cargo build` | `npm install` |
| **学习曲线** | 适中 | 简单 | 陡峭 | 适中 |

</div>

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
export GO_CODE_API_KEY=sk-ant-...
./go-code
```

---

> **翻译说明**：本站所有主要页面已完成中英文双语版本。部分深度架构内容目前仅有英文版本，中文翻译将持续更新。
