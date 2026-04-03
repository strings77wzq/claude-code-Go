## Context

官网存在 3 类问题：
1. 中英文切换 404 — 英文文件在 docs/en/ 下，VitePress root locale 期望在根目录
2. 导航不完整 — 对比 schhaohao.github.io/docs/ 缺少项目简介、项目结构、核心代码、MCP 集成
3. 品牌特色不足 — emoji 图标、缺少项目核心数字和学习收获板块

## Goals / Non-Goals

**Goals:**
- 修复中英文切换 404（移动 docs/en/ → docs/ 根）
- 完善导航结构（6 个导航组，每个中英双语）
- 品牌化 SVG 图标（Go gopher 风格，不抄袭）
- 首页增加项目核心数字和学习收获

**Non-Goals:**
- 不改变 VitePress 框架
- 不改仓库名
- 不做自定义域名

## Decisions

### 1. 英文文件移到根目录

**Decision**: docs/en/guide/* → docs/guide/*, docs/en/architecture/* → docs/architecture/*

**Rationale**: VitePress root locale 期望英文在根目录，中文在 /zh/。这是标准做法。

### 2. SVG 图标设计

**Decision**: 手绘风格 SVG 图标，Go 蓝色 (#00ADD8) 为主色调，融入 Go gopher 元素。

**Rationale**: 与 Go 社区风格一致，有品牌辨识度，不抄袭参考项目。

### 3. 导航结构

**Decision**: 6 个导航组：入门指南、架构设计、核心代码、工具系统、MCP 集成、Harness。

**Rationale**: 对标 schhaohao/docs/ 的导航深度，同时体现 Go + Python Harness 的双语特色。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 移动文件后 sidebar 路径不匹配 | 同步更新 config.ts sidebar 配置 |
| SVG 图标设计质量 | 使用简洁的几何图形，保持专业感 |
| 文档内容过多 | 先完成核心页面，Harness 页面可后续补充 |
