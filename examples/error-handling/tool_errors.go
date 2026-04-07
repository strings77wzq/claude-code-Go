package main

import (
	"context"
	"fmt"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// ToolExecutionErrorExample demonstrates how to handle tool execution errors
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
