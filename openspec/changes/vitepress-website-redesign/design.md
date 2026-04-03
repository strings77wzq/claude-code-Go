## Context

当前官网使用 MkDocs Material 主题，纯英文文档风格，缺乏产品感。需要迁移到 VitePress，实现：
- 专业产品级首页（hero + 功能卡片 + 统计数据）
- 中英文双语切换
- 深色主题 + 代码高亮 + 动画效果
- 参考 opencode.ai 和 schhaohao.github.io/docs/ 的水准

## Goals / Non-Goals

**Goals:**
- VitePress 站点完全可用，部署到 GitHub Pages
- 首页有产品感：hero 标语、安装命令、功能特性卡片
- 中英文双语内容
- 文档内容从 MkDocs 迁移过来（架构、指南、Harness）

**Non-Goals:**
- 不做自定义域名配置（保持 GitHub Pages 免费方案）
- 不做 SEO 优化（项目初期不需要）
- 不做博客/新闻板块

## Decisions

### 1. VitePress 而非 Next.js

**Decision**: 使用 VitePress 而非 Next.js 定制。

**Rationale**: VitePress 是 Markdown 驱动的静态站点生成器，适合文档+官网一体。Next.js 太重，需要 React 生态，维护成本高。VitePress 天然支持多语言 i18n。

### 2. 双语目录结构

**Decision**: 使用 VitePress 的 locales 配置，`/` 为默认语言（英文），`/zh/` 为中文。

```
docs/
├── .vitepress/
│   └── config.ts        ← 主配置，含 locales 定义
├── index.md             ← 英文首页
├── guide/               ← 英文文档
├── zh/
│   ├── index.md         ← 中文首页
│   └── guide/           ← 中文文档
└── public/
    └── logo.svg         ← 项目 logo
```

### 3. 首页设计

**Decision**: 首页使用 VitePress 的 `layout: home` 模式，包含：
- Hero: 项目名 + 标语 + 安装命令（一键复制）
- Features: 6 个功能卡片（Agent Loop、6 内置工具、权限系统、MCP 集成、SSE 流式、上下文管理）
- 底部 CTA: GitHub 链接 + 文档链接

### 4. 主题定制

**Decision**: 使用 VitePress 默认的 dark/light 主题切换，通过 CSS 变量微调品牌色。不写自定义 CSS 组件。

**Rationale**: VitePress 默认主题已经非常专业，微调即可达到 opencode.ai 的 80% 水准。

### 5. 部署方式

**Decision**: GitHub Actions 工作流改为 `npm install && npx vitepress build docs`，输出到 `docs/.vitepress/dist`，部署到 `gh-pages` 分支。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| VitePress 配置不熟悉 | 使用官方 starter 模板，最小化配置 |
| 文档内容迁移工作量大 | 先迁移首页和核心文档，Harness 文档后续补充 |
| 中文翻译质量 | 由开发者直接撰写，确保技术准确性 |
