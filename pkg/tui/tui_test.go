package tui

import (
	"context"
	"strings"
	"testing"
)

// mockAgent implements AgentInterface for testing.
type mockAgent struct {
	runFunc func(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
}

func (m *mockAgent) Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error) {
	if m.runFunc != nil {
		return m.runFunc(ctx, userInput, onTextDelta)
	}
	return "", nil
}

func (m *mockAgent) ClearHistory()           {}
func (m *mockAgent) SetModel(model string)   {}
func (m *mockAgent) Model() string           { return "test-model" }

func TestNewModel_Valid(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	// Initial state assertions
	if len(m.messages) != 0 {
		t.Errorf("expected empty messages, got %d", len(m.messages))
	}
	if m.isLoading {
		t.Error("expected isLoading to be false")
	}
	if m.quitting {
		t.Error("expected quitting to be false")
	}
	if m.provider != "test-provider" {
		t.Errorf("expected provider 'test-provider', got %q", m.provider)
	}
	if m.modelName != "test-model" {
		t.Errorf("expected modelName 'test-model', got %q", m.modelName)
	}
	if m.version != "v0.2.0" {
		t.Errorf("expected version 'v0.2.0', got %q", m.version)
	}
	if m.agent == nil {
		t.Error("expected agent to be set")
	}
	if m.input.Placeholder != "Type a message or /help for commands..." {
		t.Errorf("unexpected placeholder: %q", m.input.Placeholder)
	}
}

func TestNewModel_Debug(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", true)

	if !m.debug {
		t.Error("expected debug to be true")
	}
}

func TestNewModel_NilAgent(t *testing.T) {
	// Creating a model with a nil agent should not panic.
	m := NewModel(nil, "v0.2.0", "test-provider", "test-model", false)

	if m.agent != nil {
		t.Error("expected agent to be nil")
	}
	// View should not panic even with nil agent (agent is not called during rendering).
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestModel_Quitting(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	m.quitting = true
	view := m.View()

	if !strings.Contains(view, "Goodbye!") {
		t.Errorf("expected quitting view to contain 'Goodbye!', got %q", view)
	}
}

func TestModel_ViewContainsHeader(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	view := m.View()
	if !strings.Contains(view, "go-code v0.2.0") {
		t.Errorf("expected header to contain version, got %q", view)
	}
	if !strings.Contains(view, "test-model") {
		t.Errorf("expected header to contain model name, got %q", view)
	}
}

func TestModel_ViewWithMessages(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	m.messages = append(m.messages, message{role: "user", content: "hello"})
	m.messages = append(m.messages, message{role: "assistant", content: "hi there"})
	m.messages = append(m.messages, message{role: "system", content: "system note"})
	m.messages = append(m.messages, message{role: "error", content: "something broke"})

	view := m.View()
	if !strings.Contains(view, "hello") {
		t.Errorf("expected view to contain user message 'hello', got %q", view)
	}
	if !strings.Contains(view, "hi there") {
		t.Errorf("expected view to contain assistant message, got %q", view)
	}
	if !strings.Contains(view, "system note") {
		t.Errorf("expected view to contain system message, got %q", view)
	}
	if !strings.Contains(view, "something broke") {
		t.Errorf("expected view to contain error message, got %q", view)
	}
}

func TestModel_ViewWithLoading(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	m.isLoading = true
	view := m.View()

	if !strings.Contains(view, "Thinking...") {
		t.Errorf("expected view to show 'Thinking...' when loading, got %q", view)
	}
}

func TestModel_ViewDebug(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", true)

	view := m.View()
	if !strings.Contains(view, "DEBUG") {
		t.Errorf("expected view to contain DEBUG status bar, got %v", view)
	}
}

func TestModel_ViewWithStreamBuffer(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)

	m.streamBuffer = "streaming text"
	view := m.View()

	if !strings.Contains(view, "streaming text") {
		t.Errorf("expected view to contain stream buffer, got %q", view)
	}
}
