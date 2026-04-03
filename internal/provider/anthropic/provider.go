package anthropic

import (
	"context"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider"
)

type AnthropicProvider struct {
	client *api.Client
}

func NewProvider(apiKey, baseURL, model string) *AnthropicProvider {
	return &AnthropicProvider{
		client: api.NewClient(apiKey, baseURL, model),
	}
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

func (p *AnthropicProvider) DefaultModel() string {
	return "claude-sonnet-4-20250514"
}

func (p *AnthropicProvider) SendMessage(ctx context.Context, req *provider.Request) (*provider.Response, error) {
	apiReq := convertToAPIRequest(req)
	apiResp, err := p.client.SendMessage(ctx, apiReq)
	if err != nil {
		return nil, err
	}
	return convertToProviderResponse(apiResp), nil
}

func (p *AnthropicProvider) SendMessageStream(ctx context.Context, req *provider.Request, onTextDelta func(text string)) (*provider.Response, error) {
	apiReq := convertToAPIRequest(req)
	apiResp, err := p.client.SendMessageStream(ctx, apiReq, onTextDelta)
	if err != nil {
		return nil, err
	}
	return convertToProviderResponse(apiResp), nil
}

func convertToAPIRequest(req *provider.Request) *api.ApiRequest {
	messages := make([]api.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	tools := make([]api.ToolDefinition, len(req.Tools))
	for i, tool := range req.Tools {
		tools[i] = api.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		}
	}

	return &api.ApiRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		System:    req.System,
		Stream:    req.Stream,
		Messages:  messages,
		Tools:     tools,
	}
}

func convertToProviderResponse(apiResp *api.ApiResponse) *provider.Response {
	var content strings.Builder
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			content.WriteString(block.Text)
		}
	}

	return &provider.Response{
		ID:         apiResp.ID,
		Content:    content.String(),
		StopReason: apiResp.StopReason,
		Usage: provider.Usage{
			InputTokens:  apiResp.Usage.InputTokens,
			OutputTokens: apiResp.Usage.OutputTokens,
		},
	}
}
