package subagent

import (
	"context"
	"strings"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type fakeAPIClient struct {
	response string
}

func (f *fakeAPIClient) SendMessageStream(_ context.Context, _ *api.ApiRequest, _ func(string)) (*api.ApiResponse, error) {
	return &api.ApiResponse{
		StopReason: "end_turn",
		Content:    []api.ContentBlock{{Type: "text", Text: f.response}},
		Usage:      api.Usage{InputTokens: 100, OutputTokens: 50},
	}, nil
}

type fakeRegistry struct {
	tools map[string]tool.Tool
}

func (r *fakeRegistry) GetTool(name string) tool.Tool { return r.tools[name] }
func (r *fakeRegistry) Execute(_ context.Context, name string, _ map[string]any) tool.Result {
	return tool.Success("ok")
}
func (r *fakeRegistry) GetAllDefinitions() []tool.ToolDefinition {
	var defs []tool.ToolDefinition
	for _, t := range r.tools {
		defs = append(defs, tool.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return defs
}
func (r *fakeRegistry) GetDefinitionsByTier(_ ...tool.ToolTier) []tool.ToolDefinition {
	return r.GetAllDefinitions()
}

type allowAllPolicy struct{}

func (p *allowAllPolicy) Evaluate(string, map[string]any, bool) permission.Decision {
	return permission.Allow
}

func TestExploreSubAgentToolSet(t *testing.T) {
	tools := ExploreToolSet()
	names := toolNames(tools)
	expected := []string{"Read", "Glob", "Grep", "Tree"}
	for _, exp := range expected {
		found := false
		for _, n := range names {
			if n == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Explore tool set missing %s (has: %v)", exp, names)
		}
	}
}

func TestReviewSubAgentToolSet(t *testing.T) {
	tools := ReviewToolSet()
	names := toolNames(tools)
	for _, exp := range []string{"Read", "Diff"} {
		found := false
		for _, n := range names {
			if n == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Review tool set missing %s", exp)
		}
	}
}

func TestTestGenSubAgentToolSet(t *testing.T) {
	tools := TestGenToolSet()
	names := toolNames(tools)
	for _, exp := range []string{"Read", "Write", "Edit", "Bash"} {
		found := false
		for _, n := range names {
			if n == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("TestGen tool set missing %s", exp)
		}
	}
}

func TestRunSubAgentReturnsStructuredResult(t *testing.T) {
	client := &fakeAPIClient{response: "exploration complete"}
	reg := &fakeRegistry{tools: make(map[string]tool.Tool)}
	policy := &allowAllPolicy{}

	result, err := Run(context.Background(), TypeExplore, "find all Go files", client, reg, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
	if !strings.Contains(result, "exploration complete") {
		t.Errorf("expected result to contain response, got: %s", result)
	}
}

func TestRunSubAgentInvalidType(t *testing.T) {
	_, err := Run(context.Background(), "invalid_type", "prompt", nil, nil, nil)
	if err == nil {
		t.Error("expected error for invalid sub-agent type")
	}
}

func TestSubAgentTypeString(t *testing.T) {
	if TypeExplore.String() != "Explore" {
		t.Errorf("expected 'Explore', got '%s'", TypeExplore.String())
	}
	if TypeReview.String() != "Review" {
		t.Errorf("expected 'Review', got '%s'", TypeReview.String())
	}
	if TypeTestGen.String() != "TestGen" {
		t.Errorf("expected 'TestGen', got '%s'", TypeTestGen.String())
	}
}

func toolNames(tools map[string]tool.Tool) []string {
	var names []string
	for n := range tools {
		names = append(names, n)
	}
	return names
}
