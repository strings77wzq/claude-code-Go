## Why

项目已具备发布条件，但 README 存在多处过时信息（默认模型、项目结构、工具数量、Features），且未创建 v0.1.0 Release。这些问题会影响 fork 者、PR 贡献者和直接使用者的体验。

## What Changes

### README 更新
- 更新默认模型为 claude-sonnet-4-6-20251001
- 更新项目结构（添加 provider/, cost/, lsp/, scripts/ 目录）
- 更新工具数量为 10 个（+TodoWrite）
- 更新 Features（添加 Multi-Provider, LSP Integration, Auto-Recovery）
- 更新 Supported Providers（添加 DeepSeek, Qwen, GLM, OpenAI）

### Release 创建
- 创建 v0.1.0 tag
- 更新 CHANGELOG.md 添加 v0.1.0 发布说明

### 开发者体验优化
- 添加 .github/ISSUE_TEMPLATE/（bug report, feature request）
- 添加 .github/PULL_REQUEST_TEMPLATE.md（已存在，验证完整性）
- 更新 CONTRIBUTING.md 添加 fork + PR 流程说明

## Capabilities

### New Capabilities
- `release-v010`: v0.1.0 正式发布
- `issue-templates`: GitHub Issue 模板
- `contributing-guide`: 贡献者指南更新

### Modified Capabilities
- `readme-docs`: README 全面更新
- `changelog`: CHANGELOG v0.1.0 条目

## Impact

- 修改文件: README.md, CHANGELOG.md, CONTRIBUTING.md
- 新增文件: .github/ISSUE_TEMPLATE/bug_report.md, .github/ISSUE_TEMPLATE/feature_request.md
- 新增 tag: v0.1.0
- 不影响: Go 源代码、Python Harness
