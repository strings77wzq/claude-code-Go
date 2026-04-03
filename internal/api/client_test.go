package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSendMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") == "" {
			t.Error("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Error("missing anthropic-version header")
		}

		var req ApiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		resp := ApiResponse{
			ID:         "msg_123",
			Type:       "message",
			Role:       "assistant",
			Content:    []ContentBlock{{Type: "text", Text: "Hello"}},
			Model:      "claude-3-5-sonnet-20241022",
			StopReason: "end_turn",
			Usage:      Usage{InputTokens: 10, OutputTokens: 5},
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	resp, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "msg_123" {
		t.Errorf("expected id msg_123, got %s", resp.ID)
	}
	if len(resp.Content) != 1 || resp.Content[0].Text != "Hello" {
		t.Errorf("unexpected content: %v", resp.Content)
	}
}

func TestSendMessage_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("invalid-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 error, got: %v", err)
	}
}

func TestSendMessage_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected 403 error, got: %v", err)
	}
}

func TestSendMessage_Retry429(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		resp := ApiResponse{
			ID:         "msg_123",
			Type:       "message",
			Role:       "assistant",
			Content:    []ContentBlock{{Type: "text", Text: "Success after retry"}},
			Model:      "claude-3-5-sonnet-20241022",
			StopReason: "end_turn",
			Usage:      Usage{InputTokens: 10, OutputTokens: 5},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestSendMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err == nil {
		t.Fatal("expected error for 500")
	}
}

func TestSendMessageStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ApiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if !req.Stream {
			t.Error("expected stream=true in request")
		}

		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected http.ResponseWriter to implement http.Flusher")
		}

		fmt.Fprintf(w, "event: message_start\ndata: %s\n\n", `{"type":"message_start","message":{"id":"msg_123","role":"assistant","content":[],"model":"claude-3-5-sonnet-20241022","stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_start\ndata: %s\n\n", `{"type":"content_block_start","index":0,"content_block":{"type":"text","id":"block_1"}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_delta\ndata: %s\n\n", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_delta\ndata: %s\n\n", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" World"}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_stop\ndata: %s\n\n", `{"type":"content_block_stop","index":0}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_delta\ndata: %s\n\n", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_stop\ndata: %s\n\n", `{"type":"message_stop"}`)
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	var accumulatedText strings.Builder

	resp, err := client.SendMessageStream(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	}, func(text string) {
		accumulatedText.WriteString(text)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "msg_123" {
		t.Errorf("expected id msg_123, got %s", resp.ID)
	}
	if accumulatedText.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", accumulatedText.String())
	}
	if resp.StopReason != "end_turn" {
		t.Errorf("expected stop_reason end_turn, got %s", resp.StopReason)
	}
	if resp.Usage.OutputTokens != 5 {
		t.Errorf("expected output_tokens 5, got %d", resp.Usage.OutputTokens)
	}
}

func TestSendMessageStream_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("invalid-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessageStream(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	}, nil)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 error, got: %v", err)
	}
}

func TestSendMessageStream_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessageStream(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	}, nil)
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected 403 error, got: %v", err)
	}
}

func TestSendMessageStream_Retry429(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, _ := w.(http.Flusher)
		fmt.Fprintf(w, "event: message_start\ndata: %s\n\n", `{"type":"message_start","message":{"id":"msg_123","role":"assistant","content":[],"model":"claude-3-5-sonnet-20241022","stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
		flusher.Flush()
		fmt.Fprintf(w, "event: message_delta\ndata: %s\n\n", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`)
		flusher.Flush()
		fmt.Fprintf(w, "event: message_stop\ndata: %s\n\n", `{"type":"message_stop"}`)
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.SendMessageStream(ctx, &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if resp.ID != "msg_123" {
		t.Errorf("expected id msg_123, got %s", resp.ID)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestSendMessageStream_ToolUse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, _ := w.(http.Flusher)

		fmt.Fprintf(w, "event: message_start\ndata: %s\n\n", `{"type":"message_start","message":{"id":"msg_123","role":"assistant","content":[],"model":"claude-3-5-sonnet-20241022","stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_start\ndata: %s\n\n", `{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_123"}}`)
		flusher.Flush()

		inputJSON := `{"action":"run"}`
		escapedInputJSON, _ := json.Marshal(inputJSON)
		fmt.Fprintf(w, "event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"input_json_delta\",\"input_json\":%s}}\n\n", escapedInputJSON)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_stop\ndata: %s\n\n", `{"type":"content_block_stop","index":0}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_delta\ndata: %s\n\n", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_stop\ndata: %s\n\n", `{"type":"message_stop"}`)
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	var accumulatedText strings.Builder

	resp, err := client.SendMessageStream(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Use a tool"}},
		MaxTokens: 100,
	}, func(text string) {
		accumulatedText.WriteString(text)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(resp.Content))
	}
	if resp.Content[0].Type != "tool_use" {
		t.Errorf("expected tool_use block, got %s", resp.Content[0].Type)
	}
	if resp.Content[0].ToolUseID != "toolu_123" {
		t.Errorf("expected tool_use_id toolu_123, got %q", resp.Content[0].ToolUseID)
	}
	if resp.Content[0].Input == nil {
		t.Fatal("expected input to be parsed")
	}
	if resp.Content[0].Input["action"] == nil {
		t.Errorf("expected action in input, got %v", resp.Content[0].Input)
	}
}

func TestSendMessageStream_MultiLineData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, _ := w.(http.Flusher)

		fmt.Fprintf(w, "event: message_start\ndata: %s\n\n", `{"type":"message_start","message":{"id":"msg_123","role":"assistant","content":[],"model":"claude-3-5-sonnet-20241022","stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: content_block_start\ndata: %s\n\n", `{"type":"content_block_start","index":0,"content_block":{"type":"text","id":"block_1"}}`)
		flusher.Flush()

		multiLineData := `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Line1\nLine2\nLine3"}}`
		fmt.Fprintf(w, "event: content_block_delta\ndata: %s\n\n", multiLineData)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_delta\ndata: %s\n\n", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`)
		flusher.Flush()

		fmt.Fprintf(w, "event: message_stop\ndata: %s\n\n", `{"type":"message_stop"}`)
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	var accumulatedText strings.Builder

	resp, err := client.SendMessageStream(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	}, func(text string) {
		accumulatedText.WriteString(text)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accumulatedText.String() != "Line1\nLine2\nLine3" {
		t.Errorf("expected multi-line text, got %q", accumulatedText.String())
	}
	if len(resp.Content) != 1 || resp.Content[0].Text != "Line1\nLine2\nLine3" {
		t.Errorf("content block should have full text: %v", resp.Content)
	}
}

func TestSendMessage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.SendMessage(ctx, &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "context") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected context timeout error, got: %v", err)
	}
}

func TestSendMessage_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "not valid json")
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, "claude-3-5-sonnet-20241022")
	_, err := client.SendMessage(context.Background(), &ApiRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 100,
	})
	if err == nil {
		t.Fatal("expected error for invalid response")
	}
}
