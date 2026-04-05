## 1. P0: 多 Provider 适配层

- [x] 1.1 创建 `internal/provider/interface.go`（Provider 接口定义）
- [x] 1.2 迁移 `internal/api/` → `internal/provider/anthropic/`
- [x] 1.3 创建 `internal/provider/openai/provider.go`（OpenAI 格式 Provider）
- [x] 1.4 创建 `internal/provider/registry/registry.go`（Provider 自动选择）
- [x] 1.5 更新默认模型为 claude-sonnet-4-6-20251001
- [x] 1.6 更新 `/models` 命令显示所有支持的模型
- [x] 1.7 更新 config 支持 provider 配置

## 2. P1: 代码组织架构优化

- [x] 2.1 创建 scripts/ 目录
- [x] 2.2 移动 install.sh, install.ps1, launch.sh 到 scripts/
- [x] 2.3 删除 claw-code-parity/ owncode-analysis/ test/
- [x] 2.4 删除编译产物 go-code, bin/
- [x] 2.5 更新 .gitignore

## 3. P2: Harness 文档

- [x] 3.1 创建 harness/README.md
- [x] 3.2 说明 Harness 的作用和为什么用 Python
- [x] 3.3 添加使用指南
