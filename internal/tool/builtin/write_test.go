package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestWriteToolBasic(t *testing.T) {
	tmpDir := t.TempDir()
	writeTool := NewWriteTool(tmpDir)

	result := writeTool.Execute(context.Background(), map[string]any{
		"file_path": "notes/test.txt",
		"content":   "hello",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "notes", "test.txt"))
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("written content = %q, want hello", string(content))
	}
}

func TestWriteToolBlocksPathEscape(t *testing.T) {
	tmpDir := t.TempDir()
	writeTool := NewWriteTool(tmpDir)

	result := writeTool.Execute(context.Background(), map[string]any{
		"file_path": "../escape.txt",
		"content":   "nope",
	})

	if !result.IsError {
		t.Fatalf("expected path escape to be blocked")
	}
}

func TestWriteToolBlocksSymlinkEscape(t *testing.T) {
	tmpDir := t.TempDir()
	outsideDir := t.TempDir()
	linkPath := filepath.Join(tmpDir, "outside")
	if err := os.Symlink(outsideDir, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	writeTool := NewWriteTool(tmpDir)
	result := writeTool.Execute(context.Background(), map[string]any{
		"file_path": filepath.Join(linkPath, "escape.txt"),
		"content":   "nope",
	})

	if !result.IsError {
		t.Fatalf("expected symlink escape to be blocked")
	}
}

func TestWriteToolRefusesBinaryOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "payload.bin")
	os.WriteFile(binaryPath, []byte{0x00, 0x01, 0x02}, 0644)

	writeTool := NewWriteTool(tmpDir)
	result := writeTool.Execute(context.Background(), map[string]any{
		"file_path": binaryPath,
		"content":   "text",
	})

	if !result.IsError {
		t.Fatalf("expected binary overwrite to be blocked")
	}
}

var _ tool.Tool = (*WriteTool)(nil)
