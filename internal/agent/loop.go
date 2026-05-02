// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	}
}

// SetHooksRegistry sets the hooks registry for the agent.
func (a *Agent) SetHooksRegistry(reg *hooks.Registry) {
	a.hooksRegistry = reg
}

// LoadExternalHooks loads hook definitions from a directory and registers them.
// Invalid hook files are skipped with warnings.
func (a *Agent) LoadExternalHooks(dir string) {
	loaded, err := hooks.LoadHooksFromDir(dir)
	if err != nil {
		slog.Warn("failed to load external hooks", "dir", dir, "error", err)
		return
	}
	for _, hook := range loaded {
		if err := a.hooksRegistry.Register(hook); err != nil {
			slog.Warn("failed to register external hook", "name", hook.Name(), "error", err)
		}
	}
}

// SetPermissionPrompter sets the prompter used when a tool requires approval.
func (a *Agent) SetPermissionPrompter(prompter permission.Prompter) {
	a.permissionPrompter = prompter
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

		recoveryCtx := &RecoveryContext{
			Manager:    a.recoveryManager,
			Agent:      a,
			RetryCount: 0,
		}

		var resp *api.ApiResponse
		var err error

		err = recoveryCtx.ExecuteWithRecovery(ctx, func() error {
			resp, err = a.apiClient.SendMessageStream(ctx, req, outputCallback)
			return err
		})

		if err != nil {
			a.traceError(fmt.Sprintf("API call failed after recovery: %v", err))
			status := "failed"
			if errors.Is(err, context.Canceled) {
				status = "cancelled"
				a.traceRuntime("request_cancelled", "request context cancelled")
			}
			a.saveSession(turns, totalInputTokens, totalOutputTokens, status)
			return "", fmt.Errorf("API call failed: %w", err)
		}

		a.traceResponse(resp.StopReason, resp.Usage.InputTokens, resp.Usage.OutputTokens)

		totalInputTokens += resp.Usage.InputTokens
		totalOutputTokens += resp.Usage.OutputTokens

		if err := a.history.AddAssistantMessage(resp.Content); err != nil {
			a.traceError(fmt.Sprintf("failed to add assistant message: %v", err))
			a.saveSession(turns, totalInputTokens, totalOutputTokens, "failed")
			return "", fmt.Errorf("failed to add assistant message: %w", err)
		}

		switch resp.StopReason {
		case "end_turn", "stop_sequence":
			result := extractTextContent(resp.Content)
			a.saveSession(turns, totalInputTokens, totalOutputTokens, "completed")
			return result, nil

		case "max_tokens":
			result := extractTextContent(resp.Content) + "\n[Warning] Response was truncated (max_tokens reached)."
			a.saveSession(turns, totalInputTokens, totalOutputTokens, "completed")
			return result, nil

		case "tool_use":
			toolResults := a.executeTools(ctx, resp.Content)
			if err := a.history.AddToolResults(toolResults); err != nil {
				a.traceError(fmt.Sprintf("failed to add tool results: %v", err))
				a.saveSession(turns, totalInputTokens, totalOutputTokens, "failed")
				return "", fmt.Errorf("failed to add tool results: %w", err)
			}
			turns++
			continue

		default:
			result := extractTextContent(resp.Content)
			a.saveSession(turns, totalInputTokens, totalOutputTokens, "completed")
			return result, nil
		}
	}

	result := "[Agent loop stopped] Reached maximum turns (" + fmt.Sprintf("%d", MaxTurns) + ")."
	a.saveSession(turns, totalInputTokens, totalOutputTokens, "max_turns")
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
		"type":          "meta",
		"session_id":    a.sessionID,
		"model":         a.model,
		"start_time_ms": a.startTime.UnixMilli(),
		"end_time_ms":   0,
		"turn_count":    0,
		"input_tokens":  0,
		"output_tokens": 0,
		"status":        "running",
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

// tracePermission logs permission decisions to the trace file.
func (a *Agent) tracePermission(toolName string, decision permission.Decision, reason permission.Reason, summary string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTracePermissionWithReason(a.traceFilePath, toolName, string(decision), summary, string(reason))
}

func (a *Agent) traceRuntime(event, summary string) {
	if a.traceFilePath == "" {
		return
	}
	session.AppendTraceRuntime(a.traceFilePath, a.sessionID, event, summary)
}

func summarizePermissionInput(toolName string, input map[string]any) string {
	var summary string
	switch toolName {
	case "Bash":
		if command, ok := input["command"].(string); ok {
			summary = command
		}
	case "Read", "Write", "Edit":
		if path, ok := input["file_path"].(string); ok {
			summary = path
		}
	case "Glob", "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			summary = pattern
		}
	}
	if summary == "" {
		summary = "tool input"
	}
	return sanitizePermissionSummary(summary)
}

var permissionSecretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`),
	regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password|authorization)(\s*[:=]\s*|\s+)[^\s;&|]+`),
	regexp.MustCompile(`sk-[A-Za-z0-9_-]{8,}`),
}

func sanitizePermissionSummary(summary string) string {
	summary = strings.TrimSpace(summary)
	for _, pattern := range permissionSecretPatterns {
		summary = pattern.ReplaceAllStringFunc(summary, func(match string) string {
			if strings.HasPrefix(strings.ToLower(match), "bearer ") {
				return "Bearer [REDACTED]"
			}
			if strings.HasPrefix(match, "sk-") {
				return "sk-[REDACTED]"
			}
			parts := regexp.MustCompile(`\s*[:=]\s*|\s+`).Split(match, 2)
			if len(parts) > 0 {
				return parts[0] + "=[REDACTED]"
			}
			return "[REDACTED]"
		})
	}
	if len(summary) > 200 {
		return summary[:200] + "...(truncated)"
	}
	return summary
}

// saveSession saves the current session to disk.
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int, status string) {
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
		Status:       status,
	}

	messages := a.convertHistoryToSessionMessages()
	dir := getSessionsDir()

	if a.traceFilePath != "" {
		if err := session.AppendSessionMessages(a.traceFilePath, messages); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to append session messages: %v\n", err)
		}
		if err := session.AppendTraceStatus(a.traceFilePath, status, turnCount, inputTokens, outputTokens); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to append session status: %v\n", err)
		}
		return
	}

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
				} else if block.Type == "tool_use" {
					content += fmt.Sprintf("[tool use: %s %s]", block.Name, formatToolInput(block.Input))
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

func formatToolInput(input map[string]any) string {
	if len(input) == 0 {
		return "{}"
	}
	data, err := json.Marshal(input)
	if err != nil {
		return "{}"
	}
	return string(data)
}
