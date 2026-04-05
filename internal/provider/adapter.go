package provider

import (
	"context"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

type ApiClientAdapter struct {
	provider Provider
}

func NewApiClientAdapter(provider Provider) *ApiClientAdapter {
	return &ApiClientAdapter{provider: provider}
}

func (a *ApiClientAdapter) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	return a.provider.SendMessageStream(ctx, req, onTextDelta)
}
