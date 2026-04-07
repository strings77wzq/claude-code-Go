package integration

import (
	"testing"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestMultiTurnConversation(t *testing.T) {
	// Setup test agent
	ag := agent.New(agent.Config{
		MaxTurns: 10,
		Timeout:  30 * time.Second,
	})

	// Register mock tools
	ag.RegisterTool(tool.Definition{
		Name:        "MockRead",
		Description: "Read a file",
		Handler: func(args map[string]interface{}) (string, error) {
			return "file content", nil
		},
	})

	// Test multi-turn conversation flow
	messages := []string{
		"Read the main.go file",
		"Now edit it to add a comment",
		"Show me the diff",
	}

	for i, msg := range messages {
		t.Run(fmt.Sprintf("turn_%d", i), func(t *testing.T) {
			response, err := ag.Process(msg)
			if err != nil {
				t.Fatalf("Turn %d failed: %v", i, err)
			}
			if response == "" {
				t.Fatal("Empty response")
			}
		})
	}
}

func TestToolChain(t *testing.T) {
	ag := agent.New(agent.Config{
		MaxTurns: 5,
	})

	// Test that tools can chain together
	result, err := ag.Process("Read config.json, then edit it to add a new key")
	if err != nil {
		t.Fatalf("Tool chain failed: %v", err)
	}

	// Verify the agent used multiple tools
	if len(ag.Session.ToolCalls) < 2 {
		t.Error("Expected at least 2 tool calls in chain")
	}
}

func TestErrorRecovery(t *testing.T) {
	ag := agent.New(agent.Config{
		MaxRetries: 3,
	})

	// Register a tool that fails initially then succeeds
	callCount := 0
	ag.RegisterTool(tool.Definition{
		Name: "FlakyTool",
		Handler: func(args map[string]interface{}) (string, error) {
			callCount++
			if callCount < 3 {
				return "", fmt.Errorf("temporary error")
			}
			return "success", nil
		},
	})

	result, err := ag.Process("Use FlakyTool")
	if err != nil {
		t.Fatalf("Error recovery failed: %v", err)
	}

	if result != "success" {
		t.Error("Expected success after retries")
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}
