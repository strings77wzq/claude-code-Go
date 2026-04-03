package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type messageStartEvent struct {
	Type    string      `json:"type"`
	Message ApiResponse `json:"message"`
}

type contentBlockDeltaEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta delta  `json:"delta"`
}

type delta struct {
	Type       string `json:"type"`
	Text       string `json:"text,omitempty"`
	InputJSON  string `json:"input_json,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
}

type messageDeltaEvent struct {
	Type  string `json:"type"`
	Delta delta  `json:"delta"`
	Usage Usage  `json:"usage"`
}

func parseStreamResponse(body io.Reader, onTextDelta func(text string)) (*ApiResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	var response ApiResponse
	var contentBlocks []ContentBlock
	var currentBlockIndex int
	var currentBlockType string
	var currentInputJSON strings.Builder
	var currentEventType string
	var currentData string

	processEvent := func(eventType, data string) {
		var rawData map[string]json.RawMessage
		if err := json.Unmarshal([]byte(data), &rawData); err != nil {
			var eventData map[string]any
			if err := json.Unmarshal([]byte(data), &eventData); err != nil {
				return
			}

			evtType, ok := eventData["type"].(string)
			if !ok {
				return
			}

			switch evtType {
			case "message_start":
				var msgEvent messageStartEvent
				if err := json.Unmarshal([]byte(data), &msgEvent); err == nil {
					response = msgEvent.Message
				}

			case "content_block_start":
				var blockEvent struct {
					Type  string `json:"type"`
					Index int    `json:"index"`
					Data  struct {
						Type string `json:"type"`
						ID   string `json:"id,omitempty"`
					} `json:"content_block"`
				}
				if err := json.Unmarshal([]byte(data), &blockEvent); err == nil {
					currentBlockIndex = blockEvent.Index
					contentBlock := ContentBlock{
						Type: blockEvent.Data.Type,
						ID:   blockEvent.Data.ID,
					}
					if blockEvent.Data.Type == "tool_use" {
						contentBlock.ToolUseID = blockEvent.Data.ID
					}
					contentBlocks = append(contentBlocks, contentBlock)
					currentBlockType = blockEvent.Data.Type
				}

			case "content_block_delta":
				if deltaObj, ok := eventData["delta"].(map[string]any); ok {
					deltaType, _ := deltaObj["type"].(string)

					if deltaType == "input_json_delta" {
						if inputJSON, ok := deltaObj["input_json"].(string); ok {
							currentInputJSON.WriteString(inputJSON)
						}
					} else if deltaType == "text_delta" {
						if text, ok := deltaObj["text"].(string); ok {
							if onTextDelta != nil {
								onTextDelta(text)
							}
							if len(contentBlocks) > currentBlockIndex {
								contentBlocks[currentBlockIndex].Text += text
							}
						}
					}
				}

			case "content_block_stop":
				if currentBlockType == "tool_use" && currentInputJSON.Len() > 0 {
					if len(contentBlocks) > currentBlockIndex {
						var input map[string]any
						if err := json.Unmarshal([]byte(currentInputJSON.String()), &input); err == nil {
							contentBlocks[currentBlockIndex].Input = input
						}
					}
					currentInputJSON.Reset()
				}
				currentBlockType = ""

			case "message_delta":
				if deltaObj, ok := eventData["delta"].(map[string]any); ok {
					if stopReason, ok := deltaObj["stop_reason"].(string); ok {
						response.StopReason = stopReason
					}
				}
				if usageObj, ok := eventData["usage"].(map[string]any); ok {
					if outputTokens, ok := usageObj["output_tokens"].(float64); ok {
						response.Usage.OutputTokens = int(outputTokens)
					}
				}

			case "message_stop":
				response.Content = contentBlocks
			}
			return
		}

		evtTypeRaw, ok := rawData["type"]
		if !ok {
			return
		}
		evtType := strings.Trim(string(evtTypeRaw), `"`)

		switch evtType {
		case "message_start":
			var msgEvent messageStartEvent
			if err := json.Unmarshal([]byte(data), &msgEvent); err == nil {
				response = msgEvent.Message
			}

		case "content_block_start":
			var blockEvent struct {
				Type  string `json:"type"`
				Index int    `json:"index"`
				Data  struct {
					Type string `json:"type"`
					ID   string `json:"id,omitempty"`
				} `json:"content_block"`
			}
			if err := json.Unmarshal([]byte(data), &blockEvent); err == nil {
				currentBlockIndex = blockEvent.Index
				contentBlock := ContentBlock{
					Type: blockEvent.Data.Type,
					ID:   blockEvent.Data.ID,
				}
				if blockEvent.Data.Type == "tool_use" {
					contentBlock.ToolUseID = blockEvent.Data.ID
				}
				contentBlocks = append(contentBlocks, contentBlock)
				currentBlockType = blockEvent.Data.Type
			}

		case "content_block_delta":
			var deltaEvent contentBlockDeltaEvent
			if err := json.Unmarshal([]byte(data), &deltaEvent); err == nil {
				if deltaEvent.Delta.Text != "" && onTextDelta != nil {
					onTextDelta(deltaEvent.Delta.Text)
				}
				if deltaEvent.Delta.InputJSON != "" {
					currentInputJSON.WriteString(deltaEvent.Delta.InputJSON)
				}
				if len(contentBlocks) > currentBlockIndex {
					contentBlocks[currentBlockIndex].Text += deltaEvent.Delta.Text
				}
			}

		case "content_block_stop":
			if currentBlockType == "tool_use" && currentInputJSON.Len() > 0 {
				if len(contentBlocks) > currentBlockIndex {
					var input map[string]any
					if err := json.Unmarshal([]byte(currentInputJSON.String()), &input); err == nil {
						contentBlocks[currentBlockIndex].Input = input
					}
				}
				currentInputJSON.Reset()
			}
			currentBlockType = ""

		case "message_delta":
			var msgDeltaEvent messageDeltaEvent
			if err := json.Unmarshal([]byte(data), &msgDeltaEvent); err == nil {
				response.StopReason = msgDeltaEvent.Delta.StopReason
				response.Usage = msgDeltaEvent.Usage
			}

		case "message_stop":
			response.Content = contentBlocks
		}
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.HasPrefix(line, "event:") {
			if currentEventType != "" && currentData != "" {
				processEvent(currentEventType, currentData)
			}
			currentEventType = strings.TrimPrefix(line, "event:")
			currentEventType = strings.TrimSpace(currentEventType)
			currentData = ""
			continue
		}

		if strings.HasPrefix(line, "data:") {
			currentData = strings.TrimPrefix(line, "data:")
			currentData = strings.TrimSpace(currentData)
			continue
		}

		if line == "" {
			if currentEventType != "" && currentData != "" {
				processEvent(currentEventType, currentData)
				currentEventType = ""
				currentData = ""
			}
		}
	}

	if currentEventType != "" && currentData != "" {
		processEvent(currentEventType, currentData)
	}

	response.Content = contentBlocks
	return &response, nil
}
