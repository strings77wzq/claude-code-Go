package lsp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetReferences_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{ReferencesProvider: true}})
		case MethodTextDocumentReferences:
			writeLSPResponse(w, req.ID, []Location{
				{URI: "file:///test.go", Range: Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 5}}},
				{URI: "file:///test.go", Range: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 5, Character: 10}}},
			})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetReferences(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetReferences failed: %v", err)
	}
	if len(locs) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locs))
	}
	if locs[0].URI != "file:///test.go" {
		t.Errorf("expected uri %q, got %q", "file:///test.go", locs[0].URI)
	}
}

func TestGetReferences_NilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{ReferencesProvider: true}})
		case MethodTextDocumentReferences:
			writeLSPResponse(w, req.ID, nil)
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetReferences(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetReferences failed: %v", err)
	}
	if len(locs) != 0 {
		t.Errorf("expected 0 locations for nil result, got %d", len(locs))
	}
}

func TestGetReferences_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetReferences(context.Background(), "file:///test.go", 1, 2)
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error containing %q, got %v", "not initialized", err)
	}
}
