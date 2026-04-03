// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"fmt"

	"github.com/user/go-code/internal/api"
	"github.com/user/go-code/internal/permission"
	"github.com/user/go-code/internal/tool"
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
}

// PermissionPolicyInterface defines the interface for permission checking.
type PermissionPolicyInterface interface {
	Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision
}

// Agent is the core agent that drives the think → act → observe → think again loop.
type Agent struct {
	apiClient        ApiClientInterface
	toolRegistry     ToolRegistryInterface
	permissionPolicy PermissionPolicyInterface
	systemPrompt     string
	model            string
	maxTokens        int
	history          *History
	contextConfig    *ContextConfig
}

// NewAgent creates a new Agent with the given dependencies.
func NewAgent(
	apiClient ApiClientInterface,
	toolRegistry ToolRegistryInterface,
	permissionPolicy PermissionPolicyInterface,
	systemPrompt string,
	model string,
) *Agent {
	return &Agent{
		apiClient:        apiClient,
		toolRegistry:     toolRegistry,
		permissionPolicy: permissionPolicy,
		systemPrompt:     systemPrompt,
		model:            model,
		maxTokens:        DefaultMaxTokens,
		history:          NewHistory(),
		contextConfig:    DefaultContextConfig(),
	}
}

// Run is the main entry point for the agent. It takes user input and returns the final text response.
// The function drives the agent loop: think → act → observe → think again.
func (a *Agent) Run(ctx context.Context, userInput string, outputCallback func(string)) (string, error) {
	// Add user message to history
	if err := a.history.AddUserMessage(userInput); err != nil {
		return "", fmt.Errorf("failed to add user message: %w", err)
	}

	turns := 0
	for turns < MaxTurns {
		// Compact history if needed before building request
		CompactIfNeeded(a.history, a.contextConfig)

		// Build API request
		req := a.buildRequest()

		// Call API with streaming
		resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
		if err != nil {
			return "", fmt.Errorf("API call failed: %w", err)
		}

		// Add assistant response to history
		if err := a.history.AddAssistantMessage(resp.Content); err != nil {
			return "", fmt.Errorf("failed to add assistant message: %w", err)
		}

		// Check stop_reason
		switch resp.StopReason {
		case "end_turn", "stop_sequence":
			// LLM believes task is complete
			return extractTextContent(resp.Content), nil

		case "max_tokens":
			// Output was truncated
			return extractTextContent(resp.Content) + "\n[Warning] Response was truncated (max_tokens reached).", nil

		case "tool_use":
			// LLM wants to call tools - execute them and continue loop
			toolResults := a.executeTools(ctx, resp.Content)
			if err := a.history.AddToolResults(toolResults); err != nil {
				return "", fmt.Errorf("failed to add tool results: %w", err)
			}
			turns++
			continue

		default:
			// Unknown stop_reason - safe fallback
			return extractTextContent(resp.Content), nil
		}
	}

	// Max turns reached
	return "[Agent loop stopped] Reached maximum turns (" + fmt.Sprintf("%d", MaxTurns) + ").", nil
}

// buildRequest assembles the API request with system prompt, tools, and messages.
func (a *Agent) buildRequest() *api.ApiRequest {
	toolDefs := make([]api.ToolDefinition, 0)
	for _, td := range a.toolRegistry.GetAllDefinitions() {
		toolDefs = append(toolDefs, api.ToolDefinition{
			Name:        td.Name,
			Description: td.Description,
			InputSchema: td.InputSchema,
		})
	}

	return &api.ApiRequest{
		Model:     a.model,
		MaxTokens: a.maxTokens,
		System:    a.systemPrompt,
		Stream:    true,
		Tools:     toolDefs,
		Messages:  a.history.GetMessages(),
	}
}

// executeTools executes tool calls from the assistant's response.
// It checks permissions before each tool execution.
func (a *Agent) executeTools(ctx context.Context, content []api.ContentBlock) []api.ContentBlock {
	var toolResults []api.ContentBlock

	for _, block := range content {
		if block.Type != "tool_use" {
			continue
		}

		toolName := block.Name
		toolInput := block.Input
		toolUseID := block.ID

		// Check permission
		if !a.checkPermission(toolName, toolInput) {
			toolResults = append(toolResults, api.ContentBlock{
				Type:      "tool_result",
				ToolUseID: toolUseID,
				IsError:   true,
			})
			continue
		}

		// Execute tool
		result := a.toolRegistry.Execute(ctx, toolName, toolInput)
		toolResults = append(toolResults, api.ContentBlock{
			Type:      "tool_result",
			ToolUseID: toolUseID,
			Text:      result.Content,
			IsError:   result.IsError,
		})
	}

	return toolResults
}

// checkPermission delegates to the permission policy to determine if a tool can be executed.
func (a *Agent) checkPermission(toolName string, input map[string]any) bool {
	t := a.toolRegistry.GetTool(toolName)
	requiresPermission := t != nil && t.RequiresPermission()

	decision := a.permissionPolicy.Evaluate(toolName, input, requiresPermission)
	return decision == permission.Allow || decision == permission.Ask
}

// GetHistory returns a copy of the current conversation history.
func (a *Agent) GetHistory() *History {
	return a.history
}

// extractTextContent extracts all text from content blocks.
func extractTextContent(blocks []api.ContentBlock) string {
	var text string
	for _, block := range blocks {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return text
}
