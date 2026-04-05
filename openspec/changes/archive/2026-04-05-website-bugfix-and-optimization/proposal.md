## Why

官网存在 27 个已发现问题，分为三个优先级：P0（16 个导致 404/跳转错误）、P1（5 个影响体验）、P2（6 个影响品质）。中文页面点击标签跳转到英文或直接 404 是最严重的用户体验问题，必须立即修复。同时需要精简导航栏、清理重复文档、优化外观，使官网达到顶级开源项目水准。

## What Changes

### P0: 修复 404 和跳转错误
- 中文侧边栏链接加 /zh/ 前缀（5 处）
- 中文导航栏链接加 /zh/ 前缀（5 处）
- 创建缺失的中文页面（6 个）
- 修复英文侧边栏路径（troubleshooting/contributing）

### P1: 短期体验优化
- 精简导航栏（10+ → 5-6 个）
- 修复 Discord 占位链接
- 移除不存在的字体引用
- 清理 node_modules（从 git 中移除）
- 合并重复文档

### P2: 中期品质提升
- 自定义 VitePress 主题（Hero 终端打字机动画）
- 品牌色系统深度定制
- 移动端响应式优化
- 孤立文件整合到导航

## Capabilities

### New Capabilities
- `zh-pages-complete`: 完整的中文文档页面
- `nav-optimization`: 精简导航栏结构
- `hero-typewriter`: Hero 区域终端打字机动画
- `mobile-responsive`: 移动端响应式优化

### Modified Capabilities
- `docs-website`: 全面修复链接、导航、主题、外观

## Impact

- 修改文件: `docs/.vitepress/config.ts`（全面重构导航和侧边栏）
- 新增文件: 6 个中文文档页面
- 删除文件: 重复文档、node_modules
- 修改文件: `docs/.vitepress/theme/custom.css`（深度定制）
- 不影响: Go 源代码、Python Harness
