## Why

当前官网（VitePress 默认主题）虽然功能完整，但在交互性、设计性、引导性、完整性四个维度距离顶级开源项目（如 React、Vue、Tailwind CSS、Rust）仍有差距。需要从顶级开源项目中汲取设计灵感，打造产品级官网，同时优化 GitHub repo 的 About 描述，提升项目吸引力和转化率。

## What Changes

### 交互性升级
- Hero 区域增加终端打字机动画（模拟真实 REPL 交互）
- 功能卡片 hover 动效 + 点击展开详情
- 代码块一键复制 + 语言切换
- 滚动触发动画（fade-in, slide-up）

### 设计性升级
- 自定义 VitePress 主题（非默认主题）
- 暗色/亮色主题深度定制（非简单切换）
- 品牌色系统（Go 蓝 #00ADD8 + 辅助色）
- 响应式布局优化（移动端体验）

### 引导性升级
- 新用户引导流程（3 步快速开始）
- 交互式架构图（可点击模块查看详情）
- 使用场景卡片（选择你的角色：开发者/架构师/学生）
- "为什么选 Go" 对比表格（vs Python/Rust/TS）

### 完整性升级
- API 参考文档（所有工具、命令、配置项）
- 故障排查指南（FAQ + 常见错误）
- 版本发布日志（CHANGELOG 自动生成）
- 贡献者指南（如何开发、测试、提交 PR）

### GitHub Repo About
- 优化 About 描述，提升点击转化率
- 添加 Topics 标签，增加搜索曝光
- 添加网站链接、徽章

## Capabilities

### New Capabilities
- `interactive-hero`: 终端打字机动画 + 模拟 REPL 交互
- `custom-theme`: 自定义 VitePress 主题
- `user-onboarding`: 新用户引导流程
- `architecture-diagram`: 交互式架构图
- `api-reference`: 完整 API 参考文档
- `troubleshooting`: 故障排查指南

### Modified Capabilities
- `docs-website`: 官网全面升级（交互/设计/引导/完整）
- `github-about`: GitHub repo About 优化

## Impact

- 修改文件: `docs/.vitepress/` 主题配置和组件
- 修改文件: `docs/index.md` 首页重构
- 新增文件: `docs/api/`, `docs/troubleshooting/`, `docs/contributing/`
- 修改文件: GitHub repo About 和 Topics
