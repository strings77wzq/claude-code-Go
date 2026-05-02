package diagnostic

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"

	RedactedMarker = "[REDACTED]"
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
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			if sensitiveKey(key) {
				out[key] = RedactedMarker
				continue
			}
			out[key] = redactValue(item)
		}
		return out
	case map[string]string:
		out := make(map[string]any, len(v))
		for key, item := range v {
			if sensitiveKey(key) {
				out[key] = RedactedMarker
				continue
			}
			out[key] = redactString(item)
		}
		return out
	case []any:
		out := make([]any, 0, len(v))
		for _, item := range v {
			out = append(out, redactValue(item))
		}
		return out
	case string:
		return redactString(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return v
		}
		var decoded any
		if err := json.Unmarshal(data, &decoded); err != nil {
			return v
		}
		switch decoded.(type) {
		case map[string]any, []any:
			return redactValue(decoded)
		default:
			return decoded
		}
	}
}

func sensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(key))
	for _, token := range []string{"apikey", "authorization", "password", "secret"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return strings.Contains(normalized, "token") && !strings.Contains(normalized, "tokens")
}

var (
	bearerPattern = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)
	keyPattern    = regexp.MustCompile(`sk-[A-Za-z0-9][A-Za-z0-9._-]{8,}`)
	tokenPattern  = regexp.MustCompile(`(?i)\b(api[_-]?key|authorization|password|secret|token)[_:=.-][A-Za-z0-9._~+/=-]+`)
)

func redactString(value string) string {
	value = bearerPattern.ReplaceAllString(value, "Bearer "+RedactedMarker)
	value = keyPattern.ReplaceAllString(value, RedactedMarker)
	value = tokenPattern.ReplaceAllString(value, RedactedMarker)
	return value
}
