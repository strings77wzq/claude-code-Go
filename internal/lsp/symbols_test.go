package lsp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSymbols_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{WorkspaceSymbolProvider: true}})
		case MethodWorkspaceSymbol:
			writeLSPResponse(w, req.ID, []SymbolInformation{
				{Name: "myFunc", Kind: SymbolKindFunction, Location: Location{URI: "file:///test.go", Range: Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 10}}}},
				{Name: "MyStruct", Kind: SymbolKindClass, Location: Location{URI: "file:///test.go", Range: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 5, Character: 20}}}},
			})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	symbols, err := client.GetSymbols(context.Background(), "myFunc")
	if err != nil {
		t.Fatalf("GetSymbols failed: %v", err)
	}
	if len(symbols) != 2 {
		t.Fatalf("expected 2 symbols, got %d", len(symbols))
	}
	if symbols[0].Name != "myFunc" || symbols[0].Kind != SymbolKindFunction {
		t.Errorf("unexpected first symbol: %+v", symbols[0])
	}
	if symbols[1].Name != "MyStruct" || symbols[1].Kind != SymbolKindClass {
		t.Errorf("unexpected second symbol: %+v", symbols[1])
	}
}

func TestGetSymbols_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{WorkspaceSymbolProvider: true}})
		case MethodWorkspaceSymbol:
			writeLSPResponse(w, req.ID, []SymbolInformation{})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	symbols, err := client.GetSymbols(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("GetSymbols failed: %v", err)
	}
	if len(symbols) != 0 {
		t.Errorf("expected 0 symbols, got %d", len(symbols))
	}
}

func TestGetSymbols_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetSymbols(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error containing %q, got %v", "not initialized", err)
	}
}

func TestGetDocumentSymbols_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{Capabilities: ServerCapabilities{DocumentSymbolProvider: true}})
		case MethodTextDocumentDocumentSymbol:
			writeLSPResponse(w, req.ID, []DocumentSymbol{
				{
					Name: "myFunc", Kind: SymbolKindFunction,
					Range:          Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 10, Character: 0}},
					SelectionRange: Range{Start: Position{Line: 1, Character: 0}, End: Position{Line: 1, Character: 5}},
				},
			})
		}
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	symbols, err := client.GetDocumentSymbols(context.Background(), "file:///test.go")
	if err != nil {
		t.Fatalf("GetDocumentSymbols failed: %v", err)
	}
	if len(symbols) != 1 {
		t.Fatalf("expected 1 symbol, got %d", len(symbols))
	}
	if symbols[0].Name != "myFunc" || symbols[0].Kind != SymbolKindFunction {
		t.Errorf("unexpected symbol: %+v", symbols[0])
	}
}

func TestGetDocumentSymbols_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetDocumentSymbols(context.Background(), "file:///test.go")
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error containing %q, got %v", "not initialized", err)
	}
}
