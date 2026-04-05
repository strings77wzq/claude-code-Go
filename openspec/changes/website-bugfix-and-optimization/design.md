## Context

官网发现 27 个问题：16 个导致 404/跳转错误，5 个影响体验，6 个影响品质。

## Goals / Non-Goals

**Goals:**
- 修复所有 404 和跳转错误
- 精简导航栏
- 清理重复文档和孤立文件
- 优化外观和移动端体验

**Non-Goals:**
- 不重写整个网站
- 不添加新功能（除了 Hero 打字机动画）

## Decisions

### 1. 导航栏重构

**Decision**: 从 10+ 个顶级菜单项精简到 5 个。

```
之前: Guide | Architecture | Extensions | Tools | API Reference | Troubleshooting | Contributing | Why Go? | Resources | GitHub
之后: Guide | Architecture | API | Resources | GitHub
```

### 2. 中文页面链接

**Decision**: 所有中文侧边栏和导航链接统一加 `/zh/` 前缀。

### 3. 重复文档处理

**Decision**: 保留更完整的版本，删除重复的。

| 保留 | 删除 |
|------|------|
| `docs/extension/skills.md` | `docs/guide/skills.md` |
| `docs/api/config.md` | `docs/guide/configuration.md` |
| `docs/guide/installation.md` | `docs/guide/installation-for-agents.md` |

### 4. 孤立文件整合

**Decision**: 将孤立文件整合到导航中。

| 文件 | 整合到 |
|------|--------|
| `docs/core-code/entry-point.md` | Architecture 子菜单 |
| `docs/core-code/agent-loop-impl.md` | Architecture 子菜单 |
| `docs/guide/homebrew.md` | Guide 子菜单 |
| `docs/guide/session-management.md` | Guide 子菜单 |

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| 删除文档可能丢失内容 | 先整合到保留文件中再删除 |
| 导航变更影响 SEO | 设置 redirect 规则 |
