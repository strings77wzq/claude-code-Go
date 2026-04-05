package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggingHookName(t *testing.T) {
	hook := NewLoggingHook()
	if hook.Name() != "logging" {
		t.Errorf("Expected name 'logging', got '%s'", hook.Name())
	}
}

func TestLoggingHookPreExecute(t *testing.T) {
	hook := NewLoggingHook()
	input := map[string]any{"file": "test.txt", "limit": 100}

	err := hook.PreExecute("Read", input)
	if err != nil {
		t.Errorf("PreExecute should not return error: %v", err)
	}
}

func TestLoggingHookPostExecuteNoError(t *testing.T) {
	hook := NewLoggingHook()
	input := map[string]any{"file": "test.txt"}
	result := "file content"

	hook.PostExecute("Read", input, result, false)
}

func TestLoggingHookPostExecuteWithError(t *testing.T) {
	hook := NewLoggingHook()
	input := map[string]any{"file": "test.txt"}
	result := "error: file not found"

	hook.PostExecute("Read", input, result, true)
}

func TestAuditHookName(t *testing.T) {
	tmpDir := t.TempDir()
	hook, err := NewAuditHook(filepath.Join(tmpDir, "audit.jsonl"))
	if err != nil {
		t.Fatalf("Failed to create AuditHook: %v", err)
	}
	defer hook.Close()

	if hook.Name() != "audit" {
		t.Errorf("Expected name 'audit', got '%s'", hook.Name())
	}
}

func TestAuditHookPreExecute(t *testing.T) {
	tmpDir := t.TempDir()
	hook, err := NewAuditHook(filepath.Join(tmpDir, "audit.jsonl"))
	if err != nil {
		t.Fatalf("Failed to create AuditHook: %v", err)
	}
	defer hook.Close()

	input := map[string]any{"file": "test.txt"}
	err = hook.PreExecute("Read", input)
	if err != nil {
		t.Errorf("PreExecute should not return error: %v", err)
	}
}

func TestAuditHookPostExecute(t *testing.T) {
	tmpDir := t.TempDir()
	hook, err := NewAuditHook(filepath.Join(tmpDir, "audit.jsonl"))
	if err != nil {
		t.Fatalf("Failed to create AuditHook: %v", err)
	}
	defer hook.Close()

	input := map[string]any{"file": "test.txt"}
	result := "file content"

	err = hook.PostExecute("Read", input, result, false)
	if err != nil {
		t.Errorf("PostExecute should not return error: %v", err)
	}
}

func TestAuditHookFileWriting(t *testing.T) {
	tmpDir := t.TempDir()
	auditPath := filepath.Join(tmpDir, "audit.jsonl")
	hook, err := NewAuditHook(auditPath)
	if err != nil {
		t.Fatalf("Failed to create AuditHook: %v", err)
	}
	defer hook.Close()

	input := map[string]any{"file": "test.txt"}
	result := "file content here"

	_ = hook.PostExecute("Read", input, result, false)

	data, err := os.ReadFile(auditPath)
	if err != nil {
		t.Errorf("Failed to read audit file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Audit file should not be empty")
	}

	if !strings.Contains(string(data), "Read") {
		t.Error("Audit record should contain tool name")
	}
}

func TestAuditHookPostExecuteWithError(t *testing.T) {
	tmpDir := t.TempDir()
	hook, err := NewAuditHook(filepath.Join(tmpDir, "audit.jsonl"))
	if err != nil {
		t.Fatalf("Failed to create AuditHook: %v", err)
	}
	defer hook.Close()

	input := map[string]any{"file": "nonexistent.txt"}
	result := "error: file not found"

	err = hook.PostExecute("Read", input, result, true)
	if err != nil {
		t.Errorf("PostExecute should not return error: %v", err)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 10, "short"},
		{"verylongstring", 10, "verylongst..."},
		{"exact", 5, "exact"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
		}
	}
}

func TestTruncateMap(t *testing.T) {
	input := map[string]any{
		"short": "hi",
		"long":  "this is a very long string that should be truncated",
		"nested": map[string]any{
			"inner": "another very long string that needs truncation",
		},
		"number": 42,
	}

	result := truncateMap(input, 20)

	if result["short"] != "hi" {
		t.Errorf("short should remain unchanged, got %v", result["short"])
	}

	longStr, ok := result["long"].(string)
	if !ok {
		t.Fatal("long should be a string")
	}
	if len(longStr) > 23 { // 20 + 3 for "..."
		t.Errorf("long should be at most 23 chars, got %d: %s", len(longStr), longStr)
	}

	nested, ok := result["nested"].(map[string]any)
	if !ok {
		t.Fatal("nested should be a map")
	}
	innerStr, ok := nested["inner"].(string)
	if !ok {
		t.Fatal("inner should be a string")
	}
	if len(innerStr) > 23 { // 20 + 3 for "..."
		t.Errorf("inner should be at most 23 chars, got %d: %s", len(innerStr), innerStr)
	}

	if result["number"] != 42 {
		t.Errorf("number should remain unchanged, got %v", result["number"])
	}
}

func TestTruncateMapNil(t *testing.T) {
	result := truncateMap(nil, 10)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestAuditHookNewAuditHookNonExistentDir(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent", "subdir", "audit.jsonl")

	_, err := NewAuditHook(nonExistentPath)
	if err != nil {
		t.Errorf("NewAuditHook should create directories: %v", err)
	}
}
