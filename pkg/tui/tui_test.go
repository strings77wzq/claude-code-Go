package tui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockAgent implements AgentInterface for testing.
type mockAgent struct {
	runFunc func(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error)
	model   string
	cleared bool
}

func (m *mockAgent) Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error) {
	if m.runFunc != nil {
		return m.runFunc(ctx, userInput, onTextDelta)
	}
	return "", nil
}

func (m *mockAgent) ClearHistory() {
	m.cleared = true
}

func (m *mockAgent) SetModel(model string) {
	m.model = model
}

func (m *mockAgent) Model() string {
	if m.model == "" {
		return "test-model"
	}
	return m.model
}

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

func TestModel_InitReturnsBlinkCommand(t *testing.T) {
	m := NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	if cmd := m.Init(); cmd == nil {
		t.Fatal("expected Init to return a blink command")
	}
}

func TestModel_UpdateQuitKey(t *testing.T) {
	m := NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	got := updated.(model)
	if !got.quitting {
		t.Fatal("expected ctrl-c to mark model as quitting")
	}
	if cmd == nil {
		t.Fatal("expected ctrl-c to return quit command")
	}
}

func TestModel_UpdateEnterStartsAgent(t *testing.T) {
	agent := &mockAgent{
		runFunc: func(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error) {
			if userInput != "hello" {
				t.Fatalf("Run input = %q, want hello", userInput)
			}
			return "done", nil
		},
	}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)
	m.input.SetValue("hello")

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if !got.isLoading {
		t.Fatal("expected enter to start loading")
	}
	if len(got.messages) != 1 || got.messages[0].role != "user" || got.messages[0].content != "hello" {
		t.Fatalf("unexpected messages after enter: %#v", got.messages)
	}
	if cmd == nil {
		t.Fatal("expected enter to return agent command")
	}
}

func TestModel_UpdateEnterIgnoresEmptyInput(t *testing.T) {
	m := NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if len(got.messages) != 0 || got.isLoading {
		t.Fatalf("expected empty input to leave model idle: %#v", got)
	}
	if cmd != nil {
		t.Fatal("expected empty input to return no command")
	}
}

func TestModel_UpdateSlashCommand(t *testing.T) {
	agent := &mockAgent{}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)
	m.input.SetValue("/clear")

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(model)
	if cmd != nil {
		t.Fatal("expected /clear to return no async command")
	}
	if !agent.cleared {
		t.Fatal("expected /clear to clear agent history")
	}
	if len(got.messages) != 1 || !strings.Contains(got.messages[0].content, "cleared") {
		t.Fatalf("expected clear confirmation message, got %#v", got.messages)
	}
}

func TestModel_UpdateMessages(t *testing.T) {
	m := NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)

	updated, _ := m.Update(streamMsg{text: "partial"})
	got := updated.(model)
	if got.streamBuffer != "partial" {
		t.Fatalf("streamBuffer = %q, want partial", got.streamBuffer)
	}

	updated, _ = got.Update(connectionStatusMsg{text: "Connecting", elapsedStr: "3s"})
	got = updated.(model)
	if got.connectionMsg != "Connecting" || got.elapsedTime != "3s" {
		t.Fatalf("unexpected connection state: msg=%q elapsed=%q", got.connectionMsg, got.elapsedTime)
	}

	got.isLoading = true
	updated, _ = got.Update(doneMsg{result: "fallback"})
	got = updated.(model)
	if got.isLoading || got.streamBuffer != "" || len(got.messages) != 1 || got.messages[0].content != "partial" {
		t.Fatalf("unexpected done state with stream buffer: %#v", got)
	}

	m = NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	m.isLoading = true
	updated, _ = m.Update(doneMsg{result: "final"})
	got = updated.(model)
	if got.isLoading || len(got.messages) != 1 || got.messages[0].content != "final" {
		t.Fatalf("unexpected done state with result: %#v", got)
	}
}

func TestModel_UpdateErrors(t *testing.T) {
	m := NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	m.isLoading = true
	m.streamBuffer = "partial"

	updated, _ := m.Update(errorMsg{err: errors.New("bad credentials")})
	got := updated.(model)
	if got.isLoading || got.streamBuffer != "" || len(got.messages) != 1 || !strings.Contains(got.messages[0].content, "bad credentials") {
		t.Fatalf("unexpected error state: %#v", got)
	}

	m = NewModel(&mockAgent{}, "v0.2.0", "test-provider", "test-model", false)
	updated, _ = m.Update(errorMsg{err: context.DeadlineExceeded})
	got = updated.(model)
	if len(got.messages) != 1 || !strings.Contains(got.messages[0].content, "timed out") {
		t.Fatalf("expected timeout remediation, got %#v", got.messages)
	}
}

func TestModel_RunAgentCommandReturnsAgentMessage(t *testing.T) {
	agent := &mockAgent{
		runFunc: func(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error) {
			onTextDelta("delta")
			return "done", nil
		},
	}
	m := NewModel(agent, "v0.2.0", "test-provider", "test-model", false)
	cmd := m.runAgent("hello")
	if cmd == nil {
		t.Fatal("expected runAgent to return a command")
	}
	batch, ok := cmd().(tea.BatchMsg)
	if !ok || len(batch) != 2 {
		t.Fatalf("expected batch with spinner and agent command, got %#v", batch)
	}
	msg := batch[1]()
	switch msg.(type) {
	case streamMsg, doneMsg:
	default:
		t.Fatalf("expected stream or done message, got %#v", msg)
	}
}
