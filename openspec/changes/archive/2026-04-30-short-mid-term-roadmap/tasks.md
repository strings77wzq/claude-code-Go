## 短期任务（1-2 周）

### 1. Bash 命令语义验证

- [x] 1.1 创建 `internal/permission/bash_semantic.go`（语义验证层）
- [x] 1.2 只读命令验证（确认不修改文件系统）
- [x] 1.3 破坏性命令检测和警告（rm, mv, cp 等）
- [x] 1.4 sed/awk 写入操作验证
- [x] 1.5 路径验证（工作区边界 + symlink 解析）
- [x] 1.6 命令语义分析（管道组合、重定向、子 shell）

### 2. 测试覆盖到 80%

- [x] 2.1 hooks 模块测试（80% 覆盖）
- [x] 2.2 skills 模块测试（80% 覆盖）
- [x] 2.3 logger 模块测试（80% 覆盖）
- [x] 2.4 tui 模块测试（70% 覆盖）
- [x] 2.5 mcp 模块测试（80% 覆盖）
- [x] 2.6 provider 模块测试（80% 覆盖）
- [x] 2.7 cost 模块测试（80% 覆盖）

### 3. Release 流程

- [x] 3.1 创建 `.goreleaser.yml` 配置
- [x] 3.2 更新 GitHub Actions 添加 release 工作流
- [x] 3.3 创建 Homebrew tap 配置

### 4. 开源社区基础设施

- [x] 4.1 创建 `CONTRIBUTING.md`（贡献指南）
- [x] 4.2 创建 `SECURITY.md`（安全策略）
- [x] 4.3 创建 `CODE_OF_CONDUCT.md`（行为准则）
- [x] 4.4 创建 `.github/PULL_REQUEST_TEMPLATE.md`（PR 模板）

## 中期任务（1-2 月）

### 5. LSP 集成

- [x] 5.1 创建 `internal/lsp/client.go`（LSP 客户端）
- [x] 5.2 实现 workspace/symbols 操作
- [x] 5.3 实现 textDocument/references 操作
- [x] 5.4 实现 textDocument/publishDiagnostics 操作
- [x] 5.5 实现 textDocument/definition 操作
- [x] 5.6 实现 textDocument/hover 操作

### 6. NotebookEdit

- [x] 6.1 创建 `internal/tool/builtin/notebook.go`
- [x] 6.2 支持 .ipynb 文件读写
- [x] 6.3 支持 cell 内容编辑
- [x] 6.4 支持添加/删除 cell
- [x] 6.5 支持执行 cell

### 7. 自动恢复

- [x] 7.1 创建 `internal/agent/recovery.go`（恢复机制）
- [x] 7.2 API timeout 恢复（指数退避重试）
- [x] 7.3 rate limit 恢复（等待后重试）
- [x] 7.4 tool error 恢复（报告并继续）
- [x] 7.5 context full 恢复（压缩后重试）
