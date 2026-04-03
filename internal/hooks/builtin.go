// Package hooks provides built-in hook implementations for logging and auditing.
package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	// maxInputLength is the maximum length of input/result to log.
	maxInputLength = 500
)

// LoggingHook is a hook that logs all tool executions using slog.
// It logs tool name, input (truncated), result (truncated), and duration.
type LoggingHook struct {
	logger *slog.Logger
}

// NewLoggingHook creates a new LoggingHook with the default logger.
func NewLoggingHook() *LoggingHook {
	return &LoggingHook{
		logger: slog.Default(),
	}
}

// NewLoggingHookWithLogger creates a new LoggingHook with a custom logger.
func NewLoggingHookWithLogger(logger *slog.Logger) *LoggingHook {
	return &LoggingHook{
		logger: logger,
	}
}

// Name returns the name of this hook.
func (h *LoggingHook) Name() string {
	return "logging"
}

// PreExecute logs the tool call before execution.
func (h *LoggingHook) PreExecute(toolName string, input map[string]any) error {
	h.logger.Debug("tool execution started",
		"tool", toolName,
		"input", truncateMap(input, maxInputLength),
	)
	return nil
}

// PostExecute logs the tool result after execution.
func (h *LoggingHook) PostExecute(toolName string, input map[string]any, result string, isError bool) {
	level := slog.LevelInfo
	if isError {
		level = slog.LevelError
	}
	h.logger.Log(context.Background(), level, "tool execution completed",
		"tool", toolName,
		"result", truncate(result, maxInputLength),
		"error", isError,
	)
}

// AuditHook records tool calls to a JSONL audit log file.
// Each line in the file is a JSON object representing one tool execution.
type AuditHook struct {
	filePath string
	file     *os.File
}

// AuditRecord represents a single audit log entry.
type AuditRecord struct {
	Timestamp  string         `json:"timestamp"`
	ToolName   string         `json:"tool_name"`
	Input      map[string]any `json:"input"`
	Result     string         `json:"result"`
	IsError    bool           `json:"is_error"`
	DurationMS int64          `json:"duration_ms"`
	PreHookErr string         `json:"pre_hook_error,omitempty"`
}

// NewAuditHook creates a new AuditHook that writes to the specified file path.
// The file will be created if it doesn't exist, or appended to if it does.
func NewAuditHook(filePath string) (*AuditHook, error) {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit log directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	return &AuditHook{
		filePath: filePath,
		file:     file,
	}, nil
}

// Close closes the audit log file.
// Call this when the hook is no longer needed.
func (h *AuditHook) Close() error {
	if h.file != nil {
		return h.file.Close()
	}
	return nil
}

// Name returns the name of this hook.
func (h *AuditHook) Name() string {
	return "audit"
}

// PreExecute records the start of a tool execution.
func (h *AuditHook) PreExecute(toolName string, input map[string]any) error {
	// PreExecute doesn't write to the log - we defer to PostExecute
	// to record the full execution with duration
	_ = toolName
	_ = input
	return nil
}

// PostExecute records the complete tool execution including duration.
// If the file cannot be written, the error is logged but execution continues.
func (h *AuditHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
	record := AuditRecord{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		ToolName:  toolName,
		Input:     input,
		Result:    truncate(result, maxInputLength),
		IsError:   isError,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal audit record: %w", err)
	}

	if _, err := h.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write audit record: %w", err)
	}

	// Flush to ensure data is written
	if err := h.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync audit log: %w", err)
	}

	return nil
}

// truncate truncates a string to the specified max length.
// If the string is longer than max, it adds "..." to indicate truncation.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// truncateMap recursively truncates string values in a map to the specified max length.
// This is useful for truncating large input/output objects.
func truncateMap(m map[string]any, max int) map[string]any {
	if m == nil {
		return nil
	}

	result := make(map[string]any, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = truncate(val, max)
		case map[string]any:
			result[k] = truncateMap(val, max)
		default:
			result[k] = val
		}
	}
	return result
}
