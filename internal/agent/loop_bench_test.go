package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// setupBenchAgent creates an agent with mockToolRegistry and minimal-mock API client
// for benchmarking. The returned agent has pre-populated history and a clean state.
func setupBenchAgent(b *testing.B, numTools int) *Agent {
	b.Helper()

	toolRegistry := newMockToolRegistry()
	for i := range numTools {
		mt := &mockTool{
			name:        fmt.Sprintf("tool_%d", i),
			description: "Benchmark test tool " + fmt.Sprint(i),
			inputSchema: map[string]any{"type": "object"},
			result:      tool.Success("benchmark result"),
		}
		toolRegistry.registerTool(mt)
	}

	apiClient := &mockApiClient{
		response: &api.ApiResponse{
			ID:         "bench-id",
			Type:       "message",
			Role:       "assistant",
			Content:    []api.ContentBlock{{Type: "text", Text: "benchmark response"}},
			Model:      "test-model",
			StopReason: "end_turn",
			Usage:      api.Usage{InputTokens: 10, OutputTokens: 5},
		},
	}

	permissionPolicy := newMockPermissionPolicy()
	agent := NewAgent(apiClient, toolRegistry, permissionPolicy, "You are a helpful assistant.", "test-model")

	// Pre-populate history with a few turns (simulating an active session).
	_ = agent.history.AddUserMessage("Hello")
	_ = agent.history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Hi, how can I help?"}})
	_ = agent.history.AddUserMessage("List the files in the project")
	_ = agent.history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Here are the files..."}})
	_ = agent.history.AddUserMessage("Great, now do the task")

	return agent
}

// BenchmarkBuildRequest measures Agent.buildRequest() allocation and throughput.
// Two sub-benchmarks cover the first-turn (full disclosure) and subsequent-turn
// (progressive disclosure) code paths.
func BenchmarkBuildRequest(b *testing.B) {
	b.ReportAllocs()

	b.Run("FirstTurn", func(b *testing.B) {
		agent := setupBenchAgent(b, 20)
		agent.turnCount = 0
		agent.cacheInvalidated = true
		agent.disclosedTools = make(map[string]bool)

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			agent.turnCount = 0
			agent.cacheInvalidated = true
			clear(agent.disclosedTools)
			_ = agent.buildRequest()
		}
	})

	b.Run("SubsequentTurn", func(b *testing.B) {
		agent := setupBenchAgent(b, 20)
		agent.turnCount = 1
		agent.cacheInvalidated = false
		agent.disclosedTools = map[string]bool{
			"tool_0": true,
			"tool_1": true,
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			agent.cacheInvalidated = false
			_ = agent.buildRequest()
		}
	})
}

// BenchmarkExecuteTools measures Agent.executeTools() allocation and throughput
// with varying numbers of tool_use content blocks.
func BenchmarkExecuteTools(b *testing.B) {
	b.ReportAllocs()

	for _, tc := range []struct {
		name  string
		count int
	}{
		{"SingleTool", 1},
		{"MultipleTools", 3},
		{"ManyTools", 10},
	} {
		b.Run(tc.name, func(b *testing.B) {
			agent := setupBenchAgent(b, 20)

			blocks := make([]api.ContentBlock, tc.count)
			for i := range tc.count {
				blocks[i] = api.ContentBlock{
					Type:  "tool_use",
					ID:    fmt.Sprintf("toolu_%d", i),
					Name:  fmt.Sprintf("tool_%d", i%20),
					Input: map[string]any{"arg": "value"},
				}
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				_ = agent.executeTools(context.Background(), blocks)
			}
		})
	}
}

// BenchmarkPerformTurn measures a single performTurn cycle including buildRequest,
// recovery context, mock API call, and tracing (all no-ops except buildRequest).
func BenchmarkPerformTurn(b *testing.B) {
	b.ReportAllocs()

	agent := setupBenchAgent(b, 10)
	agent.turnCount = 0
	agent.cacheInvalidated = true
	agent.disclosedTools = make(map[string]bool)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		agent.turnCount = 0
		agent.cacheInvalidated = true
		clear(agent.disclosedTools)
		_, _, _, _ = agent.performTurn(context.Background(), nil)
	}
}
