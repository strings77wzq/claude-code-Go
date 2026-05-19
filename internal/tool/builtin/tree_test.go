package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestTreeToolName(t *testing.T) {
	tr := NewTreeTool()
	if tr.Name() != "Tree" {
		t.Errorf("expected 'Tree', got '%s'", tr.Name())
	}
}

func TestTreeToolRequiresPermission(t *testing.T) {
	tr := NewTreeTool()
	if tr.RequiresPermission() {
		t.Errorf("Tree should not require permission")
	}
}

func TestTreeToolInputSchema(t *testing.T) {
	tr := NewTreeTool()
	schema := tr.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestTreeToolExecuteHappyPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "sub"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub", "b.txt"), []byte("b"), 0644)

	tr := NewTreeTool()
	result := tr.Execute(context.Background(), map[string]any{
		"path":      tmpDir,
		"max_depth": float64(3),
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content == "" {
		t.Error("expected non-empty tree output")
	}
}

func TestTreeToolExecuteNonExistent(t *testing.T) {
	tr := NewTreeTool()
	result := tr.Execute(context.Background(), map[string]any{
		"path": "/nonexistent/path/xyz",
	})
	if !result.IsError {
		t.Errorf("expected error for non-existent path")
	}
}

func TestTreeToolExecuteNotADirectory(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "file.txt")
	os.WriteFile(tmpFile, []byte("hello"), 0644)

	tr := NewTreeTool()
	result := tr.Execute(context.Background(), map[string]any{
		"path": tmpFile,
	})
	if !result.IsError {
		t.Errorf("expected error for non-directory path")
	}
}

var _ tool.Tool = (*TreeTool)(nil)

func TestTreeToolExecuteEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	tr := NewTreeTool()
	result := tr.Execute(context.Background(), map[string]any{"path": tmpDir})
	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}
