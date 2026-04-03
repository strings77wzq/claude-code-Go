## Why

项目上线后发现多处致命问题导致用户第一眼就失败：`go.mod` 的 module path 是 `github.com/user/go-code`（占位符），33 个 Go 源文件的 import 路径全部错误，文档中 `go install` 路径大小写不一致，config/session 路径在文档中有 4 种不同写法但代码实际只用 `~/.go-code/`，安装指南缺少 Windows 支持且没有一键安装脚本。这些问题让项目无法被新用户成功安装和运行。

## What Changes

- **BREAKING**: `go.mod` module path 从 `github.com/user/go-code` 改为 `github.com/strings77wzq/claude-code-Go`
- **BREAKING**: 33 个 Go 源文件的 import 路径全局更新
- 统一所有文档中的 config 路径为 `~/.go-code/`（与 `internal/config/loader.go` 一致）
- 统一所有文档中的 session 路径为 `~/.go-code/sessions/`
- 修正 README 和 quick-start 中错误的 config 格式说明（YAML → JSON）
- 重写 README 安装指南：参考 oh-my-openagent 优雅设计，分系统/分方式
- 重写 `docs/guide/installation.md`：完整的多系统安装文档（Linux/macOS/Windows）
- 增加 Windows 交叉编译支持（Makefile + CI）
- 新增 `install.sh` 一键安装脚本
- 修正 `launch.sh` 占位符
- 修正 copyright 年份

## Capabilities

### New Capabilities
- `installation-experience`: 一键安装脚本 + 多系统安装指南 + 优雅的安装体验设计
- `config-path-consistency`: 统一的 config/session 路径文档规范

### Modified Capabilities
- `readme-docs`: README 安装指南、项目名称、clone URL、config 格式全面修正

## Impact

- 33 个 Go 源文件（import 路径变更）
- `go.mod`（module path 变更）
- README.md（安装指南重写）
- `docs/index.md`（Quick Start 修正）
- `docs/guide/installation.md`（完整重写）
- `docs/guide/quick-start.md`（config 路径统一）
- `docs/guide/configuration.md`（config 路径统一）
- `docs/guide/session-management.md`（session 路径统一）
- `docs/architecture/providers.md`（config 路径统一）
- `docs/extension/mcp.md`（MCP config 路径统一）
- `docs/extension/skills.md`（skills 路径统一）
- 所有中文对应文档（zh/ 目录）
- `Makefile`（增加 Windows 构建）
- `.github/workflows/ci.yml`（增加 Windows 构建）
- `launch.sh`（移除占位符）
- `docs/.vitepress/config.ts`（修正 copyright）
- 新增 `install.sh` 安装脚本
