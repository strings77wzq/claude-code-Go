package lsp

import (
	"context"
	"fmt"
)

type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkupKind string

const (
	MarkupKindPlainText MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)

func (c *LSPClient) GetHover(ctx context.Context, uri string, line, character int) (*Hover, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	params := HoverParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Position:     Position{Line: line, Character: character},
	}

	var result interface{}
	err := c.call(ctx, MethodTextDocumentHover, params, &result)
	if err != nil {
		return nil, fmt.Errorf("hover request failed: %w", err)
	}

	if result == nil {
		return nil, nil
	}

	switch v := result.(type) {
	case map[string]interface{}:
		hover := &Hover{}

		if contents, ok := v["contents"].(string); ok {
			hover.Contents = MarkupContent{
				Kind:  MarkupKindPlainText,
				Value: contents,
			}
		} else if contentsMap, ok := v["contents"].(map[string]interface{}); ok {
			kindStr, _ := contentsMap["kind"].(string)
			value, _ := contentsMap["value"].(string)
			hover.Contents = MarkupContent{
				Kind:  MarkupKind(kindStr),
				Value: value,
			}
		}

		if rng, ok := v["range"].(map[string]interface{}); ok {
			start, _ := rng["start"].(map[string]interface{})
			end, _ := rng["end"].(map[string]interface{})
			if start != nil && end != nil {
				hover.Range = &Range{}
				if line, ok := start["line"].(float64); ok {
					hover.Range.Start.Line = int(line)
				}
				if char, ok := start["character"].(float64); ok {
					hover.Range.Start.Character = int(char)
				}
				if line, ok := end["line"].(float64); ok {
					hover.Range.End.Line = int(line)
				}
				if char, ok := end["character"].(float64); ok {
					hover.Range.End.Character = int(char)
				}
			}
		}

		return hover, nil
	default:
		return nil, nil
	}
}
