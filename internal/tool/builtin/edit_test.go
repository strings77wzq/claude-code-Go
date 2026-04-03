package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestEditToolBasic(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	os.WriteFile(tmpFile, []byte(content), 0644)

	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":  tmpFile,
		"old_string": "world",
		"new_string": "go-code",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	newContent, _ := os.ReadFile(tmpFile)
	if string(newContent) != "hello go-code" {
		t.Errorf("expected 'hello go-code', got '%s'", string(newContent))
	}
}

func TestEditToolNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	os.WriteFile(tmpFile, []byte(content), 0644)

	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":  tmpFile,
		"old_string": "nonexistent",
		"new_string": "something",
	})

	if !result.IsError {
		t.Errorf("expected error for nonexistent string")
	}
}

func TestEditToolMultipleOccurrences(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello hello world"
	os.WriteFile(tmpFile, []byte(content), 0644)

	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":  tmpFile,
		"old_string": "hello",
		"new_string": "hi",
	})

	if !result.IsError {
		t.Errorf("expected error for multiple occurrences without replace_all")
	}
}

func TestEditToolReplaceAll(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello hello world"
	os.WriteFile(tmpFile, []byte(content), 0644)

	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":   tmpFile,
		"old_string":  "hello",
		"new_string":  "hi",
		"replace_all": true,
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	newContent, _ := os.ReadFile(tmpFile)
	if string(newContent) != "hi hi world" {
		t.Errorf("expected 'hi hi world', got '%s'", string(newContent))
	}
}

func TestEditToolMissingFilePath(t *testing.T) {
	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"old_string": "old",
		"new_string": "new",
	})

	if !result.IsError {
		t.Errorf("expected error for missing file_path")
	}
}

func TestEditToolMissingOldString(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":  tmpFile,
		"new_string": "new",
	})

	if !result.IsError {
		t.Errorf("expected error for missing old_string")
	}
}

func TestEditToolName(t *testing.T) {
	editTool := NewEditTool()
	if editTool.Name() != "Edit" {
		t.Errorf("expected 'Edit', got '%s'", editTool.Name())
	}
}

func TestEditToolRequiresPermission(t *testing.T) {
	editTool := NewEditTool()
	if !editTool.RequiresPermission() {
		t.Errorf("Edit tool should require permission")
	}
}

func TestEditToolInputSchema(t *testing.T) {
	editTool := NewEditTool()
	schema := editTool.InputSchema()

	if schema["type"] != "object" {
		t.Errorf("expected type 'object', got '%v'", schema["type"])
	}

	props := schema["properties"].(map[string]any)
	if _, ok := props["file_path"]; !ok {
		t.Error("expected 'file_path' in properties")
	}
	if _, ok := props["old_string"]; !ok {
		t.Error("expected 'old_string' in properties")
	}
	if _, ok := props["new_string"]; !ok {
		t.Error("expected 'new_string' in properties")
	}
}

func TestEditToolFileNotFound(t *testing.T) {
	editTool := NewEditTool()
	result := editTool.Execute(context.Background(), map[string]any{
		"file_path":  "/nonexistent/file.txt",
		"old_string": "old",
		"new_string": "new",
	})

	if !result.IsError {
		t.Errorf("expected error for nonexistent file")
	}
}

var _ tool.Tool = (*EditTool)(nil)
