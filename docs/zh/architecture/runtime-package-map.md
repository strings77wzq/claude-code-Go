---
title: 运行时包映射
description: 面向贡献者的 claude-code-Go 运行时行为包级别映射
---

# 运行时包映射

本页面将产品工作流映射到负责实现的 Go 包。在修改运行时行为之前请参考此页面，确保修复落在正确的层次上。

## 运行时流程

```text
cmd/go-code
  -> 配置加载器
  -> 提供者注册表
  -> 工具注册表
  -> 权限策略
  -> Agent 循环
  -> TUI 或 REPL
```

## 包所有权

| 包 | 职责 | 备注 |
| --- | --- | --- |
| `cmd/go-code` | 进程入口、命令行参数、初始化、运行时组装 | 在此保持编排逻辑；避免在 `main` 中嵌入业务逻辑 |
| `pkg/tui` | 默认的 Bubble Tea 交互式 UI | 应将斜杠命令委托给共享的命令层 |
| `pkg/tty` | 传统的基于行的 REPL | 在 TUI 达到功能完备之前，作为命令覆盖范围的参考 |
| `internal/agent` | Agent 循环、历史记录、压缩、恢复、会话跟踪钩子 | 负责停止原因分发和工具结果反馈 |
| `internal/api` | Anthropic 请求/响应类型、HTTP 客户端、SSE 解析、API 错误分类 | 尽可能将提供者中立的逻辑放在此包之外 |
| `internal/provider` | 提供者接口和提供者适配器 | Anthropic 和 OpenAI 兼容的差异在此处理 |
| `internal/config` | 配置加载和默认值 | 应为诊断输出暴露配置来源细节 |
| `internal/tool` | 工具接口、注册表和内置工具注册 | 工具定义必须保持简洁且符合 schema 验证 |
| `internal/tool/builtin` | 内置本地工具，如 read、edit、bash、tree、diff、web fetch、todo、notebook | 副作用必须由调用方进行权限控制 |
| `internal/permission` | 权限策略、规则、提示器、文件边界、Bash 语义检查 | 此包做决策；agent/工具层必须强制执行 |
| `internal/session` | JSONL 会话持久化和跟踪行辅助函数 | 需要标准化的 schema 以支持回放和调试 |
| `internal/hooks` | 工具执行前后钩子 | 预钩子失败应阻止执行 |
| `internal/skills` | 用户定义的可复用提示/工作流 | 无效的技能文件不应导致启动失败 |
| `internal/tool/mcp` | MCP 客户端、传输层和适配器 | MCP 工具必须进入相同的注册表和权限流程 |
| `internal/lsp` | LSP 客户端功能 | 可选功能；未配置服务器时应安全关闭 |
| `internal/update` | 版本检查与二进制更新 | 仅限用户触发的网络行为 |
| `internal/cost` | Token 和成本跟踪辅助 | 提供者定价假设必须有文档记录 |
| `internal/telemetry` | 可选的本地遥测原语 | 在产品/隐私策略明确前保持禁用 |

## 贡献者规则

- UI 包不应重复命令语义。
- 工具实现应执行本地验证，但权限批准应高于工具执行层。
- 提供者适配器应在错误到达 UI 代码之前将其规范化。
- 需要 socket、家目录写入或构建二进制文件的测试必须记录这些要求。
- 公开文档只能声明经过测试验证或在 `PARITY.md` 中列出的工作流。

> 英文版本: [Runtime Package Map](/architecture/runtime-package-map)
