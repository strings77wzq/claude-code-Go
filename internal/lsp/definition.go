package lsp

import (
	"context"
	"fmt"
)

type DefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

func (c *LSPClient) GetDefinition(ctx context.Context, uri string, line, character int) ([]Location, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	params := DefinitionParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Position:     Position{Line: line, Character: character},
	}

	var result interface{}
	err := c.call(ctx, MethodTextDocumentDefinition, params, &result)
	if err != nil {
		return nil, fmt.Errorf("definition request failed: %w", err)
	}

	if result == nil {
		return []Location{}, nil
	}

	switch v := result.(type) {
	case []interface{}:
		locations := make([]Location, 0, len(v))
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				loc, err := convertMapToLocation(itemMap)
				if err == nil {
					locations = append(locations, loc)
				}
			}
		}
		return locations, nil
	case map[string]interface{}:
		loc, err := convertMapToLocation(v)
		if err != nil {
			return nil, err
		}
		return []Location{loc}, nil
	default:
		return []Location{}, nil
	}
}

func convertMapToLocation(m map[string]interface{}) (Location, error) {
	loc := Location{}

	if uri, ok := m["uri"].(string); ok {
		loc.URI = uri
	}

	if rng, ok := m["range"].(map[string]interface{}); ok {
		start, _ := rng["start"].(map[string]interface{})
		end, _ := rng["end"].(map[string]interface{})
		if start != nil {
			if line, ok := start["line"].(float64); ok {
				loc.Range.Start.Line = int(line)
			}
			if char, ok := start["character"].(float64); ok {
				loc.Range.Start.Character = int(char)
			}
		}
		if end != nil {
			if line, ok := end["line"].(float64); ok {
				loc.Range.End.Line = int(line)
			}
			if char, ok := end["character"].(float64); ok {
				loc.Range.End.Character = int(char)
			}
		}
	}

	return loc, nil
}
