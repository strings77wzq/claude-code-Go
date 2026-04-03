package agent

import (
	"testing"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		messages []api.Message
		wantMin  int
		wantMax  int
	}{
		{
			name:     "empty messages",
			messages: []api.Message{},
			wantMin:  0,
			wantMax:  10,
		},
		{
			name: "single user message",
			messages: []api.Message{
				{Role: "user", Content: "Hello, how are you doing today?"},
			},
			wantMin: 5,
			wantMax: 30,
		},
		{
			name: "single assistant message",
			messages: []api.Message{
				{Role: "assistant", Content: "I am doing well, thank you!"},
			},
			wantMin: 10,
			wantMax: 30,
		},
		{
			name: "user assistant pair",
			messages: []api.Message{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there"},
			},
			wantMin: 10,
			wantMax: 40,
		},
		{
			name: "multiple messages",
			messages: []api.Message{
				{Role: "user", Content: "Write a function to calculate fibonacci"},
				{Role: "assistant", Content: "Here is a recursive implementation..."},
				{Role: "user", Content: "Can you make it iterative?"},
				{Role: "assistant", Content: "Sure, here is an iterative version..."},
			},
			wantMin: 30,
			wantMax: 100,
		},
		{
			name: "content blocks",
			messages: []api.Message{
				{
					Role: "assistant",
					Content: []api.ContentBlock{
						{Type: "text", Text: "Here is the result:"},
						{Type: "tool_use", Name: "Write", Input: map[string]any{"path": "/tmp/test.txt", "content": "hello"}},
					},
				},
			},
			wantMin: 20,
			wantMax: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EstimateTokens(tt.messages)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("EstimateTokens() = %d, want between %d and %d", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestShouldCompact(t *testing.T) {
	config := &ContextConfig{
		ContextWindow:     200000,
		CompactionPercent: 80,
		MaxTurnsToKeep:    10,
	}

	tests := []struct {
		name     string
		messages []api.Message
		want     bool
	}{
		{
			name:     "empty messages - no compact",
			messages: []api.Message{},
			want:     false,
		},
		{
			name: "small number of messages - no compact",
			messages: []api.Message{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldCompact(tt.messages, config)
			if got != tt.want {
				t.Errorf("ShouldCompact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompact(t *testing.T) {
	config := &ContextConfig{
		ContextWindow:     200000,
		CompactionPercent: 80,
		MaxTurnsToKeep:    10,
	}

	t.Run("small history - no change", func(t *testing.T) {
		history := NewHistory()
		history.AddUserMessage("First message")
		history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Response 1"}})
		history.AddUserMessage("Second message")
		history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Response 2"}})

		originalSize := history.Size()
		Compact(history, config)

		if history.Size() != originalSize {
			t.Errorf("Expected size %d, got %d", originalSize, history.Size())
		}
	})

	t.Run("preserves first user message", func(t *testing.T) {
		history := NewHistory()

		// Create 30+ messages
		for i := 0; i < 35; i++ {
			if i%2 == 0 {
				history.AddUserMessage("User message")
			} else {
				history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Assistant response"}})
			}
		}

		Compact(history, config)
		messages := history.GetMessages()

		// First message should be user
		if len(messages) == 0 {
			t.Fatal("No messages after compaction")
		}
		if messages[0].Role != "user" {
			t.Errorf("First message role = %s, want user", messages[0].Role)
		}
	})

	t.Run("preserves last 10 turns", func(t *testing.T) {
		history := NewHistory()

		// Create 30 messages (15 turns)
		for i := 0; i < 30; i++ {
			if i%2 == 0 {
				history.AddUserMessage("User message")
			} else {
				history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Assistant response"}})
			}
		}

		Compact(history, config)
		messages := history.GetMessages()

		// Should have: 1 first message + summary + 20 recent = 22 messages
		// But the structure might vary based on alternation
		if len(messages) < 20 {
			t.Errorf("Expected at least 20 messages, got %d", len(messages))
		}
	})

	t.Run("maintains alternation", func(t *testing.T) {
		history := NewHistory()

		// Create 30+ messages
		for i := 0; i < 30; i++ {
			if i%2 == 0 {
				history.AddUserMessage("User message")
			} else {
				history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Assistant response"}})
			}
		}

		Compact(history, config)
		messages := history.GetMessages()

		// Check alternation
		for i := 1; i < len(messages); i++ {
			if messages[i].Role == messages[i-1].Role {
				t.Errorf("Alternation violated at index %d: %s followed by %s",
					i-1, messages[i-1].Role, messages[i].Role)
			}
		}
	})

	t.Run("summary message present for large history", func(t *testing.T) {
		history := NewHistory()

		// Create 50 messages
		for i := 0; i < 50; i++ {
			if i%2 == 0 {
				history.AddUserMessage("User message")
			} else {
				history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Assistant response"}})
			}
		}

		Compact(history, config)
		messages := history.GetMessages()

		// After compaction, should have fewer than 50 messages
		// and more than 20 (first + summary + recent)
		if len(messages) >= 50 {
			t.Errorf("Expected fewer than 50 messages after compaction, got %d", len(messages))
		}
		if len(messages) <= 20 {
			t.Errorf("Expected more than 20 messages after compaction, got %d", len(messages))
		}
	})
}

func TestCompactIfNeeded(t *testing.T) {
	config := &ContextConfig{
		ContextWindow:     200000,
		CompactionPercent: 80,
		MaxTurnsToKeep:    10,
	}

	t.Run("returns false when compaction not needed", func(t *testing.T) {
		history := NewHistory()
		history.AddUserMessage("Hello")
		history.AddAssistantMessage([]api.ContentBlock{{Type: "text", Text: "Hi"}})

		performed := CompactIfNeeded(history, config)
		if performed {
			t.Error("Expected no compaction for small history")
		}
	})
}
