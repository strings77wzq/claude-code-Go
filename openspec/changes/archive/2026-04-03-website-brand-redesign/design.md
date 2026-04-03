## Context

官网存在导航栏中英文不跟随的 bug，且设计缺乏品牌差异化。需要修复 bug 并融合 opencode.ai、openspec.dev、claude.com/docs 的优秀设计元素，打造独特的品牌形象。

## Goals / Non-Goals

**Goals:**
- 修复导航栏多语言 bug（nav 按 locale 分离）
- 首页重构，融合多家优秀设计
- 品牌色调统一（Go 蓝 #00ADD8）
- 新增架构理念、为什么选 Go、多标签快速开始板块

**Non-Goals:**
- 不改文档内容页面（只改首页）
- 不做自定义域名
- 不改 Go 源代码

## Decisions

### 1. 导航栏多语言方案

**Decision**: 使用 VitePress 的 `locales` 配置，nav 数组按 locale 分离。

```typescript
locales: {
  root: {
    label: 'English',
    themeConfig: {
      nav: [ /* English nav */ ]
    }
  },
  zh: {
    label: '中文',
    themeConfig: {
      nav: [ /* Chinese nav */ ]
    }
  }
}
```

**Rationale**: VitePress 官方推荐方式，导航栏自动跟随语言切换。

### 2. 首页设计策略

**Decision**: 用 VitePress 的 `layout: home` frontmatter + 自定义 CSS 实现，不用 CustomHome.vue 组件。

**Rationale**: VitePress 原生 home layout 已经支持 hero + features，通过 frontmatter 的 `features` 数组可以灵活配置。新增板块通过 Markdown 内容补充，不需要自定义 Vue 组件。

### 3. 品牌色调

**Decision**: 通过 `.vitepress/theme/custom.css` 覆盖 VitePress 默认 CSS 变量。

```css
:root {
  --vp-c-brand-1: #00ADD8;  /* Go blue */
  --vp-c-brand-2: #00C853;  /* Go green accent */
}
```

### 4. 页面结构

从上到下：
1. Hero（项目名 + 标语 + 安装命令 + CTA）
2. 核心优势（3 个徽章：单二进制、零依赖、可靠性优先）
3. 架构理念（手风琴：Model 提供智能 / Harness 提供可靠性 / 可扩展生态）
4. 功能特性（6 个卡片）
5. 为什么选 Go（对比表）
6. 快速开始（多标签）
7. 核心数字
8. 页脚

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 首页内容过多导致加载慢 | 纯静态 Markdown，VitePress 构建为静态 HTML，加载极快 |
| CSS 覆盖可能破坏默认主题 | 只覆盖品牌色变量，不修改布局 CSS |
| 导航栏配置复杂 | 严格按照 VitePress locales 文档配置 |
