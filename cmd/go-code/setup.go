package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// SetupWizard runs an interactive configuration wizard.
func SetupWizard() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("  go-code Setup Wizard")
	fmt.Println("========================================")
	fmt.Println()

	// Step 1: Select provider
	provider, err := selectProvider(reader)
	if err != nil {
		return err
	}

	// Step 2: Enter API key
	apiKey, skipped, err := enterAPIKey(reader, provider)
	if err != nil {
		return err
	}
	if skipped {
		fmt.Println()
		fmt.Println("No problem! You can configure your API key later:")
		fmt.Println("  export ANTHROPIC_API_KEY=sk-ant-...")
		fmt.Println("  or create ~/.go-code/settings.json")
		fmt.Println()
		fmt.Println("Run 'go-code --setup' anytime to reconfigure.")
		return nil
	}

	// Step 3: Select model
	model := selectModel(provider)

	// Step 4: Write config
	if err := writeConfig(provider, apiKey, model); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Setup Complete!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("Config saved to: ~/.go-code/settings.json\n")
	fmt.Println()
	fmt.Println("Run 'go-code' to start!")
	return nil
}

func selectProvider(reader *bufio.Reader) (string, error) {
	fmt.Println("1. Select your LLM provider:")
	fmt.Println("   [1] Anthropic (Claude)")
	fmt.Println("   [2] OpenAI (GPT)")
	fmt.Println("   [3] Custom (OpenAI-compatible)")
	fmt.Print("\nChoose [1-3]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1", "":
		return "anthropic", nil
	case "2":
		return "openai", nil
	case "3":
		return "custom", nil
	default:
		fmt.Println("Invalid choice, defaulting to Anthropic")
		return "anthropic", nil
	}
}

func enterAPIKey(reader *bufio.Reader, provider string) (string, bool, error) {
	fmt.Println()
	fmt.Println("2. Enter your API key:")

	prompt := "API key: "
	if provider == "custom" {
		fmt.Print(prompt)
	} else {
		fmt.Print(prompt)
	}

	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return "", true, nil
	}

	// Validate format
	if err := validateAPIKey(apiKey, provider); err != nil {
		fmt.Printf("Warning: %s\n", err)
		fmt.Print("Continue anyway? (y/N): ")
		confirm, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(confirm)) != "y" {
			return "", true, nil
		}
	}

	return apiKey, false, nil
}

func validateAPIKey(apiKey, provider string) error {
	switch provider {
	case "anthropic":
		if !strings.HasPrefix(apiKey, "sk-ant-") {
			return fmt.Errorf("Anthropic API keys should start with 'sk-ant-'")
		}
	case "openai":
		if !strings.HasPrefix(apiKey, "sk-") {
			return fmt.Errorf("OpenAI API keys should start with 'sk-'")
		}
	case "custom":
		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}
	}
	return nil
}

func selectModel(provider string) string {
	fmt.Println()
	fmt.Println("3. Select model:")

	switch provider {
	case "anthropic":
		fmt.Println("   [1] claude-sonnet-4-20250514 (default)")
		fmt.Println("   [2] claude-opus-4-20250514")
		fmt.Print("\nChoose [1-2]: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.TrimSpace(input) == "2" {
			return "claude-opus-4-20250514"
		}
		return "claude-sonnet-4-20250514"
	case "openai":
		fmt.Println("   [1] gpt-4o (default)")
		fmt.Println("   [2] gpt-4o-mini")
		fmt.Print("\nChoose [1-2]: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.TrimSpace(input) == "2" {
			return "gpt-4o-mini"
		}
		return "gpt-4o"
	default:
		fmt.Print("Model name: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		model := strings.TrimSpace(input)
		if model != "" {
			return model
		}
		return "gpt-4o"
	}
}

func writeConfig(provider, apiKey, model string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	configDir := filepath.Join(usr.HomeDir, ".go-code")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	settings := map[string]string{
		"apiKey":   apiKey,
		"provider": provider,
		"model":    model,
	}

	if provider == "custom" {
		fmt.Print("Base URL (e.g., https://api.example.com): ")
		baseURL, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
		if baseURL != "" {
			settings["baseUrl"] = baseURL
		}
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(configDir, "settings.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
