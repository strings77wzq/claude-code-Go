## Context

Spinner 动画不动，安装不方便，TUI 不够美观。

## Goals / Non-Goals

**Goals:**
- 修复 spinner 动画
- 优化 TUI 对话框样式（模仿 opencode/claude-code）
- 支持 `go install` 一键安装

**Non-Goals:**
- 不改 Agent Loop 逻辑
- 不重写整个 TUI（渐进优化）

## Decisions

### 1. Spinner 修复

**Decision**: 在 `runAgent` 的 `tea.Batch` 中添加 `m.spinner.Tick`。

**Rationale**: bubbletea 需要 Tick 命令来驱动 spinner 动画帧更新。

### 2. TUI 样式优化

**Decision**: 使用 lipgloss 美化消息显示，增加分隔线和状态指示器。

**Rationale**: 模仿 opencode 的简洁风格，用户消息左对齐，助手消息带缩进。

### 3. 安装方式

**Decision**: `go install ./cmd/go-code` 安装到 $GOPATH/bin。

**Rationale**: Go 标准安装方式，与 opencode 一致。
