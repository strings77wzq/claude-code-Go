// Package session provides session persistence in JSONL format.
package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a conversation session with metadata.
type Session struct {
	ID           string    `json:"session_id"`
	Model        string    `json:"model"`
	StartTime    time.Time `json:"start_time_ms"`
	EndTime      time.Time `json:"end_time_ms"`
	TurnCount    int       `json:"turn_count"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
}

// SessionMessage represents a single message within a session.
type SessionMessage struct {
	Type      string    `json:"type"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp_ms"`
}

// sessionMetaLine is the JSONL line type for session metadata.
type sessionMetaLine struct {
	Type         string `json:"type"`
	SessionID    string `json:"session_id"`
	Model        string `json:"model"`
	StartTime    int64  `json:"start_time_ms"`
	EndTime      int64  `json:"end_time_ms"`
	TurnCount    int    `json:"turn_count"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
}

// messageLine is the JSONL line type for messages.
type messageLine struct {
	Type      string `json:"type"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp_ms"`
}

// traceRequestLine is the JSONL line type for API request traces.
type traceRequestLine struct {
	Type          string `json:"type"`
	Model         string `json:"model"`
	MessagesCount int    `json:"messages_count"`
	Timestamp     int64  `json:"timestamp_ms"`
}

// traceResponseLine is the JSONL line type for API response traces.
type traceResponseLine struct {
	Type         string `json:"type"`
	StopReason   string `json:"stop_reason"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	Timestamp    int64  `json:"timestamp_ms"`
}

// traceToolLine is the JSONL line type for tool execution traces.
type traceToolLine struct {
	Type       string      `json:"type"`
	Name       string      `json:"name"`
	Input      interface{} `json:"input"`
	Output     string      `json:"output"`
	DurationMs int64       `json:"duration_ms"`
	Timestamp  int64       `json:"timestamp_ms"`
}

// traceErrorLine is the JSONL line type for error traces.
type traceErrorLine struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp_ms"`
}

// AppendTraceRequest appends a request trace line to the session file.
func AppendTraceRequest(filepath, model string, messagesCount int) error {
	return appendTraceLine(filepath, traceRequestLine{
		Type:          "request",
		Model:         model,
		MessagesCount: messagesCount,
		Timestamp:     time.Now().UnixMilli(),
	})
}

// AppendTraceResponse appends a response trace line to the session file.
func AppendTraceResponse(filepath, stopReason string, inputTokens, outputTokens int) error {
	return appendTraceLine(filepath, traceResponseLine{
		Type:         "response",
		StopReason:   stopReason,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Timestamp:    time.Now().UnixMilli(),
	})
}

// AppendTraceTool appends a tool execution trace line to the session file.
func AppendTraceTool(filepath, name string, input interface{}, output string, durationMs int64) error {
	return appendTraceLine(filepath, traceToolLine{
		Type:       "tool",
		Name:       name,
		Input:      input,
		Output:     output,
		DurationMs: durationMs,
		Timestamp:  time.Now().UnixMilli(),
	})
}

// AppendTraceError appends an error trace line to the session file.
func AppendTraceError(filepath, message string) error {
	return appendTraceLine(filepath, traceErrorLine{
		Type:      "error",
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
	})
}

// appendTraceLine appends a JSON line to the session file.
func appendTraceLine(filepath string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal trace line: %w", err)
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open session file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write trace line: %w", err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// SaveSession saves a session and its messages to a JSONL file.
// The file is written atomically: first to a temp file, then renamed to the final path.
// Filename format: {dir}/session-{timestamp}.jsonl
func SaveSession(s *Session, messages []SessionMessage, dir string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := s.StartTime.Unix()
	filename := fmt.Sprintf("session-%d.jsonl", timestamp)
	filepath := filepath.Join(dir, filename)

	// Write to temporary file first (atomic write)
	tmpFile := filepath + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	// Write session metadata as first line
	metaLine := sessionMetaLine{
		Type:         "meta",
		SessionID:    s.ID,
		Model:        s.Model,
		StartTime:    s.StartTime.UnixMilli(),
		EndTime:      s.EndTime.UnixMilli(),
		TurnCount:    s.TurnCount,
		InputTokens:  s.InputTokens,
		OutputTokens: s.OutputTokens,
	}
	metaBytes, err := json.Marshal(metaLine)
	if err != nil {
		return fmt.Errorf("failed to marshal session metadata: %w", err)
	}
	if _, err := f.Write(metaBytes); err != nil {
		return fmt.Errorf("failed to write session metadata: %w", err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// Write each message as subsequent lines
	for _, msg := range messages {
		msgLine := messageLine{
			Type:      "message",
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp.UnixMilli(),
		}
		msgBytes, err := json.Marshal(msgLine)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
		if _, err := f.Write(msgBytes); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		if _, err := f.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	// Flush and sync to ensure all data is written
	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, filepath); err != nil {
		// Clean up temp file on failure
		os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadSession loads a session and its messages from a JSONL file.
// Invalid lines are skipped with a warning.
func LoadSession(filepath string) (*Session, []SessionMessage, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open session file: %w", err)
	}
	defer f.Close()

	var session *Session
	var messages []SessionMessage
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Try to determine the type of line
		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(line), &base); err != nil {
			fmt.Printf("Warning: skipping invalid JSON line %d: %v\n", lineNum, err)
			continue
		}

		switch base.Type {
		case "meta":
			var meta sessionMetaLine
			if err := json.Unmarshal([]byte(line), &meta); err != nil {
				fmt.Printf("Warning: skipping invalid meta line %d: %v\n", lineNum, err)
				continue
			}
			session = &Session{
				ID:           meta.SessionID,
				Model:        meta.Model,
				StartTime:    time.UnixMilli(meta.StartTime),
				EndTime:      time.UnixMilli(meta.EndTime),
				TurnCount:    meta.TurnCount,
				InputTokens:  meta.InputTokens,
				OutputTokens: meta.OutputTokens,
			}

		case "message":
			var msg messageLine
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				fmt.Printf("Warning: skipping invalid message line %d: %v\n", lineNum, err)
				continue
			}
			messages = append(messages, SessionMessage{
				Type:      msg.Type,
				Role:      msg.Role,
				Content:   msg.Content,
				Timestamp: time.UnixMilli(msg.Timestamp),
			})

		default:
			fmt.Printf("Warning: skipping unknown line type at line %d: %s\n", lineNum, base.Type)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading session file: %w", err)
	}

	if session == nil {
		return nil, nil, fmt.Errorf("no session metadata found in file")
	}

	return session, messages, nil
}
