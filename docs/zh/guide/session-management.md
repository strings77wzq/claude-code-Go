---
title: 会话管理
description: 如何在 claude-code-Go 中列出、恢复和管理已保存的会话
---

# 会话管理

claude-code-Go 会自动将您的对话会话保存到磁盘，允许您查看过去的对话并恢复中断的会话。

## 会话保存方式

会话以 **JSONL**（JSON 行）格式保存在会话目录中（默认：`~/.claude-code-go/sessions/`）。

每个会话文件包含标准化 trace 行：
- `meta`：会话 ID、模型、开始时间和初始状态。
- `request` / `response`：Provider 请求和响应摘要。
- `tool`：工具名称、输入、输出和耗时。
- `permission`：权限决策、脱敏摘要和时间戳。
- `message`：可恢复的用户和 assistant 消息。
- `status`：最终状态、轮次和 token 使用量。

文件结构示例：

```jsonl
{"type":"meta","session_id":"sess_123","model":"claude-sonnet-4-6-20251001","start_time_ms":1234567890000,"end_time_ms":0,"turn_count":0,"input_tokens":0,"output_tokens":0,"status":"running"}
{"type":"request","model":"claude-sonnet-4-6-20251001","messages_count":1,"timestamp_ms":1234567891000}
{"type":"message","role":"user","content":"你好","timestamp_ms":1234567890000}
{"type":"message","role":"assistant","content":"你好！有什么可以帮你的？","timestamp_ms":1234567895000}
{"type":"status","status":"completed","turn_count":1,"input_tokens":1000,"output_tokens":500,"timestamp_ms":1234567990000}
```

## 会话文件位置

默认情况下，会话存储在：

```
~/.claude-code-go/sessions/
```

每个会话文件命名为：`session-{时间戳}.jsonl`

时间戳是会话开始时的 Unix 时间戳。

## 列出会话

要查看所有已保存的会话，请在 REPL 中使用 `/sessions` 命令：

```
/sessions
```

这将显示按最新时间排序的会话列表，包括：
- 会话 ID
- 使用的模型
- 开始/结束时间
- 轮次数量

## 恢复会话

要恢复之前的会话，请使用会话 ID 运行 `/resume` 命令：

```
/resume <session_id>
```

例如：
```
/resume sess_123
```

这将加载会话文件中的所有消息并恢复您的对话状态，允许您从中断的地方继续。

## 回放会话

使用 `go-code replay` 可以在不调用真实 Provider 的情况下检查会话：

```bash
go-code replay latest
go-code replay sess_123
go-code replay ~/.claude-code-go/sessions/session-1234567890.jsonl
```

Replay 会打印请求、工具调用、权限决策、消息、错误和最终状态，适合调试和提交 issue 时附带证据。

## 会话元数据

每个会话存储以下元数据：

| 字段 | 描述 |
|-----|------|
| `session_id` | 会话的唯一标识符 |
| `model` | 使用的 Claude 模型（例如 claude-sonnet-4-20250514）|
| `start_time_ms` | 会话开始时间（Unix 毫秒）|
| `end_time_ms` | 会话结束时间（Unix 毫秒）|
| `turn_count` | 对话轮次数量 |
| `input_tokens` | 使用的输入 token 总数 |
| `output_tokens` | 生成的输出 token 总数 |
| `status` | `running`、`completed`、`failed` 或 `max_turns` |

## 编程访问

您还可以使用 `session` 包以编程方式访问会话：

```go
import "github.com/strings77wzq/claude-code-Go/internal/session"

// 列出所有会话
sessions, err := session.ListSessions("~/.claude-code-go/sessions")

// 加载特定会话
sess, messages, err := session.LoadSession("~/.claude-code-go/sessions/session-1234567890.jsonl")

// 回放会话 trace
events, err := session.ReplaySessionFile("~/.claude-code-go/sessions/session-1234567890.jsonl")
```

`SessionInfo` 结构体包含：
- `ID`：会话标识符
- `FilePath`：会话文件的完整路径
- `StartTime`：会话开始时间
- `EndTime`：会话结束时间
- `TurnCount`：轮次数量
- `Model`：使用的模型
