## Why

三个阻塞性问题影响用户体验：1) Spinner 动画不动，用户不知道是否在思考还是卡住了；2) 每次启动都要手动 `go build`，不像 opencode/claude-code 那样输入命令直接启动；3) TUI 对话框不够美观，输入输出模式需要模仿 opencode/claude-code 的优雅交互。

## What Changes

### TUI 修复
- 修复 Spinner 动画（添加 `m.spinner.Tick` 到 tea.Batch）
- 优化对话框样式（消息气泡、工具调用高亮、分隔线）
- 流式输出逐字渲染（不等待完整响应）

### 安装优化
- 添加 Makefile `install` 目标（`go install`）
- 更新 README 安装说明

### 官网更新
- 更新安装文档（`go install` 方式）
- 更新 TUI 截图

## Capabilities

### Modified Capabilities
- `bubbletea-tui`: 修复 spinner 动画，优化对话框样式
- `docs-website`: 更新安装文档

## Impact

- 修改文件: `pkg/tui/tui.go`, `Makefile`, `README.md`
- 修改文件: 官网安装文档
- 不影响: 核心 Agent Loop、工具系统
