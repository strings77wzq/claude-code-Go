# 权限错误处理示例

本文档展示如何在 claude-code-Go 中处理权限错误。

## 权限系统概述

claude-code-Go 使用三层权限模型：
- **ReadOnly**: 只允许读取操作
- **WorkspaceWrite**: 允许在工作区内写入
- **DangerFullAccess**: 完全访问权限（包括系统命令）

## 权限错误处理示例

```go
package main

import (
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

func main() {
	// Setup permission enforcer
	enforcer := permission.NewEnforcer(
		permission.ModeWorkspaceWrite,
		[]permission.GlobRule{
			{Pattern: "*.go", Allowed: true},
			{Pattern: "*.md", Allowed: true},
			{Pattern: "*.env", Allowed: false}, // Sensitive files
		},
	)

	// Try to read a file
	target := "config.env"
	if err := enforcer.CheckRead(target); err != nil {
		if permErr, ok := err.(*permission.Error); ok {
			fmt.Printf("Permission denied: %s\n", permErr.Message)
			fmt.Println("\nTo allow this action:")
			fmt.Println("1. Add an exception:")
			fmt.Printf("   /allow read %s\n", target)
			fmt.Println("\n2. Switch to a more permissive mode:")
			fmt.Println("   /mode DangerFullAccess")
			fmt.Println("   ⚠️  WARNING: This allows all actions!")
			fmt.Println("\n3. Update your glob rules in ~/.go-code/settings.json")
		}
		return
	}

	fmt.Println("Permission granted!")
}
```

## 处理权限错误的最佳实践

### 1. 清晰告知用户

```go
if permErr, ok := err.(*permission.Error); ok {
    fmt.Println("Permission denied:", permErr.Resource)
    fmt.Println("Run '/allow read", permErr.Resource, "' to grant access")
}
```

### 2. 提供替代方案

```go
fmt.Println("\nTo allow this action:")
fmt.Println("1. Add an exception: /allow read <file>")
fmt.Println("2. Switch mode with /mode")
fmt.Println("3. Update glob rules in settings")
```

### 3. 记录权限检查

```go
log.Infow("Permission check",
    "resource", target,
    "mode", currentMode,
    "allowed", err == nil,
)
```