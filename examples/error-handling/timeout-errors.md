# 超时错误处理示例

本文档展示如何在 claude-code-Go 中处理超时错误。

## 超时处理策略

1. **指数退避重试**
2. **将复杂任务分解为更小的步骤**
3. **使用异步执行处理长时间运行的任务**
4. **在设置中增加超时时间**

## 超时错误处理示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

func main() {
	// Configure with short timeout for demonstration
	client := api.NewClient(api.Config{
		APIKey:  "sk-test",
		Model:   "claude-sonnet-4-20250514",
		Timeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start request with retry logic
	var response string
	var err error

	for retries := 0; retries < 3; retries++ {
		if retries > 0 {
			fmt.Printf("Retry %d/3 after timeout...\n", retries)
			time.Sleep(time.Duration(retries) * 2 * time.Second)
		}

		response, err = client.SendMessage(ctx, "Complex analysis task...")
		if err == nil {
			break
		}

		if apiErr, ok := err.(*api.Error); ok && apiErr.Code == "timeout" {
			fmt.Println("Request timeout detected")

			// Check if context is still valid
			if ctx.Err() == context.DeadlineExceeded {
				// Create new context with extended timeout
				ctx, cancel = context.WithTimeout(context.Background(),
					time.Duration(10+retries*5)*time.Second)
				defer cancel()
			}
			continue
		}

		// Non-timeout error, don't retry
		break
	}

	if err != nil {
		fmt.Println("\nFailed after retries. Suggestions:")
		fmt.Println("1. Check your network connection")
		fmt.Println("2. Use a simpler prompt")
		fmt.Println("3. Increase default timeout:")
		fmt.Println("   export GO_CODE_TIMEOUT=60s")
		fmt.Println("4. Try a different model with faster response")
		return
	}

	fmt.Println("Success:", response)
}
```

## 重试模式

### 指数退避

```go
func retryWithBackoff(operation func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := operation(); err != nil {
            if !isTimeout(err) {
                return err // 非超时错误不重试
            }
            // 指数退避等待
            time.Sleep(time.Duration(1<<i) * time.Second)
            continue
        }
        return nil
    }
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

### 带上下文检查的重试

```go
for retries := 0; retries < maxRetries; retries++ {
    result, err := operation(ctx)
    if err == nil {
        return result, nil
    }
    
    if ctx.Err() == context.DeadlineExceeded {
        // 创建新上下文
        ctx, cancel = context.WithTimeout(context.Background(), newTimeout)
        defer cancel()
    }
}
```

## 用户建议

当遇到超时时，向用户提供以下建议：

1. **检查网络连接**
2. **简化提示词**
3. **增加超时设置**: `export GO_CODE_TIMEOUT=60s`
4. **尝试响应更快的模型**