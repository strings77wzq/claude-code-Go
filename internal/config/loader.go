package config

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const (
	envAPIKey      = "ANTHROPIC_API_KEY"
	envBaseURL     = "ANTHROPIC_BASE_URL"
	envLLMProvider = "LLM_PROVIDER"

	settingsFileName = "settings.json"
	configDirName    = ".go-code"
)

var ErrAPIKeyRequired = errors.New("API key is required")

type CLIOverrides struct {
	APIKey   string
	BaseURL  string
	Model    string
	Provider string
}

func Load(overrides *CLIOverrides) (*Config, error) {
	cfg := DefaultConfig()

	if err := loadUserConfig(cfg); err != nil {
		return nil, err
	}

	if err := loadProjectConfig(cfg); err != nil {
		return nil, err
	}

	loadEnvConfig(cfg)

	if overrides != nil {
		if overrides.APIKey != "" {
			cfg.APIKey = overrides.APIKey
		}
		if overrides.BaseURL != "" {
			cfg.BaseURL = overrides.BaseURL
		}
		if overrides.Model != "" {
			cfg.Model = overrides.Model
		}
		if overrides.Provider != "" {
			cfg.Provider = overrides.Provider
		}
	}

	if cfg.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	return cfg, nil
}

var currentUser = user.Current

func loadUserConfig(cfg *Config) error {
	usr, err := currentUser()
	if err != nil {
		return nil
	}

	userConfigPath := filepath.Join(usr.HomeDir, configDirName, settingsFileName)
	return loadConfigFile(userConfigPath, cfg)
}

func loadProjectConfig(cfg *Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	projectConfigPath := filepath.Join(wd, configDirName, settingsFileName)
	return loadConfigFile(projectConfigPath, cfg)
}

func loadConfigFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	if settings.APIKey != "" {
		cfg.APIKey = settings.APIKey
	}
	if settings.BaseURL != "" {
		cfg.BaseURL = settings.BaseURL
	}
	if settings.Model != "" {
		cfg.Model = settings.Model
	}
	if settings.Provider != "" {
		cfg.Provider = settings.Provider
	}

	return nil
}

func loadEnvConfig(cfg *Config) {
	if apiKey := os.Getenv(envAPIKey); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if baseURL := os.Getenv(envBaseURL); baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if provider := os.Getenv(envLLMProvider); provider != "" {
		cfg.Provider = provider
	}
}
