package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// ReasoningStrategy controls agent reasoning behavior before the main loop.
type ReasoningStrategy interface {
	// BeforeLoop is called after initializeRun but before the main agent loop.
	// It can inject additional messages into history (e.g., planning prompts).
	BeforeLoop(ctx context.Context, history *History, apiClient ApiClientInterface, model string, outputCallback func(string)) error

	// Name returns the strategy name for tracing and logging.
	Name() string
}

// ReactStrategy is the default reactive strategy — no planning, no extra calls.
type ReactStrategy struct{}

func (s ReactStrategy) BeforeLoop(ctx context.Context, history *History, apiClient ApiClientInterface, model string, outputCallback func(string)) error {
	return nil // no-op
}

func (s ReactStrategy) Name() string { return "react" }

// PlanExecuteVerifyStrategy adds an explicit planning phase before execution.
// It asks the model to create a structured plan, then proceeds with normal execution.
type PlanExecuteVerifyStrategy struct {
	planPrompt string
	plan       string // stored after planning phase
}

// NewPlanExecuteVerifyStrategy creates a new Plan-Execute-Verify strategy.
// planPrompt is the prompt sent to the model to generate a plan.
func NewPlanExecuteVerifyStrategy(planPrompt string) *PlanExecuteVerifyStrategy {
	return &PlanExecuteVerifyStrategy{
		planPrompt: planPrompt,
	}
}

func (s *PlanExecuteVerifyStrategy) BeforeLoop(ctx context.Context, history *History, apiClient ApiClientInterface, model string, outputCallback func(string)) error {
	if s.planPrompt == "" {
		s.planPrompt = defaultPlanPrompt()
	}

	// Add assistant ack first to maintain alternation (user msg was added by initializeRun).
	if err := history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "I'll create a plan first."}}); err != nil {
		return fmt.Errorf("failed to add planning ack: %w", err)
	}

	// Add planning prompt as user message.
	if err := history.AddUserMessage(s.planPrompt); err != nil {
		return fmt.Errorf("failed to add planning prompt: %w", err)
	}

	// Call API to get the plan.
	req := &api.ApiRequest{
		Model:     model,
		MaxTokens: 4096,
		Stream:    true,
		Messages:  history.GetMessages(),
	}

	resp, err := apiClient.SendMessageStream(ctx, req, outputCallback)
	if err != nil {
		// On failure, degrade gracefully — the main loop proceeds without a plan.
		// History has: user, assistant ack, user planning prompt.
		// The main loop will add the execution assistant response after the planning prompt,
		// which maintains correct alternation (the model sees the planning prompt as context
		// but proceeds to execute since there's no plan response).
		slog.Warn("planning phase API call failed, continuing without plan", "error", err)
		return nil
	}

	// Store the plan text.
	s.plan = api.ExtractTextContent(resp.Content)

	// Add plan as assistant message so the model sees it in context.
	if err := history.AddAssistantMessage(resp.Content); err != nil {
		return fmt.Errorf("failed to add plan response: %w", err)
	}

	// Add execution instruction so the model proceeds.
	if err := history.AddUserMessage("Good. Now execute the plan. Start with the first step."); err != nil {
		return fmt.Errorf("failed to add execution instruction: %w", err)
	}

	return nil
}

func (s *PlanExecuteVerifyStrategy) Name() string { return "plan-execute-verify" }

// Plan returns the stored plan text (empty if planning hasn't run or failed).
func (s *PlanExecuteVerifyStrategy) Plan() string { return s.plan }

func defaultPlanPrompt() string {
	return strings.Join([]string{
		"Before executing, create a structured plan for this task.",
		"List the steps you will take, what tools you will use, and what the expected outcome is.",
		"Be specific and actionable. Do not execute yet — just plan.",
	}, " ")
}
