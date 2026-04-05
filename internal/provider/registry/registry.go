package registry

import (
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/provider"
	"github.com/strings77wzq/claude-code-Go/internal/provider/anthropic"
	"github.com/strings77wzq/claude-code-Go/internal/provider/openai"
)

type ModelInfo struct {
	Name        string
	Provider    string
	Description string
}

var modelRegistry = []ModelInfo{
	{Name: "claude-opus-4-6-20251001", Provider: "anthropic", Description: "Most powerful model for complex reasoning"},
	{Name: "claude-sonnet-4-6-20251001", Provider: "anthropic", Description: "Balanced model for everyday tasks"},
	{Name: "claude-haiku-4-20250514", Provider: "anthropic", Description: "Fast and efficient model"},

	{Name: "gpt-4o", Provider: "openai", Description: "OpenAI's most capable model"},
	{Name: "gpt-4o-mini", Provider: "openai", Description: "Fast and affordable model"},
	{Name: "o1", Provider: "openai", Description: "Reasoning model for complex problems"},
	{Name: "o3", Provider: "openai", Description: "Advanced reasoning model"},

	{Name: "deepseek-chat", Provider: "openai", Description: "DeepSeek's chat model"},
	{Name: "deepseek-reasoner", Provider: "openai", Description: "DeepSeek's reasoning model"},

	{Name: "qwen-max", Provider: "openai", Description: "Alibaba Qwen's most capable model"},
	{Name: "qwen-plus", Provider: "openai", Description: "Alibaba Qwen's balanced model"},
	{Name: "qwen-turbo", Provider: "openai", Description: "Alibaba Qwen's fast model"},

	{Name: "glm-4-plus", Provider: "openai", Description: "Zhipu GLM's most capable model"},
	{Name: "glm-4", Provider: "openai", Description: "Zhipu GLM's balanced model"},
	{Name: "glm-4-flash", Provider: "openai", Description: "Zhipu GLM's fast model"},
}

func SelectProvider(apiKey, baseURL, modelName string) provider.Provider {
	providerName := detectProvider(modelName)

	switch providerName {
	case "openai":
		return openai.NewProvider(apiKey, baseURL, modelName)
	case "anthropic":
		fallthrough
	default:
		return anthropic.NewProvider(apiKey, baseURL, modelName)
	}
}

func detectProvider(modelName string) string {
	modelName = strings.ToLower(modelName)

	if strings.HasPrefix(modelName, "claude-") {
		return "anthropic"
	}

	openAIPrefixes := []string{"gpt-", "o1", "o3", "deepseek-", "qwen-", "glm-"}
	for _, prefix := range openAIPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			return "openai"
		}
	}

	return "anthropic"
}

func GetSupportedModels() []ModelInfo {
	result := make([]ModelInfo, len(modelRegistry))
	copy(result, modelRegistry)
	return result
}
