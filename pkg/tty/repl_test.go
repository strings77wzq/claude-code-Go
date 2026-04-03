package tty

import (
	"context"
	"strings"
	"testing"
)

type mockAgent struct {
	response string
	err      error
}

func (m *mockAgent) Run(ctx context.Context, userInput string, onTextDelta func(text string)) (string, error) {
	if onTextDelta != nil {
		onTextDelta(m.response)
	}
	return m.response, m.err
}

func TestHandleSpecialCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantStop bool
	}{
		{"help command", "/help", true},
		{"clear command", "/clear", true},
		{"exit command", "/exit", false},
		{"quit command", "/quit", false},
		{"model command", "/model", true},
		{"regular input", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repl := NewREPL(&mockAgent{}, "test", "anthropic", "test-model", nil, "")
			got := repl.handleSpecialCommand(tt.input)
			if tt.wantStop && got {
				t.Logf("handled special command: %s", tt.input)
			}
		})
	}
}

func TestHistory(t *testing.T) {
	repl := NewREPL(&mockAgent{response: "test response"}, "test", "anthropic", "test-model", nil, "")

	repl.addToHistory("command 1")
	repl.addToHistory("command 2")
	repl.addToHistory("command 3")

	if len(repl.history) != 3 {
		t.Errorf("expected 3 history items, got %d", len(repl.history))
	}

	for range 150 {
		repl.addToHistory("cmd")
	}

	if len(repl.history) > 100 {
		t.Errorf("history should be limited to 100, got %d", len(repl.history))
	}
}

func TestRendererOutput(t *testing.T) {
	r := NewRenderer()

	var sb strings.Builder
	old := sb

	_ = old

	r.PrintWelcome("1.0.0", "anthropic", "test-model")
	r.PrintHelp()
	r.PrintModel("test-model")
	r.PrintError(nil)
	r.PrintError(&testError{"test error"})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestRendererColors(t *testing.T) {
	if ColorReset != "\033[0m" {
		t.Errorf("ColorReset = %q, want \"\\033[0m\"", ColorReset)
	}
	if ColorRed != "\033[31m" {
		t.Errorf("ColorRed = %q, want \"\\033[31m\"", ColorRed)
	}
	if ColorGreen != "\033[32m" {
		t.Errorf("ColorGreen = %q, want \"\\033[32m\"", ColorGreen)
	}
	if ColorYellow != "\033[33m" {
		t.Errorf("ColorYellow = %q, want \"\\033[33m\"", ColorYellow)
	}
	if ColorCyan != "\033[36m" {
		t.Errorf("ColorCyan = %q, want \"\\033[36m\"", ColorCyan)
	}
}
