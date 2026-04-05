## Why

当前 TUI 交互是黑箱：用户输入后只显示 "Thinking..."，无法知道是网络慢、模型慢、还是卡死了。日志只输出到 stderr，被 bubbletea alt screen 覆盖，崩溃后无法回溯。缺少 debug 模式，资深工程师无法审查 AI 行为。参考 opencode/claude-code/claw-code-parity 的最佳实践，需要完整的日志、调试和回溯系统。

## What Changes

### 文件日志系统
- 创建 `internal/logger/` 包，日志写入 `~/.go-code/go-code.log`
- JSON 格式，按天轮转
- 记录 API 请求/响应、工具调用、错误、会话事件

### Debug 模式
- 添加 `--debug` 参数，开启后日志同时输出到 stderr
- TUI 底部显示 debug 状态栏（API 延迟、token 用量、工具执行详情）
- 添加 `--trace-http` 参数，记录完整 HTTP 请求/响应

### 硬超时与状态计时器
- API 请求 5 分钟硬超时
- TUI 显示实时计时：`⏳ Waiting... (2.3s)`
- 超时后返回明确错误，不阻塞后续输入

### Session Trace 增强
- Session 文件保存完整 API 交互 trace
- 新增 `/trace last` 命令查看上次交互
- 新增 `/export session` 导出当前会话

## Capabilities

### New Capabilities
- `file-logging`: 结构化文件日志（JSON 格式，~/.go-code/go-code.log）
- `debug-mode`: --debug 参数，日志输出到 stderr + TUI 状态栏
- `hard-timeout`: API 请求 5 分钟硬超时 + 实时计时器
- `session-trace`: 完整 session trace + /trace 命令

### Modified Capabilities
- `bubbletea-tui`: 增加 debug 状态栏 + 实时计时显示
- `cli-entry`: 增加 --debug 和 --trace-http 参数

## Impact

- 新增文件: `internal/logger/` 目录
- 修改文件: `cmd/go-code/main.go`（--debug 参数）
- 修改文件: `pkg/tui/tui.go`（debug 状态栏 + 计时器）
- 修改文件: `internal/agent/loop.go`（trace 记录）
- 不影响: 核心 Agent Loop 逻辑、工具系统、权限系统
