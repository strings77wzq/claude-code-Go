## Why

Phase 1 完成了核心 Agent Loop 和基础工具系统，但作为一个"好用的智能体"还缺少关键能力。Phase 2 需要补齐 Skills 自定义工作流、多 Provider 支持、Session 恢复 UI、更多内置工具，同时更新官网展示 Phase 2 的新能力，让项目从"能跑通的原型"进化为"社区认可的实用工具"。

## What Changes

### 后端新增
- **Skills 系统**：自定义命令/工作流（如 `/review-pr`、`/deploy`），社区可扩展
- **多 Provider 支持**：OpenAI、Gemini、本地模型，降低使用门槛
- **Session 恢复**：从 JSONL 加载历史，`/resume` 命令继续上次对话
- **更多内置工具**：Diff（查看文件差异）、Tree（目录树）、WebFetch（网页抓取）
- **手动压缩命令**：`/compact` 主动压缩上下文
- **自动更新**：`/update` 检查并下载最新版本

### 前端官网更新
- 首页增加 Phase 2 新功能展示卡片
- 新增 "Roadmap" 页面展示项目演进路线
- 更新功能对比表（Phase 1 vs Phase 2 vs Claude Code 官方）
- 新增 "Community" 页面展示 Skills 生态

## Capabilities

### New Capabilities
- `skills-system`: 自定义命令/工作流系统，社区可扩展
- `multi-provider`: 多模型 Provider 支持（OpenAI、Gemini、本地模型）
- `session-resume`: Session 恢复 UI，`/resume` 命令
- `additional-tools`: Diff、Tree、WebFetch 内置工具
- `manual-compaction`: 手动上下文压缩命令
- `auto-update`: 自动更新机制
- `roadmap-page`: 官网 Roadmap 页面
- `community-page`: 官网 Community 页面

### Modified Capabilities
- `docs-website`: 官网首页更新、新增页面、功能对比表更新

## Impact

- 新增 Go 模块: `internal/skills/`, `internal/provider/`
- 新增文档页面: Roadmap、Community、Skills 指南
- 修改文件: 官网首页、功能特性卡片、对比表
- 不影响: Phase 1 已完成的核心模块
