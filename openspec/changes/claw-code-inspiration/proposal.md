## Why

Claw Code (Rust, 166K stars) 和 claude-code-Go 共享同一句哲学口号——"模型提供智能，Harness 提供可靠性"，但 Claw Code 的 Harness 是一个完整的自主开发操作系统，而 claude-code-Go 的 Harness 只是进程内的安全网。通过严格对比，发现 6 个高价值、可落地的改进方向，可以让 claude-code-Go 从"教学级项目"向"工业级项目"迈进。

## What Changes

### 1. Bash 命令深度验证
- 危险命令检测和警告（rm -rf, curl | bash, sudo 等）
- 只读命令白名单（ls, cat, grep, find 等无需审批）
- 路径注入防护（防止命令参数中的路径遍历）
- sed/awk 写入操作验证

### 2. 文件边界守卫
- 二进制文件检测（防止意外修改二进制文件）
- 文件大小限制（防止写入超大文件）
- Symlink escape 检测（防止软链接跳出工作区）
- 工作区边界强制校验

### 3. TodoWrite 工具
- LLM 可主动创建/更新/完成任务列表
- 用户可见的任务进度追踪
- 支持多任务并行标记

### 4. 成本追踪
- 每次 API 调用的 token 用量和费用估算
- 会话级别的费用汇总
- 按模型定价表自动计算

### 5. PermissionEnforcer 独立模块
- 将权限检查从 Agent Loop 中解耦
- 前置拦截 + 输入级验证
- 工具级权限标注（每个 Tool 声明所需权限级别）

### 6. 工具描述约束
- 工具描述限制在 250 字符以内
- 强制要求 InputSchema 完整定义
- 工具注册时自动校验

## Capabilities

### New Capabilities
- `bash-validation`: Bash 命令深度验证（危险检测、只读白名单、路径防护）
- `file-boundary`: 文件边界守卫（二进制检测、大小限制、symlink escape）
- `todo-tool`: TodoWrite 工具
- `cost-tracking`: 成本追踪和费用估算
- `permission-enforcer`: 独立权限执行器模块

### Modified Capabilities
- `tool-system`: 工具级权限标注、描述约束校验
- `agent-loop`: 集成 PermissionEnforcer

## Impact

- 新增文件: `internal/tool/builtin/todo.go`, `internal/permission/enforcer.go`, `internal/permission/bash_validation.go`, `internal/permission/file_boundary.go`, `internal/cost/`
- 修改文件: `internal/tool/tool.go`（增加 RequiredPermissionLevel 方法）
- 修改文件: `internal/agent/loop.go`（集成 PermissionEnforcer）
- 不影响: API Client、MCP、Session、TUI
