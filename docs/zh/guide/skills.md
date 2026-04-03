---
title: Skills 系统
description: 如何创建和使用自定义 Skills
---

# Skills

Skills 是增强 agent 能力的自定义命令。当你在 REPL 中使用 `/<skill-name>` 调用 Skills 时，Skills 的提示词会被注入到 agent 的系统提示词中，为当前任务提供专门的指令。

## 什么是 Skills

Skills 本质上是存储为 JSON 文件的命名提示词。它们允许你：

- 为常见任务定义可重用的提示词
- 自定义 agent 在特定工作流中的行为
- 创建领域特定的指令

## 如何创建自定义 Skills

在 `.go-code/skills/` 目录下创建 `.json` 文件，格式如下：

```json
{
  "name": "skill-name",
  "description": "这个 Skills 的用途",
  "prompt": "注入的指令提示词",
  "examples": ["/skill-name"]
}
```

### 字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 是 | Skills 的唯一标识符 |
| `description` | 是 | 帮助中显示的简要描述 |
| `prompt` | 是 | 调用 Skills 时注入的提示词 |
| `examples` | 否 | 使用示例 |

### 示例

创建 `.go-code/skills/refactor.json`：

```json
{
  "name": "refactor",
  "description": "重构代码以提高质量",
  "prompt": "重构以下代码以提高可读性、性能和可维护性。确保所有现有测试继续通过。",
  "examples": ["/refactor"]
}
```

## 内置 Skills

默认提供以下 Skills：

| Skills | 说明 |
|------|------|
| `review-pr` | 审查拉取请求 |
| `explain-code` | 解释代码工作原理 |
| `write-tests` | 为代码编写测试 |

## 使用方法

在 REPL 中，键入 `/<skill-name>` 来调用 Skills：

```
> /review-pr
> /explain-code
> /write-tests
```

当 Skills 被调用时，它的提示词会前置到你的消息中并发送给 agent，为任务提供专门的上下文。
