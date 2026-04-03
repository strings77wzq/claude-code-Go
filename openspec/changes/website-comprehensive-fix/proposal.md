## Why

官网存在多处严重问题影响专业形象：中英文切换 404（致命 bug）、导航结构不完整（缺少项目简介、项目结构、核心代码、MCP 集成页面）、首页使用 emoji 图标缺乏品牌特色、缺少学习路径和项目核心数字展示。参考 schhaohao.github.io/docs/ 的专业水准，需要全面修复官网。

## What Changes

- **BREAKING**: 移动英文文档从 docs/en/ 到 docs/ 根目录（修复中英文切换 404）
- 更新 VitePress 配置 sidebar 路径（/en/ → /）
- 新增"入门指南"导航组：项目简介、快速开始、项目结构（中英双语）
- 新增"核心代码"导航组：入口点、Agent Loop 实现（中英双语）
- 新增"工具系统"导航组：工具概览（中英双语）
- 新增"MCP 集成"导航组：MCP 协议详解（中英双语）
- 首页用 SVG 图标替换 emoji（Go gopher、AI 芯片、盾牌等品牌特色）
- 首页增加"项目核心数字"板块
- 首页增加"学习收获"板块
- 所有页面底部增加"最后更新于"和上下页导航

## Capabilities

### New Capabilities
- `docs-introduction`: 项目简介页面（中英双语），说明项目是什么、技术栈、核心数字、学习收获
- `docs-project-structure`: 项目结构页面（中英双语），六大模块详解、依赖关系
- `docs-core-code`: 核心代码页面（中英双语），入口点解析、Agent Loop 实现细节
- `docs-tools-overview`: 工具系统概览页面（中英双语），Tool 接口设计、注册表模式、扩展指南
- `docs-mcp-integration`: MCP 集成详解页面（中英双语），协议原理、stdio transport、工具发现
- `website-branding`: 品牌化 SVG 图标系统，项目核心数字展示，学习收获板块

### Modified Capabilities
- `docs-website`: 修复中英文切换 404，完善导航结构，更新 sidebar 配置

## Impact

- 移动文件: docs/en/* → docs/ 根目录（约 6 个文件）
- 新增文件: ~12 个文档页面（中英双语）+ ~6 个 SVG 图标
- 修改文件: docs/.vitepress/config.ts, docs/index.md, docs/zh/index.md, CustomHome.vue
- 不影响: Go 源代码、Python Harness
