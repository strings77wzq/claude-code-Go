package telemetry

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConsentManager struct {
	configPath string
}

func NewConsentManager(configDir string) *ConsentManager {
	return &ConsentManager{
		configPath: filepath.Join(configDir, "telemetry-consent"),
	}
}

func (cm *ConsentManager) HasConsent() bool {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == "enabled"
}

func (cm *ConsentManager) RequestConsent() (bool, error) {
	fmt.Println("Help improve claude-code-Go by sharing anonymous usage statistics?")
	fmt.Println()
	fmt.Println("What we collect:")
	fmt.Println("  - Feature usage (which tools are used)")
	fmt.Println("  - Error counts (no stack traces or personal data)")
	fmt.Println("  - Performance metrics")
	fmt.Println()
	fmt.Println("What we DON'T collect:")
	fmt.Println("  - Your code or file contents")
	fmt.Println("  - API keys or personal information")
	fmt.Println("  - Session content or prompts")
	fmt.Println()
	fmt.Print("Enable anonymous telemetry? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	enabled := response == "y" || response == "yes"

	if err := cm.SetConsent(enabled); err != nil {
		return enabled, err
	}

	return enabled, nil
}

func (cm *ConsentManager) SetConsent(enabled bool) error {
	value := "disabled"
	if enabled {
		value = "enabled"
	}

	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(cm.configPath, []byte(value), 0644)
}
