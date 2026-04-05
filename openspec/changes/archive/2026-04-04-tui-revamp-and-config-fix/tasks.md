## 1. Phase 1: Bug 修复

- [x] 1.1 读取 ANTHROPIC_MODEL 环境变量

## 2. Phase 2: TUI 重构

- [x] 2.1 引入 bubbletea 依赖（go.mod）
- [x] 2.2 创建 pkg/tui/tui.go（完整 TUI：状态+渲染+事件处理）
- [x] 2.3 创建 pkg/tui/tui.go（渲染逻辑，View 方法）
- [x] 2.4 创建 pkg/tui/tui.go（事件处理，Update 方法）
- [x] 2.5 创建 pkg/tui/tui.go（spinner 加载动画）
- [x] 2.6 创建 pkg/tui/components/permission.go（权限审批）
- [x] 2.7 修改 cmd/go-code/main.go（使用新 TUI，保留 --legacy-repl）
- [x] 2.8 更新 README 说明新 TUI
