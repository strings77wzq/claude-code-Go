package openai

import (
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

func TestNewProvider_ValidConfig(t *testing.T) {
	p := NewProvider("sk-openai-test123", "https://api.openai.com", "gpt-4o")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}

	if got := p.Name(); got != "openai" {
		t.Fatalf("Name() = %q, want %q", got, "openai")
	}
}

func TestNewProvider_DefaultBaseURL(t *testing.T) {
	p := NewProvider("sk-openai-test123", "", "gpt-4o")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if p.baseURL != defaultBaseURL {
		t.Fatalf("baseURL = %q, want %q", p.baseURL, defaultBaseURL)
	}
}

func TestSetModel(t *testing.T) {
	p := NewProvider("sk-openai-test123", "", "gpt-4o")
	p.SetModel("gpt-4-turbo")
	if p.model != "gpt-4-turbo" {
		t.Fatalf("model = %q, want %q", p.model, "gpt-4-turbo")
	}
}

func TestNewProvider_MissingAPIKey(t *testing.T) {
	p := NewProvider("", "https://api.openai.com", "gpt-4o")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if p.apiKey != "" {
		t.Fatalf("apiKey = %q, want empty", p.apiKey)
	}
}

func TestConvertToOpenAIRequest(t *testing.T) {
	req := &api.ApiRequest{
		Model: "gpt-4o",
		System: "You are a helpful assistant.",
		Messages: []api.Message{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there"},
		},
	}

	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if openaiReq == nil {
		t.Fatal("convertToOpenAIRequest() returned nil")
	}

	if openaiReq.Model != "gpt-4o" {
		t.Fatalf("Model = %q, want %q", openaiReq.Model, "gpt-4o")
	}

	if len(openaiReq.Messages) != 3 {
		t.Fatalf("got %d messages, want 3", len(openaiReq.Messages))
	}

	if openaiReq.Messages[0].Role != "system" {
		t.Fatalf("first message role = %q, want %q", openaiReq.Messages[0].Role, "system")
	}
	if openaiReq.Messages[0].Content != "You are a helpful assistant." {
		t.Fatalf("first message content = %q, want %q", openaiReq.Messages[0].Content, "You are a helpful assistant.")
	}

	if openaiReq.Messages[1].Role != "user" {
		t.Fatalf("second message role = %q, want %q", openaiReq.Messages[1].Role, "user")
	}
	if openaiReq.Messages[1].Content != "Hello" {
		t.Fatalf("second message content = %q, want %q", openaiReq.Messages[1].Content, "Hello")
	}
}

func TestConvertToOpenAIRequest_NoSystem(t *testing.T) {
	req := &api.ApiRequest{
		Model: "gpt-4o",
		Messages: []api.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if len(openaiReq.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(openaiReq.Messages))
	}
}

func TestConvertToApiResponse(t *testing.T) {
	openaiResp := &openAIChatResponse{
		ID: "chatcmpl-123",
		Choices: []openAIChoice{
			{
				Message:      openAIMessage{Role: "assistant", Content: "Hello world"},
				FinishReason: "stop",
			},
		},
		Usage: openAIUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
		},
	}

	resp := convertToApiResponse(openaiResp)
	if resp == nil {
		t.Fatal("convertToApiResponse() returned nil")
	}

	if resp.ID != "chatcmpl-123" {
		t.Fatalf("ID = %q, want %q", resp.ID, "chatcmpl-123")
	}
	if len(resp.Content) != 1 {
		t.Fatalf("got %d content blocks, want 1", len(resp.Content))
	}
	if resp.Content[0].Text != "Hello world" {
		t.Fatalf("content text = %q, want %q", resp.Content[0].Text, "Hello world")
	}
	if resp.StopReason != "stop" {
		t.Fatalf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
	if resp.Usage.InputTokens != 10 {
		t.Fatalf("InputTokens = %d, want %d", resp.Usage.InputTokens, 10)
	}
	if resp.Usage.OutputTokens != 20 {
		t.Fatalf("OutputTokens = %d, want %d", resp.Usage.OutputTokens, 20)
	}
}

func TestConvertToApiResponse_EmptyChoices(t *testing.T) {
	openaiResp := &openAIChatResponse{
		ID:      "chatcmpl-456",
		Choices: []openAIChoice{},
	}

	resp := convertToApiResponse(openaiResp)
	if resp == nil {
		t.Fatal("convertToApiResponse() returned nil")
	}
	if len(resp.Content) != 0 {
		t.Fatalf("got %d content blocks, want 0", len(resp.Content))
	}
}
