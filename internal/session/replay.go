package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ReplayEvent struct {
	Type    string
	Summary string
}

func ReplaySessionFile(filePath string) ([]ReplayEvent, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open session file: %w", err)
	}
	defer f.Close()

	var events []ReplayEvent
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var base map[string]any
		if err := json.Unmarshal([]byte(line), &base); err != nil {
			events = append(events, ReplayEvent{Type: "invalid", Summary: err.Error()})
			continue
		}
		lineType, _ := base["type"].(string)
		events = append(events, ReplayEvent{
			Type:    lineType,
			Summary: summarizeReplayLine(lineType, base),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}
	return events, nil
}

func FormatReplay(events []ReplayEvent) string {
	var b strings.Builder
	for _, event := range events {
		if event.Summary == "" {
			continue
		}
		b.WriteString(event.Summary)
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func summarizeReplayLine(lineType string, line map[string]any) string {
	switch lineType {
	case "meta":
		return fmt.Sprintf("session %s model=%s status=%s", asString(line["session_id"]), asString(line["model"]), asString(line["status"]))
	case "message":
		return fmt.Sprintf("message %s: %s", asString(line["role"]), truncateReplay(asString(line["content"])))
	case "request":
		return fmt.Sprintf("request model=%s messages=%s", asString(line["model"]), asNumber(line["messages_count"]))
	case "response":
		return fmt.Sprintf("response stop_reason=%s input_tokens=%s output_tokens=%s", asString(line["stop_reason"]), asNumber(line["input_tokens"]), asNumber(line["output_tokens"]))
	case "tool":
		return fmt.Sprintf("tool %s duration_ms=%s output=%s", asString(line["name"]), asNumber(line["duration_ms"]), truncateReplay(asString(line["output"])))
	case "permission":
		return fmt.Sprintf("permission %s decision=%s summary=%s", asString(line["tool"]), asString(line["decision"]), truncateReplay(asString(line["summary"])))
	case "extension":
		summary := fmt.Sprintf("extension %s %s status=%s", asString(line["name"]), asString(line["event"]), asString(line["status"]))
		for _, key := range []string{"reason", "warning", "error"} {
			if value := asString(line[key]); value != "" {
				summary += fmt.Sprintf(" %s=%s", key, truncateReplay(value))
			}
		}
		if operations := asString(line["operations"]); operations != "" {
			summary += " operations=" + truncateReplay(operations)
		}
		return summary
	case "error":
		return fmt.Sprintf("error: %s", truncateReplay(asString(line["message"])))
	case "status":
		return fmt.Sprintf("status %s turns=%s input_tokens=%s output_tokens=%s", asString(line["status"]), asNumber(line["turn_count"]), asNumber(line["input_tokens"]), asNumber(line["output_tokens"]))
	default:
		return fmt.Sprintf("unknown %s", lineType)
	}
}

func FormatReplayEvidence(events []ReplayEvent) string {
	var b strings.Builder
	for _, event := range events {
		switch event.Type {
		case "message", "tool", "permission", "extension", "error", "status":
			if event.Summary == "" {
				continue
			}
			b.WriteString(event.Summary)
			b.WriteString("\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func asNumber(v any) string {
	if v == nil {
		return "0"
	}
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", n)
	default:
		return fmt.Sprint(v)
	}
}

func truncateReplay(s string) string {
	s = redactTraceString(s)
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) > 120 {
		return s[:120] + "...(truncated)"
	}
	return s
}
