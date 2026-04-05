package lsp

import (
	"context"
	"encoding/json"
	"fmt"
)

type SymbolInformation struct {
	Name          string     `json:"name"`
	Kind          SymbolKind `json:"kind"`
	Location      Location   `json:"location"`
	ContainerName string     `json:"containerName,omitempty"`
}

type SymbolKind int

const (
	SymbolKindFile        SymbolKind = 1
	SymbolKindModule      SymbolKind = 2
	SymbolKindNamespace   SymbolKind = 3
	SymbolKindPackage     SymbolKind = 4
	SymbolKindClass       SymbolKind = 5
	SymbolKindMethod      SymbolKind = 6
	SymbolKindProperty    SymbolKind = 7
	SymbolKindField       SymbolKind = 8
	SymbolKindConstructor SymbolKind = 9
	SymbolKindEnum        SymbolKind = 10
	SymbolKindInterface   SymbolKind = 11
	SymbolKindFunction    SymbolKind = 12
	SymbolKindVariable    SymbolKind = 13
	SymbolKindConstant    SymbolKind = 14
	SymbolKindString      SymbolKind = 15
	SymbolKindNumber      SymbolKind = 16
	SymbolKindBoolean     SymbolKind = 17
	SymbolKindArray       SymbolKind = 18
)

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Kind           SymbolKind       `json:"kind"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
	Detail         string           `json:"detail,omitempty"`
	Deprecated     bool             `json:"deprecated,omitempty"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type WorkspaceSymbolParams struct {
	Query string `json:"query"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

func (c *LSPClient) GetSymbols(ctx context.Context, query string) ([]SymbolInformation, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	params := WorkspaceSymbolParams{Query: query}

	var result []SymbolInformation
	err := c.call(ctx, MethodWorkspaceSymbol, params, &result)
	if err != nil {
		return nil, fmt.Errorf("workspace symbol request failed: %w", err)
	}

	return result, nil
}

func (c *LSPClient) GetDocumentSymbols(ctx context.Context, uri string) ([]DocumentSymbol, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	params := DocumentSymbolParams{
		TextDocument: TextDocumentIdentifier{URI: uri},
	}

	var result []DocumentSymbol
	err := c.call(ctx, MethodTextDocumentDocumentSymbol, params, &result)
	if err != nil {
		if err == nil && result == nil {
			var legacyResult []json.RawMessage
			errLegacy := c.call(ctx, MethodTextDocumentDocumentSymbol, params, &legacyResult)
			if errLegacy == nil && len(legacyResult) > 0 {
				symbols := make([]DocumentSymbol, 0, len(legacyResult))
				for _, raw := range legacyResult {
					var sym DocumentSymbol
					if json.Unmarshal(raw, &sym) == nil {
						symbols = append(symbols, sym)
					}
				}
				return symbols, nil
			}
		}
		return nil, fmt.Errorf("document symbol request failed: %w", err)
	}

	return result, nil
}
