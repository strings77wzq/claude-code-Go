## 1. Bash 命令深度验证

- [x] 1.1 创建 `internal/permission/bash_validation.go`（命令分类 + 验证）
- [x] 1.2 只读命令白名单自动允许
- [x] 1.3 危险命令检测和阻止
- [x] 1.4 路径注入防护

## 2. 文件边界守卫

- [x] 2.1 创建 `internal/permission/file_boundary.go`（边界检查）
- [x] 2.2 二进制文件检测
- [x] 2.3 文件大小限制（10MB）
- [x] 2.4 Symlink escape 检测

## 3. TodoWrite 工具

- [x] 3.1 创建 `internal/tool/builtin/todo.go`
- [x] 3.2 注册到工具系统
- [x] 3.3 TUI 显示任务列表

## 4. 成本追踪

- [x] 4.1 创建 `internal/cost/tracker.go`（定价表 + 费用估算）
- [x] 4.2 集成到 Agent Loop（每次 API 调用后记录）
- [x] 4.3 会话结束时显示费用汇总

## 5. PermissionEnforcer

- [x] 5.1 创建 `internal/permission/enforcer.go`（独立权限执行器）
- [x] 5.2 工具级权限标注（Tool 接口增加 RequiredPermissionLevel 方法）
- [x] 5.3 集成到 Agent Loop（替换原有 checkPermission）

## 6. 工具约束

- [x] 6.1 工具注册时校验描述长度（≤250 字符）
- [x] 6.2 工具注册时校验 InputSchema 完整性
