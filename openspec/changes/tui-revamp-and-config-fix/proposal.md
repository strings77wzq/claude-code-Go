## Why

两个阻塞性问题：1) ANTHROPIC_MODEL 环境变量未被读取，导致腾讯云 API 返回 400 错误；2) 当前 TUI 基于 bufio.Scanner，无语法高亮、无流式美化、无加载动画、无历史导航，用户体验远落后于 opencode/claude-code/codex。

## What Changes

### Phase 1: Bug 修复
- 读取 ANTHROPIC_MODEL 环境变量

### Phase 2: TUI 重构
- 引入 bubbletea 框架（charmbracelet/bubbletea + bubbles + lipgloss）
- 流式逐字输出渲染
- 工具调用加载动画（spinner 组件）
- 权限审批交互优化
- 历史记录上下键导航
- 保留旧 REPL 作为 --legacy-repl 选项

## Capabilities

### New Capabilities
- `bubbletea-tui`: 基于 bubbletea 的新 TUI
- `spinner-animation`: 工具调用加载动画
- `history-navigation`: 上下键浏览历史
- `permission-dialog`: 交互式权限审批

### Modified Capabilities
- `config-loading`: 增加 ANTHROPIC_MODEL 环境变量读取

## Impact

- 新增文件: pkg/tui/ 目录（model.go, view.go, update.go, components/）
- 修改文件: internal/config/loader.go, cmd/go-code/main.go
- 新增依赖: bubbletea, bubbles, lipgloss
- 不影响: 核心 Agent Loop、工具系统、权限系统
