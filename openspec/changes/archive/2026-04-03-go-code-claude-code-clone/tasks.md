## 1. Go 运行时 — 核心基础

- [x] 1.1 项目脚手架：go.mod, 目录结构, cmd/go-code/main.go 入口
- [x] 1.2 配置系统：Config 类型, 多源加载（CLI→env→project→user）, 默认值, API key 验证
- [x] 1.3 API 客户端：HTTP 客户端, 请求/响应类型, 阻塞调用, 429 重试

## 2. Go 运行时 — Agent 核心

- [x] 2.1 SSE 流式解析器：事件解析, text_delta 回调, input_json_delta 累积
- [x] 2.2 Agent Loop：stop_reason 调度, MAX_TURNS 限制, 流式输出回调
- [x] 2.3 历史管理：消息累积, user/assistant 严格交替, 副本保护

## 3. Go 运行时 — 工具系统

- [x] 3.1 Tool 接口 + Registry：接口定义, 注册, 查找, 执行, panic 恢复
- [x] 3.2 Bash 工具：进程执行, goroutine 读输出, 超时, 输出截断, exit code
- [x] 3.3 Read 工具：行号格式, offset/limit, 200KB 限制
- [x] 3.4 Write 工具：文件创建/覆写, 父目录自动创建
- [x] 3.5 Edit 工具：精确字符串替换, 唯一性检查, replace_all
- [x] 3.6 Glob 工具：文件模式匹配, ** 递归支持
- [x] 3.7 Grep 工具：正则内容搜索, include/exclude 过滤

## 4. Go 运行时 — 权限系统

- [x] 4.1 权限模式：ReadOnly/WorkspaceWrite/DangerFullAccess 定义
- [x] 4.2 规则匹配：glob 规则解析, bash(git:*) 格式, 文件路径匹配
- [x] 4.3 交互式审批：终端提示, y/n/a 输入, 会话记忆
- [x] 4.4 Agent Loop 集成：工具执行前权限检查

## 5. Go 运行时 — REPL 交互

- [x] 5.1 基础 REPL：欢迎界面, 输入循环, 特殊命令（/help, /clear, /exit）
- [x] 5.2 流式输出：ANSI 颜色渲染, 工具调用/结果显示
- [x] 5.3 中断处理：Ctrl+C 取消当前请求, Ctrl+D 退出

## 6. Go 运行时 — MCP 集成

- [x] 6.1 Stdio Transport：子进程启动, stdin/stdout JSON-RPC 通信
- [x] 6.2 JSON-RPC 客户端：initialize, tools/list, tools/call
- [x] 6.3 MCP Tool Adapter：远程工具适配为本地 Tool 接口, mcp__ 命名
- [x] 6.4 配置加载：settings.json 中的 mcpServers, 环境变量插值
- [x] 6.5 生命周期管理：优雅关闭, 崩溃隔离

## 7. Go 运行时 — 上下文管理

- [x] 7.1 Token 估算：基于字符数的近似计算
- [x] 7.2 自动压缩：80% 阈值触发, 保留首消息+最近10轮, 中间历史摘要
- [x] 7.3 手动压缩：用户命令触发

## 8. Go 运行时 — 集成与打磨

- [x] 8.1 完整集成：main.go 串联所有组件（config→api→tools→permission→agent→repl）
- [x] 8.2 优雅退出：SIGINT/SIGTERM 处理, MCP 子进程清理
- [x] 8.3 结构化日志：slog 集成, 错误上下文
- [x] 8.4 Go 单元测试：覆盖核心逻辑, `go test ./...` 全通过

## 9. Python Harness — Mock Server

- [x] 9.1 FastAPI Mock Anthropic API：SSE 流式响应, message_start/content_block_delta/message_delta 事件
- [x] 9.2 可配置场景：streaming_text, tool_use, multi-turn 响应序列
- [x] 9.3 请求记录：存储所有请求用于测试断言
- [x] 9.4 服务生命周期：动态端口绑定, 启动/停止

## 10. Python Harness — Parity 测试

- [x] 10.1 streaming_text 场景：验证 SSE 文本完整传递
- [x] 10.2 tool_roundtrip 场景：Read/Bash 工具调用→执行→结果→模型响应
- [x] 10.3 permission_flow 场景：allow/deny 流程验证
- [x] 10.4 mcp_integration 场景：MCP 工具发现与执行
- [x] 10.5 edit_uniqueness 场景：Edit 工具唯一性检查

## 11. Python Harness — 评估与回放

- [x] 11.1 输出质量评估：文本完整性, 工具调用正确性
- [x] 11.2 延迟监控：请求到首 token, 到完成的时间度量
- [x] 11.3 Session 回放：从 JSONL 加载, 按序重现对话
- [x] 11.4 Trace 分析：工具调用模式, 错误模式识别

## 12. 文档与 CI/CD

- [x] 12.1 MkDocs 官网：Home, Guide, Architecture, Harness 页面
- [x] 12.2 README.md：项目介绍, 安装, 快速开始
- [x] 12.3 Makefile：build, test, docs, build-all 目标
- [x] 12.4 GitHub Actions CI：Go 测试, Python harness 测试, 多平台构建
- [x] 12.5 GitHub Pages 部署：文档自动发布
