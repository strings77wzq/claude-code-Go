## Context

本项目从零构建一个 Claude Code 开源替代品。已有三个参考实现：
- **OwnCode (Java)**: 28 文件，~3000 LOC，清晰的 Agent Loop + Tool + Permission 架构
- **claw-code (Rust)**: 生产级实现，多 Provider、OAuth、Sandbox、Context Compaction
- **claw-code-parity**: 40 工具 + 67 命令的 Parity Checklist + Mock Harness

技术选型：Go 负责生产运行时（单二进制分发），Python 负责 Harness 基础设施（测试生态强），MkDocs 负责官网（Python 生态，部署简单）。

当前已有部分 Phase 1 代码完成：项目脚手架、配置加载、工具系统（6 个内置工具）、权限系统、Agent Loop、API 客户端、REPL。

## Goals / Non-Goals

**Goals:**
- 构建可发布的 Go CLI 工具，功能对标 Claude Code 核心体验
- 构建 Python Harness 测试基础设施，确保行为对齐（Parity）
- 构建 MkDocs 官网，提供完整文档
- 保持代码库可维护性，Harness 代码量 ≥ 运行时代码量

**Non-Goals:**
- 不实现 OAuth 流程（Phase 2）
- 不实现多 Provider 支持（Phase 2，先只做 Anthropic）
- 不实现 Linux Sandbox/namespace 隔离（Phase 2）
- 不实现 Session 持久化/恢复（Phase 2）
- 不实现 Context Compaction（Phase 2）

## Decisions

### 1. Go + Python 混合架构

**Decision**: Go 负责生产运行时，Python 负责 Harness + 官网。

**Rationale**: Go 单二进制分发适合 CLI 工具；Python 的 pytest/FastAPI/MkDocs 生态远强于 Go 的测试生态。

**Alternatives considered**:
- 纯 Go（包括 Harness）：Go 的测试生态够用但不如 Python 灵活，特别是 Mock HTTP 服务和参数化测试
- 纯 Rust（参考 claw-code-parity）：Rust 学习曲线陡峭，且用户已选定 Go

### 2. 不使用 anthropic-sdk-go，直接实现 HTTP 客户端

**Decision**: 使用 net/http 标准库直接与 Anthropic API 通信。

**Rationale**: 完全控制请求/响应格式，便于理解 SSE 流式协议，减少外部依赖。

**Alternatives considered**:
- 使用 anthropic-sdk-go：减少代码量，但隐藏了协议细节，调试困难

### 3. SSE 解析器自己实现

**Decision**: 自定义 SSE 解析器，不使用外部库。

**Rationale**: SSE 格式简单（data: 行 + \n\n 分隔），自己实现可控且无依赖。

### 4. 工具注册避免 import cycle

**Decision**: 使用 `internal/tool/init/` 子包注册内置工具，避免 tool → builtin → tool 的循环依赖。

**Rationale**: Go 不允许循环导入。builtin 包需要 tool 包的接口定义，tool 包的 registry 需要 builtin 包的工具实例。通过 init 子包打破循环。

### 5. 权限系统：三级模式 + 规则匹配

**Decision**: ReadOnly / WorkspaceWrite / DangerFullAccess 三级模式，支持 glob 规则匹配（如 `bash(git:*)`）。

**Rationale**: 与 Claude Code 行为一致，用户熟悉。glob 规则灵活且易于理解。

### 6. 文件编辑：精确字符串替换

**Decision**: Edit 工具使用 exact string replacement（old_string → new_string），不是行号编辑。

**Rationale**: 更安全——LLM 必须提供足够的上下文来唯一定位目标。行号编辑容易因文件变化而错位。

### 7. Mock Server 用 FastAPI

**Decision**: Python Mock Anthropic API 使用 FastAPI。

**Rationale**: FastAPI 原生支持 SSE（Server-Sent Events），类型提示好，异步性能好。

### 8. 官网用 MkDocs + Material 主题

**Decision**: MkDocs with Material for MkDocs 主题。

**Rationale**: Markdown 写文档，零前端知识即可搭建美观官网。GitHub Pages 免费部署。用户第一次做官网，这是最简路径。

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| SSE 解析边界情况 | 高 | 严格遵循 SSE 规范，Mock Server 覆盖各种边缘情况 |
| Bash 工具安全风险 | 关键 | 权限系统 + 超时 + 输出截断，用户审批每一危险操作 |
| MCP Server 崩溃 | 中 | 进程监控 + 自动重启 + 错误隔离 |
| Token 估算不准确 | 中 | 保守阈值，手动覆盖选项 |
| API 限流 | 中 | 指数退避重试，最多 3 次 |
| Go + Python 双语言维护成本 | 中 | Makefile 统一入口，CI 自动化，职责清晰分离 |
