package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReplayCommandWithFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	sessionFile := filepath.Join(tmpDir, "session-1.jsonl")
	content := `{"type":"meta","session_id":"sess_1","model":"test-model","start_time_ms":1,"status":"running"}
{"type":"message","role":"user","content":"hello","timestamp_ms":2}
{"type":"status","status":"completed","turn_count":1,"input_tokens":3,"output_tokens":4,"timestamp_ms":5}
`
	if err := os.WriteFile(sessionFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write session file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runReplayCommand([]string{sessionFile}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runReplayCommand() code = %d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "session sess_1") || !strings.Contains(stdout.String(), "status completed") {
		t.Fatalf("unexpected replay output:\n%s", stdout.String())
	}
}
