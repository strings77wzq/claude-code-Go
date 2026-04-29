package provider

import (
	"context"
	"sync"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

type ProviderFactory func(model string) Provider

type ApiClientAdapter struct {
	mu       sync.RWMutex
	provider Provider
	factory  ProviderFactory
}

func NewApiClientAdapter(provider Provider, factories ...ProviderFactory) *ApiClientAdapter {
	var factory ProviderFactory
	if len(factories) > 0 {
		factory = factories[0]
	}
	return &ApiClientAdapter{provider: provider, factory: factory}
}

func (a *ApiClientAdapter) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	a.mu.RLock()
	provider := a.provider
	a.mu.RUnlock()
	return provider.SendMessageStream(ctx, req, onTextDelta)
}

func (a *ApiClientAdapter) SetModel(model string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.factory != nil {
		a.provider = a.factory(model)
		return
	}

	if setter, ok := a.provider.(interface{ SetModel(string) }); ok {
		setter.SetModel(model)
	}
}
