// Package lsp provides Language Server Protocol client and tool integration.
package lsp

import (
	"context"
	"errors"
)

// ErrLSPUnavailable is returned when no LSP server is configured or healthy.
var ErrLSPUnavailable = errors.New("LSP server not configured or unavailable")

// LSPGate checks LSP server availability and gates LSP tool exposure.
type LSPGate struct {
	client     *LSPClient
	configured bool
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

// HealthCheck verifies the LSP server is reachable by sending an initialize request.
func (g *LSPGate) HealthCheck(ctx context.Context) error {
	if !g.configured {
		return ErrLSPUnavailable
	}
	return g.client.Initialize(ctx)
}

// GetClient returns the LSP client if configured, or an error if unavailable.
func (g *LSPGate) GetClient() (*LSPClient, error) {
	if !g.configured {
		return nil, ErrLSPUnavailable
	}
	return g.client, nil
}
