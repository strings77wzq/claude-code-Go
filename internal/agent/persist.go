package agent

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/session"
)

// saveSession saves the current session to disk.
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int, status string) {
	if a.sessionID == "" {
		return
	}

	s := &session.Session{
		ID:           a.sessionID,
		Model:        a.model,
		StartTime:    a.startTime,
		EndTime:      time.Now(),
		TurnCount:    turnCount,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Status:       status,
	}

	messages := a.convertHistoryToSessionMessages()
	dir := getSessionsDir()

	if a.traceFilePath != "" {
		if err := session.AppendSessionMessages(a.traceFilePath, messages); err != nil {
			slog.Warn("failed to append session messages", "session_id", a.sessionID, "error", err)
		}
		if err := session.AppendTraceStatus(a.traceFilePath, status, turnCount, inputTokens, outputTokens); err != nil {
			slog.Warn("failed to append session status", "session_id", a.sessionID, "error", err)
		}
		return
	}

	if err := session.SaveSession(s, messages, dir); err != nil {
		slog.Warn("failed to save session", "session_id", a.sessionID, "error", err)
	}
}

// convertHistoryToSessionMessages converts the agent's history to session messages.
func (a *Agent) convertHistoryToSessionMessages() []session.SessionMessage {
	var messages []session.SessionMessage
	historyMsg := a.history.GetMessages()

	for _, msg := range historyMsg {
		var content string
		switch c := msg.Content.(type) {
		case string:
			content = c
		case []api.ContentBlock:
			for _, block := range c {
				if block.Type == "text" {
					content += block.Text
				} else if block.Type == "tool_use" {
					content += fmt.Sprintf("[tool use: %s %s]", block.Name, formatToolInput(block.Input))
				} else if block.Type == "tool_result" {
					content += fmt.Sprintf("[tool result: %s]", block.Text)
				}
			}
		}
		messages = append(messages, session.SessionMessage{
			Type:      "message",
			Role:      msg.Role,
			Content:   content,
			Timestamp: time.Now(),
		})
	}

	return messages
}

func formatToolInput(input map[string]any) string {
	if len(input) == 0 {
		return "{}"
	}
	data, err := json.Marshal(input)
	if err != nil {
		return "{}"
	}
	return string(data)
}
