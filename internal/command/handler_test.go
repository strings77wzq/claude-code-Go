package command

import (
	"strings"
	"testing"
)

type testAgent struct {
	model     string
	cleared   bool
	compacted bool
}

func (a *testAgent) ClearHistory() {
	a.cleared = true
}

func (a *testAgent) Compact() {
	a.compacted = true
}

func (a *testAgent) Model() string {
	return a.model
}

func (a *testAgent) SetModel(model string) {
	a.model = model
}

func TestHandleUnknownCommand(t *testing.T) {
	result := Handler{}.Handle("/nope")
	if !result.Handled {
		t.Fatal("expected command to be handled")
	}
	if !strings.Contains(result.Message, "Unknown command") {
		t.Fatalf("unexpected message: %s", result.Message)
	}
}

func TestHandleModelSwitch(t *testing.T) {
	agent := &testAgent{model: "claude-sonnet-4-6-20251001"}
	handler := Handler{Agent: agent, Model: agent.model}

	result := handler.Handle("/model gpt-4o")
	if !result.Handled {
		t.Fatal("expected command to be handled")
	}
	if agent.model != "gpt-4o" {
		t.Fatalf("expected model to switch, got %s", agent.model)
	}
	if result.Model != "gpt-4o" {
		t.Fatalf("expected result model to be updated, got %s", result.Model)
	}
}

func TestHandleUnsupportedModelSwitch(t *testing.T) {
	agent := &testAgent{model: "claude-sonnet-4-6-20251001"}
	handler := Handler{Agent: agent, Model: agent.model}

	result := handler.Handle("/model made-up-model")
	if !strings.Contains(result.Message, "Unsupported model") {
		t.Fatalf("expected unsupported model message, got %s", result.Message)
	}
	if !strings.Contains(result.Message, "Keeping current model: claude-sonnet-4-6-20251001") {
		t.Fatalf("expected current model to be reported, got %s", result.Message)
	}
	if agent.model != "claude-sonnet-4-6-20251001" {
		t.Fatalf("expected model to remain unchanged, got %s", agent.model)
	}
	if result.Model != "claude-sonnet-4-6-20251001" {
		t.Fatalf("expected result model to remain current, got %s", result.Model)
	}
}

func TestHandleSessionsEmptyState(t *testing.T) {
	result := Handler{SessionsDir: t.TempDir()}.Handle("/sessions")
	if !strings.Contains(result.Message, "No sessions found") {
		t.Fatalf("expected empty sessions message, got %s", result.Message)
	}
}

func TestHandleCompact(t *testing.T) {
	agent := &testAgent{}
	result := Handler{Agent: agent}.Handle("/compact")
	if !agent.compacted {
		t.Fatal("expected compact to be called")
	}
	if !strings.Contains(result.Message, "compacted") {
		t.Fatalf("unexpected message: %s", result.Message)
	}
}

func TestHandleUpdateWithInjectedChecker(t *testing.T) {
	result := Handler{
		Version: "0.1.0",
		CheckUpdate: func(currentVersion string) (string, string, bool, error) {
			return "v0.2.0", "https://example.com/go-code", true, nil
		},
	}.Handle("/update")

	if !strings.Contains(result.Message, "Update available") {
		t.Fatalf("unexpected update message: %s", result.Message)
	}
}
