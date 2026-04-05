package anthropic

import (
	"context"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider"
)

type AnthropicProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *api.Client
}

func NewProvider(apiKey, baseURL, model string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  api.NewClient(apiKey, baseURL, model, false),
	}
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

func (p *AnthropicProvider) SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error) {
	if req.Model == "" {
		req.Model = p.model
	}
	return p.client.SendMessage(ctx, req)
}

func (p *AnthropicProvider) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	if req.Model == "" {
		req.Model = p.model
	}
	return p.client.SendMessageStream(ctx, req, onTextDelta)
}

// Ensure AnthropicProvider implements the Provider interface
var _ provider.Provider = (*AnthropicProvider)(nil)
