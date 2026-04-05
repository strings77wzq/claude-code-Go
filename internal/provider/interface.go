// Package provider defines the interface for LLM providers.
package provider

import (
	"context"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// Provider is the interface that LLM providers must implement.
type Provider interface {
	// Name returns the provider name.
	Name() string

	// SendMessage sends a non-streaming message to the provider.
	SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error)

	// SendMessageStream sends a streaming message to the provider.
	// onTextDelta is called for each text delta received.
	SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}
