## Why

Claude Code 是目前最强大的 AI 编程助手之一，但它是闭源的商业产品。本项目旨在用 Go + Python 构建一个功能对等的开源替代品——**模型提供智能，harness 提供可靠性**——让任何人都可以部署、使用、审计自己的 AI 编程助手。

参考三个成熟实现：
- **OwnCode (claude-code-java)**: 28 个 Java 文件，清晰的 Agent Loop + Tool + Permission 架构
- **claw-code (Rust)**: 生产级 Rust 实现，多 Provider、OAuth、Sandbox、Context Compaction
- **claw-code-parity**: 40 个工具 + 67 个命令的完整 Parity Checklist + Mock Harness

## What Changes

- 引入 Go 生产运行时：CLI 交互、Agent Loop、6 个内置工具（Bash/Read/Write/Edit/Glob/Grep）、权限系统、MCP 集成
- 引入 Python Harness 基础设施：Mock Anthropic API 服务、Parity 测试场景、质量评估框架、Session 回放调试
- 引入 MkDocs 官网 + 文档体系
- 引入 Makefile 统一构建/测试/文档部署流程
- 引入完整的 OpenSpec 规范体系（spec-driven schema）

## Capabilities

### New Capabilities

- `cli-repl`: 交互式终端，支持流式输出、特殊命令（/help, /clear, /exit）、ANSI 渲染
- `agent-loop`: 核心 "思考-行动-观察" 循环，stop_reason 驱动调度，MAX_TURNS 安全限制
- `tool-system`: 可扩展工具系统，Tool 接口 + Registry + 6 个内置工具（Bash/Read/Write/Edit/Glob/Grep）
- `permission-system`: 三级权限（ReadOnly/WorkspaceWrite/DangerFullAccess），规则匹配，会话记忆，交互式审批
- `api-client`: Anthropic Messages API 客户端，SSE 流式响应，429 重试，错误处理
- `mcp-integration`: MCP 协议支持，stdio transport，JSON-RPC 客户端，工具自动发现注册
- `session-management`: JSONL 会话持久化，崩溃恢复，会话恢复
- `context-management`: Token 估算，自动压缩，上下文窗口管理
- `config-system`: 多源配置加载（CLI→env→project→user），优先级合并
- `harness-mock-server`: Python Mock Anthropic API 服务，模拟 SSE 流式响应，记录请求断言
- `harness-parity-tests`: Python Parity 测试套件，覆盖 streaming_text, tool_roundtrip, permission_flow, mcp_integration
- `harness-evaluators`: 输出质量评估，工具正确性验证，延迟监控
- `harness-replay`: Session 回放调试，trace 分析
- `docs-website`: MkDocs 官网，入门指南、架构文档、API 参考，GitHub Pages 部署
- `build-ci`: Makefile 统一构建，GitHub Actions CI/CD，多平台编译

### Modified Capabilities

<!-- 无已有 spec，全部新增 -->

## Impact

- **新增 Go 模块**: `github.com/user/go-code`，约 30-40 个源文件，~3000 LOC
- **新增 Python 包**: `harness/`，约 40-50 个文件，~4000 LOC
- **新增文档**: `docs/`，MkDocs 站点
- **依赖**: Go 标准库为主，Python 依赖 FastAPI/pytest/MkDocs
- **部署**: Go 单二进制分发，Python harness 仅 CI/开发使用
