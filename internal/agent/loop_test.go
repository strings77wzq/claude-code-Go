package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type mockApiClient struct {
	response *api.ApiResponse
	err      error
	sendFunc func(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}

func (m *mockApiClient) SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, req, onTextDelta)
	}
	if m.err != nil {
		return nil, m.err
	}
	if onTextDelta != nil && m.response != nil {
		for _, block := range m.response.Content {
			if block.Type == "text" {
				onTextDelta(block.Text)
			}
		}
	}
	return m.response, nil
}

type mockToolRegistry struct {
	tools map[string]tool.Tool
}

func newMockToolRegistry() *mockToolRegistry {
	return &mockToolRegistry{
		tools: make(map[string]tool.Tool),
	}
}

func (m *mockToolRegistry) GetTool(name string) tool.Tool {
	return m.tools[name]
}

func (m *mockToolRegistry) Execute(ctx context.Context, name string, input map[string]any) tool.Result {
	t := m.tools[name]
	if t == nil {
		return tool.Error("tool not found: " + name)
	}
	return t.Execute(ctx, input)
}

func (m *mockToolRegistry) GetAllDefinitions() []tool.ToolDefinition {
	var defs []tool.ToolDefinition
	for _, t := range m.tools {
		defs = append(defs, tool.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return defs
}

func (m *mockToolRegistry) registerTool(t tool.Tool) {
	m.tools[t.Name()] = t
}

type mockPermissionPolicy struct {
	decisions map[string]permission.Decision
}

func newMockPermissionPolicy() *mockPermissionPolicy {
	return &mockPermissionPolicy{
		decisions: make(map[string]permission.Decision),
	}
}

func (m *mockPermissionPolicy) Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision {
	if decision, exists := m.decisions[toolName]; exists {
		return decision
	}
	return permission.Allow
}

func (m *mockPermissionPolicy) setDecision(toolName string, decision permission.Decision) {
	m.decisions[toolName] = decision
}

type mockTool struct {
	name        string
	description string
	inputSchema map[string]any
	result      tool.Result
}

func (m *mockTool) Name() string                { return m.name }
func (m *mockTool) Description() string         { return m.description }
func (m *mockTool) InputSchema() map[string]any { return m.inputSchema }
func (m *mockTool) RequiresPermission() bool    { return false }
func (m *mockTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (m *mockTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	return m.result
}

func TestAgent_Run_EndTurn(t *testing.T) {
	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "test-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "Hello, world!"}},
			Model:      "test-model",
			StopReason: "end_turn",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
		},
	}
	toolRegistry := newMockToolRegistry()
	permissionPolicy := newMockPermissionPolicy()
	systemPrompt := "You are a helpful assistant."
	model := "test-model"

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, systemPrompt, model)

	var output string
	callback := func(text string) { output += text }

	result, err := agent.Run(context.Background(), "Hello", callback)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result != "Hello, world!" {
		t.Errorf("Run() = %q, want %q", result, "Hello, world!")
	}
}

func TestAgent_Run_ToolUse(t *testing.T) {
	callCount := 0

	mockTool := &mockTool{
		name:        "test_tool",
		description: "A test tool",
		inputSchema: map[string]any{"type": "object"},
		result:      tool.Success("tool executed successfully"),
	}

	firstCall := true
	apiClient := &mockApiClient{
		response: nil,
		err:      nil,
		sendFunc: func(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error) {
			defer func() { firstCall = false }()
			if firstCall {
				if onTextDelta != nil {
					onTextDelta("I'll use the tool")
				}
				return &api.ApiResponse{
					ID:   "test-id",
					Type: "message",
					Role: "assistant",
					Content: []api.ContentBlock{
						{Type: "text", Text: "I'll use the tool"},
						{Type: "tool_use", ID: "toolu_123", Name: "test_tool", Input: map[string]any{"arg": "value"}},
					},
					Model:      "test-model",
					StopReason: "tool_use",
					Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
				}, nil
			}
			if onTextDelta != nil {
				onTextDelta("Tool result received")
			}
			return &api.ApiResponse{
				ID:         "test-id-2",
				Type:       "message",
				Role:       "assistant",
				Content:    []api.ContentBlock{{Type: "text", Text: "Tool result received"}},
				Model:      "test-model",
				StopReason: "end_turn",
				Usage:      api.Usage{InputTokens: 20, OutputTokens: 10},
			}, nil
		},
	}
	toolRegistry := newMockToolRegistry()
	toolRegistry.registerTool(mockTool)
	permissionPolicy := newMockPermissionPolicy()
	permissionPolicy.setDecision("test_tool", permission.Allow)

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are helpful.", "test-model")

	var outputs []string
	callback := func(text string) { outputs = append(outputs, text) }

	result, err := agent.Run(context.Background(), "Use the tool", callback)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result == "" {
		t.Error("Run() returned empty result after tool_use")
	}

	history := agent.GetHistory()
	if history.Size() != 4 {
		t.Errorf("History size = %d, want 4 (user, assistant with tool_use + text, tool result, assistant response)", history.Size())
	}

	_ = callCount
}

func TestAgent_Run_MaxTurns(t *testing.T) {
	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "test-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "Thinking..."}},
			Model:      "test-model",
			StopReason: "tool_use",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
		},
	}
	toolRegistry := newMockToolRegistry()
	permissionPolicy := newMockPermissionPolicy()

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are helpful.", "test-model")

	result, err := agent.Run(context.Background(), "Test max turns", func(string) {})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	expectedMsg := "Reached maximum turns"
	if result == "" || !strings.Contains(result, expectedMsg) {
		t.Errorf("Run() = %q, should contain %q", result, expectedMsg)
	}
}

func TestAgent_Run_UnknownStopReason(t *testing.T) {
	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "test-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "Some response"}},
			Model:      "test-model",
			StopReason: "unknown_reason",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
		},
	}
	toolRegistry := newMockToolRegistry()
	permissionPolicy := newMockPermissionPolicy()

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are helpful.", "test-model")

	result, err := agent.Run(context.Background(), "Test", func(string) {})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result != "Some response" {
		t.Errorf("Run() = %q, want %q", result, "Some response")
	}
}

func TestAgent_Run_MaxTokensTruncation(t *testing.T) {
	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "test-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "Truncated..."}},
			Model:      "test-model",
			StopReason: "max_tokens",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 8192},
		},
	}
	toolRegistry := newMockToolRegistry()
	permissionPolicy := newMockPermissionPolicy()

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are helpful.", "test-model")

	result, err := agent.Run(context.Background(), "Test", func(string) {})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	warningSubstring := "truncated"
	if !strings.Contains(result, warningSubstring) {
		t.Errorf("Run() = %q, should contain %q", result, warningSubstring)
	}
}

func TestAgent_Run_ApiError(t *testing.T) {
	apiClient := &mockApiClient{
		err: errors.New("API error"),
	}
	toolRegistry := newMockToolRegistry()
	permissionPolicy := newMockPermissionPolicy()

	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are helpful.", "test-model")

	_, err := agent.Run(context.Background(), "Test", func(string) {})
	if err == nil {
		t.Error("Run() should return error on API failure")
	}
}

func TestHistory_AddUserMessage(t *testing.T) {
	h := NewHistory()

	err := h.AddUserMessage("Hello")
	if err != nil {
		t.Errorf("AddUserMessage() error = %v", err)
	}

	if h.Size() != 1 {
		t.Errorf("History size = %d, want 1", h.Size())
	}
}

func TestHistory_AddAssistantMessage(t *testing.T) {
	h := NewHistory()

	h.AddUserMessage("Hello")
	err := h.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Hi there"}})
	if err != nil {
		t.Errorf("AddAssistantMessage() error = %v", err)
	}

	if h.Size() != 2 {
		t.Errorf("History size = %d, want 2", h.Size())
	}
}

func TestHistory_AddToolResults(t *testing.T) {
	h := NewHistory()

	h.AddUserMessage("Hello")
	h.AddAssistantMessage([]api.ContentBlock{{Type: "tool_use", ID: "toolu_1", Name: "test"}})
	err := h.AddToolResults([]api.ContentBlock{{Type: "tool_result", ToolUseID: "toolu_1", Text: "result"}})
	if err != nil {
		t.Errorf("AddToolResults() error = %v", err)
	}

	if h.Size() != 3 {
		t.Errorf("History size = %d, want 3", h.Size())
	}
}

func TestHistory_RoleAlternationError(t *testing.T) {
	h := NewHistory()

	h.AddUserMessage("Hello")
	err := h.AddUserMessage("Another user message")
	if err == nil {
		t.Error("AddUserMessage() should return error on consecutive user messages")
	}
}

func TestHistory_GetMessages(t *testing.T) {
	h := NewHistory()

	h.AddUserMessage("Hello")
	h.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Hi"}})

	messages := h.GetMessages()
	if len(messages) != 2 {
		t.Errorf("GetMessages() returned %d messages, want 2", len(messages))
	}

	messages[0].Role = "modified"
	if h.GetMessages()[0].Role != "user" {
		t.Error("GetMessages() should return a copy, not modify original")
	}
}
