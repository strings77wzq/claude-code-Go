## Context

项目核心功能已完成，需要向工业级项目迈进。对标 Claw Code 和 Claude Code。

## Goals / Non-Goals

**Goals:**
- Bash 命令语义验证（对标 Claw 1004 LOC）
- 测试覆盖到 80%
- Release 自动化
- 开源社区基础设施
- LSP 集成
- NotebookEdit
- 自动恢复

**Non-Goals:**
- 多 Agent 协调（长期目标）
- VS Code 扩展（长期目标）
- 插件市场（长期目标）

## Decisions

### 1. Bash 语义验证

**Decision**: 在现有 bash_validation.go 基础上扩展，增加命令语义分析层。

```
现有: 命令分类（ReadOnly/Write/Dangerous）
新增: 语义验证层
  ├─ 只读命令验证（确认不修改文件系统）
  ├─ 破坏性命令警告（rm, mv, cp 等）
  ├─ sed/awk 写入验证（目标路径检查）
  ├─ 路径验证（工作区边界 + symlink 解析）
  └─ 命令语义分析（管道组合、重定向、子 shell）
```

### 2. 测试覆盖

**Decision**: 每个模块补充测试，目标 80% 覆盖率。

| 模块 | 当前 | 目标 | 测试类型 |
|------|:---:|:---:|---------|
| hooks | 0% | 80% | 单元测试 |
| skills | 0% | 80% | 单元测试 |
| logger | 0% | 80% | 单元测试 |
| tui | 0% | 70% | 集成测试 |
| mcp | 0% | 80% | 单元测试 |
| provider | 0% | 80% | 单元测试 |
| cost | 0% | 80% | 单元测试 |

### 3. Release 流程

**Decision**: GoReleaser + GitHub Actions 自动发布。

```yaml
# .goreleaser.yml
builds:
  - goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
release:
  github:
    owner: strings77wzq
    name: claude-code-Go
```

### 4. LSP 集成

**Decision**: 实现 LSP client，支持 5 种操作。

```
internal/lsp/
├── client.go      ← LSP client 连接管理
├── symbols.go     ← workspace/symbols, textDocument/documentSymbol
├── references.go  ← textDocument/references
├── diagnostics.go ← textDocument/publishDiagnostics
├── definition.go  ← textDocument/definition
└── hover.go       ← textDocument/hover
```

### 5. NotebookEdit

**Decision**: 新增 NotebookEdit 工具，支持 Jupyter notebook 编辑。

```go
// 支持 .ipynb 文件
// 编辑 cell 内容
// 添加/删除 cell
// 执行 cell
```

### 6. 自动恢复

**Decision**: Agent 崩溃后自动重试，基于错误类型的 recovery recipes。

```
recovery_recipes:
  - API timeout    → retry with backoff
  - rate limit     → wait and retry
  - tool error     → report and continue
  - context full   → compact and retry
```

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| Bash 语义验证可能误判 | 提供配置覆盖机制 |
| LSP 集成增加复杂度 | 作为可选功能，默认关闭 |
| NotebookEdit 使用场景少 | 先实现基础功能，后续迭代 |
