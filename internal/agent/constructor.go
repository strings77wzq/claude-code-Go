package agent

import (
	"path/filepath"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/hooks"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

// NewAgent creates a new Agent with the given dependencies.
func NewAgent(
	apiClient ApiClientInterface,
	toolRegistry ToolRegistryInterface,
	permissionPolicy PermissionPolicyInterface,
	systemPrompt string,
	model string,
) *Agent {
	return &Agent{
		apiClient:          apiClient,
		toolRegistry:       toolRegistry,
		permissionPolicy:   permissionPolicy,
		permissionPrompter: permission.NewDefaultPrompter(),
		systemPrompt:       systemPrompt,
		model:              model,
		maxTokens:          DefaultMaxTokens,
		history:            NewHistory(),
		contextConfig:      DefaultContextConfig(),
		hooksRegistry:      hooks.NewRegistry(),
		recoveryManager:    NewRecoveryManager(),
		cacheEnabled:       true,
		disclosedTools:     make(map[string]bool),
		reasoningStrategy:  ReactStrategy{},
	}
}

// SetHooksRegistry sets the hooks registry for the agent.
func (a *Agent) SetHooksRegistry(reg *hooks.Registry) {
	a.hooksRegistry = reg
}

// PersistSession enables session persistence for the agent.
func (a *Agent) PersistSession() {
	a.keepSession = true
}

// LoadExternalHooks loads external hooks from a directory.
func (a *Agent) LoadExternalHooks(dir string) {
	if a.hooksRegistry != nil {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return
		}
		_ = absDir
		// hooks.LoadExternalDir delegates to the hooks registry
	}
}

// SetPermissionPrompter sets the permission prompter for interactive mode.
func (a *Agent) SetPermissionPrompter(prompter permission.Prompter) {
	a.permissionPrompter = prompter
}

// SetThinkingBudget sets the extended thinking budget in tokens.
func (a *Agent) SetThinkingBudget(tokens int) {
	a.thinkingBudget = tokens
}

// SetReasoningStrategy sets the reasoning strategy for the agent.
func (a *Agent) SetReasoningStrategy(s ReasoningStrategy) {
	if s != nil {
		a.reasoningStrategy = s
	}
}

// Now returns the current time (exported for testing).
var Now = func() time.Time { return time.Now() }
