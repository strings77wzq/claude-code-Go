package builtin

import (
	"context"
	"strings"
	"testing"
)

func TestDiffToolPureGoFallback(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "line1\nline2",
		"new_content": "line1\nline2_changed\nline3",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if !strings.Contains(result.Content, "+") && !strings.Contains(result.Content, "-") {
		t.Error("diff should contain additions or deletions")
	}
}

func TestDiffToolMultilineContent(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "A\nB\nC\nD\nE",
		"new_content": "A\nB\nX\nD\nE",
	})

	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if result.Content == "" || result.Content == "No differences found." {
		t.Error("expected diff to show difference")
	}
}

func TestDiffToolEmptyContent(t *testing.T) {
	d := NewDiffTool()
	result := d.Execute(context.Background(), map[string]any{
		"old_content": "",
		"new_content": "hello",
	})
	if !result.IsError {
		t.Errorf("expected error for empty old_content")
	}
}
