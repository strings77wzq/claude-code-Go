package agent

import (
	"context"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

func TestReactStrategy_BeforeLoop_NoOp(t *testing.T) {
	strategy := ReactStrategy{}
	history := NewHistory()
	history.AddUserMessage("test input")

	msgCountBefore := history.Size()

	err := strategy.BeforeLoop(context.Background(), history, nil, "test-model", func(string) {})
	if err != nil {
		t.Fatalf("ReactStrategy.BeforeLoop() error = %v", err)
	}

	if history.Size() != msgCountBefore {
		t.Errorf("ReactStrategy.BeforeLoop() modified history: size before = %d, after = %d", msgCountBefore, history.Size())
	}
}

func TestReactStrategy_Name(t *testing.T) {
	strategy := ReactStrategy{}
	if strategy.Name() != "react" {
		t.Errorf("ReactStrategy.Name() = %q, want %q", strategy.Name(), "react")
	}
}

func TestPlanExecuteVerifyStrategy_Name(t *testing.T) {
	strategy := NewPlanExecuteVerifyStrategy("")
	if strategy.Name() != "plan-execute-verify" {
		t.Errorf("PlanExecuteVerifyStrategy.Name() = %q, want %q", strategy.Name(), "plan-execute-verify")
	}
}

func TestPlanExecuteVerifyStrategy_BeforeLoop_InjectsPlanningPrompt(t *testing.T) {
	planResponse := &api.ApiResponse{
		ID:         "plan-id",
		Type:       "message",
		Role:       "assistant",
		Content:    []api.ContentBlock{{Type: "text", Text: "Step 1: Read files\nStep 2: Edit code\nStep 3: Run tests"}},
		Model:      "test-model",
		StopReason: "end_turn",
		Usage:      api.Usage{InputTokens: 50, OutputTokens: 30},
	}

	apiClient := &mockApiClient{
		response: planResponse,
	}

	history := NewHistory()
	history.AddUserMessage("build a REST API")

	strategy := NewPlanExecuteVerifyStrategy("Create a plan for this task.")

	err := strategy.BeforeLoop(context.Background(), history, apiClient, "test-model", func(string) {})
	if err != nil {
		t.Fatalf("BeforeLoop() error = %v", err)
	}

	// Should have: user msg, ack (assistant), planning prompt (user), plan (assistant), execute ack (user)
	if history.Size() != 5 {
		t.Errorf("history.Size() = %d, want 5 (user + assistant ack + planning prompt + plan + execute ack)", history.Size())
	}

	// Plan should be stored
	plan := strategy.Plan()
	if plan == "" {
		t.Error("Plan() is empty after BeforeLoop()")
	}
	if plan != "Step 1: Read files\nStep 2: Edit code\nStep 3: Run tests" {
		t.Errorf("Plan() = %q, want expected plan text", plan)
	}
}

func TestPlanExecuteVerifyStrategy_BeforeLoop_DefaultPrompt(t *testing.T) {
	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "plan-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "Plan: do things"}},
			Model:      "test-model",
			StopReason: "end_turn",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 10},
		},
	}

	history := NewHistory()
	history.AddUserMessage("do something")

	// Empty prompt should use default
	strategy := NewPlanExecuteVerifyStrategy("")

	err := strategy.BeforeLoop(context.Background(), history, apiClient, "test-model", func(string) {})
	if err != nil {
		t.Fatalf("BeforeLoop() error = %v", err)
	}

	// Should still work with default prompt
	if strategy.Plan() == "" {
		t.Error("Plan() is empty — default prompt should still produce a plan")
	}
}

func TestPlanExecuteVerifyStrategy_BeforeLoop_APIDegradesGracefully(t *testing.T) {
	apiClient := &mockApiClient{
		err: context.DeadlineExceeded,
	}

	history := NewHistory()
	history.AddUserMessage("test")

	strategy := NewPlanExecuteVerifyStrategy("plan")

	err := strategy.BeforeLoop(context.Background(), history, apiClient, "test-model", func(string) {})
	// Should NOT return error — degrade gracefully
	if err != nil {
		t.Fatalf("BeforeLoop() should degrade gracefully on API error, got error = %v", err)
	}

	// Plan should be empty
	if strategy.Plan() != "" {
		t.Errorf("Plan() = %q after API failure, want empty", strategy.Plan())
	}
}
