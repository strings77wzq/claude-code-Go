// Package tool provides the tool interface and implementations for the Claude Code clone.
package tool

import "context"

// Tool represents an executable tool that can be called by the agent.
type Tool interface {
	// Name returns the unique name of the tool.
	Name() string

	// Description returns a human-readable description of what the tool does.
	Description() string

	// InputSchema returns the JSON schema for the tool's input parameters.
	InputSchema() map[string]any

	// RequiresPermission returns true if the tool requires special permissions.
	RequiresPermission() bool

	// Execute runs the tool with the given input and returns a result.
	Execute(ctx context.Context, input map[string]any) Result
}

// Result represents the output of a tool execution.
type Result struct {
	Content string
	IsError bool
}

// Success creates a successful result.
func Success(content string) Result {
	return Result{
		Content: content,
		IsError: false,
	}
}

// Error creates an error result.
func Error(msg string) Result {
	return Result{
		Content: msg,
		IsError: true,
	}
}

// ToolDefinition represents a tool's definition for API responses.
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}
