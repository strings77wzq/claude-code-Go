---
title: 项目结构
description: 了解 claude-code-Go 的目录布局和模块组织
---

# 项目结构

本文档全面概述了 claude-code-Go 项目的目录结构，并解释每个模块的职责。

## 目录树

```
claude-code-Go/
├── cmd/go-code/              # 🚀 主入口点
│   └── main.go               # 应用启动 + 信号处理
│
├── internal/                 # 🏗️ 核心模块（私有）
│   ├── agent/                # 🧠 智能体循环 + 上下文管理
│   │   ├── loop.go          # 核心智能体执行周期
│   │   ├── history.go       # 消息历史追踪
│   │   └── compact.go       # 上下文压缩逻辑
│   │
│   ├── api/                  # 🌐 Anthropic API 客户端
│   │   ├── client.go        # Messages API 的 HTTP 客户端
│   │   ├── stream.go        # SSE token 流式处理
│   │   └── types.go         # API 请求/响应类型
│   │
│   ├── config/              # ⚙️ 配置管理
│   │   ├── loader.go        # 多源配置加载器
│   │   ├── loader_test.go   # 配置加载测试
│   │   └── types.go         # 配置结构体
│   │
│   ├── permission/          # 🛡️ 权限系统
│   │   ├── policy.go        # 权限策略引擎
│   │   ├── rules.go         # 规则定义
│   │   ├── prompter.go      # 用户权限提示
│   │   └── rules_test.go    # 权限规则测试
│   │
│   ├── session/             # 💾 会话持久化
│   │   ├── session.go       # 会话状态管理
│   │   └── session_test.go  # 会话测试
│   │
│   ├── hooks/               # 🔌 执行前后钩子
│   │   ├── hooks.go         # 钩子接口 + 注册表
│   │   └── builtin.go       # 内置钩子实现
│   │
│   └── tool/                # 🔧 工具系统
│       ├── tool.go          # 工具接口定义
│       ├── registry.go      # 工具注册 + 查询
│       ├── builtin/         # 内置工具
│       │   ├── read.go      # 文件读取工具
│       │   ├── write.go     # 文件写入工具
│       │   ├── edit.go      # 代码编辑工具
│       │   ├── glob.go      # 文件模式匹配
│       │   ├── grep.go      # 内容搜索
│       │   └── bash.go      # Shell 命令执行
│       ├── mcp/             # MCP 集成
│       │   ├── manager.go   # MCP 服务器生命周期
│       │   ├── client.go    # MCP 协议客户端
│       │   ├── adapter.go   # MCP 工具适配器
│       │   ├── config.go    # MCP 配置
│       │   └── transport.go # 传输层
│       └── init/            # 工具注册
│           └── register.go  # 内置工具注册
│
├── pkg/tty/                 # 🎨 终端 UI
│   ├── repl.go             # REPL 主循环
│   ├── renderer.go         # 终端输出渲染
│   └── repl_test.go        # REPL 测试
│
├── harness/                 # 🧪 Python 测试框架（可选）
│   ├── mock_server/        # 模拟 Anthropic API
│   │   └── server.py       # 模拟 API 服务器
│   ├── evaluators/         # 质量评估
│   │   └── evaluator.py    # 响应质量检查
│   └── replay/             # 会话回放 + 追踪
│       └── replay.py       # 调试回放工具
│
├── docs/                    # 📚 VitePress 文档
│   ├── en/                  # 英文文档
│   │   ├── guide/          # 用户指南
│   │   └── architecture/   # 架构文档
│   └── zh/                  # 中文文档
│       ├── guide/
│       └── architecture/
│
├── bin/                     # 📦 编译后的二进制文件（生成）
│
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 依赖
├── Makefile                # 构建自动化
└── README.md               # 项目自述文件
```

## 模块职责

| 模块 | 职责 | 关键文件 |
|------|------|----------|
| **cmd/go-code** | 应用入口点，信号处理，组件初始化 | `main.go` |
| **internal/agent** | 核心智能体循环执行，消息历史，上下文压缩 | `loop.go`, `history.go`, `compact.go` |
| **internal/api** | Anthropic Messages API 通信，SSE 流式传输 | `client.go`, `stream.go` |
| **internal/config** | 多源配置加载（环境变量、文件、默认值） | `loader.go`, `types.go` |
| **internal/permission** | 三层权限系统，用户审批提示 | `policy.go`, `rules.go`, `prompter.go` |
| **internal/session** | 会话状态持久化和管理 | `session.go` |
| **internal/hooks** | 执行前后钩子，用于扩展性 | `hooks.go`, `builtin.go` |
| **internal/tool** | 工具注册表，内置工具，MCP 集成 | `registry.go`, `builtin/*`, `mcp/*` |
| **pkg/tty** | 终端 REPL，输入处理，输出渲染 | `repl.go`, `renderer.go` |

## 依赖关系

模块依赖遵循从顶层到底层的单向流动：

```
                    ┌──────────────┐
                    │ cmd/go-code  │  (入口点)
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │  config  │ │   agent  │ │   tool   │
        └────┬─────┘ └─────┬─────┘ └────┬─────┘
             │             │            │
             ▼             ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │   api    │ │ session  │ │permission│
        └──────────┘ └────┬─────┘ └────┬─────┘
                           │            │
                           ▼            ▼
                     ┌──────────┐ ┌──────────┐
                     │  hooks   │ │  (via    │
                     └──────────┘ │ agent)   │
                                   └──────────┘

              ┌─────────────────────────────────┐
              │           pkg/tty               │  (使用 agent)
              └─────────────────────────────────┘
```

**核心原则：**
- 依赖是**单向的** — 无循环依赖
- `internal/*` 模块是私有的，形成核心层
- `pkg/tty` 依赖 `internal/agent` 提供智能体功能
- 所有模块最终通过 `cmd/go-code` 作为组合根连接

## 设计说明：外观模式

`internal/agent/loop.go` 中的 `AgentLoop` 充当一个**外观模式**，抽象了以下复杂性：

- API 通信
- 工具执行
- 权限检查
- 会话管理
- 钩子调用

`pkg/tty` 中的 REPL 只需要与这个单一接口交互，保持表示层简单且与核心逻辑解耦。

## 下一步

- [架构概览](../architecture/overview.md) — 深入了解系统设计
- [智能体循环详解](../architecture/agent-loop.md) — 理解核心执行周期
- [工具架构](../architecture/tools.md) — 了解工具系统和 MCP 支持