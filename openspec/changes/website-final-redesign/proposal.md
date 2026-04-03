## Why

当前官网 features 区块 6 个卡片并排，文字密度太高，视觉上像 README 而非产品官网。参考 ohmyopenagent.com、opencode.ai、cursor.com、aider.chat 等 5 个顶级 AI 编码工具官网的最佳实践，需要重构为：每个功能独立大区块、大数字指标、终端动画、标签系统、CTA 底部。

## What Changes

### 首页重构
- Hero 区域优化（安装命令框 + CTA 按钮）
- 核心数据指标放大（ohmyopenagent 风格大数字）
- 功能区块改为独立大区块（aider 风格卡片网格）
- 架构理念区块（ohmyopenagent 风格大区块 + 标签）
- 终端动画演示（CSS 动画 REPL 交互）
- CTA 底部（ohmyopenagent 风格行动号召）

### 视觉优化
- 暗色主题默认启用（GitHub-dark 配色）
- 终端风格代码块
- 滚动触发淡入动画
- 功能卡片 hover 动画

## Capabilities

### New Capabilities
- `hero-redesign`: Hero 区域重构
- `metrics-display`: 核心数据指标大数字展示
- `feature-blocks`: 功能大区块设计
- `terminal-animation`: 终端动画演示
- `cta-footer`: CTA 底部行动号召

### Modified Capabilities
- `docs-website`: 首页全面重构，视觉特效增强

## Impact

- 修改文件: docs/index.md, docs/zh/index.md, docs/.vitepress/theme/custom.css
- 新增文件: 无（纯重构）
- 不影响: Go 源代码、Python Harness、文档内容页面
