package builtin

import (
	"context"
	"strings"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type taskTestAPI struct{}

func (f *taskTestAPI) SendMessageStream(_ context.Context, _ *api.ApiRequest, _ func(string)) (*api.ApiResponse, error) {
	return &api.ApiResponse{
		StopReason: "end_turn",
		Content:    []api.ContentBlock{{Type: "text", Text: "result from sub-agent"}},
		Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
	}, nil
}

type taskTestReg struct{}

func (r *taskTestReg) GetTool(name string) tool.Tool { return nil }
func (r *taskTestReg) Execute(_ context.Context, _ string, _ map[string]any) tool.Result {
	return tool.Success("ok")
}
func (r *taskTestReg) GetAllDefinitions() []tool.ToolDefinition                      { return nil }
func (r *taskTestReg) GetDefinitionsByTier(_ ...tool.ToolTier) []tool.ToolDefinition { return nil }

type taskTestPolicy struct{}

func (p *taskTestPolicy) Evaluate(string, map[string]any, bool) permission.Decision {
	return permission.Allow
}

func TestTaskToolName(t *testing.T) {
	tk := NewTaskTool(&taskTestAPI{}, &taskTestReg{}, &taskTestPolicy{}, "test-model")
	if tk.Name() != "Task" {
		t.Errorf("expected 'Task', got '%s'", tk.Name())
	}
}

func TestTaskToolExplore(t *testing.T) {
	tk := NewTaskTool(&taskTestAPI{}, &taskTestReg{}, &taskTestPolicy{}, "test-model")
	result := tk.Execute(context.Background(), map[string]any{
		"subagent_type": "Explore",
		"prompt":        "find all Go files",
	})
	if result.IsError {
		t.Fatalf("expected success, got error: %s", result.Content)
	}
	if !strings.Contains(result.Content, "result from sub-agent") {
		t.Errorf("expected result to contain sub-agent output, got: %s", result.Content)
	}
}

func TestTaskToolInvalidType(t *testing.T) {
	tk := NewTaskTool(&taskTestAPI{}, &taskTestReg{}, &taskTestPolicy{}, "test-model")
	result := tk.Execute(context.Background(), map[string]any{
		"subagent_type": "InvalidType",
		"prompt":        "do something",
	})
	if !result.IsError {
		t.Error("expected error for invalid sub-agent type")
	}
}

func TestTaskToolMissingFields(t *testing.T) {
	tk := NewTaskTool(&taskTestAPI{}, &taskTestReg{}, &taskTestPolicy{}, "test-model")
	result := tk.Execute(context.Background(), map[string]any{})
	if !result.IsError {
		t.Error("expected error for missing fields")
	}
}

var _ tool.Tool = (*TaskTool)(nil)
