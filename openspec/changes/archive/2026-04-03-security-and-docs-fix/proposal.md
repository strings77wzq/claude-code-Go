## Why

Code Review 发现项目存在 5 个严重安全问题（路径遍历、命令注入、竞态条件、资源泄漏、Response Body 泄漏）和大量文档翻译问题（Harness 被误译为"驾驭"、Skills 被误译为"技能"、中英文不一致、缺失页面）。这些问题阻碍项目作为精品开源产品发布。

## What Changes

### Phase 1: 安全修复（5 个严重问题）
- 文件路径遍历防护（所有文件工具加 workingDir 校验）
- Bash 工具竞态条件修复（sync.WaitGroup）
- API 客户端 Response Body 泄漏修复
- MCP 进程资源泄漏修复
- 信号处理改为优雅关闭

### Phase 2: 文档修复（翻译 + 一致性）
- 全局替换中文 "驾驭系统" → "Harness"
- 全局替换中文 "技能" → "Skills"
- 统一 Agent/Provider 等术语
- 修复工具数量不一致（6 → 9）
- 补全缺失的中文页面
- 修复 README

## Capabilities

### New Capabilities
- `path-validation`: 文件操作路径验证，防止遍历攻击
- `graceful-shutdown`: 优雅关闭，清理资源后退出

### Modified Capabilities
- `bash-tool`: 修复竞态条件
- `api-client`: 修复 Response Body 泄漏
- `mcp-transport`: 修复进程资源泄漏
- `docs-website`: 修复翻译术语、一致性、缺失页面

## Impact

- 修改文件: 5 个 Go 源文件 + 10+ 个中文文档 + README
- 新增文件: 无（纯修复）
- 不影响: Agent Loop 核心逻辑、工具功能、Provider 系统
