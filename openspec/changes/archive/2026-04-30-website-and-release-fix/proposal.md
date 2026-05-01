## Why

官网存在严重问题：404 跳转、中英文不一致、Logo 显示为 `<C>` 像 C 语言项目、外观平庸。Release 页面有 untagged draft release 冲突。这些问题严重影响项目专业形象和用户信任。

## What Changes

### 官网修复
- 修复所有 404 链接（中英文侧边栏和导航栏对齐）
- 替换 Logo 为 Go 风格（Gopher + 代码括号）
- 精简导航栏为 5 个顶级菜单（Guide, Architecture, API, Resources, GitHub）
- 确保每个英文页面对应中文页面
- 自定义 VitePress 主题样式（Go 蓝配色）

### Release 修复
- 删除 untagged draft release
- 编辑 v0.1.0 release 添加完整说明

## Capabilities

### Modified Capabilities
- `docs-website`: 修复 404、对齐中英文、替换 Logo、精简导航
- `github-release`: 修复 Release 页面冲突

## Impact

- 修改文件: docs/.vitepress/config.ts
- 修改文件: docs/public/logo.svg
- 修改文件: docs/.vitepress/theme/custom.css
- 不影响: Go 源代码、Python Harness
