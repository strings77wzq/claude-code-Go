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

func TestGetHover_Success_Markdown(t *testing.T) {
	server := newHoverTestServer(t)
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	hover, err := client.GetHover(context.Background(), "file:///test.go", 5, 10)
	if err != nil {
		t.Fatalf("GetHover failed: %v", err)
	}
	if hover == nil {
		t.Fatal("expected non-nil hover result")
	}

	if hover.Contents.Kind != MarkupKindMarkdown {
		t.Errorf("expected Kind %q, got %q", MarkupKindMarkdown, hover.Contents.Kind)
	}
	if hover.Contents.Value != "## Test Hover\n\nThis is hover content." {
		t.Errorf("unexpected hover value:\ngot:  %q\nwant: %q", hover.Contents.Value, "## Test Hover\n\nThis is hover content.")
	}

	if hover.Range == nil {
		t.Fatal("expected non-nil range")
	}
	if hover.Range.Start.Line != 1 || hover.Range.Start.Character != 0 {
		t.Errorf("unexpected start position: got line=%d char=%d, want line=1 char=0",
			hover.Range.Start.Line, hover.Range.Start.Character)
	}
	if hover.Range.End.Line != 2 || hover.Range.End.Character != 15 {
		t.Errorf("unexpected end position: got line=%d char=%d, want line=2 char=15",
			hover.Range.End.Line, hover.Range.End.Character)
	}
}

func TestGetHover_Success_PlainText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lspResponder(w, r, t, map[string]interface{}{
			"contents": "plain text content",
		})
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	hover, err := client.GetHover(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetHover failed: %v", err)
	}
	if hover == nil {
		t.Fatal("expected non-nil hover result")
	}

	if hover.Contents.Kind != MarkupKindPlainText {
		t.Errorf("expected Kind %q, got %q", MarkupKindPlainText, hover.Contents.Kind)
	}
	if hover.Contents.Value != "plain text content" {
		t.Errorf("expected value %q, got %q", "plain text content", hover.Contents.Value)
	}
	if hover.Range != nil {
		t.Error("expected nil range for plain text hover")
	}
}

func TestGetHover_NilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lspResponder(w, r, t, nil)
	}))
	defer server.Close()

	client := NewLSPClient(server.URL)
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	hover, err := client.GetHover(context.Background(), "file:///test.go", 1, 2)
	if err != nil {
		t.Fatalf("GetHover failed: %v", err)
	}
	if hover != nil {
		t.Error("expected nil hover for nil result")
	}
}

func TestGetHover_JSONRPCError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}

		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{})
		case MethodTextDocumentHover:
			resp := JSONRPCResponse{
				JSONRPC: JSONRPCVersion,
				ID:      req.ID,
				Error: &JSONRPCError{
					Code:    -32603,
					Message: "Internal error",
				},
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

	_, err := client.GetHover(context.Background(), "file:///test.go", 1, 2)
	if err == nil {
		t.Fatal("expected error for JSON-RPC error response")
	}
	if !strings.Contains(err.Error(), "Internal error") {
		t.Errorf("expected error containing %q, got %v", "Internal error", err)
	}
}

func TestGetHover_NotInitialized(t *testing.T) {
	client := NewLSPClient("http://example.com")
	_, err := client.GetHover(context.Background(), "file:///test.go", 1, 2)
	if err == nil {
		t.Fatal("expected error when client is not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error containing %q, got %v", "not initialized", err)
	}
}

// newHoverTestServer creates an httptest.Server that responds to initialize and
// textDocument/hover with valid markdown hover content including a range.
func newHoverTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}

		switch req.Method {
		case MethodInitialize:
			writeLSPResponse(w, req.ID, InitializeResult{
				Capabilities: ServerCapabilities{
					HoverProvider: true,
				},
			})
		case MethodTextDocumentHover:
			result := map[string]interface{}{
				"contents": map[string]interface{}{
					"kind":  "markdown",
					"value": "## Test Hover\n\nThis is hover content.",
				},
				"range": map[string]interface{}{
					"start": map[string]interface{}{
						"line":      float64(1),
						"character": float64(0),
					},
					"end": map[string]interface{}{
						"line":      float64(2),
						"character": float64(15),
					},
				},
			}
			writeLSPResponse(w, req.ID, result)
		default:
			t.Fatalf("unexpected method: %s", req.Method)
		}
	}))
}

// lspResponder is a helper that writes a JSON-RPC response wrapping result
// as the response to req.Method. It handles both initialize and hover methods.
func lspResponder(w http.ResponseWriter, r *http.Request, t *testing.T, result interface{}) {
	t.Helper()
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatal(err)
	}
	if req.Method == MethodInitialize {
		writeLSPResponse(w, req.ID, InitializeResult{
			Capabilities: ServerCapabilities{
				HoverProvider: true,
			},
		})
		return
	}
	writeLSPResponse(w, req.ID, result)
}

// writeLSPResponse marshals result as a JSON-RPC success response and writes
// it to w with the LSP Content-Length framing.
func writeLSPResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	resp := JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  resultBytes,
	}
	respBytes, _ := json.Marshal(resp)
	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(respBytes), respBytes)
}
