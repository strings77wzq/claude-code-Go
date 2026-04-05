## Why

当前项目只支持 Anthropic 格式的 API，无法直接使用 DeepSeek、Qwen、GLM、ChatGPT、Gemini 等主流模型。默认模型版本也停留在 claude-sonnet-4，而 Anthropic 已发布 Opus 4.6 和 Sonnet 4.6。同时项目根目录杂乱（32 个条目），包含开发残留文件和编译产物，影响源码阅读体验。

## What Changes

### P0: 多 Provider 适配层
- 创建 Provider 接口，统一不同 API 格式
- 实现 Anthropic Provider（现有代码迁移）
- 实现 OpenAI-compatible Provider（支持 DeepSeek/Qwen/GLM/ChatGPT/腾讯云等）
- 更新默认模型为 claude-sonnet-4-6-20251001
- 添加模型列表和定价表

### P1: 代码组织架构优化
- 创建 scripts/ 目录，移动所有脚本
- 删除开发残留（claw-code-parity/, owncode-analysis/, test/）
- 删除编译产物（go-code, bin/）
- 更新 .gitignore

### P2: Harness 文档补充
- 创建 harness/README.md 说明 Harness 的作用
- 说明为什么用 Python（测试生态、Mock Server、数据分析）
- 添加使用指南

## Capabilities

### New Capabilities
- `multi-provider`: 多模型 Provider 适配层
- `openai-provider`: OpenAI 格式 Provider
- `model-registry`: 模型注册表和定价表
- `scripts-directory`: 脚本目录重组
- `harness-docs`: Harness 文档

### Modified Capabilities
- `api-client`: 重构为 Provider 模式
- `config`: 新增 provider 配置选项

## Impact

- 新增文件: `internal/provider/openai/`, `internal/provider/registry/`, `scripts/`, `harness/README.md`
- 修改文件: `internal/api/client.go` → 迁移到 `internal/provider/anthropic/`
- 修改文件: `cmd/go-code/main.go`（Provider 初始化逻辑）
- 删除文件: 开发残留目录和编译产物
- 不影响: Agent Loop、工具系统、权限系统
