## Why

当前官网文档存在多处问题：1) Roadmap 不准确（Phase 2 已实现仍显示"进行中"）；2) 社区页面太早，应改为反馈渠道；3) 缺少 Skills 技术教程；4) MCP 和 Skills 应合并到"扩展"目录；5) 核心代码文档像代码翻译，没有体现架构思想；6) "模型提供智能，harness 提供可靠性"设计理念没有深度展开。参考 Claude Code 的五大架构优势（Harness 优先工程、上下文窗口管理、分层权限防御、本地优先执行、MCP 通用桥接），需要全面重构官网文档，使其配得上一个资深架构师的作品。

## What Changes

### 文档重构
- **重写架构概览**：融入五大架构优势，深度阐述设计理念
- **新增设计理念页面**：完整展开"模型提供智能，harness 提供可靠性"
- **合并扩展目录**：MCP + Skills + Hooks → 统一"扩展"目录
- **新增 Skills 技术教程**：如何创建、注册、使用自定义 Skills
- **重写核心代码文档**：从代码翻译改为架构思想阐述
- **更新 Roadmap**：Phase 2 标记为已完成，Phase 3 为规划中
- **替换社区为反馈**：简单的 GitHub Issue/Discussion 反馈渠道

### 官网导航重构
```
当前:
├── 指南 (项目简介、快速开始、项目结构)
├── 架构 (概览、Agent 循环、工具)
├── 核心代码 (入口点、Agent Loop 实现)
├── 工具系统 (概览)
├── MCP 集成 (协议详解)
└── 资源 (路线图、社区)

重构后:
├── 指南 (项目简介、快速开始、项目结构)
├── 架构 (概览、设计理念、Agent Loop)
├── 扩展 (Skills 教程、MCP 协议、Hooks 系统)
├── 工具 (概览、内置工具)
└── 资源 (路线图、反馈)
```

## Capabilities

### New Capabilities
- `design-philosophy-doc`: 设计理念文档，深度展开核心哲学
- `skills-tutorial`: Skills 技术教程
- `hooks-doc`: Hooks 系统文档
- `feedback-page`: 反馈渠道页面（替换社区）
- `architecture-revamp`: 架构概览重写（融入五大优势）

### Modified Capabilities
- `docs-website`: 导航重构、核心代码文档重写、Roadmap 更新
- `mcp-integration`: 从架构目录迁移到扩展目录

## Impact

- 新增文件: 设计理念、Skills 教程、Hooks 文档、反馈页面
- 重写文件: 架构概览、核心代码文档、Roadmap
- 移动文件: MCP 文档从 architecture/ 到 extension/
- 删除文件: 社区页面、旧的核心代码文档
- 不影响: Go 源代码、Python Harness
