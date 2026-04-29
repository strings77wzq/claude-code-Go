package provider

import (
	"context"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

type fakeProvider struct {
	name  string
	model string
}

func (p *fakeProvider) Name() string { return p.name }

func (p *fakeProvider) SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error) {
	return &api.ApiResponse{Model: p.model}, nil
}

func (p *fakeProvider) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	return &api.ApiResponse{Model: p.model}, nil
}

func (p *fakeProvider) SetModel(model string) {
	p.model = model
}

func TestApiClientAdapterSetModelUsesFactory(t *testing.T) {
	adapter := NewApiClientAdapter(&fakeProvider{name: "anthropic", model: "claude"}, func(model string) Provider {
		return &fakeProvider{name: "openai", model: model}
	})

	adapter.SetModel("gpt-4o")
	resp, err := adapter.SendMessageStream(context.Background(), &api.ApiRequest{}, nil)
	if err != nil {
		t.Fatalf("SendMessageStream() error = %v", err)
	}
	if resp.Model != "gpt-4o" {
		t.Fatalf("model = %s, want gpt-4o", resp.Model)
	}
}

func TestApiClientAdapterSetModelFallsBackToProviderSetter(t *testing.T) {
	initial := &fakeProvider{name: "anthropic", model: "claude"}
	adapter := NewApiClientAdapter(initial)

	adapter.SetModel("claude-sonnet-4-6-20251001")
	if initial.model != "claude-sonnet-4-6-20251001" {
		t.Fatalf("provider model = %s", initial.model)
	}
}
