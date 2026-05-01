## Why

项目核心功能已完成（~90%），但距离"开源社区成熟项目"还有差距。需要系统性地补齐安全验证、测试覆盖、开源基础设施，以及中期的高级功能（LSP、NotebookEdit、自动恢复），对标 Claw Code 和 Claude Code 的工业级水准。

## What Changes

### 短期（1-2 周）
1. **Bash 命令语义验证** — 对标 Claw Code 的 bash_validation.rs（1004 LOC），实现命令语义分析、只读验证、破坏性命令警告、路径验证、命令语义分析
2. **测试覆盖到 80%** — hooks/skills/logger/tui/mcp/provider 全部补充测试
3. **Release 流程** — GoReleaser 配置、自动打 tag、多平台二进制发布、Homebrew tap
4. **CONTRIBUTING/SECURITY** — 贡献指南、代码规范、PR 模板、安全策略、漏洞报告流程

### 中期（1-2 月）
1. **LSP 集成** — Language Server Protocol 支持，symbols/references/diagnostics/definition/hover 5 种操作
2. **NotebookEdit** — Jupyter notebook 编辑支持
3. **自动恢复** — Agent 崩溃后自动重试，recovery recipes，stale branch 检测

## Capabilities

### New Capabilities
- `bash-semantic-validation`: Bash 命令语义级验证（对标 Claw 1004 LOC）
- `test-coverage-80`: 测试覆盖提升到 80%
- `release-automation`: GoReleaser + 自动发布
- `community-infrastructure`: CONTRIBUTING + SECURITY + CODE_OF_CONDUCT
- `lsp-integration`: LSP 协议支持
- `notebook-edit`: Jupyter notebook 编辑
- `auto-recovery`: Agent 自动恢复机制

## Impact

- 新增文件: `internal/permission/bash_semantic.go`, `.goreleaser.yml`, `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `.github/PULL_REQUEST_TEMPLATE.md`, `internal/lsp/`, `internal/tool/builtin/notebook.go`, `internal/agent/recovery.go`
- 修改文件: 所有缺测试的模块
- 不影响: 核心 Agent Loop 逻辑、API Client
