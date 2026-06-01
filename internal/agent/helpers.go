package agent

import (
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/sanitize"
)

func summarizePermissionInput(toolName string, input map[string]any) string {
	var summary string
	switch toolName {
	case "Bash":
		if command, ok := input["command"].(string); ok {
			summary = command
		}
	case "Read", "Write", "Edit":
		if path, ok := input["file_path"].(string); ok {
			summary = path
		}
	case "Glob", "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			summary = pattern
		}
	}
	if summary == "" {
		summary = "tool input"
	}
	return sanitizePermissionSummary(summary)
}

func sanitizePermissionSummary(summary string) string {
	return sanitize.RedactStringTruncated(strings.TrimSpace(summary), 200)
}
