## Context

Claw Code 的 Harness 是操作系统级别的自主开发系统，而 claude-code-Go 的 Harness 是进程内的安全网。对比发现 6 个可落地的改进方向。

## Goals / Non-Goals

**Goals:**
- Bash 命令深度验证（危险检测、只读白名单、路径防护）
- 文件边界守卫（二进制检测、大小限制、symlink escape）
- TodoWrite 工具
- 成本追踪
- PermissionEnforcer 独立模块
- 工具描述约束

**Non-Goals:**
- 不做多 Agent 协调（Task/Worker/Team 系统）
- 不做 LSP 集成
- 不做 OAuth PKCE
- 不做 Cron/定时任务

## Decisions

### 1. Bash 验证策略

**Decision**: 在 Bash 工具执行前增加验证层，而非在 Agent Loop 中验证。

```
Bash 工具 → BashValidator → 分类命令 → 决策
  ├─ 只读命令 (ls, cat, grep, find, wc, head, tail) → 自动允许
  ├─ 危险命令 (rm -rf, curl | bash, sudo, dd) → 拒绝 + 警告
  ├─ 写入命令 (sed -i, awk, tee) → 路径验证 + 权限检查
  └─ 其他命令 → 正常权限审批流程
```

### 2. 文件边界守卫

**Decision**: 在 Write/Edit 工具中集成边界检查。

- 文件扩展名黑名单：`.exe`, `.bin`, `.so`, `.dylib`
- 文件大小限制：10MB
- Symlink 解析后检查是否在工作区内
- 工作区边界：启动时确定，不可更改

### 3. TodoWrite 工具

**Decision**: 简单的内存任务列表，不持久化。

```go
type TodoItem struct {
    ID       int
    Content  string
    Status   string // "pending", "in_progress", "completed"
}
```

### 4. 成本追踪

**Decision**: 按模型定价表估算费用，不依赖 API 返回的实际费用。

```
claude-sonnet-4: $3/MTok input, $15/MTok output
claude-opus-4:   $15/MTok input, $75/MTok output
hunyuan-2.0:     ¥0.5/MTok input, ¥2.0/MTok output
```

### 5. PermissionEnforcer

**Decision**: 独立模块，在 Agent Loop 的 checkPermission 中调用。

```
Agent.checkPermission()
  → PermissionEnforcer.Evaluate()
    → 1. 工具级权限标注检查
    → 2. Bash 命令验证
    → 3. 文件边界检查
    → 4. 原始 Policy 评估
```

### 6. 工具描述约束

**Decision**: 工具注册时自动校验描述长度和 InputSchema 完整性。

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| Bash 验证可能误判 | 提供配置覆盖机制 |
| 文件边界守卫影响正常操作 | 默认宽松，可配置严格模式 |
| 成本追踪不准确 | 标注为估算值，非精确计费 |
