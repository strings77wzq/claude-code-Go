package builtin

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestBashToolExecute(t *testing.T) {
	tmpDir := t.TempDir()
	bashTool := NewBashTool(tmpDir)

	result := bashTool.Execute(context.Background(), map[string]any{
		"command": "echo hello",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
	if result.Content != "hello\n" {
		t.Errorf("expected 'hello\\n', got '%s'", result.Content)
	}
}

func TestBashToolExitCode(t *testing.T) {
	tmpDir := t.TempDir()
	bashTool := NewBashTool(tmpDir)

	result := bashTool.Execute(context.Background(), map[string]any{
		"command": "exit 1",
	})

	if !result.IsError {
		t.Errorf("expected error for exit code 1")
	}
}

func TestBashToolTimeout(t *testing.T) {
	tmpDir := t.TempDir()
	bashTool := NewBashTool(tmpDir)

	start := time.Now()
	result := bashTool.Execute(context.Background(), map[string]any{
		"command": "sleep 5",
		"timeout": 1,
	})
	elapsed := time.Since(start)

	if !result.IsError {
		t.Errorf("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Errorf("timeout should be around 1 second, got %v", elapsed)
	}
}

func TestBashToolOutputTruncation(t *testing.T) {
	tmpDir := t.TempDir()
	bashTool := NewBashTool(tmpDir)

	largeOutput := make([]byte, 150*1024)
	for i := range largeOutput {
		largeOutput[i] = 'a'
	}

	result := bashTool.Execute(context.Background(), map[string]any{
		"command": "echo '" + string(largeOutput) + "'",
	})

	if len(result.Content) > 110*1024 {
		t.Errorf("output should be truncated to ~100KB, got %d bytes", len(result.Content))
	}
}

func TestBashToolWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)

	tmpFile := filepath.Join(subDir, "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	bashTool := NewBashTool(tmpDir)
	result := bashTool.Execute(context.Background(), map[string]any{
		"command": "cat subdir/test.txt",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}
	if result.Content != "test" {
		t.Errorf("expected 'test\\n', got '%s'", result.Content)
	}
}

func TestBashToolName(t *testing.T) {
	bashTool := NewBashTool(".")
	if bashTool.Name() != "Bash" {
		t.Errorf("expected 'Bash', got '%s'", bashTool.Name())
	}
}

func TestBashToolRequiresPermission(t *testing.T) {
	bashTool := NewBashTool(".")
	if !bashTool.RequiresPermission() {
		t.Errorf("Bash tool should require permission")
	}
}

func TestBashToolInputSchema(t *testing.T) {
	bashTool := NewBashTool(".")
	schema := bashTool.InputSchema()

	if schema["type"] != "object" {
		t.Errorf("expected type 'object', got '%v'", schema["type"])
	}

	props := schema["properties"].(map[string]any)
	if _, ok := props["command"]; !ok {
		t.Error("expected 'command' in properties")
	}
	if _, ok := props["timeout"]; !ok {
		t.Error("expected 'timeout' in properties")
	}
}

func TestBashToolMissingCommand(t *testing.T) {
	bashTool := NewBashTool(".")
	result := bashTool.Execute(context.Background(), map[string]any{})

	if !result.IsError {
		t.Errorf("expected error for missing command")
	}
}

var _ tool.Tool = (*BashTool)(nil)
