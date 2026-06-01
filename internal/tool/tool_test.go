package tool

import (
	"context"
	"strings"
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

type panicTool struct{}

func (p *panicTool) Name() string                { return "panic_tool" }
func (p *panicTool) Description() string         { return "a tool that panics for registry recovery tests" }
func (p *panicTool) InputSchema() map[string]any { return map[string]any{"type": "object"} }
func (p *panicTool) RequiresPermission() bool    { return false }
func (p *panicTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (p *panicTool) Execute(_ context.Context, _ map[string]any) Result {
	panic("boom")
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

// largeOutputTool returns content exceeding maxResultSize to test truncation.
type largeOutputTool struct{}

func (l *largeOutputTool) Name() string { return "large_output" }
func (l *largeOutputTool) Description() string {
	return "a tool that returns more than 100KB of output"
}
func (l *largeOutputTool) InputSchema() map[string]any { return map[string]any{"type": "object"} }
func (l *largeOutputTool) RequiresPermission() bool    { return false }
func (l *largeOutputTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (l *largeOutputTool) Execute(_ context.Context, _ map[string]any) Result {
	buf := make([]byte, 150*1024)
	for i := range buf {
		buf[i] = 'x'
	}
	return Success(string(buf))
}

func TestRegistryTruncatesLargeOutput(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(&largeOutputTool{}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result := reg.Execute(context.Background(), "large_output", nil)
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
	if len(result.Content) > 110*1024 {
		t.Errorf("output should be truncated to ~100KB, got %d bytes", len(result.Content))
	}
	if !strings.Contains(result.Content, "[truncated") {
		t.Error("truncated output should contain a truncation marker")
	}
}

func TestRegistryExecuteRecoversPanickingTool(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(&panicTool{}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	result := reg.Execute(context.Background(), "panic_tool", nil)

	if !result.IsError {
		t.Fatalf("expected recovered panic to return an error result, got %#v", result)
	}
	if result.Content == "" || !strings.Contains(result.Content, "panic_tool") || !strings.Contains(result.Content, "panic recovered") {
		t.Fatalf("expected structured panic recovery message with tool name, got %q", result.Content)
	}
}

type tieredTool struct {
	name string
	tier ToolTier
}

func (t *tieredTool) Name() string                { return t.name }
func (t *tieredTool) Description() string         { return "a tool with a tier" }
func (t *tieredTool) InputSchema() map[string]any { return map[string]any{"type": "object"} }
func (t *tieredTool) RequiresPermission() bool    { return false }
func (t *tieredTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (t *tieredTool) Execute(_ context.Context, _ map[string]any) Result {
	return Success("tiered executed")
}
func (t *tieredTool) Tier() ToolTier { return t.tier }

func TestToolTierConstants(t *testing.T) {
	if TierCore == "" {
		t.Error("TierCore should not be empty")
	}
	if TierExtension == "" {
		t.Error("TierExtension should not be empty")
	}
	if TierMCP == "" {
		t.Error("TierMCP should not be empty")
	}
	if TierCore == TierExtension {
		t.Error("TierCore and TierExtension should be different")
	}
}

func TestRegistryGetDefinitionsByTier(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(&tieredTool{name: "core_tool", tier: TierCore}); err != nil {
		t.Fatal(err)
	}
	if err := reg.Register(&tieredTool{name: "ext_tool", tier: TierExtension}); err != nil {
		t.Fatal(err)
	}

	coreDefs := reg.GetDefinitionsByTier(TierCore)
	if len(coreDefs) != 1 {
		t.Fatalf("expected 1 core definition, got %d", len(coreDefs))
	}
	if coreDefs[0].Name != "core_tool" {
		t.Errorf("unexpected tool name: %s", coreDefs[0].Name)
	}

	allDefs := reg.GetDefinitionsByTier(TierCore, TierExtension)
	if len(allDefs) != 2 {
		t.Fatalf("expected 2 definitions for core+extension, got %d", len(allDefs))
	}
}

func TestToolDefaultsToCoreTier(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&stubTool{}) // stubTool doesn't implement TieredTool

	defs := reg.GetDefinitionsByTier(TierCore)
	if len(defs) != 1 {
		t.Fatalf("expected 1 core definition for default tier, got %d", len(defs))
	}
}

var _ Tool = (*stubTool)(nil)
var _ Tool = (*tieredTool)(nil)
