## Context

TUI 交互是黑箱，日志只输出到 stderr 被 alt screen 覆盖，无 debug 模式，无硬超时。

## Goals / Non-Goals

**Goals:**
- 文件日志（JSON 格式，~/.go-code/go-code.log）
- --debug 参数开启详细日志
- TUI 实时计时 + debug 状态栏
- 5 分钟硬超时
- Session 完整 trace

**Non-Goals:**
- 不改 Agent Loop 核心逻辑
- 不改工具系统
- 不做日志轮转（后续迭代）

## Decisions

### 1. 日志架构

**Decision**: 使用标准库 `log/slog` + `os.File` 写入文件。

```
internal/logger/
├── logger.go     ← 初始化，创建 ~/.go-code/go-code.log
└── trace.go      ← HTTP trace 记录
```

**Rationale**: 项目已使用 slog，无需引入新依赖。

### 2. Debug 模式

**Decision**: `--debug` 参数同时启用文件日志和 stderr 输出。

**Rationale**: 简单直接，与 opencode 的 `-d` 一致。

### 3. 超时策略

**Decision**: API 请求使用 `context.WithTimeout`，默认 5 分钟。

**Rationale**: 腾讯云 Coding Plan 模型可能较慢，5 分钟是合理上限。

### 4. TUI 计时器

**Decision**: 在 `runAgent` 中启动独立 ticker，每 500ms 更新状态。

**Rationale**: 不阻塞消息读取，与现有 ticker 架构一致。

### 5. Session Trace

**Decision**: 在 Agent Loop 中记录 request/response 到 session JSONL 文件。

**Rationale**: 完整 trace 支持事后排查，与 claw-code-parity 一致。

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| 日志文件过大 | 后续添加按天轮转，当前手动清理 |
| Debug 模式性能影响 | 仅在 --debug 时记录 HTTP trace |
| TUI 状态栏占用空间 | 仅在 --debug 时显示 |
