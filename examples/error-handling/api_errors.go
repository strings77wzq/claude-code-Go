package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// APIErrorHandlingExample demonstrates how to handle API errors
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
