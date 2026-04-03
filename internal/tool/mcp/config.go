package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// McpServerConfig represents the configuration for an MCP server.
type McpServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// InterpolateEnvVars replaces ${VAR} placeholders with environment variable values.
func InterpolateEnvVars(env map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range env {
		result[k] = envVarPattern.ReplaceAllStringFunc(v, func(match string) string {
			varName := match[2 : len(match)-1]
			val := os.Getenv(varName)
			if val != "" {
				return val
			}
			return match
		})
	}
	return result
}

// LoadMcpConfigs loads MCP server configurations from a JSON file.
func LoadMcpConfigs(settingsPath string) (map[string]McpServerConfig, error) {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var configs map[string]McpServerConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse settings file: %w", err)
	}

	for name, config := range configs {
		if config.Env != nil {
			interpolated := InterpolateEnvVars(config.Env)
			config.Env = interpolated
			configs[name] = config
		}
	}

	return configs, nil
}

// GetDefaultMcpConfigPath returns the default path for MCP config.
func GetDefaultMcpConfigPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return ""
	}
	return strings.Join([]string{home, ".config", "go-code", "mcp.json"}, string(os.PathSeparator))
}
