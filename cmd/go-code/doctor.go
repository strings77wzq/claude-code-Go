package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/strings77wzq/claude-code-Go/internal/config"
	"github.com/strings77wzq/claude-code-Go/internal/provider/registry"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
	toolinit "github.com/strings77wzq/claude-code-Go/internal/tool/init"
)

type doctorStatus string

const (
	doctorPass doctorStatus = "PASS"
	doctorFail doctorStatus = "FAIL"
	doctorSkip doctorStatus = "SKIP"
)

type DoctorOptions struct {
	HomeDir    string
	WorkingDir string
	Offline    bool
}

type doctorCheck struct {
	Name        string
	Status      doctorStatus
	Detail      string
	Remediation string
}

func runDoctorCommand(args []string, stdout, stderr io.Writer) int {
	opts := DoctorOptions{}
	for _, arg := range args {
		switch arg {
		case "--offline", "--no-network":
			opts.Offline = true
		case "-h", "--help":
			fmt.Fprintln(stdout, "Usage: go-code doctor [--offline]")
			fmt.Fprintln(stdout, "")
			fmt.Fprintln(stdout, "Validate local config, session paths, tools, and docs.")
			return 0
		default:
			fmt.Fprintf(stderr, "Unknown doctor option: %s\n", arg)
			fmt.Fprintln(stderr, "Usage: go-code doctor [--offline]")
			return 2
		}
	}
	return RunDoctor(stdout, opts)
}

func RunDoctor(w io.Writer, opts DoctorOptions) int {
	if opts.HomeDir == "" {
		if home, err := os.UserHomeDir(); err == nil {
			opts.HomeDir = home
		}
	}
	if opts.WorkingDir == "" {
		if wd, err := os.Getwd(); err == nil {
			opts.WorkingDir = wd
		}
	}

	checks := []doctorCheck{
		checkBinary(),
		checkWorkingDir(opts.WorkingDir),
	}

	cfg, cfgSource, cfgCheck := checkConfig(opts.HomeDir, opts.WorkingDir)
	checks = append(checks, cfgCheck)
	checks = append(checks, checkProvider(cfg, cfgSource, opts.Offline))
	checks = append(checks, checkSessionDir(opts.HomeDir))
	checks = append(checks, checkTools(opts.WorkingDir))
	checks = append(checks, checkDocs(opts.WorkingDir))

	fmt.Fprintln(w, "go-code doctor")
	fmt.Fprintln(w, "==============")
	for _, check := range checks {
		fmt.Fprintf(w, "[%s] %s: %s\n", check.Status, check.Name, check.Detail)
		if check.Remediation != "" {
			fmt.Fprintf(w, "      fix: %s\n", check.Remediation)
		}
	}

	for _, check := range checks {
		if check.Status == doctorFail {
			return 1
		}
	}
	return 0
}

func checkBinary() doctorCheck {
	return doctorCheck{
		Name:   "binary",
		Status: doctorPass,
		Detail: "go-code " + version,
	}
}

func checkWorkingDir(workingDir string) doctorCheck {
	if workingDir == "" {
		return doctorCheck{
			Name:        "working directory",
			Status:      doctorFail,
			Detail:      "could not determine current working directory",
			Remediation: "run go-code from an accessible project directory",
		}
	}
	if info, err := os.Stat(workingDir); err != nil {
		return doctorCheck{
			Name:        "working directory",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("%s is not accessible: %v", workingDir, err),
			Remediation: "check the path and filesystem permissions",
		}
	} else if !info.IsDir() {
		return doctorCheck{
			Name:        "working directory",
			Status:      doctorFail,
			Detail:      workingDir + " is not a directory",
			Remediation: "run go-code from a project directory",
		}
	}
	return doctorCheck{
		Name:   "working directory",
		Status: doctorPass,
		Detail: workingDir,
	}
}

func checkConfig(homeDir, workingDir string) (*config.Config, string, doctorCheck) {
	cfg := config.DefaultConfig()
	source := "defaults"

	userPath := filepath.Join(homeDir, ".go-code", "settings.json")
	if loaded, err := loadDoctorSettings(userPath, cfg); err != nil {
		return cfg, userPath, doctorCheck{
			Name:        "configuration",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("failed to read %s: %v", userPath, err),
			Remediation: "fix or remove the invalid settings file",
		}
	} else if loaded {
		source = userPath
	}

	projectPath := filepath.Join(workingDir, ".go-code", "settings.json")
	if loaded, err := loadDoctorSettings(projectPath, cfg); err != nil {
		return cfg, projectPath, doctorCheck{
			Name:        "configuration",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("failed to read %s: %v", projectPath, err),
			Remediation: "fix or remove the invalid project settings file",
		}
	} else if loaded {
		source = projectPath
	}

	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
		source = "ANTHROPIC_API_KEY"
	}
	if baseURL := os.Getenv("ANTHROPIC_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if model := os.Getenv("ANTHROPIC_MODEL"); model != "" {
		cfg.Model = model
	}
	if provider := os.Getenv("LLM_PROVIDER"); provider != "" {
		cfg.Provider = provider
	}

	if cfg.APIKey == "" {
		return cfg, source, doctorCheck{
			Name:        "configuration",
			Status:      doctorFail,
			Detail:      "API key is missing; checked project settings, user settings, and ANTHROPIC_API_KEY",
			Remediation: "run go-code --setup, export ANTHROPIC_API_KEY, or create ~/.go-code/settings.json",
		}
	}
	if resolved, err := registry.ResolveConfig(cfg.Provider, cfg.BaseURL, cfg.Model, cfg.APIKey); err == nil {
		cfg.Provider = resolved.Provider
		cfg.BaseURL = resolved.BaseURL
		cfg.Model = resolved.Model
	}

	return cfg, source, doctorCheck{
		Name:   "configuration",
		Status: doctorPass,
		Detail: fmt.Sprintf("source=%s provider=%s model=%s base_url=%s", source, cfg.Provider, cfg.Model, cfg.BaseURL),
	}
}

func loadDoctorSettings(path string, cfg *config.Config) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	var settings config.Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return true, err
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
	return true, nil
}

func checkProvider(cfg *config.Config, source string, offline bool) doctorCheck {
	if cfg == nil || cfg.APIKey == "" {
		return doctorCheck{
			Name:        "provider",
			Status:      doctorFail,
			Detail:      "provider cannot be validated without an API key",
			Remediation: "fix the configuration check first",
		}
	}
	resolved, err := registry.ResolveConfig(cfg.Provider, cfg.BaseURL, cfg.Model, cfg.APIKey)
	if err != nil {
		return doctorCheck{
			Name:        "provider",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("%v from %s", err, source),
			Remediation: "set provider, model, baseUrl, and API key in settings.json or environment variables",
		}
	}
	if offline {
		return doctorCheck{
			Name:   "provider",
			Status: doctorSkip,
			Detail: fmt.Sprintf("provider=%s model=%s network probe skipped by --offline", resolved.Provider, resolved.Model),
		}
	}
	return doctorCheck{
		Name:   "provider",
		Status: doctorPass,
		Detail: fmt.Sprintf("provider=%s model=%s base_url=%s configuration is syntactically valid", resolved.Provider, resolved.Model, resolved.BaseURL),
	}
}

func checkSessionDir(homeDir string) doctorCheck {
	if homeDir == "" {
		return doctorCheck{
			Name:        "session directory",
			Status:      doctorFail,
			Detail:      "could not determine user home directory",
			Remediation: "set HOME or run in a normal user environment",
		}
	}
	sessionDir := filepath.Join(homeDir, ".claude-code-go", "sessions")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return doctorCheck{
			Name:        "session directory",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("%s is not writable: %v", sessionDir, err),
			Remediation: "fix directory permissions or remove the blocking path",
		}
	}
	probe := filepath.Join(sessionDir, ".doctor-write-test")
	if err := os.WriteFile(probe, []byte("ok"), 0644); err != nil {
		return doctorCheck{
			Name:        "session directory",
			Status:      doctorFail,
			Detail:      fmt.Sprintf("%s cannot be written: %v", sessionDir, err),
			Remediation: "fix directory permissions",
		}
	}
	_ = os.Remove(probe)
	return doctorCheck{
		Name:   "session directory",
		Status: doctorPass,
		Detail: sessionDir,
	}
}

func checkTools(workingDir string) doctorCheck {
	registry := tool.NewRegistry()
	if err := toolinit.RegisterBuiltinTools(registry, workingDir); err != nil {
		return doctorCheck{
			Name:        "tools",
			Status:      doctorFail,
			Detail:      err.Error(),
			Remediation: "check built-in tool registration and schemas",
		}
	}
	return doctorCheck{
		Name:   "tools",
		Status: doctorPass,
		Detail: fmt.Sprintf("%d built-in tools registered", len(registry.GetAllDefinitions())),
	}
}

func checkDocs(workingDir string) doctorCheck {
	required := []string{
		"README.md",
		filepath.Join("docs", "zh", "guide", "quick-start.md"),
	}
	for _, rel := range required {
		path := filepath.Join(workingDir, rel)
		if _, err := os.Stat(path); err != nil {
			return doctorCheck{
				Name:        "documentation",
				Status:      doctorFail,
				Detail:      fmt.Sprintf("missing %s: %v", rel, err),
				Remediation: "restore docs or run from the repository root",
			}
		}
	}
	return doctorCheck{
		Name:   "documentation",
		Status: doctorPass,
		Detail: "README.md and Chinese quick start found",
	}
}
