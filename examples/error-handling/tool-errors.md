# 工具执行错误处理示例

本文档展示如何在 claude-code-Go 中处理工具执行错误。

## 常见工具错误类型

| 错误类型 | 原因 | 解决方案 |
|----------|------|----------|
| `ErrCommandFailed` | 命令返回非零退出码 | 检查命令语法和资源 |
| `ErrTimeout` | 工具执行超时 | 使用更简单的命令或增加超时 |
| `ErrInvalidInput` | 参数缺失或无效 | 检查工具参数要求 |
| `ErrNotFound` | 目标文件/资源不存在 | 先检查文件是否存在 |

## 工具错误处理示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func main() {
	// Create tool registry
	registry := tool.NewRegistry()

	// Create execution context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to execute a tool
	result, err := registry.Execute(ctx, "Bash", map[string]interface{}{
		"command": "ls -la /nonexistent",
	})

	if err != nil {
		if toolErr, ok := err.(*tool.ExecutionError); ok {
			fmt.Printf("Tool execution failed: %s\n", toolErr.Message)

			switch toolErr.Type {
			case tool.ErrCommandFailed:
				fmt.Println("\nThe command returned a non-zero exit code.")
				fmt.Println("Check the command syntax and ensure resources exist.")

			case tool.ErrTimeout:
				fmt.Println("\nCommand timed out.")
				fmt.Println("Tips:")
				fmt.Println("- Use simpler commands")
				fmt.Println("- Add timeout: /timeout 30s")
				fmt.Println("- Run long tasks in background")

			case tool.ErrInvalidInput:
				fmt.Println("\nInvalid tool input.")
				fmt.Println("Check the tool's required parameters:")
				fmt.Printf("  %s\n", toolErr.ToolSchema)
			}
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
		return
	}

	fmt.Println("Result:", result)
}
```

## 最佳实践

### 1. 始终检查工具结果

```go
result, err := registry.Execute(ctx, toolName, input)
if err != nil {
    // 处理错误
    return fmt.Errorf("tool %s failed: %w", toolName, err)
}
// 使用结果
processResult(result)
```

### 2. 使用合适的超时

```go
ctx, cancel := context.WithTimeout(context.Background(), appropriateTimeout)
defer cancel()
```

### 3. 验证输入参数

```go
if err := validateInput(input); err != nil {
    return fmt.Errorf("invalid input: %w", err)
}
```