package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestGrepToolName(t *testing.T) {
	g := NewGrepTool()
	if g.Name() != "Grep" {
		t.Errorf("expected 'Grep', got '%s'", g.Name())
	}
}

func TestGrepToolRequiresPermission(t *testing.T) {
	g := NewGrepTool()
	if g.RequiresPermission() {
		t.Errorf("Grep should not require permission")
	}
}

func TestGrepToolInputSchema(t *testing.T) {
	g := NewGrepTool()
	schema := g.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestGrepToolExecuteHappyPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello world\nfoo bar\nhello again"), 0644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	g := NewGrepTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "hello",
		"path":    ".",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}

func TestGrepToolExecuteNoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("foo"), 0644)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	g := NewGrepTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "zzz_not_found",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content != "No matches found" {
		t.Errorf("expected 'No matches found', got '%s'", result.Content)
	}
}

func TestGrepToolExecuteEmptyPattern(t *testing.T) {
	g := NewGrepTool()
	result := g.Execute(context.Background(), map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for empty pattern")
	}
}

func TestGrepToolExecuteInvalidRegex(t *testing.T) {
	g := NewGrepTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "[invalid",
	})
	if !result.IsError {
		t.Errorf("expected error for invalid regex")
	}
}

var _ tool.Tool = (*GrepTool)(nil)

func TestGrepToolExecutePathNotFound(t *testing.T) {
	g := NewGrepTool()
	result := g.Execute(context.Background(), map[string]any{
		"pattern": "foo",
		"path":    "/nonexistent_test_path_xyz",
	})
	// Walk skips non-existent root silently (callback returns nil on error),
	// so grep returns "No matches found" rather than an error
	if result.IsError {
		t.Fatalf("unexpected error: %s", result.Content)
	}
}
