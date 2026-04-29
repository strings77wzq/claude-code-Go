package registry

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/provider"
	"github.com/strings77wzq/claude-code-Go/internal/provider/anthropic"
	"github.com/strings77wzq/claude-code-Go/internal/provider/openai"
)

type ModelInfo struct {
	Name               string
	Provider           string
	Description        string
	Deprecated         bool
	DeprecationMessage string
}

type ResolvedConfig struct {
	Provider     string
	BaseURL      string
	Model        string
	APIKeySource string
}

var modelRegistry = []ModelInfo{
	// Anthropic
	{Name: "claude-opus-4-6-20251001", Provider: "anthropic", Description: "Most powerful model for complex reasoning"},
	{Name: "claude-sonnet-4-6-20251001", Provider: "anthropic", Description: "Balanced model for everyday tasks"},
	{Name: "claude-haiku-4-20250514", Provider: "anthropic", Description: "Fast and efficient model"},

	// OpenAI
	{Name: "gpt-4o", Provider: "openai", Description: "OpenAI's most capable model"},
	{Name: "gpt-4o-mini", Provider: "openai", Description: "Fast and affordable model"},
	{Name: "o1", Provider: "openai", Description: "Reasoning model for complex problems"},
	{Name: "o3", Provider: "openai", Description: "Advanced reasoning model"},

	// DeepSeek (primary models)
	{Name: "deepseek-v4-pro", Provider: "openai", Description: "DeepSeek's latest flagship model"},
	{Name: "deepseek-v4-flash", Provider: "openai", Description: "DeepSeek's fast reasoning model"},

	// DeepSeek (deprecated aliases)
	{Name: "deepseek-chat", Provider: "openai", Description: "[Deprecated] Use deepseek-v4-pro instead", Deprecated: true, DeprecationMessage: "deepseek-chat is deprecated, use deepseek-v4-pro instead"},
	{Name: "deepseek-reasoner", Provider: "openai", Description: "[Deprecated] Use deepseek-v4-flash instead", Deprecated: true, DeprecationMessage: "deepseek-reasoner is deprecated, use deepseek-v4-flash instead"},

	// Qwen
	{Name: "qwen-max", Provider: "openai", Description: "Alibaba Qwen's most capable model"},
	{Name: "qwen-plus", Provider: "openai", Description: "Alibaba Qwen's balanced model"},
	{Name: "qwen-turbo", Provider: "openai", Description: "Alibaba Qwen's fast model"},

	// GLM
	{Name: "glm-4-plus", Provider: "openai", Description: "Zhipu GLM's most capable model"},
	{Name: "glm-4", Provider: "openai", Description: "Zhipu GLM's balanced model"},
	{Name: "glm-4-flash", Provider: "openai", Description: "Zhipu GLM's fast model"},

	// MiMo (OpenAI-compatible)
	{Name: "mimo-v2.5-pro", Provider: "openai", Description: "MiMo's flagship model (OpenAI-compatible)"},
}

func SelectProvider(apiKey, baseURL, modelName string) provider.Provider {
	providerName := DetectProvider(modelName)
	return SelectProviderFor(providerName, apiKey, baseURL, modelName)
}

func SelectProviderFor(providerName, apiKey, baseURL, modelName string) provider.Provider {
	switch providerName {
	case "openai":
		return openai.NewProvider(apiKey, baseURL, modelName)
	case "anthropic":
		fallthrough
	default:
		return anthropic.NewProvider(apiKey, baseURL, modelName)
	}
}

func DetectProvider(modelName string) string {
	modelName = strings.ToLower(modelName)

	if strings.HasPrefix(modelName, "claude-") {
		return "anthropic"
	}

	openAIPrefixes := []string{"gpt-", "o1", "o3", "deepseek-", "qwen-", "glm-", "mimo-"}
	for _, prefix := range openAIPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			return "openai"
		}
	}

	return "anthropic"
}

func ResolveConfig(providerName, baseURL, modelName, apiKey string) (*ResolvedConfig, error) {
	providerName = strings.ToLower(strings.TrimSpace(providerName))
	baseURL = strings.TrimSpace(baseURL)
	modelName = strings.TrimSpace(modelName)
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return nil, fmt.Errorf("provider API key is required")
	}
	if modelName == "" {
		return nil, fmt.Errorf("model is required")
	}
	if providerName == "" {
		providerName = DetectProvider(modelName)
	}
	if !IsSupportedProvider(providerName) {
		return nil, fmt.Errorf("unsupported provider %q; supported providers: %s", providerName, strings.Join(SupportedProviders(), ", "))
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL(providerName)
	}
	if parsed, err := url.ParseRequestURI(baseURL); err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid base URL %q", baseURL)
	}

	if info, ok := LookupModel(modelName); ok {
		if info.Deprecated && info.DeprecationMessage != "" {
			slog.Warn(info.DeprecationMessage, "model", modelName)
		}
		if info.Provider != providerName {
			return nil, fmt.Errorf("model %q belongs to provider %q, not %q", modelName, info.Provider, providerName)
		}
	} else {
		slog.Warn("model not in verified registry, proceeding with inferred provider",
			"model", modelName,
			"provider", providerName,
		)
	}

	return &ResolvedConfig{
		Provider:     providerName,
		BaseURL:      baseURL,
		Model:        modelName,
		APIKeySource: "configuration",
	}, nil
}

func LookupModel(modelName string) (ModelInfo, bool) {
	for _, info := range modelRegistry {
		if info.Name == modelName {
			return info, true
		}
	}
	return ModelInfo{}, false
}

func IsKnownModel(modelName string) bool {
	_, ok := LookupModel(modelName)
	return ok
}

func IsSupportedProvider(providerName string) bool {
	for _, supported := range SupportedProviders() {
		if supported == providerName {
			return true
		}
	}
	return false
}

func SupportedProviders() []string {
	return []string{"anthropic", "openai"}
}

func DefaultBaseURL(providerName string) string {
	switch providerName {
	case "openai":
		return "https://api.openai.com"
	default:
		return "https://api.anthropic.com"
	}
}

func GetSupportedModels() []ModelInfo {
	result := make([]ModelInfo, len(modelRegistry))
	copy(result, modelRegistry)
	return result
}
