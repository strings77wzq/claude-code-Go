---
title: 钩子系统
description: 钩子系统技术深度解析 — 钩子接口、内置实现、自定义钩子创建、注册和执行流程
---

# 钩子系统

go-code 实现了钩子系统，提供工具执行前后的回调，用于监控和扩展工具行为。本文提供全面的技术概述。

## 什么是钩子？

**钩子**是在工具执行前后触发的回调。它们提供以下机制：

- **监控** — 跟踪工具调用和结果
- **日志** — 记录详细的执行信息
- **审计** — 维护符合合规要求的审计跟踪
- **指标** — 收集性能和使用统计
- **通知** — 对特定工具操作发出警报
- **验证** — 预检查输入或验证输出

## 钩子接口

钩子接口定义了两个方法：

```go
type Hook interface {
    // Name returns the unique identifier for this hook.
    Name() string

    // PreExecute is called before a tool is executed.
    // It receives the tool name and input parameters.
    // Returning an error will prevent tool execution.
    PreExecute(toolName string, input map[string]any) error

    // PostExecute is called after a tool has been executed.
    // It receives the tool name, input parameters, execution result, and whether the result is an error.
    // Errors from PostExecute are logged but do not affect the tool result.
    PostExecute(toolName string, input map[string]any, result string, isError bool) error
}
```

### PreExecute

`PreExecute` 方法在工具运行之前调用。用例包括：
- 验证输入参数
- 检查权限
- 启动性能计时器
- 记录即将执行的日志
- 根据条件阻止执行

如果 `PreExecute` 返回错误，工具执行将被**阻止**，错误将传播给调用者。

### PostExecute

`PostExecute` 方法在工具完成后调用。用例包括：
- 记录执行结果
- 记录指标和持续时间
- 写入审计记录
- 处理结果
- 发送通知

`PostExecute` 的错误会被**记录但忽略** — 它们不影响工具结果或 agent 的执行流程。

## 钩子注册表

`Registry` 管理钩子的注册和执行：

```go
type Registry struct {
    mu        sync.RWMutex
    preHooks  []Hook
    postHooks []Hook
}
```

### 关键方法

```go
// Register adds a hook to the registry.
func (r *Registry) Register(hook Hook) error

// RunPreHooks executes all pre-execute hooks in order.
// Returns the first error encountered (stops execution).
func (r *Registry) RunPreHooks(toolName string, input map[string]any) error

// RunPostHooks executes all post-execute hooks in order.
// Errors are logged but do not stop execution of subsequent hooks.
func (r *Registry) RunPostHooks(toolName string, input map[string]any, result string, isError bool)
```

### 注册流程

```
┌─────────────────────────────────────────────────────────────────────┐
│                     钩子注册流程                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   1. 创建钩子实例 (LoggingHook, AuditHook, 自定义)                 │
│         │                                                           │
│         ▼                                                           │
│   2. 调用 registry.Register(hook)                                  │
│         │                                                           │
│         ▼                                                           │
│   3. 验证无重复名称                                                │
│         │                                                           │
│         ▼                                                           │
│   4. 将钩子添加到 preHooks 和 postHooks 切片                       │
│         │                                                           │
│         ▼                                                           │
│   5. 钩子准备好在所有工具调用时执行                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 执行流程

```
┌─────────────────────────────────────────────────────────────────────┐
│                     钩子执行流程                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   Tool.Execute(input) 调用                                          │
│         │                                                           │
│         ▼                                                           │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │ RunPreHooks(toolName, input)                                │   │
│   │   - Hook1.PreExecute() → OK                                 │   │
│   │   - Hook2.PreExecute() → OK                                 │   │
│   │   - Hook3.PreExecute() → ERROR → STOP                      │   │
│   └─────────────────────────────────────────────────────────────┘   │
│         │                                                           │
│         ▼ (如果所有 pre-hooks 通过)                                 │
│   工具执行 (Read, Write, Edit, Glob, Grep, Bash, MCP)              │
│         │                                                           │
│         ▼                                                           │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │ RunPostHooks(toolName, input, result, isError)               │   │
│   │   - Hook1.PostExecute() → OK                                 │   │
│   │   - Hook2.PostExecute() → OK                                │   │
│   │   - Hook3.PostExecute() → ERROR (已记录，被忽略)           │   │
│   └─────────────────────────────────────────────────────────────┘   │
│         │                                                           │
│         ▼                                                           │
│   将工具结果返回给 agent                                           │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 内置钩子

go-code 提供两个内置钩子：

### LoggingHook

`LoggingHook` 使用 Go 的 `slog` 包记录所有工具执行：

```go
type LoggingHook struct {
    logger *slog.Logger
}
```

**功能：**
- 记录工具名称、输入（截断）和结果（截断）
- 使用 DEBUG 级别进行 pre-execute，成功时使用 INFO，失败时使用 ERROR
- 可配置的截断长度（默认：500 字符）

**用法：**
```go
hook := hooks.NewLoggingHook()              // 使用默认 logger
hook := hooks.NewLoggingHookWithLogger(l)   // 使用自定义 logger
registry.Register(hook)
```

### AuditHook

`AuditHook` 将工具执行记录到 JSONL 审计日志文件：

```go
type AuditHook struct {
    filePath string
    file     *os.File
}

type AuditRecord struct {
    Timestamp  string         `json:"timestamp"`
    ToolName   string         `json:"tool_name"`
    Input      map[string]any `json:"input"`
    Result     string         `json:"result"`
    IsError    bool           `json:"is_error"`
    DurationMS int64          `json:"duration_ms"`
    PreHookErr string         `json:"pre_hook_error,omitempty"`
}
```

**功能：**
- 每行写入一条 JSON 记录（JSONL 格式）
- 记录时间戳、工具名称、输入、结果、错误状态
- 自动创建目录和文件
- 线程安全，带有适当的文件同步

**用法：**
```go
hook, err := hooks.NewAuditHook("/var/log/go-code/audit.jsonl")
if err != nil {
    log.Fatal(err)
}
defer hook.Close()
registry.Register(hook)
```

## 创建自定义钩子

要创建自定义钩子，实现 `Hook` 接口：

```go
import "log/slog"

type MetricsHook struct {
    counters map[string]int64
    mu       sync.Mutex
}

func (h *MetricsHook) Name() string {
    return "metrics"
}

func (h *MetricsHook) PreExecute(toolName string, input map[string]any) error {
    h.mu.Lock()
    h.counters[toolName]++
    h.mu.Unlock()
    return nil
}

func (h *MetricsHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    // 可以记录计时、成功率等
    return nil
}
```

### 示例：通知钩子

```go
type NotificationHook struct {
    channels []chan string
}

func (h *NotificationHook) Name() string {
    return "notification"
}

func (h *NotificationHook) PreExecute(toolName string, input map[string]any) error {
    if toolName == "Bash" {
        for _, ch := range h.channels {
            ch <- fmt.Sprintf("Executing dangerous command: %v", input["command"])
        }
    }
    return nil
}

func (h *NotificationHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    return nil
}
```

### 示例：验证钩子

```go
type ValidationHook struct {
    allowedPaths []string
}

func (h *ValidationHook) Name() string {
    return "validation"
}

func (h *ValidationHook) PreExecute(toolName string, input map[string]any) error {
    if toolName == "Write" || toolName == "Edit" {
        if path, ok := input["filePath"].(string); ok {
            if !h.isPathAllowed(path) {
                return fmt.Errorf("path not allowed: %s", path)
            }
        }
    }
    return nil
}

func (h *ValidationHook) isPathAllowed(path string) bool {
    for _, allowed := range h.allowedPaths {
        if strings.HasPrefix(path, allowed) {
            return true
        }
    }
    return false
}

func (h *ValidationHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
    return nil
}
```

## 用例

### 日志和监控

```go
// 跟踪所有工具使用
registry.Register(hooks.NewLoggingHook())

// 查询日志：
grep "tool execution" /var/log/app.log
```

### 合规审计

```go
// 记录所有操作以符合合规要求
hook, _ := hooks.NewAuditHook("/audit/tool-access.jsonl")
registry.Register(hook)

// 查看审计日志：
cat /audit/tool-access.jsonl | jq '.'
```

### 性能指标

```go
// 跟踪工具使用和性能
type PerfHook struct{}

func (h *PerfHook) Name() string { return "perf" }
func (h *PerfHook) PreExecute(t string, i map[string]any) error { /* start timer */ return nil }
func (h *PerfHook) PostExecute(t string, i map[string]any, r string, e bool) error {
    // 记录持续时间
    return nil
}
registry.Register(&PerfHook{})
```

### 安全验证

```go
// 阻止危险操作
type SecurityHook struct{}

func (h *SecurityHook) Name() string { return "security" }
func (h *SecurityHook) PreExecute(t string, i map[string]any) error {
    if t == "Bash" && strings.Contains(i["command"].(string), "rm -rf") {
        return errors.New("blocked dangerous command")
    }
    return nil
}
func (h *SecurityHook) PostExecute(t string, i map[string]any, r string, e bool) error {
    return nil
}
registry.Register(&SecurityHook{})
```

### 通知

```go
// 对特定事件发出警报
notifHook := &NotificationHook{channels: []chan string{alertChannel}}
registry.Register(notifHook)
```

## 架构概述

```
┌─────────────────────────────────────────────────────────────────────┐
│                        钩子系统架构                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    工具执行                                   │   │
│   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐       │   │
│   │  │  PreExecute │───▶│    工具    │───▶│ PostExecute │       │   │
│   │  │    钩子     │    │   (Read,   │    │    钩子     │       │   │
│   │  │             │    │  Write,    │    │             │       │   │
│   │  │ - Logging   │    │  Edit,     │    │ - Logging   │       │   │
│   │  │ - Audit     │    │  Glob,     │    │ - Audit     │       │   │
│   │  │ - Custom    │    │  Grep,     │    │ - Custom    │       │   │
│   │  └─────────────┘    │  Bash,     │    └─────────────┘       │   │
│   │                      │  MCP)      │                          │   │
│   │                      └─────────────┘                          │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│                              ▼                                       │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    钩子注册表                                │   │
│   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │   │
│   │  │  Register  │    │ RunPreHooks │    │RunPostHooks │      │   │
│   │  └─────────────┘    └─────────────┘    └─────────────┘      │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 相关文档

- [技能系统](./skills.md) — 自定义 agent 行为的命名提示词
- [MCP 集成](./mcp.md) — 外部工具集成
- [工具系统概述](../tools/overview.md) — 工具接口和注册表
- [Agent 循环实现](../core-code/agent-loop-impl.md) — 工具执行流程

---

<div class="nav-prev-next">

- [MCP 集成](./mcp.md) ←
- → [扩展概述](./overview.md)

</div>