package anthropic

import (
	"testing"
)

func TestNewProvider_ValidConfig(t *testing.T) {
	p := NewProvider("sk-ant-test123", "https://api.anthropic.com", "claude-sonnet-4-20250514")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}

	if got := p.Name(); got != "anthropic" {
		t.Fatalf("Name() = %q, want %q", got, "anthropic")
	}
}

func TestNewProvider_DefaultBaseURL(t *testing.T) {
	p := NewProvider("sk-ant-test123", "", "claude-sonnet-4-20250514")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if p.baseURL != "https://api.anthropic.com" {
		t.Fatalf("baseURL = %q, want %q", p.baseURL, "https://api.anthropic.com")
	}
}

func TestSetModel(t *testing.T) {
	p := NewProvider("sk-ant-test123", "", "claude-sonnet-4-20250514")
	p.SetModel("claude-opus-4-20250514")
	if p.model != "claude-opus-4-20250514" {
		t.Fatalf("model = %q, want %q", p.model, "claude-opus-4-20250514")
	}
}

func TestNewProvider_MissingAPIKey(t *testing.T) {
	p := NewProvider("", "https://api.anthropic.com", "claude-sonnet-4-20250514")
	if p == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if p.apiKey != "" {
		t.Fatalf("apiKey = %q, want empty", p.apiKey)
	}
}
