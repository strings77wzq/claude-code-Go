## Why

浅色主题切换无效（appearance 配置 API 错误），底部 CTA 区块与 Hero 按钮功能完全重复，造成页面冗余。

## What Changes

### 主题切换修复
- 修正 `appearance: true`（VitePress 正确 API）
- 启用深色/浅色/跟随系统三种模式
- 完善浅色主题下的自定义组件样式

### CTA 清理
- 删除底部冗余 CTA 区块
- 页面以"快速开始"代码块自然结尾

## Capabilities

### Modified Capabilities
- `docs-website`: 主题切换修复，CTA 清理

## Impact

- 修改文件: docs/.vitepress/config.ts
- 修改文件: docs/index.md, docs/zh/index.md
- 修改文件: docs/.vitepress/theme/custom.css
- 不影响: 英文文档内容、Go 源代码
