// Package config provides configuration loading for go-code.
package config

// Config holds the runtime configuration for go-code.
type Config struct {
	APIKey     string
	BaseURL    string
	Model      string
	MaxTokens  int
	WorkingDir string
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		BaseURL:   "https://api.anthropic.com",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 8192,
	}
}

// Settings represents the JSON config file structure.
type Settings struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}
