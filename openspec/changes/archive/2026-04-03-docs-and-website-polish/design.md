## Context

用户反馈文档表述不够高级，网页视觉效果不够 hack 风格。需要在不改变内容结构的前提下，提升文档的专业感和官网的视觉冲击力。

## Goals / Non-Goals

**Goals:**
- 文档表述升级：工程语言更精准，架构图更清晰
- 官网视觉特效：终端窗口、打字机效果、暗色主题
- 保持内容准确性，不为了"高级"而牺牲可读性

**Non-Goals:**
- 不改 Go 源代码
- 不改官网内容结构
- 不做 React/Vue 组件重写（VitePress 原生能力优先）

## Decisions

### 1. 文档表述升级策略

**Decision**: 用工程师的语言，不用营销语言。

**Before**: "claude-code-Go 是一个终端驱动的 AI 编程助手"
**After**: "claude-code-Go 是一个基于 Agent Loop 架构的 AI 编码代理，通过 stop_reason 状态机驱动工具调用链"

**Rationale**: 目标读者是开发者，他们喜欢精准的工程语言，不喜欢空洞的营销词。

### 2. 终端窗口效果

**Decision**: 用纯 CSS + 少量 JS 实现终端窗口效果，不依赖外部库。

```
┌─────────────────────────────────────────────┐
│ $ claude-code-Go                            │
│ > 帮我重构这个模块                           │
│ 🔄 Agent 思考中...                           │
│ 🛠️ 调用工具: Read → main.go                 │
│ ✓ 文件已读取                                 │
│ 🔄 Agent 继续思考...                         │
│ 🛠️ 调用工具: Edit → main.go                 │
│ ✓ 重构完成                                   │
└─────────────────────────────────────────────┘
```

**Rationale**: 纯 CSS 实现简单，加载快，不需要 React/Vue 组件。

### 3. 暗色主题

**Decision**: 默认暗色主题，通过 CSS 变量覆盖 VitePress 默认样式。

```css
:root {
  --vp-c-bg: #0d1117;
  --vp-c-bg-soft: #161b22;
  --vp-c-text-1: #e6edf3;
  --vp-c-brand-1: #00ADD8;
}
```

**Rationale**: Hacker 风格默认暗色，开发者习惯。

### 4. 打字机效果

**Decision**: 用 CSS animation 实现，不依赖 JS 库。

```css
@keyframes typewriter {
  from { width: 0 }
  to { width: 100% }
}
```

**Rationale**: 纯 CSS 性能好，不需要 JS。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| CSS 动画可能影响性能 | 只用 transform 和 opacity 动画，GPU 加速 |
| 暗色主题可能影响可读性 | 保持足够的对比度（4.5:1 以上） |
| 终端窗口效果可能太复杂 | 用纯 CSS，不依赖 JS 库 |
