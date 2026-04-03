## Context

ANTHROPIC_MODEL 环境变量未被读取。当前 TUI 基于 bufio.Scanner，体验差。

## Goals / Non-Goals

**Goals:**
- 修复 ANTHROPIC_MODEL 环境变量读取
- 用 bubbletea 重构 TUI
- 保留旧 REPL 作为 --legacy-repl 选项

**Non-Goals:**
- 不改 Agent Loop 逻辑
- 不改工具系统

## Decisions

### 1. bubbletea 框架

**Decision**: 使用 charmbracelet/bubbletea + bubbles + lipgloss。

**Rationale**: Go 生态最成熟的 TUI 框架，15K+ stars，opencode 也使用类似架构。

### 2. 兼容性

**Decision**: 保留旧 REPL 作为 --legacy-repl 标志。

**Rationale**: 防止新 TUI 有兼容性问题，用户可回退。

### 3. 架构

**Decision**: pkg/tui/ 独立目录，不修改 pkg/tty/。

**Rationale**: 新旧 TUI 并存，互不影响。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 新依赖增加构建复杂度 | bubbletea 是纯 Go，无 CGO |
| TUI 重构工作量大 | 分阶段：先基础框架，再逐步加组件 |
