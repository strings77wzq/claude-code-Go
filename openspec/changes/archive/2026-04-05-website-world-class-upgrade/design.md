## Context

当前官网使用 VitePress 默认主题，功能完整但缺乏顶级开源项目的交互体验和设计感。需要从 React、Vue、Tailwind、Rust 等顶级项目中汲取设计灵感。

## Goals / Non-Goals

**Goals:**
- Hero 区域终端打字机动画
- 自定义 VitePress 主题
- 新用户引导流程
- 交互式架构图
- API 参考文档 + 故障排查
- GitHub About 优化

**Non-Goals:**
- 不重写整个网站（基于 VitePress 升级）
- 不使用 React/Vue 自定义框架
- 不做自定义域名（后续迭代）

## Decisions

### 1. 交互性：终端打字机动画

**Decision**: 使用 CSS animation + JavaScript 实现终端打字机效果。

```
Hero 区域显示模拟终端：
┌──────────────────────────────────────┐
│ $ go-code                            │
│ > 帮我写个 HTTP 服务器              │  ← 打字机效果
│ ⡿ Thinking...                       │  ← Spinner 动画
│ 🛠️ Tool call: Write → main.go       │
│ ✓ File written                       │
│ > _                                   │  ← 光标闪烁
└──────────────────────────────────────┘
```

### 2. 设计性：自定义主题

**Decision**: 基于 VitePress 默认主题深度定制，不从头写。

```css
:root {
  --vp-c-brand-1: #00ADD8;  /* Go 蓝 */
  --vp-c-brand-2: #00875a;  /* Go 绿 */
  --vp-c-bg: #0d1117;       /* GitHub 暗色 */
  --vp-c-bg-soft: #161b22;
}
```

### 3. 引导性：角色选择

**Decision**: 首页增加"选择你的角色"板块。

```
你是哪种开发者？

┌──────────┐  ┌──────────┐  ┌──────────┐
│ 👨‍💻 全栈   │  │ 🏗️ 架构师 │  │ 🎓 学生   │
│ 快速上手  │  │ 深度架构  │  │ 学习原理  │
└──────────┘  └──────────┘  └──────────┘
```

### 4. GitHub About

**Decision**: 简洁有力，突出差异化优势。

```
About:
Model provides intelligence, Harness provides reliability.
A production-grade AI coding assistant in pure Go — with multi-provider support,
runtime model switching, and full agent loop.

Website: https://strings77wzq.github.io/claude-code-Go/
Topics: ai-agent, claude-code, go, coding-assistant, mcp, llm, terminal
```

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| 自定义主题可能影响性能 | 只定制 CSS，不添加重 JS |
| 打字机动画可能卡顿 | 使用 CSS animation，GPU 加速 |
| 内容过多导致页面臃肿 | 分页加载，按需渲染 |
