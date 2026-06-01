package lsp

import (
	"context"
	"encoding/json"
	"testing"
)

func TestGetDiagnostics_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetDiagnostics(context.Background(), "file:///test.go")
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
}

func TestGetDiagnostics_EmptyCache(t *testing.T) {
	// Clean up any residual state
	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.diags = make(map[string][]Diagnostic)
	globalDiagnosticsCache.mu.Unlock()

	client := NewLSPClient("http://example.com")
	// Manually mark as initialized to bypass network dependency
	client.mu.Lock()
	client.initialized = true
	client.mu.Unlock()

	diags, err := client.GetDiagnostics(context.Background(), "file:///nonexistent.go")
	if err != nil {
		t.Fatalf("GetDiagnostics failed: %v", err)
	}
	if len(diags) != 0 {
		t.Errorf("expected empty diags, got %d", len(diags))
	}
}

func TestGetDiagnostics_WithCache(t *testing.T) {
	expected := []Diagnostic{
		{
			Range:    Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 10}},
			Severity: SeverityError,
			Message:  "test error",
			Source:   "test",
		},
	}

	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.diags["file:///test.go"] = expected
	globalDiagnosticsCache.mu.Unlock()

	client := NewLSPClient("http://example.com")
	client.mu.Lock()
	client.initialized = true
	client.mu.Unlock()

	diags, err := client.GetDiagnostics(context.Background(), "file:///test.go")
	if err != nil {
		t.Fatalf("GetDiagnostics failed: %v", err)
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Message != "test error" {
		t.Errorf("expected message %q, got %q", "test error", diags[0].Message)
	}
	if diags[0].Severity != SeverityError {
		t.Errorf("expected severity %d, got %d", SeverityError, diags[0].Severity)
	}
}

func TestHandleDiagnosticsNotification(t *testing.T) {
	params := PublishDiagnosticsParams{
		URI:     "file:///test.go",
		Version: 1,
		Diagnostics: []Diagnostic{
			{Range: Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 1}}, Message: "err1", Severity: SeverityError},
			{Range: Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 1}}, Message: "warn1", Severity: SeverityWarning},
		},
	}

	rawParams, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	client := NewLSPClient("http://example.com")
	if err := client.handleDiagnosticsNotification(rawParams); err != nil {
		t.Fatalf("handleDiagnosticsNotification failed: %v", err)
	}

	// Verify cache was updated
	globalDiagnosticsCache.mu.RLock()
	diags, exists := globalDiagnosticsCache.diags["file:///test.go"]
	globalDiagnosticsCache.mu.RUnlock()

	if !exists {
		t.Fatal("expected diagnostics to be cached")
	}
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(diags))
	}
}

func TestHandleDiagnosticsNotification_InvalidJSON(t *testing.T) {
	client := NewLSPClient("http://example.com")
	err := client.handleDiagnosticsNotification(json.RawMessage(`invalid`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestHandleDiagnosticsNotification_ListenerDelivery(t *testing.T) {
	// Clean up listeners
	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.diags = make(map[string][]Diagnostic)
	globalDiagnosticsCache.listeners = nil
	globalDiagnosticsCache.mu.Unlock()

	client := NewLSPClient("http://example.com")
	ch := client.SubscribeDiagnostics()
	defer client.UnsubscribeDiagnostics(ch)

	params := PublishDiagnosticsParams{
		URI: "file:///listener.go",
		Diagnostics: []Diagnostic{
			{Message: "listener test", Range: Range{Start: Position{}, End: Position{}}},
		},
	}
	rawParams, _ := json.Marshal(params)
	if err := client.handleDiagnosticsNotification(rawParams); err != nil {
		t.Fatalf("handleDiagnosticsNotification failed: %v", err)
	}

	select {
	case diags := <-ch:
		if len(diags) != 1 || diags[0].Message != "listener test" {
			t.Errorf("unexpected diags via listener: %+v", diags)
		}
	default:
		t.Fatal("expected diagnostic to be delivered to listener channel")
	}
}

func TestSubscribeUnsubscribeDiagnostics(t *testing.T) {
	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.listeners = nil
	globalDiagnosticsCache.mu.Unlock()

	client := NewLSPClient("http://example.com")

	ch1 := client.SubscribeDiagnostics()
	ch2 := client.SubscribeDiagnostics()

	if ch1 == nil || ch2 == nil {
		t.Fatal("expected non-nil channels")
	}

	globalDiagnosticsCache.mu.RLock()
	listenerCount := len(globalDiagnosticsCache.listeners)
	globalDiagnosticsCache.mu.RUnlock()
	if listenerCount != 2 {
		t.Errorf("expected 2 listeners, got %d", listenerCount)
	}

	client.UnsubscribeDiagnostics(ch1)

	globalDiagnosticsCache.mu.RLock()
	listenerCount = len(globalDiagnosticsCache.listeners)
	globalDiagnosticsCache.mu.RUnlock()
	if listenerCount != 1 {
		t.Errorf("expected 1 listener after unsubscribe, got %d", listenerCount)
	}

	// Verify ch1 is closed
	_, ok := <-ch1
	if ok {
		t.Error("expected ch1 to be closed after unsubscribe")
	}

	// Clean up
	client.UnsubscribeDiagnostics(ch2)
}

func TestClearDiagnostics(t *testing.T) {
	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.diags["file:///test.go"] = []Diagnostic{{Message: "test"}}
	globalDiagnosticsCache.mu.Unlock()

	client := NewLSPClient("http://example.com")
	client.ClearDiagnostics("file:///test.go")

	globalDiagnosticsCache.mu.RLock()
	_, exists := globalDiagnosticsCache.diags["file:///test.go"]
	globalDiagnosticsCache.mu.RUnlock()

	if exists {
		t.Error("expected diagnostics to be cleared")
	}
}
