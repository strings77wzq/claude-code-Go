package registry

import (
	"strings"
	"testing"
)

func TestResolveConfigDetectsProviderFromKnownModel(t *testing.T) {
	resolved, err := ResolveConfig("", "", "gpt-4o", "test-key")
	if err != nil {
		t.Fatalf("ResolveConfig() error = %v", err)
	}
	if resolved.Provider != "openai" {
		t.Fatalf("provider = %s, want openai", resolved.Provider)
	}
	if resolved.BaseURL != "https://api.openai.com" {
		t.Fatalf("baseURL = %s, want OpenAI default", resolved.BaseURL)
	}
}

func TestResolveConfigRejectsUnknownProvider(t *testing.T) {
	_, err := ResolveConfig("unknown", "https://example.test", "custom-model", "test-key")
	if err == nil {
		t.Fatal("expected unsupported provider error")
	}
	if !strings.Contains(err.Error(), "unsupported provider") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveConfigRejectsProviderModelMismatch(t *testing.T) {
	_, err := ResolveConfig("anthropic", "", "gpt-4o", "test-key")
	if err == nil {
		t.Fatal("expected provider/model mismatch error")
	}
	if !strings.Contains(err.Error(), "belongs to provider") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveConfigAllowsExplicitOpenAICompatibleModel(t *testing.T) {
	resolved, err := ResolveConfig("openai", "https://api.deepseek.com", "vendor-custom-model", "test-key")
	if err != nil {
		t.Fatalf("ResolveConfig() error = %v", err)
	}
	if resolved.Provider != "openai" || resolved.Model != "vendor-custom-model" {
		t.Fatalf("unexpected resolved config: %#v", resolved)
	}
}

func TestIsKnownModel(t *testing.T) {
	if !IsKnownModel("claude-sonnet-4-6-20251001") {
		t.Fatal("expected known Claude model")
	}
	if IsKnownModel("vendor-custom-model") {
		t.Fatal("expected custom model to be unknown to runtime switching")
	}
}
