package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetDefinition_SingleLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DefinitionProvider: true}})
		case MethodTextDocumentDefinition:
			writeLSPResponse(w, req.ID, map[string]interface{}{
				"uri": "file:///test.go",
				"range": map[string]interface{}{
					"start": map[string]interface{}{"line": float64(10), "character": float64(5)},
					"end":   map[string]interface{}{"line": float64(10), "character": float64(15)},
				},
			})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetDefinition(context.Background(), "file:///test.go", 5, 10)
	if err != nil {
		t.Fatalf("GetDefinition failed: %v", err)
	}
	if len(locs) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locs))
	}
	if locs[0].URI != "file:///test.go" {
		t.Errorf("expected uri %q, got %q", "file:///test.go", locs[0].URI)
	}
	if locs[0].Range.Start.Line != 10 || locs[0].Range.Start.Character != 5 {
		t.Errorf("unexpected start: got line=%d char=%d", locs[0].Range.Start.Line, locs[0].Range.Start.Character)
	}
	if locs[0].Range.End.Line != 10 || locs[0].Range.End.Character != 15 {
		t.Errorf("unexpected end: got line=%d char=%d", locs[0].Range.End.Line, locs[0].Range.End.Character)
	}
}

func TestGetDefinition_MultipleLocations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DefinitionProvider: true}})
		case MethodTextDocumentDefinition:
			writeLSPResponse(w, req.ID, []interface{}{
				map[string]interface{}{
					"uri": "file:///a.go",
					"range": map[string]interface{}{
						"start": map[string]interface{}{"line": float64(1), "character": float64(0)},
						"end":   map[string]interface{}{"line": float64(1), "character": float64(5)},
					},
				},
				map[string]interface{}{
					"uri": "file:///b.go",
					"range": map[string]interface{}{
						"start": map[string]interface{}{"line": float64(2), "character": float64(0)},
						"end":   map[string]interface{}{"line": float64(2), "character": float64(5)},
					},
				},
			})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetDefinition(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetDefinition failed: %v", err)
	}
	if len(locs) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locs))
	}
}

func TestGetDefinition_NilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DefinitionProvider: true}})
		case MethodTextDocumentDefinition:
			writeLSPResponse(w, req.ID, nil)
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetDefinition(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetDefinition failed: %v", err)
	}
	if locs == nil || len(locs) != 0 {
		t.Errorf("expected empty slice, got %v", locs)
	}
}

func TestGetDefinition_DefaultCase(t *testing.T) {
	// When result is neither a slice nor a map (e.g., a string), should return empty slice.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DefinitionProvider: true}})
		case MethodTextDocumentDefinition:
			writeLSPResponse(w, req.ID, "unexpected_string_result")
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	locs, err := client.GetDefinition(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetDefinition failed: %v", err)
	}
	if len(locs) != 0 {
		t.Errorf("expected 0 locations for unexpected type, got %d", len(locs))
	}
}

func TestGetDefinition_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetDefinition(context.Background(), "file:///test.go", 1, 2)
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error containing %q, got %v", "not initialized", err)
	}
}

func TestGetDefinition_JSONRPCError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DefinitionProvider: true}})
		case MethodTextDocumentDefinition:
			resp := JSONRPCResponse{
				JSONRPC: JSONRPCVersion,
				ID:      req.ID,
				Error:   &JSONRPCError{Code: -32601, Message: "Method not found"},
			}
			respBytes, _ := json.Marshal(resp)
			fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(respBytes), respBytes)
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	_, err := client.GetDefinition(context.Background(), "file:///test.go", 1, 2)
	if err == nil {
		t.Fatal("expected error for JSON-RPC error")
	}
	if !strings.Contains(err.Error(), "Method not found") {
		t.Errorf("expected error containing %q, got %v", "Method not found", err)
	}
}

func TestConvertMapToLocation_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		wantURI  string
		wantZero bool
	}{
		{
			name:     "no uri",
			input:    map[string]interface{}{},
			wantURI:  "",
			wantZero: true,
		},
		{
			name: "no range",
			input: map[string]interface{}{
				"uri": "file:///test.go",
			},
			wantURI: "file:///test.go",
		},
		{
			name: "range without start",
			input: map[string]interface{}{
				"uri": "file:///test.go",
				"range": map[string]interface{}{
					"end": map[string]interface{}{"line": float64(5), "character": float64(0)},
				},
			},
			wantURI: "file:///test.go",
		},
		{
			name: "range without end",
			input: map[string]interface{}{
				"uri": "file:///test.go",
				"range": map[string]interface{}{
					"start": map[string]interface{}{"line": float64(0), "character": float64(0)},
				},
			},
			wantURI: "file:///test.go",
		},
		{
			name: "non-float line values",
			input: map[string]interface{}{
				"uri": "file:///test.go",
				"range": map[string]interface{}{
					"start": map[string]interface{}{"line": "one", "character": float64(0)},
					"end":   map[string]interface{}{"line": "five", "character": "ten"},
				},
			},
			wantURI: "file:///test.go",
		},
		{
			name: "full valid range",
			input: map[string]interface{}{
				"uri": "file:///test.go",
				"range": map[string]interface{}{
					"start": map[string]interface{}{"line": float64(0), "character": float64(1)},
					"end":   map[string]interface{}{"line": float64(10), "character": float64(20)},
				},
			},
			wantURI:  "file:///test.go",
			wantZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := convertMapToLocation(tt.input)
			if err != nil {
				t.Fatalf("convertMapToLocation returned error: %v", err)
			}
			if loc.URI != tt.wantURI {
				t.Errorf("uri = %q, want %q", loc.URI, tt.wantURI)
			}
			if tt.wantZero {
				if loc.Range.Start.Line != 0 || loc.Range.Start.Character != 0 ||
					loc.Range.End.Line != 0 || loc.Range.End.Character != 0 {
					t.Errorf("expected zero range, got %+v", loc.Range)
				}
			}
		})
	}
}
