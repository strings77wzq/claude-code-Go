package eval

import (
	"context"
	"errors"
	"testing"
)

func TestRunner_Run_AllPass(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "Hello! I used Read and Grep tools.", []string{"Read", "Grep"}, 100, 50, nil
	}

	cases := []Case{
		{Name: "greet", Prompt: "hello", ExpectText: "Hello"},
		{Name: "tools", Prompt: "search", ExpectTools: []string{"Read", "Grep"}},
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "test-model", cases)

	if report.Passed != 2 {
		t.Errorf("Passed = %d, want 2", report.Passed)
	}
	if report.Failed != 0 {
		t.Errorf("Failed = %d, want 0", report.Failed)
	}
	if report.Model != "test-model" {
		t.Errorf("Model = %q, want %q", report.Model, "test-model")
	}
}

func TestRunner_Run_TextMismatch(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "Goodbye!", nil, 10, 5, nil
	}

	cases := []Case{
		{Name: "check", Prompt: "hi", ExpectText: "Hello"},
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "model", cases)

	if report.Passed != 0 {
		t.Errorf("Passed = %d, want 0 (text mismatch)", report.Passed)
	}
	if report.Failed != 1 {
		t.Errorf("Failed = %d, want 1", report.Failed)
	}
}

func TestRunner_Run_ToolMismatch(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "done", []string{"Read"}, 10, 5, nil
	}

	cases := []Case{
		{Name: "tools", Prompt: "search", ExpectTools: []string{"Read", "Grep"}},
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "model", cases)

	if report.Passed != 0 {
		t.Errorf("Passed = %d, want 0 (tool mismatch)", report.Passed)
	}
}

func TestRunner_Run_ExecutionError(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "", nil, 0, 0, errors.New("API timeout")
	}

	cases := []Case{
		{Name: "fail", Prompt: "do something"},
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "model", cases)

	if report.Passed != 0 {
		t.Errorf("Passed = %d, want 0 (execution error)", report.Passed)
	}
	if report.Results[0].Error != "API timeout" {
		t.Errorf("Error = %q, want %q", report.Results[0].Error, "API timeout")
	}
}

func TestRunner_Run_NoChecks(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "anything", nil, 10, 5, nil
	}

	cases := []Case{
		{Name: "no-checks", Prompt: "hello"}, // no ExpectText or ExpectTools
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "model", cases)

	if report.Passed != 1 {
		t.Errorf("Passed = %d, want 1 (no checks = pass if no error)", report.Passed)
	}
}

func TestRunner_Run_TokenCounting(t *testing.T) {
	exec := func(ctx context.Context, prompt string) (string, []string, int, int, error) {
		return "ok", nil, 100, 200, nil
	}

	cases := []Case{
		{Name: "a", Prompt: "x"},
		{Name: "b", Prompt: "y"},
	}

	runner := NewRunner(exec)
	report := runner.Run(context.Background(), "model", cases)

	if report.TotalTok != 600 { // (100+200) * 2
		t.Errorf("TotalTok = %d, want 600", report.TotalTok)
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct {
		name     string
		haystack []string
		needles  []string
		want     bool
	}{
		{"exact", []string{"A", "B"}, []string{"A", "B"}, true},
		{"subset", []string{"A", "B", "C"}, []string{"A", "B"}, true},
		{"missing", []string{"A"}, []string{"A", "B"}, false},
		{"empty needles", []string{"A"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsAll(tt.haystack, tt.needles); got != tt.want {
				t.Errorf("containsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
