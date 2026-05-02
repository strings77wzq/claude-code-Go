package diagnostic

import (
	"strings"
	"testing"
)

func TestDiagnosticFormatRedactsMetadata(t *testing.T) {
	d := Diagnostic{
		Component: "mcp",
		Severity:  SeverityError,
		Code:      "mcp.launch.blocked",
		Summary:   "launch blocked",
		Detail:    "Authorization: Bearer test-token",
		Retryable: false,
		Metadata: map[string]any{
			"server":        "fixture",
			"authorization": "Bearer metadata-token",
			"env": map[string]any{
				"API_TOKEN": "secret-token",
			},
		},
	}

	formatted := d.Format()

	for _, want := range []string{"ERROR", "mcp.launch.blocked", "launch blocked", "[REDACTED]"} {
		if !strings.Contains(formatted, want) {
			t.Fatalf("expected %q in formatted diagnostic:\n%s", want, formatted)
		}
	}
	for _, leaked := range []string{"test-token", "metadata-token", "secret-token"} {
		if strings.Contains(formatted, leaked) {
			t.Fatalf("diagnostic leaked %q:\n%s", leaked, formatted)
		}
	}
}

func TestDiagnosticTraceFieldsAreRedacted(t *testing.T) {
	d := Diagnostic{
		Component: "skills",
		Severity:  SeverityWarn,
		Code:      "skills.invalid",
		Summary:   "invalid skill",
		Metadata: map[string]any{
			"file":   "bad.json",
			"apiKey": "sk-test-secret-value",
		},
	}

	fields := d.TraceFields()

	if fields["component"] != "skills" || fields["severity"] != "WARN" || fields["code"] != "skills.invalid" {
		t.Fatalf("unexpected trace fields: %#v", fields)
	}
	if metadata, ok := fields["metadata"].(map[string]any); !ok || metadata["apiKey"] != RedactedMarker {
		t.Fatalf("expected redacted metadata, got %#v", fields["metadata"])
	}
}
