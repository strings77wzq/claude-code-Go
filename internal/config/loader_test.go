package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestLoadFromEnvVars(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "env-api-key")
	os.Setenv("ANTHROPIC_BASE_URL", "https://custom.api.anthropic.com")
	defer os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Unsetenv("ANTHROPIC_BASE_URL")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIKey != "env-api-key" {
		t.Errorf("expected APIKey 'env-api-key', got '%s'", cfg.APIKey)
	}
	if cfg.BaseURL != "https://custom.api.anthropic.com" {
		t.Errorf("expected BaseURL 'https://custom.api.anthropic.com', got '%s'", cfg.BaseURL)
	}
}

func TestCLIOverridesEnvVars(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "env-api-key")
	os.Setenv("ANTHROPIC_BASE_URL", "https://env.api.anthropic.com")
	defer os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Unsetenv("ANTHROPIC_BASE_URL")

	overrides := &CLIOverrides{
		APIKey:  "cli-api-key",
		BaseURL: "https://cli.api.anthropic.com",
		Model:   "claude-opus",
	}

	cfg, err := Load(overrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIKey != "cli-api-key" {
		t.Errorf("expected APIKey 'cli-api-key', got '%s'", cfg.APIKey)
	}
	if cfg.BaseURL != "https://cli.api.anthropic.com" {
		t.Errorf("expected BaseURL 'https://cli.api.anthropic.com', got '%s'", cfg.BaseURL)
	}
	if cfg.Model != "claude-opus" {
		t.Errorf("expected Model 'claude-opus', got '%s'", cfg.Model)
	}
}

func TestConfigFileLoading(t *testing.T) {
	tmpDir := t.TempDir()

	userDir := filepath.Join(tmpDir, "user")
	os.MkdirAll(filepath.Join(userDir, ".go-code"), 0755)
	userConfigPath := filepath.Join(userDir, ".go-code", "settings.json")
	userSettings := Settings{
		APIKey:  "user-api-key",
		BaseURL: "https://user.api.anthropic.com",
		Model:   "user-model",
	}
	userData, _ := json.Marshal(userSettings)
	os.WriteFile(userConfigPath, userData, 0644)

	projectDir := filepath.Join(tmpDir, "project")
	os.MkdirAll(filepath.Join(projectDir, ".go-code"), 0755)
	projectConfigPath := filepath.Join(projectDir, ".go-code", "settings.json")
	projectSettings := Settings{
		APIKey:  "project-api-key",
		BaseURL: "https://project.api.anthropic.com",
		Model:   "project-model",
	}
	projectData, _ := json.Marshal(projectSettings)
	os.WriteFile(projectConfigPath, projectData, 0644)

	origWorkingDir := mustGetwd()
	origUser := currentUser
	defer func() {
		currentUser = origUser
		os.Chdir(origWorkingDir)
	}()

	os.Chdir(projectDir)
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: userDir}, nil
	}

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIKey != "project-api-key" {
		t.Errorf("expected APIKey 'project-api-key', got '%s'", cfg.APIKey)
	}
	if cfg.BaseURL != "https://project.api.anthropic.com" {
		t.Errorf("expected BaseURL 'https://project.api.anthropic.com', got '%s'", cfg.BaseURL)
	}
	if cfg.Model != "project-model" {
		t.Errorf("expected Model 'project-model', got '%s'", cfg.Model)
	}
}

func TestPriorityChain(t *testing.T) {
	tmpDir := t.TempDir()

	userDir := filepath.Join(tmpDir, "user")
	os.MkdirAll(filepath.Join(userDir, ".go-code"), 0755)
	userConfigPath := filepath.Join(userDir, ".go-code", "settings.json")
	userSettings := Settings{
		APIKey:  "user-key",
		BaseURL: "https://user.com",
		Model:   "user-model",
	}
	userData, _ := json.Marshal(userSettings)
	os.WriteFile(userConfigPath, userData, 0644)

	projectDir := filepath.Join(tmpDir, "project")
	os.MkdirAll(filepath.Join(projectDir, ".go-code"), 0755)
	projectConfigPath := filepath.Join(projectDir, ".go-code", "settings.json")
	projectSettings := Settings{
		APIKey:  "project-key",
		BaseURL: "https://project.com",
		Model:   "project-model",
	}
	projectData, _ := json.Marshal(projectSettings)
	os.WriteFile(projectConfigPath, projectData, 0644)

	origWorkingDir := mustGetwd()
	origUser := currentUser
	defer func() {
		currentUser = origUser
		os.Chdir(origWorkingDir)
	}()

	os.Chdir(projectDir)
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: userDir}, nil
	}

	os.Setenv("ANTHROPIC_API_KEY", "env-key")
	os.Setenv("ANTHROPIC_BASE_URL", "https://env.com")
	defer os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Unsetenv("ANTHROPIC_BASE_URL")

	overrides := &CLIOverrides{
		APIKey:  "cli-key",
		BaseURL: "https://cli.com",
		Model:   "cli-model",
	}

	cfg, err := Load(overrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIKey != "cli-key" {
		t.Errorf("expected CLI to override: APIKey 'cli-key', got '%s'", cfg.APIKey)
	}
	if cfg.BaseURL != "https://cli.com" {
		t.Errorf("expected CLI to override: BaseURL 'https://cli.com', got '%s'", cfg.BaseURL)
	}
	if cfg.Model != "cli-model" {
		t.Errorf("expected CLI to override: Model 'cli-model', got '%s'", cfg.Model)
	}
}

func TestAPIKeyValidation(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load(nil)
	if err == nil {
		t.Fatalf("expected error for empty API key, got cfg: %+v", cfg)
	}
	if err != ErrAPIKeyRequired {
		t.Errorf("expected ErrAPIKeyRequired, got: %v", err)
	}
}

func TestDefaults(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_BASE_URL")

	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	os.MkdirAll(filepath.Join(projectDir, ".go-code"), 0755)

	userDir := filepath.Join(tmpDir, "user")
	os.MkdirAll(filepath.Join(userDir, ".go-code"), 0755)

	origWorkingDir := mustGetwd()
	origUser := currentUser
	defer func() {
		currentUser = origUser
		os.Chdir(origWorkingDir)
	}()

	os.Chdir(projectDir)
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: userDir}, nil
	}

	settings := Settings{
		APIKey: "test-key",
	}
	data, _ := json.Marshal(settings)
	os.WriteFile(filepath.Join(userDir, ".go-code", "settings.json"), data, 0644)

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BaseURL != "https://api.anthropic.com" {
		t.Errorf("expected default BaseURL, got '%s'", cfg.BaseURL)
	}
	if cfg.Model != "claude-sonnet-4-20250514" {
		t.Errorf("expected default Model, got '%s'", cfg.Model)
	}
	if cfg.MaxTokens != 8192 {
		t.Errorf("expected default MaxTokens 8192, got %d", cfg.MaxTokens)
	}
}

func mustGetwd() string {
	wd, _ := os.Getwd()
	return wd
}
