## Context

当前 `/model` 命令只显示当前模型，不能切换。`api.Client` 的 model 字段在初始化后不可变。

## Goals / Non-Goals

**Goals:**
- `/model` 无参数显示当前模型，带参数切换模型
- `/models` 列出支持的模型（含 Coding Plan 模型）
- 切换后立即生效，无需重启

**Non-Goals:**
- 不改 Provider 切换（需要重启，涉及 API 格式差异）
- 不改配置文件持久化（运行时切换不写文件）

## Decisions

### 1. `/model` 命令设计

**Decision**: `/model` 无参数显示当前模型，`/model <name>` 切换模型。

```
go-code> /model
Current model: claude-sonnet-4-20250514

go-code> /model hunyuan-2.0-instruct
Model switched to: hunyuan-2.0-instruct
```

**Rationale**: 与 `/compact`、`/update` 等命令风格一致。

### 2. `/models` 命令设计

**Decision**: `/models` 列出支持的模型，分组显示。

```
go-code> /models
Available models:
  Anthropic:
    claude-sonnet-4-20250514 (default)
    claude-opus-4-20250514
    claude-haiku-4-20250514
  Tencent Coding Plan:
    tc-code-latest (Auto)
    hunyuan-2.0-instruct
    minimax-m2.5
    kimi-k2.5
    glm-5
```

**Rationale**: 用户需要知道可用模型列表，特别是 Coding Plan 的 model 参数值。

### 3. 模型切换实现

**Decision**: `api.Client` 增加 `SetModel(model string)` 方法，通过 mutex 保护。

```go
func (c *Client) SetModel(model string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.model = model
}
```

**Rationale**: 简单直接，不需要重建连接或会话。

### 4. 模型列表硬编码

**Decision**: 在 REPL 中硬编码常用模型列表。

**Rationale**: 模型列表变化不频繁，硬编码简单可靠。后续可从配置文件读取。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 切换模型后历史消息格式不兼容 | 不同模型的消息格式由 Agent Loop 统一处理 |
| 无效模型名导致 API 错误 | 切换时不验证，API 返回错误时提示用户 |
