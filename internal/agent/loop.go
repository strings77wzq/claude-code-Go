// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/hooks"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/session"
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
	traceFilePath    string
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
	a.initTraceFile()

	var totalInputTokens, totalOutputTokens int

	if err := a.history.AddUserMessage(userInput); err != nil {
		a.traceError(fmt.Sprintf("failed to add user message: %v", err))
		return "", fmt.Errorf("failed to add user message: %w", err)
	}

	turns := 0

	for turns < MaxTurns {
		CompactIfNeeded(a.history, a.contextConfig)

		req := a.buildRequest()
		a.traceRequest(req.Model, len(req.Messages))

		resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
		if err != nil {
			a.traceError(fmt.Sprintf("API call failed: %v", err))
			a.saveSession(turns, totalInputTokens, totalOutputTokens)
			return "", fmt.Errorf("API call failed: %w", err)
		}

		a.traceResponse(resp.StopReason, resp.Usage.InputTokens, resp.Usage.OutputTokens)

		totalInputTokens += resp.Usage.InputTokens
		totalOutputTokens += resp.Usage.OutputTokens

		if err := a.history.AddAssistantMessage(resp.Content); err != nil {
			a.traceError(fmt.Sprintf("failed to add assistant message: %v", err))
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
				a.traceError(fmt.Sprintf("failed to add tool results: %v", err))
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

		startTime := time.Now()

		if !a.checkPermission(toolName, toolInput) {
			a.traceTool(toolName, toolInput, "permission denied", time.Since(startTime).Milliseconds())
			toolResults = append(toolResults, api.ContentBlock{
				Type:      "tool_result",
				ToolUseID: toolUseID,
				IsError:   true,
			})
			continue
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

// SetModel updates the model used by the agent at runtime.
func (a *Agent) SetModel(model string) {
	a.model = model
	if api, ok := a.apiClient.(interface{ SetModel(string) }); ok {
		api.SetModel(model)
	}
}

// Model returns the current model name.
func (a *Agent) Model() string {
	return a.model
}

// TraceFilePath returns the path to the trace file for the current session.
func (a *Agent) TraceFilePath() string {
	return a.traceFilePath
}

// ClearHistory resets the conversation history.
func (a *Agent) ClearHistory() {
	a.history = NewHistory()
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

// initTraceFile initializes the trace file with meta information and returns its path.
func (a *Agent) initTraceFile() string {
	if a.sessionID == "" {
		return ""
	}

	dir := getSessionsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}

	timestamp := a.startTime.Unix()
	filename := fmt.Sprintf("session-%d.jsonl", timestamp)
	filepath := filepath.Join(dir, filename)
	a.traceFilePath = filepath

	metaLine := map[string]interface{}{
		"type":       "meta",
		"session_id": a.sessionID,
		"model":      a.model,
		"start_ms":   a.startTime.UnixMilli(),
	}

	data, err := json.Marshal(metaLine)
	if err != nil {
		return ""
	}

	f, err := os.Create(filepath)
	if err != nil {
		return ""
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return ""
	}
	if _, err := f.WriteString("\n"); err != nil {
		return ""
	}

	return filepath
}

// traceRequest logs the request to the trace file.
func (a *Agent) traceRequest(model string, messagesCount int) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceRequest(a.traceFilePath, model, messagesCount)
}

// traceResponse logs the response to the trace file.
func (a *Agent) traceResponse(stopReason string, inputTokens, outputTokens int) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceResponse(a.traceFilePath, stopReason, inputTokens, outputTokens)
}

// traceTool logs tool execution to the trace file.
func (a *Agent) traceTool(name string, input interface{}, output string, durationMs int64) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceTool(a.traceFilePath, name, input, output, durationMs)
}

// traceError logs an error to the trace file.
func (a *Agent) traceError(message string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceError(a.traceFilePath, message)
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
