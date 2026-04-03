---
title: 智能体循环实现
description: 深入解析 Run() 方法 —— stop_reason 分派、历史管理、工具执行和会话持久化
---

# 智能体循环实现

智能体循环是 go-code 的核心执行引擎，实现位于 `internal/agent/loop.go`。本文档详细说明 `Run()` 方法及其关键组件。

## Run() 方法概述

```go
func (a *Agent) Run(ctx context.Context, userInput string, 
                   outputCallback func(string)) (string, error)
```

方法接收：
- `ctx` — 取消上下文
- `userInput` — 用户消息
- `outputCallback` — 接收流式文本的函数

返回最终文本响应或错误。

## 执行流程

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Run() 方法流程                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. 初始化会话（生成 ID，记录开始时间）                              │
│  2. 将用户消息添加到历史记录                                        │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    FOR 循环 (MaxTurns)                       │   │
│  │                                                              │   │
│  │  3. 必要时压缩历史记录                                       │   │
│  │  4. 构建 API 请求（系统 + 工具 + 消息）                      │   │
│  │  5. 发送到 API 并接收流式响应                                │   │
│  │  6. 将助手消息添加到历史记录                                  │   │
│  │  7. 根据 stop_reason 分派：                                 │   │
│  │     - end_turn / stop_sequence → 返回文本                  │   │
│  │     - max_tokens → 返回文本 + 警告                          │   │
│  │     - tool_use → 执行工具，添加结果，继续                    │   │
│  │     - default → 返回文本                                    │   │
│  │  8. 退出时保存会话                                           │   │
│  │                                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 会话初始化

```go
func (a *Agent) Run(ctx context.Context, userInput string, 
                   outputCallback func(string)) (string, error) {
    a.sessionID = generateSessionID()
    a.startTime = time.Now()

    var totalInputTokens, totalOutputTokens int

    if err := a.history.AddUserMessage(userInput); err != nil {
        return "", fmt.Errorf("failed to add user message: %w", err)
    }
```

每个会话获得：
- 唯一 ID：`sess_<时间戳>`
- 开始时间戳用于会话跟踪
- 用于监控的令牌计数器

## 主循环

```go
turns := 0

for turns < MaxTurns {
    // 步骤 1：必要时压缩历史记录
    CompactIfNeeded(a.history, a.contextConfig)

    // 步骤 2：构建并发送 API 请求
    req := a.buildRequest()
    resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
```

### 历史记录压缩

```go
CompactIfNeeded(a.history, a.contextConfig)
```

随着对话增长，历史记录会定期压缩以保持在令牌限制内。这由压缩模块处理，以防止上下文溢出。

### 请求构建

```go
func (a *Agent) buildRequest() *api.ApiRequest {
    toolDefs := make([]api.ToolDefinition, 0)
    for _, td := range a.toolRegistry.GetAllDefinitions() {
        toolDefs = append(toolDefs, api.ToolDefinition{
            Name:        td.Name,
            Description: td.Description,
            InputSchema: td.InputSchema,
        })
    }

    return &api.ApiRequest{
        Model:     a.model,
        MaxTokens: a.maxTokens,
        System:    a.systemPrompt,
        Stream:    true,
        Tools:     toolDefs,
        Messages:  a.history.GetMessages(),
    }
}
```

每个请求包含：
- 模型标识符
- 最大令牌限制
- 系统提示词
- 来自注册表的工具定义
- 当前对话历史

## Stop Reason 分派

API 响应包含 `stop_reason` 字段，决定下一步操作：

### end_turn / stop_sequence

```go
switch resp.StopReason {
case "end_turn", "stop_sequence":
    result := extractTextContent(resp.Content)
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

模型认为任务已完成。将文本内容作为最终响应返回。

### max_tokens

```go
case "max_tokens":
    result := extractTextContent(resp.Content) + 
              "\n[Warning] Response was truncated (max_tokens reached)."
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

响应被截断。返回部分内容并附带警告。

### tool_use

```go
case "tool_use":
    toolResults := a.executeTools(ctx, resp.Content)
    if err := a.history.AddToolResults(toolResults); err != nil {
        a.saveSession(turns, totalInputTokens, totalOutputTokens)
        return "", fmt.Errorf("failed to add tool results: %w", err)
    }
    turns++
    continue
```

模型希望调用工具。执行所有请求的工具，将结果添加到历史记录，然后继续循环。

### unknown / default

```go
default:
    result := extractTextContent(resp.Content)
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

意外 stop_reason 的后备处理。返回生成的任何内容。

## MAX_TURNS 安全限制

```go
// MaxTurns is the maximum number of agent loop iterations to prevent infinite loops.
const MaxTurns = 50
```

循环强制执行最多 50 次迭代以防止：
- 无限工具执行循环
- 无效的模型振荡
- 资源耗尽

如果达到限制：

```go
result := "[Agent loop stopped] Reached maximum turns (" + 
          fmt.Sprintf("%d", MaxTurns) + ")."
a.saveSession(turns, totalInputTokens, totalOutputTokens)
return result, nil
```

## 历史记录管理

### 消息类型

```go
// 添加用户消息
if err := a.history.AddUserMessage(userInput); err != nil {
    return "", fmt.Errorf("failed to add user message: %w", err)
}

// 添加助手响应
if err := a.history.AddAssistantMessage(resp.Content); err != nil {
    return "", fmt.Errorf("failed to add assistant message: %w", err)
}

// 添加工具结果
if err := a.history.AddToolResults(toolResults); err != nil {
    return "", fmt.Errorf("failed to add tool results: %w", err)
}
```

### API 消息格式

Anthropic API 期望以下格式的消息：

```json
{
  "messages": [
    { "role": "user", "content": "..." },
    { "role": "assistant", "content": "..." },
    { "role": "tool", "tool_use_id": "...", "content": "..." }
  ]
}
```

## 工具执行与权限检查

```go
func (a *Agent) executeTools(ctx context.Context, 
                             content []api.ContentBlock) []api.ContentBlock {
    var toolResults []api.ContentBlock

    for _, block := range content {
        if block.Type != "tool_use" {
            continue
        }

        toolName := block.Name
        toolInput := block.Input
        toolUseID := block.ID

        // 权限检查
        if !a.checkPermission(toolName, toolInput) {
            toolResults = append(toolResults, api.ContentBlock{
                Type:      "tool_result",
                ToolUseID: toolUseID,
                IsError:   true,
            })
            continue
        }

        // 执行 pre-hooks
        if a.hooksRegistry != nil {
            if err := a.hooksRegistry.RunPreHooks(toolName, toolInput); err != nil {
                toolResults = append(toolResults, api.ContentBlock{
                    Type:      "tool_result",
                    ToolUseID: toolUseID,
                    Text:      "pre-hook error: " + err.Error(),
                    IsError:   true,
                })
                continue
            }
        }

        // 执行工具
        result := a.toolRegistry.Execute(ctx, toolName, toolInput)

        // 执行 post-hooks
        if a.hooksRegistry != nil {
            a.hooksRegistry.RunPostHooks(toolName, toolInput, 
                                         result.Content, result.IsError)
        }

        // 收集结果
        toolResults = append(toolResults, api.ContentBlock{
            Type:      "tool_result",
            ToolUseID: toolUseID,
            Text:      result.Content,
            IsError:   result.IsError,
        })
    }

    return toolResults
}
```

### 权限检查

```go
func (a *Agent) checkPermission(toolName string, 
                                input map[string]any) bool {
    t := a.toolRegistry.GetTool(toolName)
    requiresPermission := t != nil && t.RequiresPermission()

    decision := a.permissionPolicy.Evaluate(toolName, input, requiresPermission)
    return decision == permission.Allow || decision == permission.Ask
}
```

如果工具的 `RequiresPermission()` 返回 `true`，权限策略将评估是允许、拒绝还是询问用户。

## 会话持久化

```go
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int) {
    if a.sessionID == "" {
        return
    }

    s := &session.Session{
        ID:           a.sessionID,
        Model:        a.model,
        StartTime:    a.startTime,
        EndTime:      time.Now(),
        TurnCount:    turnCount,
        InputTokens:  inputTokens,
        OutputTokens: outputTokens,
    }

    messages := a.convertHistoryToSessionMessages()
    dir := getSessionsDir()

    if err := session.SaveSession(s, messages, dir); err != nil {
        fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
    }
}
```

会话保存到 `~/.claude-code-go/sessions/`，包含：
- 会话 ID 和时间戳
- 使用的模型
- 回合数和令牌使用量
- 完整对话历史

## 相关文档

- [入口点详解](entry-point.md) — main.go 初始化顺序
- [工具系统概述](../tools/overview.md) — 工具接口和注册表
- [架构概述](../architecture/overview.md) — 系统组件

---

<div class="nav-prev-next">

- [入口点详解](entry-point.md) ←
- → [工具系统概述](../tools/overview.md)

</div>