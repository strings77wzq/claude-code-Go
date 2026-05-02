package lsp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLSPGateUnavailable(t *testing.T) {
	gate := NewLSPGate("")

	if gate.IsAvailable() {
		t.Error("gate without server URL should not be available")
	}

	if err := gate.HealthCheck(context.Background()); err != ErrLSPUnavailable {
		t.Errorf("expected ErrLSPUnavailable, got %v", err)
	}

	client, err := gate.GetClient()
	if err != ErrLSPUnavailable {
		t.Errorf("expected ErrLSPUnavailable, got %v", err)
	}
	if client != nil {
		t.Error("client should be nil when unavailable")
	}
}

func TestLSPGateDiagnosticUnavailable(t *testing.T) {
	diag := NewLSPGate("").Diagnostic()

	if diag.Component != "lsp" || diag.Code != "lsp.unavailable" {
		t.Fatalf("unexpected diagnostic: %#v", diag)
	}
	if diag.Severity != "WARN" {
		t.Fatalf("severity = %s, want WARN", diag.Severity)
	}
}

func TestLSPGateConfigured(t *testing.T) {
	// Configured gate reports available (health depends on server reachability)
	gate := NewLSPGate("http://localhost:8080/lsp")

	if !gate.IsAvailable() {
		t.Error("gate with server URL should be available")
	}

	client, err := gate.GetClient()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil when configured")
	}
}

func TestLSPGateHealthCheckWithTraceUnavailable(t *testing.T) {
	traceFile := t.TempDir() + "/trace.jsonl"
	gate := NewLSPGate("")

	err := gate.HealthCheckWithTrace(context.Background(), traceFile)
	if err != ErrLSPUnavailable {
		t.Fatalf("expected ErrLSPUnavailable, got %v", err)
	}

	trace, err := os.ReadFile(traceFile)
	if err != nil {
		t.Fatalf("failed to read trace file: %v", err)
	}
	traceText := string(trace)
	for _, want := range []string{
		`"type":"extension"`,
		`"name":"lsp"`,
		`"event":"health_check"`,
		`"status":"unavailable"`,
	} {
		if !strings.Contains(traceText, want) {
			t.Fatalf("expected %s in trace:\n%s", want, traceText)
		}
	}
}

func TestLSPGateHealthCheckWithTraceSuccess(t *testing.T) {
	traceFile := t.TempDir() + "/trace.jsonl"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gateLSPResponder(w, r, t, InitializeResult{
			Capabilities: ServerCapabilities{
				HoverProvider:              true,
				DefinitionProvider:         true,
				ReferencesProvider:         true,
				DocumentSymbolProvider:     true,
				WorkspaceSymbolProvider:    true,
				PublishDiagnosticsProvider: true,
			},
			ServerInfo: &ServerInfo{Name: "fixture-lsp", Version: "1.0.0"},
		})
	}))
	defer server.Close()
	gate := NewLSPGate(server.URL)

	if err := gate.HealthCheckWithTrace(context.Background(), traceFile); err != nil {
		t.Fatalf("HealthCheckWithTrace() error = %v", err)
	}

	traceText := readGateTrace(t, traceFile)
	for _, want := range []string{
		`"type":"extension"`,
		`"name":"lsp"`,
		`"event":"health_check"`,
		`"status":"available"`,
		`"server":"fixture-lsp"`,
	} {
		if !strings.Contains(traceText, want) {
			t.Fatalf("expected %s in trace:\n%s", want, traceText)
		}
	}
}

func TestLSPGateHealthCheckWithTraceFailure(t *testing.T) {
	traceFile := t.TempDir() + "/trace.jsonl"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	gate := NewLSPGate(server.URL)

	err := gate.HealthCheckWithTrace(context.Background(), traceFile)
	if err == nil {
		t.Fatal("expected health check error")
	}

	traceText := readGateTrace(t, traceFile)
	for _, want := range []string{
		`"type":"extension"`,
		`"name":"lsp"`,
		`"event":"health_check"`,
		`"status":"error"`,
		`"error":"initialize failed: server returned status 500"`,
	} {
		if !strings.Contains(traceText, want) {
			t.Fatalf("expected %s in trace:\n%s", want, traceText)
		}
	}
}

func TestLSPGateAdvertisesOperationsOnlyAfterHealthCheckPasses(t *testing.T) {
	unconfigured := NewLSPGate("")
	if got := unconfigured.AdvertisedOperations(); len(got) != 0 {
		t.Fatalf("unconfigured gate advertised operations: %v", got)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gateLSPResponder(w, r, t, InitializeResult{
			Capabilities: ServerCapabilities{
				HoverProvider:              true,
				DefinitionProvider:         true,
				ReferencesProvider:         true,
				DocumentSymbolProvider:     true,
				WorkspaceSymbolProvider:    true,
				PublishDiagnosticsProvider: true,
			},
		})
	}))
	defer server.Close()
	gate := NewLSPGate(server.URL)
	if got := gate.AdvertisedOperations(); len(got) != 0 {
		t.Fatalf("gate advertised operations before health check: %v", got)
	}
	if err := gate.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
	got := strings.Join(gate.AdvertisedOperations(), ",")
	for _, want := range []string{"diagnostics", "symbols", "definitions", "references", "hover"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in advertised operations %q", want, got)
		}
	}
}

func readGateTrace(t *testing.T, traceFile string) string {
	t.Helper()
	trace, err := os.ReadFile(traceFile)
	if err != nil {
		t.Fatalf("failed to read trace file: %v", err)
	}
	return string(trace)
}

func gateLSPResponder(w http.ResponseWriter, r *http.Request, t *testing.T, result InitializeResult) {
	t.Helper()
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatal(err)
	}
	if req.Method != MethodInitialize {
		t.Fatalf("unexpected method: %s", req.Method)
	}
	writeLSPResponse(w, req.ID, result)
}
