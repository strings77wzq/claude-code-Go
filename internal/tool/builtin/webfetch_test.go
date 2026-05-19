package builtin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestWebFetchToolName(t *testing.T) {
	w := NewWebFetchTool()
	if w.Name() != "WebFetch" {
		t.Errorf("expected 'WebFetch', got '%s'", w.Name())
	}
}

func TestWebFetchToolRequiresPermission(t *testing.T) {
	w := NewWebFetchTool()
	if !w.RequiresPermission() {
		t.Errorf("WebFetch should require permission")
	}
}

func TestWebFetchToolInputSchema(t *testing.T) {
	w := NewWebFetchTool()
	schema := w.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestWebFetchToolExecuteHappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><p>Hello World</p></body></html>"))
	}))
	defer srv.Close()

	wf := NewWebFetchTool()
	result := wf.Execute(context.Background(), map[string]any{
		"url": srv.URL,
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content == "" {
		t.Error("expected non-empty response content")
	}
}

func TestWebFetchToolExecute404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	wf := NewWebFetchTool()
	result := wf.Execute(context.Background(), map[string]any{
		"url": srv.URL,
	})

	if !result.IsError {
		t.Errorf("expected error for 404")
	}
}

func TestWebFetchToolExecuteInvalidURL(t *testing.T) {
	wf := NewWebFetchTool()
	result := wf.Execute(context.Background(), map[string]any{
		"url": "://invalid-url",
	})

	if !result.IsError {
		t.Errorf("expected error for invalid URL")
	}
}

func TestWebFetchToolExecuteTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	wf := NewWebFetchTool()
	result := wf.Execute(ctx, map[string]any{
		"url": srv.URL,
	})

	if !result.IsError {
		t.Errorf("expected timeout error")
	}
}

func TestWebFetchToolExecuteEmptyURL(t *testing.T) {
	wf := NewWebFetchTool()
	result := wf.Execute(context.Background(), map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for empty URL")
	}
}

var _ tool.Tool = (*WebFetchTool)(nil)
