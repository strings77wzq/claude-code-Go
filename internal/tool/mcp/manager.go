package mcp

import (
	"fmt"
	"log"

	"github.com/user/go-code/internal/tool"
)

// McpManager manages MCP client connections and tool adapters.
type McpManager struct {
	clients  map[string]*McpClient
	adapters map[string]*McpToolAdapter
}

// NewMcpManager creates a new MCP manager.
func NewMcpManager() *McpManager {
	return &McpManager{
		clients:  make(map[string]*McpClient),
		adapters: make(map[string]*McpToolAdapter),
	}
}

// InitializeAndRegister initializes MCP servers from configs and registers their tools.
func (m *McpManager) InitializeAndRegister(configs map[string]McpServerConfig, registry *tool.Registry) error {
	for name, config := range configs {
		if err := m.initializeServer(name, config, registry); err != nil {
			log.Printf("Failed to initialize MCP server %q: %v", name, err)
		}
	}
	return nil
}

// initializeServer initializes a single MCP server.
func (m *McpManager) initializeServer(name string, config McpServerConfig, registry *tool.Registry) error {
	transport := NewStdioTransport(config.Command, config.Args, config.Env)
	if err := transport.Start(); err != nil {
		return fmt.Errorf("failed to start transport: %w", err)
	}

	client := NewMcpClient(transport)
	if err := client.Initialize(); err != nil {
		transport.Close()
		return fmt.Errorf("failed to initialize: %w", err)
	}

	tools, err := client.ListTools()
	if err != nil {
		transport.Close()
		return fmt.Errorf("failed to list tools: %w", err)
	}

	m.clients[name] = client

	for _, toolInfo := range tools {
		adapter, err := newMcpToolAdapter(name, toolInfo, client)
		if err != nil {
			log.Printf("Failed to create adapter for tool %q: %v", toolInfo.Name, err)
			continue
		}

		if err := registry.Register(adapter); err != nil {
			log.Printf("Failed to register tool %q: %v", adapter.Name(), err)
			continue
		}

		m.adapters[adapter.Name()] = adapter
	}

	return nil
}

// Close closes all MCP client connections.
func (m *McpManager) Close() error {
	var lastErr error
	for name, client := range m.clients {
		if err := client.transport.Close(); err != nil {
			if lastErr == nil {
				lastErr = fmt.Errorf("failed to close client %q: %w", name, err)
			}
		}
	}
	m.clients = make(map[string]*McpClient)
	m.adapters = make(map[string]*McpToolAdapter)
	return lastErr
}
