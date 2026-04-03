// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/go-code/internal/api"
	"github.com/user/go-code/internal/hooks"
	"github.com/user/go-code/internal/permission"
	"github.com/user/go-code/internal/session"
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
	hooksRegistry    *hooks.Registry
	systemPrompt     string
	model            string
	maxTokens        int
	history          *History
	contextConfig    *ContextConfig
	sessionID        string
	startTime        time.Time
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
		hooksRegistry:    hooks.NewRegistry(),
	}
}

// SetHooksRegistry sets the hooks registry for the agent.
func (a *Agent) SetHooksRegistry(reg *hooks.Registry) {
	a.hooksRegistry = reg
}

// Run is the main entry point for the agent. It takes user input and returns the final text response.
func (a *Agent) Run(ctx context.Context, userInput string, outputCallback func(string)) (string, error) {
	a.sessionID = generateSessionID()
	a.startTime = time.Now()

	var totalInputTokens, totalOutputTokens int

	if err := a.history.AddUserMessage(userInput); err != nil {
		return "", fmt.Errorf("failed to add user message: %w", err)
	}

	turns := 0

	for turns < MaxTurns {
		CompactIfNeeded(a.history, a.contextConfig)

		req := a.buildRequest()

		resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
		if err != nil {
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return "", fmt.Errorf("API call failed: %w", err)
		}

		totalInputTokens += resp.Usage.InputTokens
		totalOutputTokens += resp.Usage.OutputTokens

		if err := a.history.AddAssistantMessage(resp.Content); err != nil {
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return "", fmt.Errorf("failed to add assistant message: %w", err)
		}

		switch resp.StopReason {
		case "end_turn", "stop_sequence":
			result := extractTextContent(resp.Content)
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return result, nil

		case "max_tokens":
			result := extractTextContent(resp.Content) + "\n[Warning] Response was truncated (max_tokens reached)."
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return result, nil

		case "tool_use":
			toolResults := a.executeTools(ctx, resp.Content)
			if err := a.history.AddToolResults(toolResults); err != nil {
				a.saveSession(turns, totalInputTokens, totalOutputTokens)
				return "", fmt.Errorf("failed to add tool results: %w", err)
			}
			turns++
			continue

		default:
			result := extractTextContent(resp.Content)
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return result, nil
		}
	}

	result := "[Agent loop stopped] Reached maximum turns (" + fmt.Sprintf("%d", MaxTurns) + ")."
	a.saveSession(turns, totalInputTokens, totalOutputTokens)
	return result, nil
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
func (a *Agent) executeTools(ctx context.Context, content []api.ContentBlock) []api.ContentBlock {
	var toolResults []api.ContentBlock

	for _, block := range content {
		if block.Type != "tool_use" {
			continue
		}

		toolName := block.Name
		toolInput := block.Input
		toolUseID := block.ID

		if !a.checkPermission(toolName, toolInput) {
			toolResults = append(toolResults, api.ContentBlock{
				Type:      "tool_result",
				ToolUseID: toolUseID,
				IsError:   true,
			})
			continue
		}

		if a.hooksRegistry != nil {
			if err := a.hooksRegistry.RunPreHooks(toolName, toolInput); err != nil {
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

		if a.hooksRegistry != nil {
			a.hooksRegistry.RunPostHooks(toolName, toolInput, result.Content, result.IsError)
		}

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

// generateSessionID creates a unique session identifier.
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixMilli())
}

// getSessionsDir returns the directory for storing session files.
func getSessionsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".claude-code-go/sessions"
	}
	return filepath.Join(homeDir, ".claude-code-go", "sessions")
}

// saveSession saves the current session to disk.
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int) {
	if a.sessionID == "" {
		return
	}

	s := &session.Session{
		ID:           a.sessionID,
		Model:        a.model,
		StartTime:    a.startTime,
		EndTime:      time.Now(),
		TurnCount:    turnCount,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
	}

	messages := a.convertHistoryToSessionMessages()
	dir := getSessionsDir()

	if err := session.SaveSession(s, messages, dir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
	}
}

// convertHistoryToSessionMessages converts the agent's history to session messages.
func (a *Agent) convertHistoryToSessionMessages() []session.SessionMessage {
	var messages []session.SessionMessage
	historyMsg := a.history.GetMessages()

	for _, msg := range historyMsg {
		var content string
		switch c := msg.Content.(type) {
		case string:
			content = c
		case []api.ContentBlock:
			for _, block := range c {
				if block.Type == "text" {
					content += block.Text
				} else if block.Type == "tool_result" {
					content += fmt.Sprintf("[tool result: %s]", block.Text)
				}
			}
		}
		messages = append(messages, session.SessionMessage{
			Type:      "message",
			Role:      msg.Role,
			Content:   content,
			Timestamp: time.Now(),
		})
	}

	return messages
}
