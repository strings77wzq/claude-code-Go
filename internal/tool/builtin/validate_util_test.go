package builtin

import (
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

func TestResolvePathInsideWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	path, err := ResolvePath("test.txt", tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path == "" {
		t.Error("expected non-empty resolved path")
	}
}

func TestResolvePathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := ResolvePath("../outside.txt", tmpDir)
	if err == nil {
		t.Errorf("expected error for path traversal")
	}
}

func TestResolvePathBlockedSystem(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := ResolvePath("/dev/sda", tmpDir)
	if err == nil {
		t.Errorf("expected error for blocked system path")
	}
}

func TestValidatePathInsideWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	err := ValidatePath("test.txt", tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidatePathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	err := ValidatePath("../outside.txt", tmpDir)
	if err == nil {
		t.Errorf("expected error for path traversal")
	}
}

func TestValidatePathBlockedDev(t *testing.T) {
	tmpDir := t.TempDir()
	err := ValidatePath("/dev/null", tmpDir)
	if err == nil {
		t.Errorf("expected error for blocked /dev path")
	}
}

// Ensure validate.go uses the permission package
var _ = permission.LevelReadOnly

func TestValidatePathEmpty(t *testing.T) {
	// Empty path joined with working dir produces a valid workspace path,
	// so ValidatePath does not reject it
	err := ValidatePath("", t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error for empty path: %v", err)
	}
}
