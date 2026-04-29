package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewLSPClient(t *testing.T) {
	c := NewLSPClient("http://example.com")
	if c == nil {
		t.Fatal("NewLSPClient returned nil")
	}
	if c.serverURL != "http://example.com" {
		t.Errorf("expected serverURL %q, got %q", "http://example.com", c.serverURL)
	}
	if c.httpClient == nil {
		t.Error("expected httpClient to be non-nil")
	}
	if c.nextID != 1 {
		t.Errorf("expected nextID to be 1, got %d", c.nextID)
	}
	if c.initialized {
		t.Error("expected new client to not be initialized")
	}
}

func TestInitialize_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}

		if req.Method != MethodInitialize {
			t.Errorf("expected method %q, got %q", MethodInitialize, req.Method)
		}

		result := InitializeResult{
			Capabilities: ServerCapabilities{
				HoverProvider:      true,
				DefinitionProvider: true,
			},
			ServerInfo: &ServerInfo{
				Name:    "test-lsp-server",
				Version: "1.0.0",
			},
		}
		resultBytes, _ := json.Marshal(result)
		resp := JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Result:  resultBytes,
		}
		respBytes, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/vscode-jsonrpc")
		fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(respBytes), respBytes)
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if !client.IsInitialized() {
		t.Error("expected client to be marked as initialized")
	}

	caps := client.GetCapabilities()
	if caps.HoverProvider != true {
		t.Error("expected HoverProvider to be true")
	}
	if caps.DefinitionProvider != true {
		t.Error("expected DefinitionProvider to be true")
	}
}

func TestInitialize_UnreachableServer(t *testing.T) {
	client := NewLSPClient("http://127.0.0.1:1")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := client.Initialize(ctx)
	if err == nil {
		t.Fatal("expected error when connecting to unreachable server")
	}
}

func TestInitialize_AlreadyInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}

		resultBytes, _ := json.Marshal(InitializeResult{})
		resp := JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Result:  resultBytes,
		}
		respBytes, _ := json.Marshal(resp)

		fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(respBytes), respBytes)
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)

	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("first Initialize should succeed: %v", err)
	}

	err := client.Initialize(context.Background())
	if err == nil {
		t.Fatal("expected error when initializing an already-initialized client")
	}
	if !strings.Contains(err.Error(), "already initialized") {
		t.Errorf("expected error containing %q, got %v", "already initialized", err)
	}
}

func TestInitialize_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	err := client.Initialize(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
	if !strings.Contains(err.Error(), "server returned status") {
		t.Errorf("expected error containing %q, got %v", "server returned status", err)
	}
}
