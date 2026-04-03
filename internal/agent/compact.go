// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"fmt"
	"strings"

	"github.com/user/go-code/internal/api"
)

const (
	// DefaultContextWindow is the default context window for claude-sonnet-4
	DefaultContextWindow = 200000

	// CompactionThreshold triggers compaction when token count exceeds this percentage
	CompactionThreshold = 80

	// MaxTurnsToKeep is the number of turns (message pairs) to keep at the end
	MaxTurnsToKeep = 10

	// TokensPerChar is the heuristic ratio for token estimation
	TokensPerChar = 4
)

// ContextConfig holds configuration for context management
type ContextConfig struct {
	ContextWindow     int
	CompactionPercent int
	MaxTurnsToKeep    int
}

// DefaultContextConfig returns the default context configuration
func DefaultContextConfig() *ContextConfig {
	return &ContextConfig{
		ContextWindow:     DefaultContextWindow,
		CompactionPercent: CompactionThreshold,
		MaxTurnsToKeep:    MaxTurnsToKeep,
	}
}

// EstimateTokens estimates the total token count for a list of messages.
// Uses a simple heuristic of ~4 characters per token.
func EstimateTokens(messages []api.Message) int {
	totalChars := 0

	for _, msg := range messages {
		// Add overhead for role field (e.g., "role": "user")
		totalChars += len(msg.Role) + 10

		// Process content based on type
		switch content := msg.Content.(type) {
		case string:
			totalChars += len(content)
		case []api.ContentBlock:
			for _, block := range content {
				switch block.Type {
				case "text":
					totalChars += len(block.Text)
				case "tool_use":
					// Add overhead for tool use
					totalChars += len(block.Name) + 50
					// Estimate input JSON size
					if block.Input != nil {
						// Rough estimate of JSON string length
						totalChars += 100
					}
				case "tool_result":
					totalChars += len(block.Text) + 50
				}
			}
		}
	}

	// Convert characters to tokens (4 chars ≈ 1 token)
	return totalChars / TokensPerChar
}

// ShouldCompact returns true if the estimated token count exceeds the threshold.
func ShouldCompact(messages []api.Message, config *ContextConfig) bool {
	if config == nil {
		config = DefaultContextConfig()
	}

	estimatedTokens := EstimateTokens(messages)
	threshold := (config.ContextWindow * config.CompactionPercent) / 100

	return estimatedTokens > threshold
}

// Compact compacts the history to fit within the token limit.
// Strategy:
//   - Keep the first user message (initial context)
//   - Keep the most recent 10 turns (20 messages)
//   - Summarize everything in between
//   - Maintain user/assistant alternation
func Compact(history *History, config *ContextConfig) {
	if config == nil {
		config = DefaultContextConfig()
	}

	messages := history.GetMessages()
	if len(messages) <= 20 {
		// Not enough messages to need compaction
		return
	}

	// Build the compacted message list
	var compacted []api.Message
	var lastRole string

	// 1. Keep the first user message (initial context)
	// Find first user message
	firstUserIdx := -1
	for i, msg := range messages {
		if msg.Role == "user" {
			firstUserIdx = i
			break
		}
	}

	if firstUserIdx >= 0 {
		compacted = append(compacted, messages[firstUserIdx])
		lastRole = messages[firstUserIdx].Role
		messages = messages[firstUserIdx+1:]
	}

	// 2. Keep the most recent 10 turns (20 messages)
	// But first, check if we have enough messages
	if len(messages) <= 20 {
		// Just add all remaining messages
		for _, msg := range messages {
			if msg.Role != lastRole {
				compacted = append(compacted, msg)
				lastRole = msg.Role
			}
		}
	} else {
		// Split into middle and recent
		middleCount := len(messages) - 20
		middle := messages[:middleCount]
		recent := messages[middleCount:]

		// 3. Summarize the middle messages
		summary := summarizeMessages(middle)
		summaryRole := "user"
		if lastRole == "user" {
			summaryRole = "assistant"
		}
		summaryMsg := api.Message{
			Role:    summaryRole,
			Content: summary,
		}
		compacted = append(compacted, summaryMsg)
		lastRole = summaryRole

		// 4. Add recent messages
		for _, msg := range recent {
			if msg.Role != lastRole {
				compacted = append(compacted, msg)
				lastRole = msg.Role
			}
		}
	}

	// Replace history with compacted messages
	history.replaceMessages(compacted)
}

// summarizeMessages creates a summary of omitted messages.
func summarizeMessages(messages []api.Message) string {
	if len(messages) == 0 {
		return ""
	}

	var summaryParts []string

	for _, msg := range messages {
		role := msg.Role
		var contentText string

		switch content := msg.Content.(type) {
		case string:
			contentText = truncateText(content, 200)
		case []api.ContentBlock:
			contentText = extractBlockText(content)
			contentText = truncateText(contentText, 200)
		}

		if contentText != "" {
			summaryParts = append(summaryParts, fmt.Sprintf("[%s]: %s", role, contentText))
		}
	}

	summary := strings.Join(summaryParts, "\n")
	return fmt.Sprintf("Previous conversation summary: %s", summary)
}

// extractBlockText extracts text content from content blocks.
func extractBlockText(blocks []api.ContentBlock) string {
	var texts []string
	for _, block := range blocks {
		if block.Type == "text" {
			texts = append(texts, block.Text)
		}
	}
	return strings.Join(texts, " ")
}

// truncateText truncates text to a maximum length.
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// CompactIfNeeded compacts history if token count exceeds threshold.
// Returns true if compaction was performed.
func CompactIfNeeded(history *History, config *ContextConfig) bool {
	messages := history.GetMessages()
	if !ShouldCompact(messages, config) {
		return false
	}

	Compact(history, config)
	return true
}
