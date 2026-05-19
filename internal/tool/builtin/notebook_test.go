package builtin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func createTestNotebook(t *testing.T, tmpDir string) string {
	t.Helper()
	nb := Notebook{
		NBFormat:   4,
		NBFormatMA: 5,
		Metadata:   json.RawMessage(`{}`),
		Cells: []Cell{
			{
				CellType: "code",
				Source:   []string{"print('hello')"},
				Metadata: json.RawMessage(`{}`),
				Outputs:  json.RawMessage(`[]`),
			},
			{
				CellType: "markdown",
				Source:   "# Title",
				Metadata: json.RawMessage(`{}`),
			},
		},
	}
	path := filepath.Join(tmpDir, "test.ipynb")
	data, _ := json.MarshalIndent(nb, "", "  ")
	os.WriteFile(path, data, 0644)
	return path
}

func TestNotebookToolName(t *testing.T) {
	n := NewNotebookTool(".")
	if n.Name() != "NotebookEdit" {
		t.Errorf("expected 'NotebookEdit', got '%s'", n.Name())
	}
}

func TestNotebookToolRequiresPermission(t *testing.T) {
	n := NewNotebookTool(".")
	if !n.RequiresPermission() {
		t.Errorf("NotebookEdit should require permission")
	}
}

func TestNotebookToolInputSchema(t *testing.T) {
	n := NewNotebookTool(".")
	schema := n.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("expected type 'object'")
	}
}

func TestNotebookToolRead(t *testing.T) {
	tmpDir := t.TempDir()
	nbPath := createTestNotebook(t, tmpDir)

	n := NewNotebookTool(tmpDir)
	result := n.Execute(context.Background(), map[string]any{
		"file_path": nbPath,
		"operation": "read",
	})

	if result.IsError {
		t.Fatalf("expected success reading notebook, got error: %s", result.Content)
	}
}

func TestNotebookToolReadNonExistent(t *testing.T) {
	n := NewNotebookTool(".")
	result := n.Execute(context.Background(), map[string]any{
		"file_path": "/nonexistent/notebook.ipynb",
		"operation": "read",
	})

	if !result.IsError {
		t.Errorf("expected error for non-existent notebook")
	}
}

func TestNotebookToolMissingOperation(t *testing.T) {
	n := NewNotebookTool(".")
	result := n.Execute(context.Background(), map[string]any{
		"file_path": "test.ipynb",
	})
	if !result.IsError {
		t.Errorf("expected error for missing operation")
	}
}

func TestNotebookToolInvalidPath(t *testing.T) {
	n := NewNotebookTool(t.TempDir())
	result := n.Execute(context.Background(), map[string]any{
		"file_path": "../outside.ipynb",
		"operation": "read",
	})
	if !result.IsError {
		t.Errorf("expected error for path outside workspace")
	}
}

var _ tool.Tool = (*NotebookTool)(nil)
