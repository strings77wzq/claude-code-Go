## Context

项目 `go-code` 是 Claude Code 的 Go+Python 仿写实现。当前存在三类问题：
1. **致命级**：`go.mod` module path 为占位符 `github.com/user/go-code`，导致 `go build` / `go install` 全部失败
2. **严重级**：文档中 config/session 路径有 4 种不同写法，但代码实际只使用 `~/.go-code/`
3. **质量级**：安装体验不够优雅，缺少 Windows 支持和一键安装脚本

## Goals / Non-Goals

**Goals:**
- 让新用户能在 30 秒内成功安装并运行 go-code
- 所有文档中的路径、命令、URL 完全一致且可执行
- 支持 Linux (amd64)、macOS (amd64/arm64)、Windows (amd64)
- 安装体验对标 oh-my-openagent 的优雅程度

**Non-Goals:**
- 不改变现有代码的 config/session 实现逻辑（文档适配代码，而非代码适配文档）
- 不引入新的外部依赖
- 不改变二进制名称（保持 `go-code`）

## Decisions

### Decision 1: Module Path 统一为 `github.com/strings77wzq/claude-code-Go`

**Why**: GitHub 仓库 URL 是 `strings77wzq/claude-code-Go`，module path 必须匹配才能让 `go install` 工作。

**Alternatives considered**:
- 改为 `github.com/strings77wzq/go-code` — 但 repo 名已定为 `claude-code-Go`，改 repo 名影响所有已有链接
- 保持 `github.com/user/go-code` — 无法被外部用户 `go install`

### Decision 2: 文档 Config 路径统一为 `~/.go-code/`

**Why**: `internal/config/loader.go:17` 定义 `configDirName = ".go-code"`，实际路径是 `~/.go-code/settings.json`。文档必须与代码一致。

**Alternatives considered**:
- 改代码为 `~/.config/go-code/` — 更大改动，引入回归风险
- 保持文档现状 — 用户照着做会失败

### Decision 3: Config 格式文档统一为 JSON

**Why**: 代码使用 `encoding/json` 解析 `settings.json`，不支持 YAML。README 和 quick-start 中提到的 `config.yaml` 是错误的。

### Decision 4: 安装脚本参考 oh-my-openagent 的 "For Humans" 模式

**Why**: oh-my-openagent 的安装体验极佳——一条 curl 命令搞定，同时提供多种方式。我们采用相同策略：
- 一键安装脚本（curl pipe bash）
- `go install` 方式
- 源码编译方式
- 预编译二进制下载

### Decision 5: 项目对外品牌为 `go-code`，副标题说明是 Claude Code 仿写

**Why**: 用户说项目就叫 go-code。副标题 "Claude Code 的 Go+Python 仿写" 清晰传达项目定位。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| `go install` 需要 repo 有 release tag | 提供 `@latest` 和 `@main` 两种用法说明 |
| Windows 预编译二进制需要实际构建 | CI 增加 windows/amd64 构建，release 时附带 |
| install.sh 脚本需要维护 | 保持脚本简单：检测 OS → 下载对应二进制 → 放入 PATH |
| 全局 import 路径替换可能遗漏 | 使用 `grep -r` 验证零残留 `github.com/user/go-code` |
