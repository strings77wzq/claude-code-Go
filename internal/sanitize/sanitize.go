// Package sanitize provides token and secret redaction for trace output,
// permission summaries, and diagnostic reports.
package sanitize

import (
	"regexp"
	"strings"
)

// RedactedMarker is used to replace sensitive values in output.
const RedactedMarker = "[REDACTED]"

var (
	bearerPattern = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)
	apiKeyPattern = regexp.MustCompile(`sk-[A-Za-z0-9][A-Za-z0-9._-]{8,}`)
	tokenPattern  = regexp.MustCompile(`(?i)\b(api[_-]?key|authorization|password|secret|token)[_:=.-][A-Za-z0-9._~+/=-]+`)
	// Extended patterns for permission summaries (covers inline assignments like "key=value" too).
	inlineSecretPattern = regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password|authorization)(\s*[:=]\s*|\s+)[^\s;&|]+`)
)

// RedactString replaces known secret patterns (API keys, bearer tokens, passwords)
// in the given string with RedactedMarker.
func RedactString(s string) string {
	s = bearerPattern.ReplaceAllStringFunc(s, func(match string) string {
		if strings.HasPrefix(strings.ToLower(match), "bearer ") {
			return "Bearer " + RedactedMarker
		}
		return RedactedMarker
	})
	s = apiKeyPattern.ReplaceAllString(s, RedactedMarker)
	s = tokenPattern.ReplaceAllString(s, RedactedMarker)
	s = inlineSecretPattern.ReplaceAllStringFunc(s, func(match string) string {
		if strings.HasPrefix(strings.ToLower(match), "bearer ") {
			return "Bearer " + RedactedMarker
		}
		if strings.HasPrefix(match, "sk-") {
			return "sk-" + RedactedMarker
		}
		parts := regexp.MustCompile(`\s*[:=]\s*|\s+`).Split(match, 2)
		if len(parts) > 0 {
			return parts[0] + "=" + RedactedMarker
		}
		return RedactedMarker
	})
	return s
}

// RedactStringTruncated is like RedactString but truncates the result to maxLen
// characters, appending "...(truncated)" when truncated.
func RedactStringTruncated(s string, maxLen int) string {
	s = RedactString(s)
	if len(s) > maxLen {
		return s[:maxLen] + "...(truncated)"
	}
	return s
}

// SensitiveKey reports whether the given key names a field whose value
// should be redacted in structured output.
func SensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, token := range []string{"apikey", "authorization", "secret", "password", "api_key"} {
		if strings.Contains(lower, token) {
			return true
		}
	}
	// Check for "token" but exclude safe plural forms like "input_tokens",
	// "output_tokens" which are common non-sensitive trace fields.
	if strings.Contains(lower, "token") && !strings.Contains(lower, "tokens") {
		return true
	}
	return false
}

// RedactValue recursively redacts sensitive values in structured data.
// Supported types: map[string]any, []any, string. Other types are returned
// unchanged.
func RedactValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for key, item := range val {
			if SensitiveKey(key) {
				out[key] = RedactedMarker
				continue
			}
			out[key] = RedactValue(item)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = RedactValue(item)
		}
		return out
	case string:
		return RedactString(val)
	default:
		return v
	}
}
