package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/session"
)

// GetHistory returns a copy of the current conversation history.
func (a *Agent) GetHistory() *History {
	return a.history
}

// SetModel updates the model used by the agent at runtime.
func (a *Agent) SetModel(model string) {
	a.model = model
	if api, ok := a.apiClient.(interface{ SetModel(string) }); ok {
		api.SetModel(model)
	}
}

// Model returns the current model name.
func (a *Agent) Model() string {
	return a.model
}

// TraceFilePath returns the path to the trace file for the current session.
func (a *Agent) TraceFilePath() string {
	return a.traceFilePath
}

// ClearHistory resets the conversation history.
func (a *Agent) ClearHistory() {
	a.history = NewHistory()
}

// generateSessionID creates a unique session identifier.
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixMilli())
}

// getSessionsDir returns the directory for storing session files.
func getSessionsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".claude-code-go/sessions"
	}
	return filepath.Join(homeDir, ".claude-code-go", "sessions")
}

// initTraceFile initializes the trace file with meta information and returns its path.
func (a *Agent) initTraceFile() string {
	if a.sessionID == "" {
		return ""
	}

	dir := getSessionsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}

	timestamp := a.startTime.Unix()
	filename := fmt.Sprintf("session-%d.jsonl", timestamp)
	path := filepath.Join(dir, filename)
	a.traceFilePath = path

	metaLine := map[string]interface{}{
		"type":          "meta",
		"session_id":    a.sessionID,
		"model":         a.model,
		"start_time_ms": a.startTime.UnixMilli(),
		"end_time_ms":   0,
		"turn_count":    0,
		"input_tokens":  0,
		"output_tokens": 0,
		"status":        "running",
	}

	data, err := json.Marshal(metaLine)
	if err != nil {
		return ""
	}

	f, err := os.Create(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return ""
	}
	if _, err := f.WriteString("\n"); err != nil {
		return ""
	}

	return path
}

// traceRequest logs the request to the trace file.
func (a *Agent) traceRequest(model string, messagesCount int) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceRequest(a.traceFilePath, model, messagesCount)
}

// traceResponse logs the response to the trace file.
func (a *Agent) traceResponse(stopReason string, inputTokens, outputTokens int) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceResponse(a.traceFilePath, stopReason, inputTokens, outputTokens)
}

// traceTool logs tool execution to the trace file.
func (a *Agent) traceTool(name string, input interface{}, output string, durationMs int64) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceTool(a.traceFilePath, name, input, output, durationMs)
}

// traceError logs an error to the trace file.
func (a *Agent) traceError(message string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceError(a.traceFilePath, message)
}

// tracePermission logs permission decisions to the trace file.
func (a *Agent) tracePermission(toolName string, decision permission.Decision, reason permission.Reason, summary string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTracePermissionWithReason(a.traceFilePath, toolName, string(decision), summary, string(reason))
}

func (a *Agent) traceThinking(text string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceRuntime(a.traceFilePath, a.sessionID, "thinking", text)
}

func (a *Agent) traceRuntime(event, summary string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceRuntime(a.traceFilePath, a.sessionID, event, summary)
}
