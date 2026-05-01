package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSaveAndLoadRoundtrip(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "session-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a session
	session := &Session{
		ID:           "sess_123",
		Model:        "claude-sonnet-4-20250514",
		StartTime:    time.UnixMilli(1234567890),
		EndTime:      time.UnixMilli(1234567990),
		TurnCount:    5,
		InputTokens:  1000,
		OutputTokens: 500,
	}

	// Create messages
	messages := []SessionMessage{
		{Type: "message", Role: "user", Content: "hello", Timestamp: time.UnixMilli(1234567890)},
		{Type: "message", Role: "assistant", Content: "Hi! How can I help?", Timestamp: time.UnixMilli(1234567891)},
		{Type: "message", Role: "user", Content: "Write a function", Timestamp: time.UnixMilli(1234567892)},
		{Type: "message", Role: "assistant", Content: "Here is a function...", Timestamp: time.UnixMilli(1234567893)},
	}

	// Save the session
	err = SaveSession(session, messages, tmpDir)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Find the saved file
	files, err := filepath.Glob(filepath.Join(tmpDir, "session-*.jsonl"))
	if err != nil {
		t.Fatalf("failed to find session files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 session file, got %d", len(files))
	}

	// Load the session
	loadedSession, loadedMessages, err := LoadSession(files[0])
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	// Verify session metadata
	if loadedSession.ID != session.ID {
		t.Errorf("session ID mismatch: got %s, want %s", loadedSession.ID, session.ID)
	}
	if loadedSession.Model != session.Model {
		t.Errorf("session model mismatch: got %s, want %s", loadedSession.Model, session.Model)
	}
	if loadedSession.TurnCount != session.TurnCount {
		t.Errorf("turn count mismatch: got %d, want %d", loadedSession.TurnCount, session.TurnCount)
	}
	if loadedSession.InputTokens != session.InputTokens {
		t.Errorf("input tokens mismatch: got %d, want %d", loadedSession.InputTokens, session.InputTokens)
	}
	if loadedSession.OutputTokens != session.OutputTokens {
		t.Errorf("output tokens mismatch: got %d, want %d", loadedSession.OutputTokens, session.OutputTokens)
	}

	// Verify messages
	if len(loadedMessages) != len(messages) {
		t.Errorf("message count mismatch: got %d, want %d", len(loadedMessages), len(messages))
	}
	for i, msg := range messages {
		if loadedMessages[i].Role != msg.Role {
			t.Errorf("message %d role mismatch: got %s, want %s", i, loadedMessages[i].Role, msg.Role)
		}
		if loadedMessages[i].Content != msg.Content {
			t.Errorf("message %d content mismatch: got %s, want %s", i, loadedMessages[i].Content, msg.Content)
		}
	}
}

func TestAtomicWrite(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "session-atomic-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	session := &Session{
		ID:        "sess_atomic",
		Model:     "test-model",
		StartTime: time.Now(),
		EndTime:   time.Now(),
		TurnCount: 1,
	}

	messages := []SessionMessage{
		{Type: "message", Role: "user", Content: "test", Timestamp: time.Now()},
	}

	// Save the session
	err = SaveSession(session, messages, tmpDir)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Check that the temp file no longer exists
	tmpFile := filepath.Join(tmpDir, "session-"+string(rune(session.StartTime.Unix()))+".jsonl.tmp")
	if _, err := os.Stat(tmpFile); err == nil {
		t.Error("temp file still exists after successful save")
	}

	// Check that the actual file exists
	files, err := filepath.Glob(filepath.Join(tmpDir, "session-*.jsonl"))
	if err != nil {
		t.Fatalf("failed to find session files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 session file, got %d", len(files))
	}

	// Verify file content is complete (not empty)
	info, err := os.Stat(files[0])
	if err != nil {
		t.Fatalf("failed to stat session file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("session file is empty")
	}
}

func TestLoadSessionWithInvalidLines(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "session-invalid-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with invalid lines
	testFile := filepath.Join(tmpDir, "session-test.jsonl")
	content := `{"type":"meta","session_id":"sess_123","model":"test","start_time_ms":1234567890,"end_time_ms":1234567990,"turn_count":1,"input_tokens":100,"output_tokens":50}
invalid json line
{"type":"message","role":"user","content":"hello","timestamp_ms":1234567890}
{"invalid":"json"}
{"type":"message","role":"assistant","content":"hi","timestamp_ms":1234567891}
unknown type line
{"type":"message","role":"user","content":"bye","timestamp_ms":1234567892}
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Load the session - should skip invalid lines
	session, messages, err := LoadSession(testFile)
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	// Verify session was loaded
	if session.ID != "sess_123" {
		t.Errorf("session ID mismatch: got %s, want sess_123", session.ID)
	}

	// Verify valid messages were loaded (3 out of 5 - skipping invalid lines)
	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}

	// Verify message content
	if len(messages) > 0 && messages[0].Content != "hello" {
		t.Errorf("first message content mismatch: got %s, want hello", messages[0].Content)
	}
	if len(messages) > 1 && messages[1].Content != "hi" {
		t.Errorf("second message content mismatch: got %s, want hi", messages[1].Content)
	}
	if len(messages) > 2 && messages[2].Content != "bye" {
		t.Errorf("third message content mismatch: got %s, want bye", messages[2].Content)
	}
}

func TestSaveSessionCreatesDirectory(t *testing.T) {
	// Use a non-existent nested directory
	tmpDir, err := os.MkdirTemp("", "session-dir-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	newDir := filepath.Join(tmpDir, "sessions", "2024", "01")

	session := &Session{
		ID:        "sess_newdir",
		Model:     "test-model",
		StartTime: time.Now(),
		EndTime:   time.Now(),
		TurnCount: 1,
	}

	messages := []SessionMessage{
		{Type: "message", Role: "user", Content: "test", Timestamp: time.Now()},
	}

	// Save should create the directory
	err = SaveSession(session, messages, newDir)
	if err != nil {
		t.Fatalf("failed to save session to new directory: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}

	// Verify file was created
	files, err := filepath.Glob(filepath.Join(newDir, "session-*.jsonl"))
	if err != nil {
		t.Fatalf("failed to find session files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 session file, got %d", len(files))
	}
}

func TestTraceReplayAndResumeSchema(t *testing.T) {
	tmpDir := t.TempDir()
	sessionFile := filepath.Join(tmpDir, "session-123.jsonl")

	meta := `{"type":"meta","session_id":"sess_trace","model":"test-model","start_time_ms":1234567890,"end_time_ms":0,"turn_count":0,"input_tokens":0,"output_tokens":0,"status":"running"}`
	if err := os.WriteFile(sessionFile, []byte(meta+"\n"), 0644); err != nil {
		t.Fatalf("failed to write trace file: %v", err)
	}

	if err := AppendTraceRequest(sessionFile, "test-model", 2); err != nil {
		t.Fatalf("failed to append request: %v", err)
	}
	if err := AppendTraceTool(sessionFile, "Read", map[string]any{"file_path": "README.md"}, "ok", 4); err != nil {
		t.Fatalf("failed to append tool: %v", err)
	}
	if err := AppendTracePermission(sessionFile, "Bash", "Deny", "rm -rf /"); err != nil {
		t.Fatalf("failed to append permission: %v", err)
	}
	if err := AppendTraceExtension(sessionFile, "lsp", "health_check", "unavailable", map[string]any{
		"reason": "not configured",
	}); err != nil {
		t.Fatalf("failed to append extension: %v", err)
	}
	if err := AppendSessionMessages(sessionFile, []SessionMessage{
		{Type: "message", Role: "user", Content: "hello", Timestamp: time.UnixMilli(1234567891)},
		{Type: "message", Role: "assistant", Content: "[tool result: ok]", Timestamp: time.UnixMilli(1234567892)},
	}); err != nil {
		t.Fatalf("failed to append messages: %v", err)
	}
	if err := AppendTraceStatus(sessionFile, "completed", 1, 10, 5); err != nil {
		t.Fatalf("failed to append status: %v", err)
	}

	loaded, messages, err := LoadSession(sessionFile)
	if err != nil {
		t.Fatalf("LoadSession() error = %v", err)
	}
	if loaded.Status != "completed" || loaded.TurnCount != 1 {
		t.Fatalf("loaded status/turns = %s/%d", loaded.Status, loaded.TurnCount)
	}
	if len(messages) != 2 {
		t.Fatalf("messages = %d, want 2", len(messages))
	}

	events, err := ReplaySessionFile(sessionFile)
	if err != nil {
		t.Fatalf("ReplaySessionFile() error = %v", err)
	}
	replay := FormatReplay(events)
	for _, want := range []string{"session sess_trace", "request model=test-model", "tool Read", "permission Bash", "extension lsp health_check status=unavailable", "status completed"} {
		if !strings.Contains(replay, want) {
			t.Fatalf("replay missing %q:\n%s", want, replay)
		}
	}

	sessions, err := ListSessions(tmpDir)
	if err != nil {
		t.Fatalf("ListSessions() error = %v", err)
	}
	if len(sessions) != 1 || sessions[0].TurnCount != 1 || sessions[0].Status != "completed" {
		t.Fatalf("unexpected session info: %#v", sessions)
	}
}

func TestTraceRedactsSecretsInStoredOutputAndReplay(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "session-redact.jsonl")

	if err := AppendTraceTool(sessionFile, "ProviderCall", map[string]any{
		"api_key":       "sk-ant-test-secret-value",
		"Authorization": "Bearer provider-token-value",
	}, "authorization: Bearer output-token-value", 1); err != nil {
		t.Fatalf("failed to append tool: %v", err)
	}

	raw, err := os.ReadFile(sessionFile)
	if err != nil {
		t.Fatalf("failed to read trace file: %v", err)
	}
	rawText := string(raw)
	for _, leaked := range []string{"sk-ant-test-secret-value", "provider-token-value", "output-token-value"} {
		if strings.Contains(rawText, leaked) {
			t.Fatalf("trace leaked secret %q:\n%s", leaked, rawText)
		}
	}
	if !strings.Contains(rawText, "[REDACTED]") {
		t.Fatalf("expected redaction marker in trace:\n%s", rawText)
	}

	events, err := ReplaySessionFile(sessionFile)
	if err != nil {
		t.Fatalf("ReplaySessionFile() error = %v", err)
	}
	replay := FormatReplay(events)
	for _, leaked := range []string{"sk-ant-test-secret-value", "provider-token-value", "output-token-value"} {
		if strings.Contains(replay, leaked) {
			t.Fatalf("replay leaked secret %q:\n%s", leaked, replay)
		}
	}
}
