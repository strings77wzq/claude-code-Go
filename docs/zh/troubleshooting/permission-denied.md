---
title: 权限被拒绝
description: 理解并解决 claude-code-Go 的权限错误
---

# 权限被拒绝

理解并解决权限错误。

## 默认安全模型

`go-code` 默认以 `WorkspaceWrite` 模式启动。此模式允许读取工作区文件而不提示，并在写入、编辑、Bash 执行、网络访问或其他副作用操作之前请求确认。

当出现提示时:

- 选择 `y` 进行一次性批准。
- 选择 `a` 在当前会话中记住相同的操作。
- 选择 `n` 拒绝操作。工具不会执行，agent 会收到错误结果。

权限决策会记录到会话跟踪中，包含工具名称、决策、摘要和时间戳。密钥和完整工具负载不会存储在权限审计条目中。

## 常见场景

### 读取敏感文件

```
> Read .env

❌ 权限被拒绝: .env 匹配被阻止的模式
```

**解决方案**:
1. 一次性访问: `/allow read .env`
2. 模式访问: `/allow read *.env`
3. 切换模式: `/mode ReadOnly`（可读取任何内容）

### 写入系统目录

```
> Write /etc/config "data"

❌ 权限被拒绝: /etc/* 被阻止
```

**解决方案**:
1. 改用用户目录
2. 切换到 DangerFullAccess 模式（不推荐）
3. 在 claude-code-Go 外部手动使用 sudo

### 执行危险命令

```
> Bash rm -rf /

⚠️ 需要权限: 需要 DangerFullAccess 模式
```

**解决方案**:
1. 确认您确实要执行此操作
2. 使用 `/allow` 进行一次性批准
3. 如需自动化，切换模式

## 权限级别

| 级别 | 可读取 | 可写入 | 可执行 |
|------|--------|--------|--------|
| ReadOnly | ✅ 全部 | ❌ 无 | ❌ 无 |
| WorkspaceWrite | ✅ 全部 | ✅ 工作区 | ✅ 安全命令 |
| DangerFullAccess | ✅ 全部 | ✅ 全部 | ✅ 全部 |

## 自定义规则

添加到 `~/.go-code/settings.json`:

```json
{
  "rules": [
    {"pattern": "*.secret", "allowed": false},
    {"pattern": "docs/*", "allowed": true},
    {"pattern": "*.tmp", "allowed": true}
  ]
}
```

## 会话记忆

记住权限设置:

```
> /remember allow read *.log
> /remember allow bash go test
```

## 调试

查看当前权限:

```
> /mode
当前模式: WorkspaceWrite

> /rules
当前规则:
- *.env → 拒绝
- *.go → 允许
- * → 询问
```

> 英文版本: [Permission Denied](/troubleshooting/permission-denied.md)
