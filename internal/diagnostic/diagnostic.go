package diagnostic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/sanitize"
)

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"

	RedactedMarker = sanitize.RedactedMarker
)

type Diagnostic struct {
	Component string
	Severity  Severity
	Code      string
	Summary   string
	Detail    string
	Retryable bool
	Metadata  map[string]any
}

func (d Diagnostic) Format() string {
	parts := []string{
		string(d.Severity),
		d.Component,
		d.Code,
		redactString(d.Summary),
	}
	if d.Detail != "" {
		parts = append(parts, redactString(d.Detail))
	}
	if len(d.Metadata) > 0 {
		keys := make([]string, 0, len(d.Metadata))
		for key := range d.Metadata {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		metadata := make([]string, 0, len(keys))
		redacted := redactValue(d.Metadata).(map[string]any)
		for _, key := range keys {
			metadata = append(metadata, fmt.Sprintf("%s=%v", key, redacted[key]))
		}
		parts = append(parts, strings.Join(metadata, " "))
	}
	return strings.Join(parts, " ")
}

func (d Diagnostic) TraceFields() map[string]any {
	return map[string]any{
		"component": d.Component,
		"severity":  string(d.Severity),
		"code":      d.Code,
		"summary":   redactString(d.Summary),
		"detail":    redactString(d.Detail),
		"retryable": d.Retryable,
		"metadata":  redactValue(d.Metadata),
	}
}

func redactValue(value any) any {
	return sanitize.RedactValue(value)
}

func sensitiveKey(key string) bool {
	return sanitize.SensitiveKey(key)
}

func redactString(value string) string {
	return sanitize.RedactString(value)
}
