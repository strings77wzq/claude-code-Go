package init

import (
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestRegisterBuiltinTools(t *testing.T) {
	registry := tool.NewRegistry()
	err := RegisterBuiltinTools(registry, t.TempDir())
	if err != nil {
		t.Fatalf("RegisterBuiltinTools returned error: %v", err)
	}

	defs := registry.GetAllDefinitions()
	expectedCount := 11
	if len(defs) != expectedCount {
		t.Errorf("expected %d tools, got %d", expectedCount, len(defs))
	}

	// Verify specific well-known tools are present by name
	names := make(map[string]bool)
	for _, d := range defs {
		names[d.Name] = true
	}

	expectedNames := []string{"Bash", "Read", "Write", "Edit", "Glob", "Grep", "Diff", "Tree", "WebFetch", "TodoWrite", "NotebookEdit"}
	for _, n := range expectedNames {
		if !names[n] {
			t.Errorf("expected tool %q not registered", n)
		}
	}
}

func TestRegisterBuiltinToolsDoubleRegister(t *testing.T) {
	registry := tool.NewRegistry()
	workingDir := t.TempDir()

	err := RegisterBuiltinTools(registry, workingDir)
	if err != nil {
		t.Fatalf("first RegisterBuiltinTools returned error: %v", err)
	}

	err = RegisterBuiltinTools(registry, workingDir)
	if err == nil {
		t.Fatal("expected error on double register, got nil")
	}
}
