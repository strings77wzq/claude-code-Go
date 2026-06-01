package openai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/provider"
)

// ---- Helpers ----

func testProvider(serverURL string) *OpenAIProvider {
	return NewProvider("sk-test-key", serverURL, "gpt-4o")
}

// ---- Provider constructor & config ----

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

func TestNewProvider_MissingAPIKey(t *testing.T) {
	p := NewProvider("", "https://api.openai.com", "gpt-4o")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if p.apiKey != "" {
		t.Fatalf("apiKey = %q, want empty", p.apiKey)
	}
}

func TestSetModel(t *testing.T) {
	p := NewProvider("sk-openai-test123", "", "gpt-4o")
	p.SetModel("gpt-4-turbo")
	if p.model != "gpt-4-turbo" {
		t.Fatalf("model = %q, want %q", p.model, "gpt-4-turbo")
	}
}

func TestName(t *testing.T) {
	p := NewProvider("sk-test", "http://example.com", "gpt-4o")
	if got := p.Name(); got != "openai" {
		t.Errorf("Name() = %q, want %q", got, "openai")
	}
}

// ---- setHeaders ----

func TestSetHeaders(t *testing.T) {
	p := testProvider("http://test.com")
	req, err := http.NewRequest("POST", "http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	p.setHeaders(req)

	if got := req.Header.Get("Authorization"); got != "Bearer sk-test-key" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer sk-test-key")
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

// ---- convertToOpenAIRequest ----

func TestConvertToOpenAIRequest_ContentBlocks(t *testing.T) {
	req := &api.ApiRequest{
		Model: "gpt-4o",
		Messages: []api.Message{
			{
				Role: "user",
				Content: []api.ContentBlock{
					{Type: "text", Text: "Hello "},
					{Type: "text", Text: "world"},
					{Type: "tool_use", Text: "should-be-ignored"},
				},
			},
		},
	}
	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if len(openaiReq.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(openaiReq.Messages))
	}
	// Only "text" type blocks are concatenated; non-text blocks are skipped
	if openaiReq.Messages[0].Content != "Hello world" {
		t.Errorf("content = %q, want %q", openaiReq.Messages[0].Content, "Hello world")
	}
}

func TestConvertToOpenAIRequest_CachedSystemPrompt(t *testing.T) {
	req := &api.ApiRequest{
		Model:  "gpt-4o",
		System: api.CachedSystemPrompt{Text: "You are a bot.", Type: "text"},
		Messages: []api.Message{
			{Role: "user", Content: "Hello"},
		},
	}
	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if len(openaiReq.Messages) != 2 {
		t.Fatalf("got %d messages, want 2 (system + user)", len(openaiReq.Messages))
	}
	if openaiReq.Messages[0].Role != "system" || openaiReq.Messages[0].Content != "You are a bot." {
		t.Errorf("system message = %+v", openaiReq.Messages[0])
	}
}

func TestConvertToOpenAIRequest_NilSystem(t *testing.T) {
	req := &api.ApiRequest{
		Model:    "gpt-4o",
		System:   nil,
		Messages: []api.Message{{Role: "user", Content: "Hello"}},
	}
	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if len(openaiReq.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(openaiReq.Messages))
	}
}

func TestConvertToOpenAIRequest_StreamPreserved(t *testing.T) {
	req := &api.ApiRequest{
		Model:    "gpt-4o",
		Stream:   true,
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}
	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if !openaiReq.Stream {
		t.Error("Stream field should be preserved from request")
	}
}

func TestConvertToOpenAIRequest_EmptyContent(t *testing.T) {
	req := &api.ApiRequest{
		Model: "gpt-4o",
		Messages: []api.Message{
			{Role: "user", Content: ""},
			{Role: "user", Content: []api.ContentBlock{}},
		},
	}
	openaiReq := convertToOpenAIRequest(req, "gpt-4o")
	if len(openaiReq.Messages) != 2 {
		t.Fatalf("got %d messages, want 2", len(openaiReq.Messages))
	}
	if openaiReq.Messages[0].Content != "" {
		t.Errorf("empty string content = %q, want empty", openaiReq.Messages[0].Content)
	}
	if openaiReq.Messages[1].Content != "" {
		t.Errorf("empty content block = %q, want empty", openaiReq.Messages[1].Content)
	}
}

func TestConvertToOpenAIRequest_ModelOverride(t *testing.T) {
	req := &api.ApiRequest{
		Model:    "custom-model",
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}
	// The model parameter passed to convertToOpenAIRequest should be used when req.Model is empty
	// Here req.Model is non-empty, but the function still uses the passed model
	openaiReq := convertToOpenAIRequest(req, "provider-default-model")
	if openaiReq.Model != "provider-default-model" {
		t.Errorf("Model = %q, want %q", openaiReq.Model, "provider-default-model")
	}
}

// ---- parseStreamResponse ----

func TestParseStreamResponse(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		wantText string
		wantStop string
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "single token",
			data:     "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"},\"finish_reason\":null}]}\n",
			wantText: "Hello",
			wantLen:  1,
		},
		{
			name:     "multiple tokens",
			data:     "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"},\"finish_reason\":null}]}\ndata: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\" world\"},\"finish_reason\":null}]}\n",
			wantText: "Hello world",
			wantLen:  2,
		},
		{
			name:     "with stop reason",
			data:     "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Done\"},\"finish_reason\":\"stop\"}]}\n",
			wantText: "Done",
			wantStop: "stop",
			wantLen:  1,
		},
		{
			name:    "done marker",
			data:    "data: [DONE]\n",
			wantLen: 0,
		},
		{
			name:    "empty data prefix",
			data:    "data:\n",
			wantLen: 0,
		},
		{
			name:    "no data prefix",
			data:    "event: foo\n",
			wantLen: 0,
		},
		{
			name:     "carriage return suffix",
			data:     "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hi\"}}]}\r\n",
			wantText: "Hi",
			wantLen:  1,
		},
		{
			name:    "empty response body",
			data:    "",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := testProvider("http://test.com")
			var callbackTexts []string
			resp, err := p.parseStreamResponse(strings.NewReader(tt.data), func(text string) {
				callbackTexts = append(callbackTexts, text)
			})
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseStreamResponse() error: %v", err)
			}
			if len(resp.Content) != tt.wantLen {
				t.Errorf("got %d content blocks, want %d", len(resp.Content), tt.wantLen)
			}
			if tt.wantLen > 0 {
				var fullText string
				for _, block := range resp.Content {
					fullText += block.Text
				}
				if fullText != tt.wantText {
					t.Errorf("full text = %q, want %q", fullText, tt.wantText)
				}
			}
			if tt.wantStop != "" && resp.StopReason != tt.wantStop {
				t.Errorf("StopReason = %q, want %q", resp.StopReason, tt.wantStop)
			}
			// Verify callback was called for each block
			if len(callbackTexts) != tt.wantLen {
				t.Errorf("onTextDelta called %d times, want %d", len(callbackTexts), tt.wantLen)
			}
		})
	}
}

func TestParseStreamResponse_NilCallback(t *testing.T) {
	p := testProvider("http://test.com")
	data := "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"},\"finish_reason\":null}]}\n"
	resp, err := p.parseStreamResponse(strings.NewReader(data), nil)
	if err != nil {
		t.Fatalf("parseStreamResponse() error: %v", err)
	}
	if len(resp.Content) != 1 || resp.Content[0].Text != "Hello" {
		t.Errorf("content = %v", resp.Content)
	}
}

func TestParseStreamResponse_InvalidJSON(t *testing.T) {
	p := testProvider("http://test.com")
	data := "data: {invalid json}\n"
	resp, err := p.parseStreamResponse(strings.NewReader(data), nil)
	if err != nil {
		t.Fatalf("parseStreamResponse() error: %v", err)
	}
	// Invalid JSON chunks are silently skipped
	if len(resp.Content) != 0 {
		t.Errorf("expected 0 content blocks, got %d", len(resp.Content))
	}
}

// ---- SendMessage ----

func TestSendMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		if r.Header.Get("Authorization") != "Bearer sk-test-key" {
			t.Error("wrong authorization header")
		}

		// Verify request body
		var reqBody openAIChatRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if reqBody.Model != "gpt-4o" {
			t.Errorf("model = %q, want %q", reqBody.Model, "gpt-4o")
		}
		if reqBody.Stream {
			t.Error("stream should be false for non-streaming")
		}
		if len(reqBody.Messages) != 1 || reqBody.Messages[0].Content != "Hi" {
			t.Errorf("messages = %+v", reqBody.Messages)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openAIChatResponse{
			ID: "chatcmpl-test",
			Choices: []openAIChoice{
				{Message: openAIMessage{Role: "assistant", Content: "Hello!"}, FinishReason: "stop"},
			},
			Usage: openAIUsage{PromptTokens: 10, CompletionTokens: 5},
		})
	}))
	defer server.Close()

	p := testProvider(server.URL)
	resp, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Model: "gpt-4o",
		Messages: []api.Message{
			{Role: "user", Content: "Hi"},
		},
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	if resp.ID != "chatcmpl-test" {
		t.Errorf("ID = %q, want %q", resp.ID, "chatcmpl-test")
	}
	if len(resp.Content) != 1 || resp.Content[0].Text != "Hello!" {
		t.Errorf("content = %+v", resp.Content)
	}
	if resp.StopReason != "stop" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
	if resp.Usage.InputTokens != 10 || resp.Usage.OutputTokens != 5 {
		t.Errorf("usage = %+v", resp.Usage)
	}
}

func TestSendMessage_ModelFromProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody openAIChatRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Model != "gpt-4o" {
			t.Errorf("model = %q, want %q (from provider's default)", reqBody.Model, "gpt-4o")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openAIChatResponse{
			ID: "chatcmpl-1",
			Choices: []openAIChoice{
				{Message: openAIMessage{Role: "assistant", Content: "OK"}, FinishReason: "stop"},
			},
		})
	}))
	defer server.Close()

	p := NewProvider("sk-test", server.URL, "gpt-4o")
	// Request without explicit model — should use provider's default model
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
}

func TestSendMessage_WithSystemPrompt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody openAIChatRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.Messages) != 2 {
			t.Fatalf("got %d messages, want 2 (system + user)", len(reqBody.Messages))
		}
		if reqBody.Messages[0].Role != "system" || reqBody.Messages[0].Content != "Be helpful." {
			t.Errorf("system message = %+v", reqBody.Messages[0])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openAIChatResponse{
			ID: "chatcmpl-sys",
			Choices: []openAIChoice{
				{Message: openAIMessage{Role: "assistant", Content: "Sure!"}, FinishReason: "stop"},
			},
		})
	}))
	defer server.Close()

	p := testProvider(server.URL)
	resp, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Model:    "gpt-4o",
		System:   "Be helpful.",
		Messages: []api.Message{{Role: "user", Content: "Help"}},
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if resp.Content[0].Text != "Sure!" {
		t.Errorf("content = %q, want %q", resp.Content[0].Text, "Sure!")
	}
}

func TestSendMessage_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var classifiedErr *provider.ClassifiedError
	if !errors.As(err, &classifiedErr) {
		t.Fatalf("error is %T, want *provider.ClassifiedError", err)
	}
	if classifiedErr.Kind != provider.ErrorAuth {
		t.Errorf("error kind = %q, want %q", classifiedErr.Kind, provider.ErrorAuth)
	}
}

func TestSendMessage_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid model"}`))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var classifiedErr *provider.ClassifiedError
	if !errors.As(err, &classifiedErr) {
		t.Fatalf("error is %T, want *provider.ClassifiedError", err)
	}
	if classifiedErr.Kind != provider.ErrorInvalidRequest {
		t.Errorf("error kind = %q, want %q", classifiedErr.Kind, provider.ErrorInvalidRequest)
	}
}

func TestSendMessage_ServerError(t *testing.T) {
	var attempt int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// All attempts (1 initial + maxRetries retries) should fail
	if attempt != maxRetries+1 {
		t.Errorf("expected %d attempts, got %d", maxRetries+1, attempt)
	}
	if !strings.Contains(err.Error(), "server") {
		t.Errorf("error = %v, want server error", err)
	}
}

func TestSendMessage_RetryOn429ThenSuccess(t *testing.T) {
	var attempt int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt++
		if attempt == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openAIChatResponse{
			ID: "chatcmpl-retry",
			Choices: []openAIChoice{
				{Message: openAIMessage{Role: "assistant", Content: "After retry"}, FinishReason: "stop"},
			},
		})
	}))
	defer server.Close()

	p := testProvider(server.URL)
	resp, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}
	if resp.Content[0].Text != "After retry" {
		t.Errorf("content = %q, want %q", resp.Content[0].Text, "After retry")
	}
	if attempt != 2 {
		t.Errorf("expected 2 attempts, got %d", attempt)
	}
}

func TestSendMessage_ContextCancelDuringRetry(t *testing.T) {
	var attempt int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before entering the retry loop

	p := testProvider(server.URL)
	_, err := p.SendMessage(ctx, &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestSendMessage_DecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "decode") {
		t.Errorf("error = %v, want decode error", err)
	}
}

func TestSendMessage_NonStandardError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// 402 Unprocessable (non-standard, non-retryable)
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(`{"error": "insufficient_quota"}`))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessage(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var classifiedErr *provider.ClassifiedError
	if !errors.As(err, &classifiedErr) {
		t.Fatalf("error is %T, want *provider.ClassifiedError", err)
	}
	if classifiedErr.StatusCode != http.StatusPaymentRequired {
		t.Errorf("status code = %d, want %d", classifiedErr.StatusCode, http.StatusPaymentRequired)
	}
}

// ---- SendMessageStream ----

func TestSendMessageStream_Success(t *testing.T) {
	streamData := "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"},\"finish_reason\":null}]}\n" +
		"data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\" world\"},\"finish_reason\":null}]}\n" +
		"data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"\"},\"finish_reason\":\"stop\"}]}\n" +
		"data: [DONE]\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk-test-key" {
			t.Error("wrong authorization header")
		}

		var reqBody openAIChatRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if !reqBody.Stream {
			t.Error("stream should be true for streaming request")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte(streamData))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	var deltas []string
	resp, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Model:    "gpt-4o",
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, func(text string) {
		deltas = append(deltas, text)
	})
	if err != nil {
		t.Fatalf("SendMessageStream() error: %v", err)
	}
	// Only 2 content blocks: "Hello" and " world" (empty delta content is skipped)
	if len(resp.Content) != 2 {
		t.Errorf("got %d content blocks, want 2", len(resp.Content))
	}
	// Check deltas received
	if len(deltas) != 2 {
		t.Errorf("onTextDelta called %d times, want 2", len(deltas))
	}
	if resp.StopReason != "stop" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
}

func TestSendMessageStream_ModelFromProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody openAIChatRequest
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Model != "gpt-4o" {
			t.Errorf("model = %q, want %q (from provider)", reqBody.Model, "gpt-4o")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: [DONE]\n"))
	}))
	defer server.Close()

	p := NewProvider("sk-test", server.URL, "gpt-4o")
	_, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err != nil {
		t.Fatalf("SendMessageStream() error: %v", err)
	}
}

func TestSendMessageStream_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var classifiedErr *provider.ClassifiedError
	if !errors.As(err, &classifiedErr) {
		t.Fatalf("error is %T, want *provider.ClassifiedError", err)
	}
	if classifiedErr.Kind != provider.ErrorAuth {
		t.Errorf("error kind = %q, want %q", classifiedErr.Kind, provider.ErrorAuth)
	}
}

func TestSendMessageStream_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var classifiedErr *provider.ClassifiedError
	if !errors.As(err, &classifiedErr) {
		t.Fatalf("error is %T, want *provider.ClassifiedError", err)
	}
	if classifiedErr.Kind != provider.ErrorInvalidRequest {
		t.Errorf("error kind = %q, want %q", classifiedErr.Kind, provider.ErrorInvalidRequest)
	}
}

func TestSendMessageStream_ServerErrorThenSuccess(t *testing.T) {
	var attempt int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt++
		if attempt == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: [DONE]\n"))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err != nil {
		t.Fatalf("SendMessageStream() error: %v", err)
	}
	if attempt != 2 {
		t.Errorf("expected 2 attempts, got %d", attempt)
	}
}

func TestSendMessageStream_ContextCancelDuringRetry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := testProvider(server.URL)
	_, err := p.SendMessageStream(ctx, &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestSendMessageStream_NilCallback(t *testing.T) {
	streamData := "data: {\"id\":\"1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\ndata: [DONE]\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte(streamData))
	}))
	defer server.Close()

	p := testProvider(server.URL)
	// Passing nil callback should not panic
	resp, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err != nil {
		t.Fatalf("SendMessageStream() error: %v", err)
	}
	if len(resp.Content) != 1 {
		t.Errorf("got %d content blocks, want 1", len(resp.Content))
	}
}

func TestSendMessageStream_NetworkError(t *testing.T) {
	// A server that closes immediately before sending any response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("server does not support hijacking")
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer server.Close()

	p := testProvider(server.URL)
	_, err := p.SendMessageStream(context.Background(), &api.ApiRequest{
		Messages: []api.Message{{Role: "user", Content: "Hi"}},
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---- convertToApiResponse ----

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

func TestConvertToApiResponse_EmptyContentString(t *testing.T) {
	openaiResp := &openAIChatResponse{
		ID: "chatcmpl-empty",
		Choices: []openAIChoice{
			{Message: openAIMessage{Role: "assistant", Content: ""}},
		},
	}
	resp := convertToApiResponse(openaiResp)
	if len(resp.Content) != 1 {
		t.Fatalf("got %d content blocks, want 1", len(resp.Content))
	}
	if resp.Content[0].Text != "" {
		t.Errorf("content text = %q, want empty", resp.Content[0].Text)
	}
}

func TestConvertToApiResponse_ZeroUsage(t *testing.T) {
	openaiResp := &openAIChatResponse{
		ID: "chatcmpl-usage",
		Choices: []openAIChoice{
			{Message: openAIMessage{Role: "assistant", Content: "Hi"}, FinishReason: "stop"},
		},
	}
	resp := convertToApiResponse(openaiResp)
	if resp.Usage.InputTokens != 0 || resp.Usage.OutputTokens != 0 {
		t.Errorf("usage = %+v, want zero values", resp.Usage)
	}
}

// ---- Interface compliance ----

func TestProviderImplementsInterface(t *testing.T) {
	p := testProvider("http://test.com")
	var _ provider.Provider = p
}
