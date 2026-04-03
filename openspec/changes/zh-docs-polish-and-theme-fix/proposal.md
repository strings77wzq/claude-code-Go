## Why

中文文档存在严重翻译腔问题，表述不自然，不符合中文母语开发者的表达习惯。部分专有名词（如 Roadmap、Skills、Harness）被错误翻译。同时，暗色主题切换后浅色页面消失，缺少"跟随系统"选项。

## What Changes

### 中文文档重构
- 全面检查所有中文文档，重构翻译腔表述
- 专有名词保留英文：Roadmap、Skills、Harness、Agent Loop、Provider、MCP、Hooks、Token
- 使用中文开发者社区惯用表达

### 主题切换修复
- 修复暗色/浅色主题切换
- 增加"跟随系统"选项
- 确保两种主题都可用

## Capabilities

### Modified Capabilities
- `docs-website`: 中文文档重构，主题切换修复

## Impact

- 修改文件: 所有 docs/zh/**/*.md 文件
- 修改文件: docs/.vitepress/config.ts（主题配置）
- 不影响: 英文文档、Go 源代码
