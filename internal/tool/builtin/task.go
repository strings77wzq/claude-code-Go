package builtin

import (
	"context"
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/agent/subagent"
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// TaskTool enables the main agent to spawn specialized sub-agents.
type TaskTool struct {
	apiClient  apiClient
	toolReg    toolReg
	permPolicy permPolicy
}

type apiClient interface {
	SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}

type toolReg interface {
	GetTool(name string) tool.Tool
	Execute(ctx context.Context, name string, input map[string]any) tool.Result
	GetAllDefinitions() []tool.ToolDefinition
	GetDefinitionsByTier(tiers ...tool.ToolTier) []tool.ToolDefinition
}

type permPolicy interface {
	Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision
}

// NewTaskTool creates a TaskTool with agent dependencies for sub-agent execution.
func NewTaskTool(api apiClient, reg toolReg, pol permPolicy) tool.Tool {
	return &TaskTool{apiClient: api, toolReg: reg, permPolicy: pol}
}

func (t *TaskTool) Name() string { return "Task" }
func (t *TaskTool) Description() string {
	return "Spawn a specialized sub-agent (Explore/Review/TestGen) to handle a specific task and return structured results."
}
func (t *TaskTool) Tier() tool.ToolTier { return tool.TierExtension }

func (t *TaskTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"subagent_type": map[string]any{
				"type":        "string",
				"description": "Sub-agent type: Explore (code exploration), Review (code review), or TestGen (test generation)",
				"enum":        []string{"Explore", "Review", "TestGen"},
			},
			"prompt": map[string]any{
				"type":        "string",
				"description": "The task description for the sub-agent",
			},
		},
		"required": []string{"subagent_type", "prompt"},
	}
}

func (t *TaskTool) RequiresPermission() bool { return false }
func (t *TaskTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}

func (t *TaskTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	subType, _ := input["subagent_type"].(string)
	prompt, _ := input["prompt"].(string)
	if subType == "" || prompt == "" {
		return tool.Error("subagent_type and prompt are required")
	}

	st := subagent.Type(subType)
	switch st {
	case subagent.TypeExplore, subagent.TypeReview, subagent.TypeTestGen:
	default:
		return tool.Error(fmt.Sprintf("unknown subagent_type: %s (valid: Explore, Review, TestGen)", subType))
	}

	result, err := subagent.Run(ctx, st, prompt, t.apiClient, t.toolReg, t.permPolicy)
	if err != nil {
		return tool.Error(fmt.Sprintf("sub-agent failed: %v", err))
	}

	return tool.Success(result)
}
