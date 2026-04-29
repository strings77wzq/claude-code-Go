package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMcpConfigsValid(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp.json")

	content := `{
		"my-server": {
			"command": "/usr/bin/python3",
			"args": ["-m", "some_server"],
			"env": {
				"KEY": "value"
			}
		}
	}`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	configs, err := LoadMcpConfigs(configPath)
	if err != nil {
		t.Fatalf("LoadMcpConfigs returned error: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}

	cfg, ok := configs["my-server"]
	if !ok {
		t.Fatal("expected key 'my-server' in configs")
	}

	if cfg.Command != "/usr/bin/python3" {
		t.Errorf("expected command '/usr/bin/python3', got '%s'", cfg.Command)
	}

	if len(cfg.Args) != 2 || cfg.Args[0] != "-m" || cfg.Args[1] != "some_server" {
		t.Errorf("unexpected args: %v", cfg.Args)
	}

	if cfg.Env == nil || cfg.Env["KEY"] != "value" {
		t.Errorf("unexpected env: %v", cfg.Env)
	}
}

func TestLoadMcpConfigsFileNotFound(t *testing.T) {
	_, err := LoadMcpConfigs("/no/such/path/mcp.json")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read settings file") {
		t.Errorf("expected 'failed to read settings file' in error, got: %v", err)
	}
}

func TestLoadMcpConfigsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(configPath, []byte("{invalid json}"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadMcpConfigs(configPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse settings file") {
		t.Errorf("expected 'failed to parse settings file' in error, got: %v", err)
	}
}

func TestInterpolateEnvVars(t *testing.T) {
	const testVar = "TEST_INTERPOLATE_VALUE"
	const testVal = "resolved_value"
	os.Setenv(testVar, testVal)
	defer os.Unsetenv(testVar)

	env := map[string]string{
		"MY_KEY": "${" + testVar + "}",
	}

	result := InterpolateEnvVars(env)

	if result["MY_KEY"] != testVal {
		t.Errorf("expected '%s', got '%s'", testVal, result["MY_KEY"])
	}
}

func TestInterpolateEnvVarsUnsetVariable(t *testing.T) {
	env := map[string]string{
		"MY_KEY": "${SOME_UNSET_VAR_XYZ}",
	}

	result := InterpolateEnvVars(env)
	expected := "${SOME_UNSET_VAR_XYZ}"
	if result["MY_KEY"] != expected {
		t.Errorf("expected '%s', got '%s'", expected, result["MY_KEY"])
	}
}

func TestInterpolateEnvVarsNilMap(t *testing.T) {
	result := InterpolateEnvVars(nil)
	if result == nil {
		t.Fatal("expected non-nil result from nil input")
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

func TestGetDefaultMcpConfigPath(t *testing.T) {
	path := GetDefaultMcpConfigPath()
	if path == "" {
		t.Fatal("GetDefaultMcpConfigPath returned empty string")
	}
	if !strings.HasSuffix(path, ".config/go-code/mcp.json") {
		t.Errorf("unexpected config path suffix: %s", path)
	}
}
