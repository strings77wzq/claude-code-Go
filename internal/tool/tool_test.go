package tool

import (
	"context"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

// stubTool is a minimal Tool implementation for testing.
type stubTool struct{}

func (s *stubTool) Name() string                { return "stub" }
func (s *stubTool) Description() string         { return "a stub tool for testing" }
func (s *stubTool) InputSchema() map[string]any { return map[string]any{"type": "object"} }
func (s *stubTool) RequiresPermission() bool    { return false }
func (s *stubTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (s *stubTool) Execute(_ context.Context, _ map[string]any) Result {
	return Success("stub executed")
}

func TestSuccessAndErrorConstructors(t *testing.T) {
	// Happy path: Success creates a non-error Result.
	ok := Success("ok")
	if ok.Content != "ok" {
		t.Errorf("expected Content 'ok', got %q", ok.Content)
	}
	if ok.IsError {
		t.Error("expected IsError to be false")
	}

	// Happy path: Error creates an error Result.
	fail := Error("something went wrong")
	if fail.Content != "something went wrong" {
		t.Errorf("expected Content 'something went wrong', got %q", fail.Content)
	}
	if !fail.IsError {
		t.Error("expected IsError to be true")
	}
}

func TestRegistryToolNotFound(t *testing.T) {
	// Error path: Executing a tool that is not registered returns an error.
	reg := NewRegistry()
	result := reg.Execute(context.Background(), "nonexistent", nil)

	if !result.IsError {
		t.Fatal("expected an error for missing tool, got success")
	}
	expected := "tool not found: nonexistent"
	if result.Content != expected {
		t.Errorf("expected %q, got %q", expected, result.Content)
	}
}

func TestRegistryRegisterAndGet(t *testing.T) {
	// Happy path: Register a tool and retrieve it.
	reg := NewRegistry()
	stub := &stubTool{}
	if err := reg.Register(stub); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	got := reg.GetTool("stub")
	if got == nil {
		t.Fatal("GetTool returned nil")
	}
	if got.Name() != "stub" {
		t.Errorf("expected Name 'stub', got %q", got.Name())
	}
}

func TestRegistryExecuteRegisteredTool(t *testing.T) {
	// Happy path: Execute a registered tool via the registry.
	reg := NewRegistry()
	if err := reg.Register(&stubTool{}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result := reg.Execute(context.Background(), "stub", nil)
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	if result.Content != "stub executed" {
		t.Errorf("expected 'stub executed', got %q", result.Content)
	}
}

var _ Tool = (*stubTool)(nil)
