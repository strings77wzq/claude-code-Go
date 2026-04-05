package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider"
)

const (
	defaultBaseURL = "https://api.openai.com"
	maxRetries     = 3
	retryDelayBase = time.Second
)

type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) SendMessage(ctx context.Context, req *api.ApiRequest) (*api.ApiResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}

	openaiReq := convertToOpenAIRequest(req, model)
	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.setHeaders(httpReq)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryDelayBase * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := p.httpClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusInternalServerError {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limited or server error (%d)", resp.StatusCode)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, fmt.Errorf("unauthorized (401): invalid API key")
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		defer resp.Body.Close()
		var openaiResp openAIChatResponse
		if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return convertToApiResponse(&openaiResp), nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (p *OpenAIProvider) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}

	openaiReq := convertToOpenAIRequest(req, model)
	openaiReq.Stream = true

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.setHeaders(httpReq)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryDelayBase * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := p.httpClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusInternalServerError {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limited or server error (%d)", resp.StatusCode)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, fmt.Errorf("unauthorized (401): invalid API key")
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		defer resp.Body.Close()
		return p.parseStreamResponse(resp.Body, onTextDelta)
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func (p *OpenAIProvider) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
}

func (p *OpenAIProvider) parseStreamResponse(body io.Reader, onTextDelta func(text string)) (*api.ApiResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	var contentBlocks []api.ContentBlock
	var stopReason string
	var outputTokens int

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if !strings.HasPrefix(line, "data:") {
			continue
		}

		dataContent := strings.TrimPrefix(line, "data:")
		dataContent = strings.TrimSpace(dataContent)

		if dataContent == "" || dataContent == "[DONE]" {
			continue
		}

		var chunk openAIStreamChunk
		if err := json.Unmarshal([]byte(dataContent), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			if chunk.Choices[0].Delta.Content != "" {
				contentBlocks = append(contentBlocks, api.ContentBlock{
					Type: "text",
					Text: chunk.Choices[0].Delta.Content,
				})
				if onTextDelta != nil {
					onTextDelta(chunk.Choices[0].Delta.Content)
				}
			}

			if chunk.Choices[0].FinishReason != "" {
				stopReason = chunk.Choices[0].FinishReason
			}
		}
	}

	return &api.ApiResponse{
		ID:         "stream",
		Type:       "message",
		Role:       "assistant",
		Content:    contentBlocks,
		StopReason: stopReason,
		Usage: api.Usage{
			InputTokens:  0,
			OutputTokens: outputTokens,
		},
	}, nil
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Stream      bool            `json:"stream,omitempty"`
	Tools       []openAITool    `json:"tools,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAITool struct {
	Type     string      `json:"type"`
	Function interface{} `json:"function"`
}

type openAIChatResponse struct {
	ID      string         `json:"id"`
	Choices []openAIChoice `json:"choices"`
	Usage   openAIUsage    `json:"usage"`
}

type openAIChoice struct {
	Message      openAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Delta        openAIDelta   `json:"delta"`
}

type openAIDelta struct {
	Content string `json:"content"`
}

type openAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type openAIStreamChunk struct {
	ID      string         `json:"id"`
	Choices []openAIChoice `json:"choices"`
}

func convertToOpenAIRequest(req *api.ApiRequest, model string) *openAIChatRequest {
	messages := make([]openAIMessage, 0, len(req.Messages))

	if req.System != "" {
		messages = append(messages, openAIMessage{
			Role:    "system",
			Content: req.System,
		})
	}

	for _, msg := range req.Messages {
		content := ""
		switch c := msg.Content.(type) {
		case string:
			content = c
		case []api.ContentBlock:
			var sb strings.Builder
			for _, block := range c {
				if block.Type == "text" {
					sb.WriteString(block.Text)
				}
			}
			content = sb.String()
		}
		messages = append(messages, openAIMessage{
			Role:    msg.Role,
			Content: content,
		})
	}

	return &openAIChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   req.Stream,
	}
}

func convertToApiResponse(openaiResp *openAIChatResponse) *api.ApiResponse {
	if len(openaiResp.Choices) == 0 {
		return &api.ApiResponse{
			ID:      openaiResp.ID,
			Type:    "message",
			Role:    "assistant",
			Content: []api.ContentBlock{},
		}
	}

	content := openaiResp.Choices[0].Message.Content
	contentBlocks := []api.ContentBlock{
		{
			Type: "text",
			Text: content,
		},
	}

	return &api.ApiResponse{
		ID:         openaiResp.ID,
		Type:       "message",
		Role:       "assistant",
		Content:    contentBlocks,
		StopReason: openaiResp.Choices[0].FinishReason,
		Usage: api.Usage{
			InputTokens:  openaiResp.Usage.PromptTokens,
			OutputTokens: openaiResp.Usage.CompletionTokens,
		},
	}
}

var _ provider.Provider = (*OpenAIProvider)(nil)
