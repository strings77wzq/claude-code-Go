## Context

appearance 配置使用了错误的 API 格式导致主题切换无效。底部 CTA 与 Hero 按钮功能重复。

## Goals / Non-Goals

**Goals:**
- 修复主题切换（深色/浅色/跟随系统）
- 删除冗余 CTA
- 完善浅色主题样式

**Non-Goals:**
- 不改文档内容
- 不改 Go 源代码

## Decisions

### 1. appearance 配置

**Decision**: `appearance: true`

**Rationale**: VitePress 文档明确说明 `true` 启用切换，默认跟随系统，用户可手动切换。之前的 `{ label: '自动', value: 'auto' }` 不是有效 API。

### 2. CTA 删除

**Decision**: 完全删除底部 CTA 区块。

**Rationale**: Hero 已有"快速开始"+"查看源码"按钮，底部 CTA 功能完全重复。页面以"快速开始"代码块结尾更自然。

### 3. 浅色主题样式

**Decision**: 终端窗口在浅色主题下保持深色背景（终端本身就是深色）。

**Rationale**: 终端模拟器传统上就是深色，浅色主题下保持深色终端窗口更符合用户预期。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 浅色主题下自定义组件样式不一致 | 测试所有自定义组件在两种主题下的表现 |
| 删除 CTA 后页面结尾突兀 | 以"快速开始"代码块自然过渡到 footer |
