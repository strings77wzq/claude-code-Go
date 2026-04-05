---
title: REPL 命令参考
description: go-code 所有 REPL 命令的完整参考
---

# REPL 命令参考

go-code 提供了一组斜杠命令来控制 REPL 会话。这些命令允许您获取帮助、管理会话、切换模型以及控制应用程序行为。

## 命令列表

| 命令 | 描述 | 示例 |
|---------|-------------|---------|
| [`/help`](#help) | 显示帮助信息 | `/help` |
| [`/clear`](#clear) | 清除会话历史 | `/clear` |
| [`/model`](#model) | 显示或切换模型 | `/model claude-opus-4-20250514` |
| [`/models`](#models) | 列出可用模型 | `/models` |
| [`/sessions`](#sessions) | 列出已保存的会话 | `/sessions` |
| [`/resume`](#resume) | 恢复会话 | `/resume session-id` |
| [`/compact`](#compact) | 压缩会话上下文 | `/compact` |
| [`/update`](#update) | 检查并应用更新 | `/update` |
| [`/exit`](#exit) | 退出应用程序 | `/exit` |
| [`/skills`](#skills) | 列出可用的技能 | `/skills` |

---

## /help

显示所有可用命令的帮助信息。

### 使用方法

```
/help
```

### 描述

显示所有可用 REPL 命令的摘要及简要描述。这是快速了解 go-code 功能的最佳方式。

### 示例输出

```
可用命令：
  /help       - 显示此帮助信息
  /clear      - 清除会话历史
  /model      - 显示/切换当前模型
  /models     - 列出可用模型
  /sessions   - 列出已保存的会话
  /resume     - 恢复会话
  /compact    - 压缩会话上下文
  /update     - 检查更新
  /exit       - 退出应用程序
  /skills     - 列出可用技能

输入 /<command> 来使用命令。
```

---

## /clear

清除会话历史。

### 使用方法

```
/clear
```

### 描述

清除当前会话的所有会话历史。代理将丢失之前消息的上下文，但会话保持活动状态。这对于在不结束会话的情况下重新开始很有用。

### 行为

- 清除代理的消息历史
- 不影响已保存的会话
- 保留当前模型和设置

### 示例

```
go-code> /clear
Conversation history cleared
```

---

## /model

显示或更改当前模型。

### 使用方法

```
/model              # 显示当前模型
/model <model-name> # 切换到不同的模型
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `model-name` | string | 否 | 要切换到的模型 |

### 描述

无参数时，显示当前活动的模型。带模型名称时，切换到该模型以进行后续请求。

### 可用模型

**Anthropic：**
- `claude-sonnet-4-20250514`（默认）
- `claude-opus-4-20250514`
- `claude-haiku-4-20250514`

**腾讯云 Coding 计划：**
- `tc-code-latest`（自动）
- `hunyuan-2.0-instruct`
- `hunyuan-2.0-thinking`
- `minimax-m2.5`
- `kimi-k2.5`
- `glm-5`
- `hunyuan-t1`
- `hunyuan-turbos`

### 示例

**显示当前模型：**
```
go-code> /model
Current model: claude-sonnet-4-20250514
```

**切换到不同模型：**
```
go-code> /model claude-opus-4-20250514
Model switched to: claude-opus-4-20250514
```

---

## /models

列出所有可用模型及描述。

### 使用方法

```
/models
```

### 描述

显示可与 go-code 一起使用的所有模型的完整列表，按提供商组织。显示 Anthropic 模型和腾讯云 Coding 计划模型。

### 示例输出

```
可用模型：

  Anthropic:
    claude-sonnet-4-20250514 (默认)
    claude-opus-4-20250514
    claude-haiku-4-20250514

  腾讯云 Coding 计划:
    tc-code-latest (自动)
    hunyuan-2.0-instruct
    hunyuan-2.0-thinking
    minimax-m2.5
    kimi-k2.5
    glm-5
    hunyuan-t1
    hunyuan-turbos

切换模型: /model <model-name>
```

---

## /sessions

列出所有已保存的会话。

### 使用方法

```
/sessions
```

### 描述

显示会话目录中的所有已保存会话。每个会话显示其 ID、模型、轮次计数和开始时间。会话会自动保存，可以在以后恢复。

### 示例输出

```
可用会话：
  abc123  model=claude-sonnet-4-20250514 turns=5 started=2026-04-05 10:30:00
  def456  model=claude-opus-4-20250514 turns=12 started=2026-04-04 15:45:30
  ghi789  model=claude-sonnet-4-20250514 turns=3 started=2026-04-03 09:20:15
```

### 参见

- [`/resume`](#resume) — 恢复特定会话

---

## /resume

恢复之前的会话。

### 使用方法

```
/resume <session-id>
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `session-id` | string | 是 | 要恢复的会话 ID |

### 描述

加载之前会话的会话历史并从中断处继续。会话必须存在于会话目录中。

### 示例

**恢复会话：**
```
go-code> /resume abc123
Resumed session abc123 with 10 messages
Session model: claude-sonnet-4-20250514
```

**无效会话：**
```
go-code> /resume nonexistent
Session not found: nonexistent
```

---

## /compact

压缩会话上下文。

### 使用方法

```
/compact
```

### 描述

触发长会话的上下文压缩。这通过总结旧消息同时保留关键信息来减少会话的内存占用。对于很长的会话很有用，可以保持性能。

### 行为

- 总结旧的会话消息
- 减少 API 调用的 token 数量
- 保留重要上下文

### 示例

```
go-code> /compact
Conversation compacted
```

---

## /update

检查并应用更新。

### 使用方法

```
/update
```

### 描述

连接到发布服务器检查是否有新版本可用。如果有更新可用，提示下载并替换当前二进制文件。

### 行为

1. 从 GitHub releases 检查最新版本
2. 与当前版本比较
3. 如果有可用更新，提示确认
4. 下载并替换二进制文件
5. 需要重启以应用

### 示例

**没有可用更新：**
```
go-code> /update
Already up to date (v0.1.0)
```

**有可用更新：**
```
go-code> /update
Update available: v0.1.0 -> v0.1.1
Download and replace binary? [y/N]: y
Update successful. Please restart go-code.
```

---

## /exit

退出应用程序。

### 使用方法

```
/exit
# 或
/quit
```

### 描述

结束当前 REPL 会话并退出应用程序。会话在退出前会自动保存。

### 示例

```
go-code> /exit

Goodbye!
```

---

## /skills

列出所有可用的技能。

### 使用方法

```
/skills
```

### 描述

显示已配置的所有自定义技能。技能是可以用 `/skillname` 调用的自定义命令。每个技能显示其名称和描述。

### 示例输出

```
可用技能：
  /brainstorming - 用于任何创造性工作之前
  /debugging    - 系统化调试工作流
  /refactor     - 智能重构工作流
  /tdd          - 测试驱动开发工作流
```

---

## 技巧

### 命令历史

使用上下箭头键在当前会话中导航之前的命令。

### 部分输入

对于 `/model`，您可以在命令后直接输入模型名称：
```
go-code> /model claude-opus-4-20250514
```

### Tab 补全

基本 REPL 不支持 Tab 补全。使用 TUI 模式获得增强功能。

---

## 相关文档

- [配置指南](../guide/configuration.md) — 模型配置
- [会话管理](../guide/session-management.md) — 会话持久化
- [技能系统](../extension/skills.md) — 自定义命令