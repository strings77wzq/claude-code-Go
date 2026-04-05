## 1. 文件日志系统

- [x] 1.1 创建 `internal/logger/logger.go`（初始化 + 文件写入）
- [x] 1.2 记录关键事件（API 请求、工具执行、错误、会话）
- [x] 1.3 创建 `internal/logger/trace.go`（HTTP trace 记录）

## 2. Debug 模式

- [x] 2.1 main.go 添加 `--debug` 和 `--trace-http` 参数
- [x] 2.2 Debug 模式日志同时输出到 stderr
- [x] 2.3 TUI debug 状态栏（API 延迟、token 用量）

## 3. 硬超时与计时器

- [x] 3.1 API 客户端 5 分钟硬超时（context.WithTimeout）
- [x] 3.2 TUI 实时计时显示（`⏳ Waiting... (2.3s)`）
- [x] 3.3 超时后明确错误信息

## 4. Session Trace

- [x] 4.1 Agent Loop 记录 request/response 到 session 文件
- [x] 4.2 新增 `/trace last` 命令
- [x] 4.3 新增 `/export session` 命令
