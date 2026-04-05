package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type Diagnostic struct {
	Range    Range                          `json:"range"`
	Severity Severity                       `json:"severity,omitempty"`
	Source   string                         `json:"source,omitempty"`
	Message  string                         `json:"message"`
	Code     interface{}                    `json:"code,omitempty"`
	Related  []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

type Severity int

const (
	SeverityError       Severity = 1
	SeverityWarning     Severity = 2
	SeverityInformation Severity = 3
	SeverityHint        Severity = 4
)

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     int          `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type diagnosticsCache struct {
	mu        sync.RWMutex
	diags     map[string][]Diagnostic
	listeners []chan []Diagnostic
}

var globalDiagnosticsCache = &diagnosticsCache{
	diags: make(map[string][]Diagnostic),
}

func (c *LSPClient) GetDiagnostics(ctx context.Context, uri string) ([]Diagnostic, error) {
	c.mu.Lock()
	if !c.initialized {
		c.mu.Unlock()
		return nil, fmt.Errorf("client not initialized")
	}
	c.mu.Unlock()

	globalDiagnosticsCache.mu.RLock()
	diags, exists := globalDiagnosticsCache.diags[uri]
	globalDiagnosticsCache.mu.RUnlock()

	if exists {
		return diags, nil
	}

	return []Diagnostic{}, nil
}

func (c *LSPClient) handleDiagnosticsNotification(params json.RawMessage) error {
	var publishParams PublishDiagnosticsParams
	if err := json.Unmarshal(params, &publishParams); err != nil {
		return fmt.Errorf("failed to unmarshal diagnostics: %w", err)
	}

	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.diags[publishParams.URI] = publishParams.Diagnostics
	globalDiagnosticsCache.mu.Unlock()

	globalDiagnosticsCache.mu.RLock()
	listeners := make([]chan []Diagnostic, len(globalDiagnosticsCache.listeners))
	copy(listeners, globalDiagnosticsCache.listeners)
	globalDiagnosticsCache.mu.RUnlock()

	for _, listener := range listeners {
		select {
		case listener <- publishParams.Diagnostics:
		default:
		}
	}

	return nil
}

func (c *LSPClient) SubscribeDiagnostics() chan []Diagnostic {
	ch := make(chan []Diagnostic, 1)
	globalDiagnosticsCache.mu.Lock()
	globalDiagnosticsCache.listeners = append(globalDiagnosticsCache.listeners, ch)
	globalDiagnosticsCache.mu.Unlock()
	return ch
}

func (c *LSPClient) UnsubscribeDiagnostics(ch chan []Diagnostic) {
	globalDiagnosticsCache.mu.Lock()
	defer globalDiagnosticsCache.mu.Unlock()
	for i, listener := range globalDiagnosticsCache.listeners {
		if listener == ch {
			globalDiagnosticsCache.listeners = append(globalDiagnosticsCache.listeners[:i], globalDiagnosticsCache.listeners[i+1:]...)
			close(ch)
			return
		}
	}
}

func (c *LSPClient) ClearDiagnostics(uri string) {
	globalDiagnosticsCache.mu.Lock()
	delete(globalDiagnosticsCache.diags, uri)
	globalDiagnosticsCache.mu.Unlock()
}
