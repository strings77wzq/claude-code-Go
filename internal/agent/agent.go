// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/hooks"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// MaxTurns is the maximum number of agent loop iterations to prevent infinite loops.
const MaxTurns = 50

// DefaultMaxTokens is the default max tokens for API requests.
const DefaultMaxTokens = 8192

// ApiClientInterface defines the interface for API communication.
type ApiClientInterface interface {
	SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}

// ToolRegistryInterface defines the interface for tool management.
type ToolRegistryInterface interface {
	GetTool(name string) tool.Tool
	Execute(ctx context.Context, name string, input map[string]any) tool.Result
	GetAllDefinitions() []tool.ToolDefinition
	GetDefinitionsByTier(tiers ...tool.ToolTier) []tool.ToolDefinition
}

// PermissionPolicyInterface defines the interface for permission checking.
type PermissionPolicyInterface interface {
	Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision
}

type permissionMemoryInterface interface {
	RememberDecision(toolName string, input map[string]any, decision permission.Decision)
}

type permissionDetailedPolicyInterface interface {
	EvaluateDetailed(toolName string, input map[string]any, requiresPermission bool) permission.Evaluation
}

// Agent is the core agent that drives the think → act → observe → think again loop.
type Agent struct {
	apiClient          ApiClientInterface
	toolRegistry       ToolRegistryInterface
	permissionPolicy   PermissionPolicyInterface
	permissionPrompter permission.Prompter
	hooksRegistry      *hooks.Registry
	systemPrompt       string
	model              string
	maxTokens          int
	history            *History
	contextConfig      *ContextConfig
	sessionID          string
	startTime          time.Time
	traceFilePath      string
	recoveryManager    *RecoveryManager
	keepSession        bool
	turnCount          int
	cacheEnabled       bool
	disclosedTools     map[string]bool
	cacheInvalidated   bool
	thinkingBudget     int
	reasoningStrategy  ReasoningStrategy
}
