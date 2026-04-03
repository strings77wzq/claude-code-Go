## Why

项目上线后发现多处质量问题影响专业形象：README 中项目名称和 clone URL 与实际仓库不一致，文档命令过时；官网 hero 区域存在视觉冗余（"claude-code-Go / Claude Code in Go" 重复）；官网底部缺少 CTA 和使用场景展示；部分 AI 新技术（Session Persistence、Hooks）未实现。这些问题让项目看起来不够专业，影响作为个人技术 IP 的长期打造。

## What Changes

- **BREAKING**: README 项目名称从 `go-code` 改为 `claude-code-Go`
- 修复 README 中错误的 clone URL（`github.com/user/go-code` → `github.com/strings77wzq/claude-code-Go`）
- 更新 README 中过时的文档命令（mkdocs → vitepress）
- 更新 README 项目结构描述（体现 en/zh 双语结构）
- 优化官网 hero 区域，去除视觉冗余
- 官网首页增加 CTA 底部区域和使用场景展示
- 补充 Session Persistence 实现（JSONL 格式）
- 补充 Hooks 系统（pre/post tool execution）

## Capabilities

### New Capabilities
- `session-persistence`: JSONL 会话持久化，崩溃恢复，会话恢复
- `hooks-system`: 工具执行前后钩子，支持 pre/post hooks

### Modified Capabilities
- `docs-website`: 官网首页 hero 优化，增加 CTA 和使用场景
- `readme-docs`: README 全面修复（名称、URL、命令、结构）

## Impact

- 修改文件: README.md, docs/index.md, docs/zh/index.md
- 新增文件: internal/session/（会话持久化）, internal/hooks/（钩子系统）
- 不影响: Go 核心逻辑（agent loop, tool, api, permission）
