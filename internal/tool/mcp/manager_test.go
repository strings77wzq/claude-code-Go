package mcp

import (
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestNewMcpManager(t *testing.T) {
	mgr := NewMcpManager()
	if mgr == nil {
		t.Fatal("NewMcpManager returned nil")
	}

	if len(mgr.clients) != 0 {
		t.Errorf("expected empty clients map, got %d entries", len(mgr.clients))
	}
	if len(mgr.adapters) != 0 {
		t.Errorf("expected empty adapters map, got %d entries", len(mgr.adapters))
	}
}

func TestMcpManagerInitializeInvalidConfig(t *testing.T) {
	mgr := NewMcpManager()
	registry := tool.NewRegistry()

	configs := map[string]McpServerConfig{
		"nonexistent": {
			Command: "/no/such/binary/ever",
			Args:    []string{},
			Env:     nil,
		},
	}

	// InitializeAndRegister should return nil (it logs errors and continues)
	err := mgr.InitializeAndRegister(configs, registry)
	if err != nil {
		t.Fatalf("InitializeAndRegister returned unexpected error: %v", err)
	}

	// The server should not have been registered (no clients or adapters)
	if len(mgr.clients) != 0 {
		t.Errorf("expected 0 clients after failed init, got %d", len(mgr.clients))
	}
	if len(mgr.adapters) != 0 {
		t.Errorf("expected 0 adapters after failed init, got %d", len(mgr.adapters))
	}
}

func TestMcpManagerCloseEmpty(t *testing.T) {
	mgr := NewMcpManager()
	err := mgr.Close()
	if err != nil {
		t.Fatalf("Close on empty manager returned error: %v", err)
	}
}
