package agent

import (
	"context"
	"errors"
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// initializeRun sets up the agent for a new user interaction.
func (a *Agent) initializeRun(userInput string) error {
	if a.history == nil {
		a.history = NewHistory()
	}
	a.startTime = Now()
	a.sessionID = generateSessionID()
	a.turnCount = 0
	a.disclosedTools = make(map[string]bool)
	a.cacheInvalidated = true

	a.history.AddUserMessage(userInput)
	a.traceFilePath = a.initTraceFile()

	return nil
}

// performTurn executes one agent turn: build request, call API with recovery, trace response.
func (a *Agent) performTurn(ctx context.Context, outputCallback func(string)) (*api.ApiResponse, int, int, error) {
	req := a.buildRequest()
	a.traceRequest(req.Model, len(req.Messages))

	recoveryCtx := &RecoveryContext{
		Manager:    a.recoveryManager,
		Agent:      a,
		RetryCount: 0,
	}

	var resp *api.ApiResponse
	var err error

	err = recoveryCtx.ExecuteWithRecovery(ctx, func() error {
		resp, err = a.apiClient.SendMessageStream(ctx, req, outputCallback)
		return err
	})

	if err != nil {
		a.traceError(fmt.Sprintf("API call failed after recovery: %v", err))
		if errors.Is(err, context.Canceled) {
			a.traceRuntime("request_cancelled", "request context cancelled")
		}
		return nil, 0, 0, fmt.Errorf("API call failed: %w", err)
	}

	a.traceResponse(resp.StopReason, resp.Usage.InputTokens, resp.Usage.OutputTokens)
	return resp, resp.Usage.InputTokens, resp.Usage.OutputTokens, nil
}

// finishRun saves the session with final status.
func (a *Agent) finishRun(turnCount, inputTokens, outputTokens int, status string) {
	a.saveSession(turnCount, inputTokens, outputTokens, status)
}

// handleThinking processes a thinking response from the model.
func (a *Agent) handleThinking(resp *api.ApiResponse) error {
	if len(resp.Content) > 0 {
		for _, block := range resp.Content {
			if block.Type == "thinking" {
				a.traceThinking(block.Text)
			}
		}
	}
	return a.history.AddAssistantMessage(resp.Content)
}

// handleToolUse executes tool calls from the response.
func (a *Agent) handleToolUse(ctx context.Context, resp *api.ApiResponse) error {
	toolResults := a.executeTools(ctx, resp.Content)
	if err := a.history.AddToolResults(toolResults); err != nil {
		a.traceError(fmt.Sprintf("failed to add tool results: %v", err))
		return fmt.Errorf("failed to add tool results: %w", err)
	}
	a.turnCount++
	return nil
}

// handleEndTurn extracts the final text content from the response.
func (a *Agent) handleEndTurn(resp *api.ApiResponse) string {
	return api.ExtractTextContent(resp.Content)
}

// Run is the main entry point for the agent. It takes user input and returns the final text response.
func (a *Agent) Run(ctx context.Context, userInput string, outputCallback func(string)) (string, error) {
	if err := a.initializeRun(userInput); err != nil {
		return "", err
	}

	// Apply reasoning strategy (e.g., planning before execution)
	if err := a.reasoningStrategy.BeforeLoop(ctx, a.history, a.apiClient, a.model, outputCallback); err != nil {
		a.traceError(fmt.Sprintf("reasoning strategy failed: %v", err))
		a.finishRun(0, 0, 0, "failed")
		return "", fmt.Errorf("reasoning strategy failed: %w", err)
	}

	var totalInputTokens, totalOutputTokens int

	for a.turnCount < MaxTurns {
		CompactIfNeeded(a.history, a.contextConfig)

		resp, inTokens, outTokens, err := a.performTurn(ctx, outputCallback)
		totalInputTokens += inTokens
		totalOutputTokens += outTokens

		if err != nil {
			status := "failed"
			if errors.Is(err, context.Canceled) {
				status = "cancelled"
			}
			a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, status)
			return "", err
		}

		if err := a.history.AddAssistantMessage(resp.Content); err != nil {
			a.traceError(fmt.Sprintf("failed to add assistant message: %v", err))
			a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "failed")
			return "", fmt.Errorf("failed to add assistant message: %w", err)
		}

		switch resp.StopReason {
		case "thinking":
			if err := a.handleThinking(resp); err != nil {
				a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "failed")
				return "", err
			}
			continue

		case "end_turn", "stop_sequence":
			a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "completed")
			return a.handleEndTurn(resp), nil

		case "max_tokens":
			a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "completed")
			return a.handleEndTurn(resp) + "\n[Warning] Response was truncated (max_tokens reached).", nil

		case "tool_use":
			if err := a.handleToolUse(ctx, resp); err != nil {
				a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "failed")
				return "", err
			}
			continue

		default:
			a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "completed")
			return a.handleEndTurn(resp), nil
		}
	}

	result := "[Agent loop stopped] Reached maximum turns (" + fmt.Sprintf("%d", MaxTurns) + ")."
	a.finishRun(a.turnCount, totalInputTokens, totalOutputTokens, "max_turns")
	return result, nil
}
