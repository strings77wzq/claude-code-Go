package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// TimeoutErrorExample demonstrates how to handle timeout errors
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
