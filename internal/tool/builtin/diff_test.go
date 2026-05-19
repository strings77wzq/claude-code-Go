package builtin

import (
	"context"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestDiffToolName(t *testing.T) {
	d := NewDiffTool()
	if d.Name() != "Diff" {
		t.Errorf("expected 'Diff', got '%s'", d.Name())
	}
}

func TestDiffToolRequiresPermission(t *testing.T) {
	d := NewDiffTool()
	if d.RequiresPermission() {
		t.Errorf("Diff should not require permission")
	}
}

func TestDiffToolInputSchema(t *testing.T) {
	d := NewDiffTool()
	schema := d.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestDiffToolExecuteHappyPath(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "line1\nline2\nline3",
		"new_content": "line1\nline2 modified\nline3",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
}

func TestDiffToolExecuteIdentical(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "hello world",
		"new_content": "hello world",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content != "No differences found." {
		t.Errorf("expected 'No differences found.', got '%s'", result.Content)
	}
}

func TestDiffToolExecuteMissingOldContent(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"new_content": "hello",
	})
	if !result.IsError {
		t.Errorf("expected error for missing old_content")
	}
}

func TestDiffToolExecuteMissingNewContent(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "hello",
	})
	if !result.IsError {
		t.Errorf("expected error for missing new_content")
	}
}

var _ tool.Tool = (*DiffTool)(nil)
