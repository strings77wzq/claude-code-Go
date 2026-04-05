## Context

用户等待 API 响应时无任何反馈，错误信息不分类，不支持非交互模式。

## Goals / Non-Goals

**Goals:**
- 连接超时反馈（3s/30s/5min 三级反馈）
- 详细错误分类（401/403/429/500/timeout/network）
- 非交互模式（-p 参数）

**Non-Goals:**
- 不改 Agent Loop 逻辑
- 不改工具系统

## Decisions

### 1. 连接超时反馈

**Decision**: 使用 goroutine + channel 实现分级反馈。

```
API 请求 goroutine:
  ├─ 0-3s:    正常等待（spinner 旋转）
  ├─ 3-30s:   发送 "Connecting to API..." 消息
  ├─ 30s-5min: 发送 "Still connecting..." 消息
  └─ 5min:    超时错误
```

**Rationale**: 不阻塞主 UI 线程，用户可以随时 Ctrl+C 取消。

### 2. 错误分类

**Decision**: 在 API Client 层统一错误分类，返回结构化错误。

```go
type APIError struct {
    Code    int    // HTTP status code
    Type    string // "auth", "rate_limit", "server", "timeout", "network"
    Message string // 用户友好的错误信息
}
```

**Rationale**: TUI 层只需显示 `APIError.Message`，不需要理解 HTTP 状态码。

### 3. 非交互模式

**Decision**: 使用 `flag` 包解析 `-p` 参数，直接调用 Agent.Run() 后退出。

```
go-code -p "explain this code"
  → 加载配置 → 创建 Agent → Run() → 输出结果 → 退出
```

**Rationale**: 简单直接，与 opencode 的 `-p` 参数一致。

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| 超时反馈可能干扰正常流式输出 | 只在首次 API 请求时显示，流式开始后隐藏 |
| 非交互模式权限自动批准 | 仅用于调试，文档中明确说明 |
| 错误分类可能遗漏新错误码 | 默认 fallback 到通用错误信息 |
