package sanitize

import (
	"strings"
	"testing"
)

func TestRedactString_Empty(t *testing.T) {
	got := RedactString("")
	if got != "" {
		t.Errorf("RedactString('') = %q, want ''", got)
	}
}

func TestRedactString_NoSecrets(t *testing.T) {
	input := "hello world this is safe"
	got := RedactString(input)
	if got != input {
		t.Errorf("RedactString(%q) = %q, want unchanged", input, got)
	}
}

func TestRedactString_BearerToken(t *testing.T) {
	got := RedactString("Authorization: Bearer sk-ant-test12345")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected redaction, got %q", got)
	}
	if strings.Contains(got, "sk-ant-test12345") {
		t.Errorf("secret leaked: %q", got)
	}
}

func TestRedactString_SkPrefix(t *testing.T) {
	got := RedactString("key=sk-abc123def456ghi789jkl")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected redaction, got %q", got)
	}
}

func TestRedactString_APIKeyInline(t *testing.T) {
	got := RedactString("api_key=my-secret-key-value")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected redaction, got %q", got)
	}
	if strings.Contains(got, "my-secret-key-value") {
		t.Errorf("secret leaked: %q", got)
	}
}

func TestRedactString_PasswordAssignment(t *testing.T) {
	got := RedactString("password := super secret 123")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected redaction, got %q", got)
	}
}

func TestRedactString_CaseInsensitive(t *testing.T) {
	got := RedactString("API_KEY=SOMEVALUE")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected case-insensitive redaction, got %q", got)
	}
}

func TestRedactString_TokenInQuery(t *testing.T) {
	got := RedactString("token=ghp_1234567890abcdef")
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected token redaction, got %q", got)
	}
}

func TestRedactStringTruncated_Short(t *testing.T) {
	got := RedactStringTruncated("safe", 10)
	if got != "safe" {
		t.Errorf("RedactStringTruncated('safe', 10) = %q, want 'safe'", got)
	}
}

func TestRedactStringTruncated_ExceedsLimit(t *testing.T) {
	got := RedactStringTruncated("this is a very long safe string that exceeds the limit", 10)
	if !strings.Contains(got, "...(truncated)") {
		t.Errorf("expected truncation marker, got %q", got)
	}
}

func TestRedactStringTruncated_WithSecret(t *testing.T) {
	got := RedactStringTruncated("Bearer sk-test-token-value-here", 100)
	if strings.Contains(got, "sk-test-token-value-here") {
		t.Errorf("secret leaked after redaction: %q", got)
	}
	if !strings.Contains(got, RedactedMarker) {
		t.Errorf("expected redaction marker, got %q", got)
	}
}

func TestSensitiveKey_Table(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"api_key", true},
		{"API_KEY", true},
		{"authorization", true},
		{"Authorization", true},
		{"secret", true},
		{"password", true},
		{"my_password_field", true},
		{"token", true},
		{"access_token", true},
		{"tokens", false},
		{"input_tokens", false},
		{"output_tokens", false},
		{"name", false},
		{"description", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := SensitiveKey(tt.key)
			if got != tt.want {
				t.Errorf("SensitiveKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestRedactValue_MapSensitiveKey(t *testing.T) {
	input := map[string]any{
		"safe":   "hello",
		"apiKey": "sk-test123",
	}
	got := RedactValue(input).(map[string]any)
	if got["safe"] != "hello" {
		t.Errorf("safe field changed: %v", got["safe"])
	}
	if got["apiKey"] != RedactedMarker {
		t.Errorf("apiKey not redacted: got %v", got["apiKey"])
	}
}

func TestRedactValue_NestedMap(t *testing.T) {
	input := map[string]any{
		"env": map[string]any{
			"SECRET": "should be hidden",
		},
	}
	got := RedactValue(input).(map[string]any)
	env := got["env"].(map[string]any)
	if env["SECRET"] != RedactedMarker {
		t.Errorf("nested secret not redacted: %v", env["SECRET"])
	}
}

func TestRedactValue_StringSlice(t *testing.T) {
	input := []any{"hello", "Bearer mytoken", "world"}
	got := RedactValue(input).([]any)
	if got[0] != "hello" {
		t.Errorf("first element changed: %v", got[0])
	}
	if !strings.Contains(got[1].(string), RedactedMarker) {
		t.Errorf("second element not redacted: %v", got[1])
	}
}

func TestRedactValue_NonStringTypes(t *testing.T) {
	tests := []any{42, 3.14, true, nil}
	for _, v := range tests {
		got := RedactValue(v)
		if got != v {
			t.Errorf("RedactValue(%#v) = %#v, want unchanged", v, got)
		}
	}
}
