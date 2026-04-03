---
layout: home
title: claude-code-Go
titleTemplate: Go 语言实现的 AI 编程助手

hero:
  name: claude-code-Go
  text: 模型提供智能，harness 提供可靠性
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

features:
  - icon: ⚡
    title: 单二进制部署
    details: 零依赖，一个文件跑遍全平台——Linux、macOS、Windows。无需运行时、无需虚拟环境、无需 node_modules。
  - icon: 🔒
    title: 可靠性优先
    details: 权限控制、超时保护、会话持久化。Harness 保障安全，让模型专注于智能。
  - icon: 🔌
    title: 可扩展生态
    details: MCP 协议、Hooks、Skills——通过自定义工具扩展能力，无缝集成你的工作流。

stats:
  - label: 源代码文件
    value: 50+
  - label: 模块数量
    value: 8
  - label: 内置工具
    value: 6
  - label: 代码行数
    value: 7,000+
---

## 架构设计理念

::: details 🧠 模型提供智能
LLM 负责：理解意图、决策使用哪个工具、解释结果、规划下一步。它是系统的大脑。
:::

::: details 🛡️ Harness 提供可靠性
运行时负责：权限控制、超时保护、输出截断、会话持久化、错误恢复。它是让系统达到生产级别的安全网。
:::

::: details 🔌 可扩展生态
MCP 协议用于外部工具发现，Hooks 用于执行前后回调，简洁的 Tool 接口——只需实现一个接口即可添加新能力。
:::

## 功能特性

| 功能 | 说明 |
|------|------|
| 🔄 Agent Loop | 基于 stop_reason 驱动的「思考→行动→观察」自主循环 |
| 🛠️ 6 大内置工具 | Read、Write、Edit、Glob、Grep、Bash——开箱即用的完整工具集 |
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