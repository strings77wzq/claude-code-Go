## Context

中文文档翻译腔严重，专有名词被错误翻译。主题切换功能异常。

## Goals / Non-Goals

**Goals:**
- 中文文档自然流畅，符合开发者表达习惯
- 专有名词保留英文
- 主题切换正常，支持跟随系统

**Non-Goals:**
- 不改英文文档
- 不改 Go 源代码

## Decisions

### 1. 专有名词保留

**Decision**: 以下术语不翻译，保留英文：
- Roadmap（开发路线图）
- Skills（技能系统）
- Harness（可靠性框架）
- Agent Loop（智能体循环）
- Provider（模型提供商）
- MCP（Model Context Protocol）
- Hooks（钩子系统）
- Token（令牌）

### 2. 主题配置

**Decision**: VitePress 默认支持 `appearance: 'force-dark'` 或 `appearance: true`（跟随系统 + 手动切换）。

```typescript
themeConfig: {
  appearance: true, // 启用切换，默认跟随系统
}
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 中文重构工作量大 | 分优先级，先改首页和核心文档 |
| 主题切换可能影响现有样式 | 测试两种主题的对比度 |
