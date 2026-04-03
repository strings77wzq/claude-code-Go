## Why

当前官网文档表述偏平淡，不够"高级"。网页视觉效果太"干净"，缺少 hack 风格的特效（终端窗口效果、代码动画、打字机效果等）。作为面向开发者的 AI 编程工具，官网应该体现技术感和 hacker 文化。

## What Changes

### 文档表述升级
- 架构概览：用更精准的工程语言，增加架构图和状态机图
- 设计理念：用类比和对比加深理解（vs 传统脚本、vs 全栈框架）
- 核心代码文档：增加"为什么这样设计"的深度解释

### 官网视觉特效
- 首页 Hero 区域增加终端窗口效果（模拟 REPL 交互动画）
- 代码块增加打字机效果（typewriter animation）
- 功能卡片增加 hover 动画效果
- 增加 ASCII art 风格的架构图
- 整体色调调整为暗色主题（hacker 风格）

### VitePress 主题定制
- 自定义 CSS 变量实现暗色主题
- 增加终端风格的代码块样式
- 增加页面加载动画
- 增加滚动触发的淡入动画

## Capabilities

### New Capabilities
- `terminal-animation`: 终端窗口效果 + REPL 动画
- `typewriter-effect`: 代码打字机效果
- `scroll-animations`: 滚动触发淡入动画
- `dark-theme`: 暗色主题默认启用

### Modified Capabilities
- `docs-website`: 文档表述升级，视觉特效增强

## Impact

- 修改文件: docs/.vitepress/theme/custom.css, docs/index.md, docs/zh/index.md
- 新增文件: docs/.vitepress/theme/components/ (动画组件)
- 重写文档: 架构概览、设计理念、核心代码文档
- 不影响: Go 源代码、Python Harness、API
