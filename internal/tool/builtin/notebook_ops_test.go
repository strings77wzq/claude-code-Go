package builtin

import (
	"context"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestNotebookToolEditCell(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "edit",
		"cell_index": float64(0),
		"source":     "print('updated')",
	})

	if result.IsError {
		t.Fatalf("expected success editing cell, got error: %s", result.Content)
	}
}

func TestNotebookToolEditCellOutOfRange(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "edit",
		"cell_index": float64(99),
		"source":     "bad",
	})

	if !result.IsError {
		t.Errorf("expected error for out-of-range cell index")
	}
}

func TestNotebookToolAddCell(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "add_cell",
		"cell_index": float64(1),
		"source":     "print('new')",
		"cell_type":  "code",
	})

	if result.IsError {
		t.Fatalf("expected success adding cell, got error: %s", result.Content)
	}
}

func TestNotebookToolDeleteCell(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "delete_cell",
		"cell_index": float64(0),
	})

	if result.IsError {
		t.Fatalf("expected success deleting cell, got error: %s", result.Content)
	}
}

func TestNotebookToolExecuteMarkdownCell(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "execute",
		"cell_index": float64(1),
	})

	if !result.IsError {
		t.Errorf("expected error for executing markdown cell")
	}
}

func TestNotebookToolUnknownOperation(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path": nbPath,
		"operation": "invalid_op",
	})

	if !result.IsError {
		t.Errorf("expected error for unknown operation")
	}
}

func TestNotebookToolEditCellMissingSource(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path":  nbPath,
		"operation":  "edit",
		"cell_index": float64(0),
	})

	if !result.IsError {
		t.Errorf("expected error for missing source in edit")
	}
}

func TestNotebookToolDeleteCellMissingIndex(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path": nbPath,
		"operation": "delete_cell",
	})

	if !result.IsError {
		t.Errorf("expected error for missing cell_index in delete")
	}
}

var _ tool.Tool = (*NotebookTool)(nil)
