# API 错误处理示例

本文档展示如何在 claude-code-Go 中处理 API 错误。

## 常见 API 错误类型

| 错误码 | 原因 | 解决方案 |
|--------|------|----------|
| `rate_limit_exceeded` | 请求过于频繁 | 等待后重试 |
| `invalid_api_key` | API 密钥错误 | 检查配置 |
| `timeout` | 请求超时 | 增加超时时间或简化请求 |
| `context_length_exceeded` | 上下文过长 | 压缩会话 |
| `model_unavailable` | 模型暂时不可用 | 尝试其他模型 |

## 错误处理示例代码

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

func main() {
	client := api.NewClient(api.Config{
		APIKey:  "sk-test",
		Model:   "claude-sonnet-4-20250514",
		Timeout: 30 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt API call
	response, err := client.SendMessage(ctx, "Hello")
	if err != nil {
		// Handle specific error types
		if apiErr, ok := err.(*api.Error); ok {
			switch apiErr.Code {
			case "rate_limit_exceeded":
				fmt.Println("Rate limit hit. Waiting before retry...")
				time.Sleep(60 * time.Second)
				// Retry logic here

			case "invalid_api_key":
				fmt.Println("Invalid API key. Please check your configuration.")
				fmt.Println("Run: go-code --setup")

			case "timeout":
				fmt.Println("Request timeout. Try:")
				fmt.Println("1. Check your internet connection")
				fmt.Println("2. Increase timeout: export GO_CODE_TIMEOUT=60s")

			case "context_length_exceeded":
				fmt.Println("Context too long. Try:")
				fmt.Println("1. Start a new session")
				fmt.Println("2. Use a model with larger context window")
				fmt.Println("3. Clear old messages: /compact")

			default:
				fmt.Printf("API Error: %s\n", apiErr.Message)
			}
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
		return
	}

	fmt.Println("Success:", response)
}
```

## 最佳实践

### 1. 分类处理错误

```go
if apiErr, ok := err.(*api.Error); ok {
    switch apiErr.Code {
    case "rate_limit_exceeded":
        time.Sleep(60 * time.Second)
        return retry()
    case "timeout":
        return retryWithLongerTimeout()
    default:
        return fmt.Errorf("API error: %w", err)
    }
}
```

### 2. 用户友好的错误信息

```go
func UserFriendlyError(err error) string {
    if apiErr, ok := err.(*api.Error); ok {
        switch apiErr.Code {
        case "rate_limit_exceeded":
            return "You've hit the rate limit. Please wait a minute before trying again."
        case "timeout":
            return "The request took too long. Try simplifying your question or check your connection."
        }
    }
    return fmt.Sprintf("An error occurred: %v", err)
}
```

### 3. 重试策略

```go
for i := 0; i < 3; i++ {
    result, err := operation()
    if err == nil {
        return result, nil
    }
    if !isTimeout(err) {
        return nil, err // Don't retry non-timeout errors
    }
    time.Sleep(time.Duration(i+1) * time.Second)
}
```