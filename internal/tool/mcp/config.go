package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/diagnostic"
)

type McpLaunchPolicy struct {
	AllowedCommands []string `json:"allowedCommands,omitempty"`
	WorkingDir      string   `json:"workingDir,omitempty"`
	InheritEnv      bool     `json:"inheritEnv,omitempty"`
}

// McpServerConfig represents the configuration for an MCP server.
type McpServerConfig struct {
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Env          map[string]string `json:"env"`
	LaunchPolicy McpLaunchPolicy   `json:"launchPolicy,omitempty"`
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

func ValidateLaunchPolicy(serverName string, config McpServerConfig) []diagnostic.Diagnostic {
	policy := config.LaunchPolicy
	allowed := policy.AllowedCommands
	if len(allowed) == 0 {
		allowed = defaultAllowedCommands()
	}
	command := filepath.Base(config.Command)
	for _, candidate := range allowed {
		if command == candidate || config.Command == candidate {
			return nil
		}
	}
	return []diagnostic.Diagnostic{{
		Component: "mcp",
		Severity:  diagnostic.SeverityError,
		Code:      "mcp.launch.command_not_allowed",
		Summary:   "MCP launch command is not allowed",
		Detail:    fmt.Sprintf("server %q command %q is outside the launch policy allowlist", serverName, config.Command),
		Metadata:  SanitizedConfigMetadata(serverName, config),
	}}
}

func SanitizedConfigMetadata(serverName string, config McpServerConfig) map[string]any {
	env := make(map[string]any, len(config.Env))
	for key, value := range config.Env {
		if isSensitiveEnvKey(key) {
			env[key] = diagnostic.RedactedMarker
			continue
		}
		env[key] = value
	}
	return map[string]any{
		"server":  serverName,
		"command": config.Command,
		"args":    config.Args,
		"env":     env,
	}
}

func defaultAllowedCommands() []string {
	return []string{"node", "npx", "uvx", "python", "python3", "go", filepath.Base(os.Args[0])}
}

func isSensitiveEnvKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(key))
	for _, token := range []string{"apikey", "authorization", "password", "secret"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return strings.Contains(normalized, "token") && !strings.Contains(normalized, "tokens")
}
