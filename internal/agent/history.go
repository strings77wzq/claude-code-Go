// Package agent provides the core agent loop for the Claude Code clone.
package agent

import (
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/api"
)

// History manages the conversation history with strict user/assistant alternation.
type History struct {
	messages []api.Message
}

// NewHistory creates a new empty History.
func NewHistory() *History {
	return &History{
		messages: make([]api.Message, 0),
	}
}

// AddUserMessage adds a user message with text content.
// Returns error if the previous message was also a user message.
func (h *History) AddUserMessage(text string) error {
	if err := h.checkRoleAlternation("user"); err != nil {
		return err
	}
	h.messages = append(h.messages, api.Message{
		Role:    "user",
		Content: text,
	})
	return nil
}

// AddAssistantMessage adds an assistant message with content blocks.
// Returns error if the previous message was also an assistant message.
func (h *History) AddAssistantMessage(blocks []api.ContentBlock) error {
	if err := h.checkRoleAlternation("assistant"); err != nil {
		return err
	}
	h.messages = append(h.messages, api.Message{
		Role:    "assistant",
		Content: blocks,
	})
	return nil
}

// AddToolResults adds tool results as a user message.
// Each tool result block contains the tool_use_id to reference the original tool_use.
// Returns error if the previous message was also a user message.
func (h *History) AddToolResults(toolResults []api.ContentBlock) error {
	if err := h.checkRoleAlternation("user"); err != nil {
		return err
	}
	h.messages = append(h.messages, api.Message{
		Role:    "user",
		Content: toolResults,
	})
	return nil
}

// GetMessages returns a copy of the current messages for API requests.
func (h *History) GetMessages() []api.Message {
	result := make([]api.Message, len(h.messages))
	copy(result, h.messages)
	return result
}

// Clear resets the history to empty.
func (h *History) Clear() {
	h.messages = h.messages[:0]
}

// Size returns the number of messages in history.
func (h *History) Size() int {
	return len(h.messages)
}

// replaceMessages replaces all messages in history with the given slice.
// This is used by compaction to update history after summarizing messages.
func (h *History) replaceMessages(messages []api.Message) {
	h.messages = make([]api.Message, len(messages))
	copy(h.messages, messages)
}

// checkRoleAlternation ensures user/assistant alternation is maintained.
func (h *History) checkRoleAlternation(newRole string) error {
	if len(h.messages) > 0 {
		lastMsg := h.messages[len(h.messages)-1]
		if lastMsg.Role == newRole {
			return fmt.Errorf("cannot add consecutive '%s' messages: messages must alternate between 'user' and 'assistant'", newRole)
		}
	}
	return nil
}
