package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDoctorHealthyOffline(t *testing.T) {
	homeDir := t.TempDir()
	workspace := t.TempDir()
	writeDoctorConfig(t, homeDir, `{"apiKey":"sk-ant-test","model":"claude-test"}`)
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
		"[PASS] session directory",
		"[PASS] tools",
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

func TestRunDoctorMissingAPIKey(t *testing.T) {
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
