package agent

import (
	"context"
	"strings"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
)

// executeTools executes tool calls from the assistant's response.
func (a *Agent) executeTools(ctx context.Context, content []api.ContentBlock) []api.ContentBlock {
	var toolResults []api.ContentBlock

	for _, block := range content {
		if block.Type != "tool_use" {
			continue
		}

		toolName := block.Name
		toolInput := block.Input
		toolUseID := block.ID

		startTime := time.Now()

		if allowed, decision, reason := a.checkPermissionDetailed(toolName, toolInput); !allowed {
			a.tracePermission(toolName, decision, reason, summarizePermissionInput(toolName, toolInput))
			a.traceTool(toolName, toolInput, "permission denied", time.Since(startTime).Milliseconds())
			toolResults = append(toolResults, api.ContentBlock{
				Type:      "tool_result",
				ToolUseID: toolUseID,
				Text:      "Permission denied for tool: " + toolName + "\nReason: " + string(reason) + "\nRe-run with --permission-mode danger-full-access to grant all permissions",
				IsError:   true,
			})
			continue
		} else {
			a.tracePermission(toolName, decision, reason, summarizePermissionInput(toolName, toolInput))
		}

		if a.hooksRegistry != nil {
			if err := a.hooksRegistry.RunPreHooks(toolName, toolInput); err != nil {
				a.traceTool(toolName, toolInput, "pre-hook error: "+err.Error(), time.Since(startTime).Milliseconds())
				toolResults = append(toolResults, api.ContentBlock{
					Type:      "tool_result",
					ToolUseID: toolUseID,
					Text:      "pre-hook error: " + err.Error(),
					IsError:   true,
				})
				continue
			}
		}

		result := a.toolRegistry.Execute(ctx, toolName, toolInput)
		if result.IsError && strings.Contains(result.Content, "panic recovered") {
			a.traceRuntime("tool_panic_recovered", result.Content)
		}

		if a.hooksRegistry != nil {
			a.hooksRegistry.RunPostHooks(toolName, toolInput, result.Content, result.IsError)
		}

		a.traceTool(toolName, toolInput, result.Content, time.Since(startTime).Milliseconds())

		toolResults = append(toolResults, api.ContentBlock{
			Type:      "tool_result",
			ToolUseID: toolUseID,
			Text:      result.Content,
			IsError:   result.IsError,
		})
	}

	return toolResults
}

// checkPermission delegates to the permission policy and prompter to determine if a tool can be executed.
func (a *Agent) checkPermission(toolName string, input map[string]any) (bool, permission.Decision) {
	allowed, decision, _ := a.checkPermissionDetailed(toolName, input)
	return allowed, decision
}

func (a *Agent) checkPermissionDetailed(toolName string, input map[string]any) (bool, permission.Decision, permission.Reason) {
	t := a.toolRegistry.GetTool(toolName)
	requiresPermission := t != nil && t.RequiresPermission()

	evaluation := permission.Evaluation{
		Decision: a.permissionPolicy.Evaluate(toolName, input, requiresPermission),
		Reason:   permission.ReasonRequiresApproval,
	}
	if detailed, ok := a.permissionPolicy.(permissionDetailedPolicyInterface); ok {
		evaluation = detailed.EvaluateDetailed(toolName, input, requiresPermission)
	}
	decision := evaluation.Decision
	switch decision {
	case permission.Allow, permission.AllowOnce, permission.AllowForSession:
		return true, decision, evaluation.Reason
	case permission.Deny:
		return false, decision, evaluation.Reason
	case permission.Ask:
		if a.permissionPrompter == nil {
			return false, permission.Deny, permission.ReasonRequiresApproval
		}
		promptDecision := a.permissionPrompter.Decide(toolName, input, "tool requires approval")
		switch promptDecision {
		case permission.Allow, permission.AllowOnce:
			return true, promptDecision, evaluation.Reason
		case permission.AllowForSession:
			if memory, ok := a.permissionPolicy.(permissionMemoryInterface); ok {
				memory.RememberDecision(toolName, input, permission.AllowForSession)
			}
			return true, promptDecision, evaluation.Reason
		default:
			return false, permission.Deny, evaluation.Reason
		}
	default:
		return false, permission.Deny, evaluation.Reason
	}
}
