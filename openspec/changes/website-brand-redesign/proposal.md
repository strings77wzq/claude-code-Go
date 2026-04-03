## Why

官网存在两个严重问题：1) 中文页面导航栏仍显示英文（Guide/Architecture/Core Code/Tools/MCP），侧边栏中文但导航栏不跟随语言切换——明显的产品级 bug；2) 官网设计缺乏品牌差异化，当前首页结构抄袭 schhaohao/docs 的功能卡片布局，没有体现 claude-code-Go 的独特定位（"模型提供智能，harness 提供可靠性"）。参考 opencode.ai、openspec.dev、claude.com/docs 的优秀设计，需要全面重塑品牌。

## What Changes

- **BUG FIX**: 导航栏按 locale 分离中英文，中文页面显示中文导航
- 首页重构：融合 opencode.ai 的 Hero + 安装命令、openspec.dev 的徽章式核心优势、claude.com 的手风琴功能展示
- 新增"架构理念"板块：展示 Model + Harness 的设计哲学
- 新增"为什么选 Go"板块：对比 Python/Rust 实现的优势
- 新增"快速开始"多标签板块：go install / 源码编译 / 预编译二进制
- 品牌色调统一：Go 蓝 #00ADD8 为主色，暗色主题背景
- 页脚增加中英文翻译

## Capabilities

### New Capabilities
- `nav-localization`: 导航栏多语言支持，nav 数组按 locale 分离
- `hero-redesign`: Hero 区域重构，安装命令组件，CTA 按钮
- `architecture-philosophy`: 架构理念展示板块（Model + Harness）
- `why-go-section`: 为什么选 Go 板块
- `quickstart-tabs`: 多标签快速开始板块

### Modified Capabilities
- `docs-website`: 首页全面重构，导航栏 bug 修复，品牌色调统一

## Impact

- 修改文件: docs/.vitepress/config.ts（nav 按 locale 分离）
- 修改文件: docs/index.md, docs/zh/index.md（首页全面重构）
- 新增文件: docs/.vitepress/theme/CustomHome.vue（自定义首页组件）
- 不影响: Go 源代码、Python Harness、文档内容页面
