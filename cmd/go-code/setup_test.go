package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestWriteConfigDeepSeekUsesOpenAICompatibleSettings(t *testing.T) {
	tmpDir := t.TempDir()
	origUser := currentSetupUser
	currentSetupUser = func() (*user.User, error) {
		return &user.User{HomeDir: tmpDir}, nil
	}
	defer func() { currentSetupUser = origUser }()

	if err := writeConfig("deepseek", "test-key", "deepseek-v4-pro"); err != nil {
		t.Fatalf("writeConfig failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, ".go-code", "settings.json"))
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("settings is invalid JSON: %v", err)
	}
	if got["provider"] != "openai" {
		t.Fatalf("provider = %q, want openai", got["provider"])
	}
	if got["model"] != "deepseek-v4-pro" {
		t.Fatalf("model = %q, want deepseek-v4-pro", got["model"])
	}
	if got["baseUrl"] != "https://api.deepseek.com" {
		t.Fatalf("baseUrl = %q, want https://api.deepseek.com", got["baseUrl"])
	}
}

func TestDefaultModelForProvider(t *testing.T) {
	tests := []struct {
		provider string
		want     string
	}{
		{"deepseek", "deepseek-v4-pro"},
		{"openai", "gpt-5.2"},
		{"anthropic", "claude-sonnet-4-6"},
		{"custom-openai", "gpt-5.2"},
	}
	for _, tt := range tests {
		if got := defaultModelForProvider(tt.provider); got != tt.want {
			t.Fatalf("defaultModelForProvider(%q) = %q, want %q", tt.provider, got, tt.want)
		}
	}
}
