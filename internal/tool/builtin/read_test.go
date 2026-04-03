package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestReadToolBasic(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\n"
	os.WriteFile(tmpFile, []byte(content), 0644)

	readTool := NewReadTool()
	result := readTool.Execute(context.Background(), map[string]any{
		"file_path": tmpFile,
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	expected := "1: line1\n2: line2\n3: line3\n"
	if result.Content != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result.Content)
	}
}

func TestReadToolOffsetLimit(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\nline4\nline5\n"
	os.WriteFile(tmpFile, []byte(content), 0644)

	readTool := NewReadTool()

	result := readTool.Execute(context.Background(), map[string]any{
		"file_path": tmpFile,
		"offset":    1,
		"limit":     2,
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	expected := "2: line2\n3: line3\n"
	if result.Content != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result.Content)
	}
}

func TestReadToolFileNotFound(t *testing.T) {
	readTool := NewReadTool()
	result := readTool.Execute(context.Background(), map[string]any{
		"file_path": "/nonexistent/file.txt",
	})

	if !result.IsError {
		t.Errorf("expected error for nonexistent file")
	}
}

func TestReadToolDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	readTool := NewReadTool()
	result := readTool.Execute(context.Background(), map[string]any{
		"file_path": tmpDir,
	})

	if !result.IsError {
		t.Errorf("expected error for directory")
	}
}

func TestReadToolMissingFilePath(t *testing.T) {
	readTool := NewReadTool()
	result := readTool.Execute(context.Background(), map[string]any{})

	if !result.IsError {
		t.Errorf("expected error for missing file_path")
	}
}

func TestReadToolName(t *testing.T) {
	readTool := NewReadTool()
	if readTool.Name() != "Read" {
		t.Errorf("expected 'Read', got '%s'", readTool.Name())
	}
}

func TestReadToolRequiresPermission(t *testing.T) {
	readTool := NewReadTool()
	if readTool.RequiresPermission() {
		t.Errorf("Read tool should not require permission")
	}
}

func TestReadToolInputSchema(t *testing.T) {
	readTool := NewReadTool()
	schema := readTool.InputSchema()

	if schema["type"] != "object" {
		t.Errorf("expected type 'object', got '%v'", schema["type"])
	}

	props := schema["properties"].(map[string]any)
	if _, ok := props["file_path"]; !ok {
		t.Error("expected 'file_path' in properties")
	}
	if _, ok := props["offset"]; !ok {
		t.Error("expected 'offset' in properties")
	}
	if _, ok := props["limit"]; !ok {
		t.Error("expected 'limit' in properties")
	}
}

var _ tool.Tool = (*ReadTool)(nil)
