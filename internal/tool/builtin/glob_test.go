package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestGlobToolName(t *testing.T) {
	g := NewGlobTool()
	if g.Name() != "Glob" {
		t.Errorf("expected 'Glob', got '%s'", g.Name())
	}
}

func TestGlobToolRequiresPermission(t *testing.T) {
	g := NewGlobTool()
	if g.RequiresPermission() {
		t.Errorf("Glob should not require permission")
	}
}

func TestGlobToolInputSchema(t *testing.T) {
	g := NewGlobTool()
	schema := g.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
	req, _ := schema["required"].([]string)
	if len(req) != 1 || req[0] != "pattern" {
		t.Errorf("expected 'pattern' in required")
	}
}

func TestGlobToolExecuteHappyPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("text"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "c.go"), []byte("package test"), 0644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	g := NewGlobTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "*.go",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content != "a.go\nc.go" {
		t.Errorf("expected 'a.go\\nc.go', got '%s'", result.Content)
	}
}

func TestGlobToolExecuteRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "a.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub", "b.go"), []byte("package sub"), 0644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	g := NewGlobTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "**/*.go",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}

func TestGlobToolExecuteEmptyPattern(t *testing.T) {
	g := NewGlobTool()
	result := g.Execute(context.Background(), map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for empty pattern")
	}
}

func TestGlobToolExecuteNoMatches(t *testing.T) {
	g := NewGlobTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "nonexistent*.xyz",
	})
	if result.IsError {
		t.Fatalf("expected success with 'No matches found', got error: %s", result.Content)
	}
	if result.Content != "No matches found" {
		t.Errorf("expected 'No matches found', got '%s'", result.Content)
	}
}

var _ tool.Tool = (*GlobTool)(nil)
