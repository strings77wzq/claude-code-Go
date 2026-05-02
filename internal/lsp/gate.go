// Package lsp provides Language Server Protocol client and tool integration.
package lsp

import (
	"context"
	"errors"

	"github.com/strings77wzq/claude-code-Go/internal/diagnostic"
	"github.com/strings77wzq/claude-code-Go/internal/session"
)

// ErrLSPUnavailable is returned when no LSP server is configured or healthy.
var ErrLSPUnavailable = errors.New("LSP server not configured or unavailable")

// LSPGate checks LSP server availability and gates LSP tool exposure.
type LSPGate struct {
	client     *LSPClient
	configured bool
	healthy    bool
}

// NewLSPGate creates an LSP gate. If serverURL is empty, all LSP features are unavailable.
func NewLSPGate(serverURL string) *LSPGate {
	if serverURL == "" {
		return &LSPGate{configured: false}
	}
	return &LSPGate{
		client:     NewLSPClient(serverURL),
		configured: true,
	}
}

// IsAvailable returns true if an LSP server is configured.
func (g *LSPGate) IsAvailable() bool {
	return g.configured
}

func (g *LSPGate) Diagnostic() diagnostic.Diagnostic {
	if !g.configured {
		return diagnostic.Diagnostic{
			Component: "lsp",
			Severity:  diagnostic.SeverityWarn,
			Code:      "lsp.unavailable",
			Summary:   "LSP server is not configured",
			Retryable: true,
		}
	}
	if !g.healthy {
		return diagnostic.Diagnostic{
			Component: "lsp",
			Severity:  diagnostic.SeverityWarn,
			Code:      "lsp.unhealthy",
			Summary:   "LSP server is configured but not healthy",
			Retryable: true,
		}
	}
	return diagnostic.Diagnostic{
		Component: "lsp",
		Severity:  diagnostic.SeverityInfo,
		Code:      "lsp.available",
		Summary:   "LSP server is available",
	}
}

// HealthCheck verifies the LSP server is reachable by sending an initialize request.
func (g *LSPGate) HealthCheck(ctx context.Context) error {
	if !g.configured {
		return ErrLSPUnavailable
	}
	if g.client.IsInitialized() {
		g.healthy = true
		return nil
	}
	if err := g.client.Initialize(ctx); err != nil {
		g.healthy = false
		return err
	}
	g.healthy = true
	return nil
}

// HealthCheckWithTrace records the LSP health outcome as a non-fatal extension event.
func (g *LSPGate) HealthCheckWithTrace(ctx context.Context, traceFile string) error {
	err := g.HealthCheck(ctx)
	fields := map[string]interface{}{}
	status := "available"
	if err != nil {
		if errors.Is(err, ErrLSPUnavailable) {
			status = "unavailable"
		} else {
			status = "error"
			fields["error"] = err.Error()
		}
	} else if info := g.client.GetServerInfo(); info != nil {
		fields["server"] = info.Name
		if info.Version != "" {
			fields["server_version"] = info.Version
		}
	}
	if g.healthy {
		fields["operations"] = g.AdvertisedOperations()
	}
	if traceErr := session.AppendTraceExtension(traceFile, "lsp", "health_check", status, fields); traceErr != nil && err == nil {
		return traceErr
	}
	return err
}

// AdvertisedOperations returns LSP operations that are safe to expose to callers.
func (g *LSPGate) AdvertisedOperations() []string {
	if !g.configured || !g.healthy || !g.client.IsInitialized() {
		return nil
	}
	caps := g.client.GetCapabilities()
	operations := make([]string, 0, 5)
	if capabilityEnabled(caps.PublishDiagnosticsProvider) {
		operations = append(operations, "diagnostics")
	}
	if capabilityEnabled(caps.WorkspaceSymbolProvider) || capabilityEnabled(caps.DocumentSymbolProvider) {
		operations = append(operations, "symbols")
	}
	if capabilityEnabled(caps.DefinitionProvider) {
		operations = append(operations, "definitions")
	}
	if capabilityEnabled(caps.ReferencesProvider) {
		operations = append(operations, "references")
	}
	if capabilityEnabled(caps.HoverProvider) {
		operations = append(operations, "hover")
	}
	return operations
}

// GetClient returns the LSP client if configured, or an error if unavailable.
func (g *LSPGate) GetClient() (*LSPClient, error) {
	if !g.configured {
		return nil, ErrLSPUnavailable
	}
	return g.client, nil
}

func capabilityEnabled(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}
