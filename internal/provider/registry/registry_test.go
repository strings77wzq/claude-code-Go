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

func TestDeepSeekV4ModelsResolve(t *testing.T) {
	tests := []struct {
		model    string
		wantProv string
	}{
		{"deepseek-v4-pro", "openai"},
		{"deepseek-v4-flash", "openai"},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			resolved, err := ResolveConfig("", "", tt.model, "test-key")
			if err != nil {
				t.Fatalf("ResolveConfig() error = %v", err)
			}
			if resolved.Provider != tt.wantProv {
				t.Fatalf("provider = %s, want %s", resolved.Provider, tt.wantProv)
			}
			if resolved.Model != tt.model {
				t.Fatalf("model = %s, want %s", resolved.Model, tt.model)
			}
		})
	}
}

func TestLegacyDeepSeekModelsResolve(t *testing.T) {
	tests := []struct {
		model    string
		wantProv string
	}{
		{"deepseek-chat", "openai"},
		{"deepseek-reasoner", "openai"},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			// Legacy models should still resolve (they're in the registry as deprecated)
			info, ok := LookupModel(tt.model)
			if !ok {
				t.Fatalf("LookupModel(%q) = false, want true (deprecated alias)", tt.model)
			}
			if !info.Deprecated {
				t.Fatalf("expected %q to be marked deprecated", tt.model)
			}
			if info.Provider != tt.wantProv {
				t.Fatalf("provider = %s, want %s", info.Provider, tt.wantProv)
			}
			if info.DeprecationMessage == "" {
				t.Fatal("expected non-empty deprecation message")
			}

			// ResolveConfig should also work
			resolved, err := ResolveConfig("", "", tt.model, "test-key")
			if err != nil {
				t.Fatalf("ResolveConfig() error = %v", err)
			}
			if resolved.Provider != tt.wantProv {
				t.Fatalf("provider = %s, want %s", resolved.Provider, tt.wantProv)
			}
			if resolved.Model != tt.model {
				t.Fatalf("model = %s, want %s", resolved.Model, tt.model)
			}
		})
	}
}

func TestMiMoModelResolution(t *testing.T) {
	// MiMo should be in the registry as known model
	if !IsKnownModel("mimo-v2.5-pro") {
		t.Fatal("expected mimo-v2.5-pro to be a known model")
	}

	info, ok := LookupModel("mimo-v2.5-pro")
	if !ok {
		t.Fatal("LookupModel(mimo-v2.5-pro) = false, want true")
	}
	if info.Provider != "openai" {
		t.Fatalf("provider = %s, want openai", info.Provider)
	}

	resolved, err := ResolveConfig("", "", "mimo-v2.5-pro", "test-key")
	if err != nil {
		t.Fatalf("ResolveConfig() error = %v", err)
	}
	if resolved.Provider != "openai" {
		t.Fatalf("provider = %s, want openai", resolved.Provider)
	}
}

func TestUnknownModelPassthrough(t *testing.T) {
	tests := []struct {
		model    string
		wantProv string
	}{
		{"gpt-5-turbo", "openai"},
		{"my-custom-model", "anthropic"},
		{"deepseek-v5-ultra", "openai"},
		{"mimo-v3", "openai"},
		{"claude-5-20261001", "anthropic"},
		{"qwen-3-max", "openai"},
		{"glm-5", "openai"},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			// Model should NOT be in the registry
			if IsKnownModel(tt.model) {
				t.Fatalf("expected %q to be unknown", tt.model)
			}

			// But ResolveConfig should allow it (no error)
			resolved, err := ResolveConfig("", "", tt.model, "test-key")
			if err != nil {
				t.Fatalf("ResolveConfig() error = %v", err)
			}
			if resolved.Provider != tt.wantProv {
				t.Fatalf("provider = %s, want %s", resolved.Provider, tt.wantProv)
			}
			if resolved.Model != tt.model {
				t.Fatalf("model = %s, want %s", resolved.Model, tt.model)
			}
		})
	}
}

func TestDetectProviderNewPrefixes(t *testing.T) {
	tests := []struct {
		model    string
		wantProv string
	}{
		{"mimo-v2.5-pro", "openai"},
		{"mimo-v3-experimental", "openai"},
		{"deepseek-v4-pro", "openai"},
		{"deepseek-v4-flash", "openai"},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			got := DetectProvider(tt.model)
			if got != tt.wantProv {
				t.Fatalf("DetectProvider(%q) = %s, want %s", tt.model, got, tt.wantProv)
			}
		})
	}
}

func TestSupportedProviders(t *testing.T) {
	providers := SupportedProviders()
	if len(providers) != 2 {
		t.Fatalf("expected 2 supported providers, got %d", len(providers))
	}
	hasAnthropic := false
	hasOpenAI := false
	for _, p := range providers {
		switch p {
		case "anthropic":
			hasAnthropic = true
		case "openai":
			hasOpenAI = true
		}
	}
	if !hasAnthropic {
		t.Fatal("expected anthropic in supported providers")
	}
	if !hasOpenAI {
		t.Fatal("expected openai in supported providers")
	}
}

func TestDefaultBaseURL(t *testing.T) {
	if DefaultBaseURL("openai") != "https://api.openai.com" {
		t.Fatalf("openai default URL mismatch")
	}
	if DefaultBaseURL("anthropic") != "https://api.anthropic.com" {
		t.Fatalf("anthropic default URL mismatch")
	}
	if DefaultBaseURL("unknown") != "https://api.anthropic.com" {
		t.Fatalf("unknown provider should default to anthropic")
	}
}

func TestLookupModel(t *testing.T) {
	// Known model
	info, ok := LookupModel("gpt-4o")
	if !ok {
		t.Fatal("expected to find gpt-4o")
	}
	if info.Provider != "openai" {
		t.Fatalf("provider = %s, want openai", info.Provider)
	}
	if info.Deprecated {
		t.Fatal("gpt-4o should not be deprecated")
	}

	// Deprecated model
	info, ok = LookupModel("deepseek-chat")
	if !ok {
		t.Fatal("expected to find deepseek-chat (deprecated)")
	}
	if !info.Deprecated {
		t.Fatal("expected deepseek-chat to be deprecated")
	}

	// Unknown model
	_, ok = LookupModel("nonexistent-model")
	if ok {
		t.Fatal("expected no lookup for nonexistent model")
	}
}

func TestProviderProfileIsTransportIndependent(t *testing.T) {
	profile := ProfileForModel("gpt-4o")

	if profile.Provider != "openai" || profile.Model != "gpt-4o" {
		t.Fatalf("unexpected profile: %#v", profile)
	}
	if profile.Transport != "" {
		t.Fatalf("provider profile should not embed transport implementation, got %q", profile.Transport)
	}
	if len(profile.Capabilities) == 0 {
		t.Fatalf("expected capabilities in profile: %#v", profile)
	}
}

func TestProviderProfileDiagnostic(t *testing.T) {
	profile := ProfileForModel("vendor-custom-model")
	diag := profile.Diagnostic()

	if diag.Component != "provider" || diag.Code != "provider.profile" {
		t.Fatalf("unexpected diagnostic: %#v", diag)
	}
	fields := diag.TraceFields()
	metadata := fields["metadata"].(map[string]any)
	if metadata["model"] != "vendor-custom-model" {
		t.Fatalf("expected model metadata, got %#v", metadata)
	}
}

func TestGetSupportedModelsFiltered(t *testing.T) {
	models := GetSupportedModels()
	foundDeprecated := false
	foundNew := false
	for _, m := range models {
		if m.Name == "deepseek-chat" || m.Name == "deepseek-reasoner" {
			foundDeprecated = true
		}
		if m.Name == "deepseek-v4-pro" || m.Name == "deepseek-v4-flash" {
			foundNew = true
		}
		if m.Name == "mimo-v2.5-pro" {
			if m.Provider != "openai" {
				t.Fatalf("mimo-v2.5-pro provider = %s, want openai", m.Provider)
			}
		}
	}
	if !foundNew {
		t.Fatal("expected new DeepSeek models in GetSupportedModels")
	}
	if !foundDeprecated {
		t.Fatal("expected deprecated models in GetSupportedModels (they're still in registry)")
	}
}
