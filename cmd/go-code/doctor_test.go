package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDoctorHealthyOffline(t *testing.T) {
	isolateDoctorProviderEnv(t)
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test","model":"claude-sonnet-4-6-20251001"}`)
	writeDoctorDocs(t, workspace)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out.String())
	}
	output := out.String()
	for _, want := range []string{
		"[PASS] binary",
		"[PASS] configuration",
		"[SKIP] provider",
		"[PASS] provider profile",
		"[PASS] session directory",
		"[PASS] tools",
		"[SKIP] mcp",
		"[SKIP] lsp",
		"[SKIP] hooks",
		"[SKIP] skills",
		"[PASS] documentation",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in doctor output:\n%s", want, output)
		}
	}
	if strings.Contains(output, "sk-ant-test") {
		t.Fatalf("doctor output leaked API key:\n%s", output)
	}
}

func isolateDoctorProviderEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_BASE_URL", "")
	t.Setenv("ANTHROPIC_MODEL", "")
	t.Setenv("LLM_PROVIDER", "")
	t.Setenv("GO_CODE_API_KEY", "")
	t.Setenv("GO_CODE_BASE_URL", "")
	t.Setenv("GO_CODE_MODEL", "")
	t.Setenv("GO_CODE_PROVIDER", "")
	t.Setenv("DEEPSEEK_API_KEY", "")
}

func TestRunDoctorMissingAPIKey(t *testing.T) {
	isolateDoctorProviderEnv(t)
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorDocs(t, workspace)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code == 0 {
		t.Fatalf("expected non-zero exit code for missing API key")
	}
	output := out.String()
	if !strings.Contains(output, "[FAIL] configuration") {
		t.Fatalf("expected configuration failure:\n%s", output)
	}
	if !strings.Contains(output, "ANTHROPIC_API_KEY") {
		t.Fatalf("expected remediation to mention ANTHROPIC_API_KEY:\n%s", output)
	}
}

func TestRunDoctorInvalidSessionDirectory(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)

	blockingPath := filepath.Join(homeDir, ".claude-code-go")
	if err := os.WriteFile(blockingPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create blocking path: %v", err)
	}

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code == 0 {
		t.Fatalf("expected non-zero exit code for invalid session directory")
	}
	output := out.String()
	if !strings.Contains(output, "[FAIL] session directory") {
		t.Fatalf("expected session directory failure:\n%s", output)
	}
}

func TestRunDoctorOfflineSkipsProviderProbe(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "[SKIP] provider:") || !strings.Contains(out.String(), "network probe skipped by --offline") {
		t.Fatalf("expected provider skip in output:\n%s", out.String())
	}
}

func TestRunDoctorExtensionDiagnosticsWithFixtures(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)
	writeDoctorMcpConfig(t, homeDir, `{"local":{"command":"echo","args":["ok"]}}`)
	writeDoctorSkill(t, homeDir, "valid.json", `{"name":"review","description":"Review code","prompt":"Review this."}`)
	writeDoctorSkill(t, homeDir, "invalid.json", `{"description":"missing name"}`)
	writeDoctorHooksDir(t, homeDir)
	t.Setenv("GO_CODE_LSP_URL", "http://127.0.0.1:65535")

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out.String())
	}
	output := out.String()
	for _, want := range []string{
		"[PASS] mcp: 1 server config(s) found",
		"[SKIP] lsp: configured via GO_CODE_LSP_URL; health check skipped by --offline",
		"[PASS] hooks: hooks directory readable",
		"[PASS] skills: 1 skill(s) loaded, 1 warning(s)",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected %q in doctor output:\n%s", want, output)
		}
	}
}

func TestRunDoctorInvalidMcpConfigFails(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)
	writeDoctorMcpConfig(t, homeDir, `{not-json}`)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    true,
	})

	if code == 0 {
		t.Fatalf("expected non-zero exit code for invalid MCP config")
	}
	output := out.String()
	if !strings.Contains(output, "[FAIL] mcp") || !strings.Contains(output, "failed to parse") {
		t.Fatalf("expected MCP parse failure:\n%s", output)
	}
}

func TestRunDoctorLSPHealthPassesWhenConfiguredOnline(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)
	server := newDoctorLSPServer(t)
	defer server.Close()
	t.Setenv("GO_CODE_LSP_URL", server.URL)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    false,
	})

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "[PASS] lsp: configured via GO_CODE_LSP_URL and health check passed") {
		t.Fatalf("expected LSP health pass:\n%s", out.String())
	}
}

func TestRunDoctorLSPHealthFailureIsNonFatalToOtherChecks(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test"}`)
	writeDoctorDocs(t, workspace)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	t.Setenv("GO_CODE_LSP_URL", server.URL)

	var out bytes.Buffer
	code := RunDoctor(&out, DoctorOptions{
		HomeDir:    homeDir,
		WorkingDir: workspace,
		Offline:    false,
	})

	if code == 0 {
		t.Fatalf("expected non-zero exit code for failing LSP health check")
	}
	output := out.String()
	if !strings.Contains(output, "[FAIL] lsp") || !strings.Contains(output, "[PASS] tools") {
		t.Fatalf("expected LSP failure while other checks continue:\n%s", output)
	}
}

func writeDoctorConfig(t *testing.T, homeDir, body string) {
	t.Helper()
	configDir := filepath.Join(homeDir, ".go-code")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "settings.json"), []byte(body), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
}

func writeDoctorDocs(t *testing.T, workspace string) {
	t.Helper()
	quickStartDir := filepath.Join(workspace, "docs", "zh", "guide")
	if err := os.MkdirAll(quickStartDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatalf("failed to write README: %v", err)
	}
	if err := os.WriteFile(filepath.Join(quickStartDir, "quick-start.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatalf("failed to write quick start: %v", err)
	}
}

func writeDoctorMcpConfig(t *testing.T, homeDir, body string) {
	t.Helper()
	configDir := filepath.Join(homeDir, ".config", "go-code")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create MCP config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "mcp.json"), []byte(body), 0644); err != nil {
		t.Fatalf("failed to write MCP config: %v", err)
	}
}

func writeDoctorSkill(t *testing.T, homeDir, name, body string) {
	t.Helper()
	skillsDir := filepath.Join(homeDir, ".go-code", "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("failed to create skills dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillsDir, name), []byte(body), 0644); err != nil {
		t.Fatalf("failed to write skill: %v", err)
	}
}

func writeDoctorHooksDir(t *testing.T, homeDir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(homeDir, ".go-code", "hooks"), 0755); err != nil {
		t.Fatalf("failed to create hooks dir: %v", err)
	}
}

func newDoctorLSPServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		result := map[string]any{
			"capabilities": map[string]any{
				"hoverProvider":      true,
				"definitionProvider": true,
			},
			"serverInfo": map[string]any{
				"name":    "doctor-lsp",
				"version": "test",
			},
		}
		resultBytes, _ := json.Marshal(result)
		respBytes, _ := json.Marshal(map[string]any{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"result":  json.RawMessage(resultBytes),
		})
		w.Header().Set("Content-Type", "application/vscode-jsonrpc")
		fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(respBytes), respBytes)
	}))
}
