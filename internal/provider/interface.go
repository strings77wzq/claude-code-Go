// Package provider defines the interface for LLM providers.
package provider

import (
	"context"
)

// Request represents a request to an LLM provider.
type Request struct {
	Model     string
	MaxTokens int
	System    string
	Stream    bool
	Messages  []Message
	Tools     []ToolDefinition
}

// Message represents a chat message.
type Message struct {
	Role    string
	Content string
}

// ToolDefinition defines a tool that can be called.
type ToolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any
}

// Response represents a response from an LLM provider.
type Response struct {
	ID         string
	Content    string
	StopReason string
	Usage      Usage
}

// Usage represents token usage statistics.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// Provider is the interface that LLM providers must implement.
type Provider interface {
	// Name returns the provider name.
	Name() string

	// DefaultModel returns the default model for this provider.
	DefaultModel() string

	// SendMessage sends a non-streaming message to the provider.
	SendMessage(ctx context.Context, req *Request) (*Response, error)

	// SendMessageStream sends a streaming message to the provider.
	// onTextDelta is called for each text delta received.
	SendMessageStream(ctx context.Context, req *Request, onTextDelta func(text string)) (*Response, error)
}
