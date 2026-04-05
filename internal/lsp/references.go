package lsp

import (
	"context"
	"fmt"
)

type ReferenceParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	Context      ReferenceContext       `json:"context"`
}

type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

func (c *LSPClient) GetReferences(ctx context.Context, uri string, line, character int) ([]Location, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	params := ReferenceParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
		Position:     Position{Line: line, Character: character},
		Context:      ReferenceContext{IncludeDeclaration: true},
	}

	var result []Location
	err := c.call(ctx, MethodTextDocumentReferences, params, &result)
	if err != nil {
		return nil, fmt.Errorf("references request failed: %w", err)
	}

	if result == nil {
		return []Location{}, nil
	}

	return result, nil
}
