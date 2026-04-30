package lsp

import (
	"context"
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
