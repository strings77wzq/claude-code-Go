package mcp

import (
	"encoding/json"
	"fmt"
)

// McpToolInfo represents a tool definition from the MCP server.
type McpToolInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// McpClient wraps the transport and provides high-level MCP operations.
type McpClient struct {
	transport *StdioTransport
	nextID    int
}

// NewMcpClient creates a new MCP client with the given transport.
func NewMcpClient(transport *StdioTransport) *McpClient {
	return &McpClient{
		transport: transport,
		nextID:    1,
	}
}

// Initialize sends the initialize JSON-RPC request, receives server info,
// and sends the notifications/initialized message.
func (c *McpClient) Initialize() error {
	// Send initialize request
	params := map[string]any{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]any{},
		"clientInfo": map[string]any{
			"name":    "go-code",
			"version": "0.1.0",
		},
	}

	id := c.nextID
	c.nextID++

	if err := c.transport.SendRequest("initialize", params, id); err != nil {
		return fmt.Errorf("failed to send initialize request: %w", err)
	}

	// Read response
	resp, err := c.transport.ReadResponse()
	if err != nil {
		return fmt.Errorf("failed to read initialize response: %w", err)
	}

	// Check for JSON-RPC error
	if errMsg, ok := resp["error"]; ok {
		return fmt.Errorf("initialize error: %v", errMsg)
	}

	// Verify ID matches
	if respID, ok := resp["id"].(float64); !ok || int(respID) != id {
		return fmt.Errorf("response ID mismatch")
	}

	// Send notifications/initialized
	notifParams := map[string]any{}
	if err := c.transport.SendRequest("notifications/initialized", notifParams, 0); err != nil {
		return fmt.Errorf("failed to send initialized notification: %w", err)
	}

	return nil
}

// ListTools sends the tools/list request and returns tool definitions.
func (c *McpClient) ListTools() ([]McpToolInfo, error) {
	params := map[string]any{}

	id := c.nextID
	c.nextID++

	if err := c.transport.SendRequest("tools/list", params, id); err != nil {
		return nil, fmt.Errorf("failed to send tools/list request: %w", err)
	}

	// Read response
	resp, err := c.transport.ReadResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to read tools/list response: %w", err)
	}

	// Check for JSON-RPC error
	if errMsg, ok := resp["error"]; ok {
		return nil, fmt.Errorf("tools/list error: %v", errMsg)
	}

	// Extract result
	result, ok := resp["result"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid result format")
	}

	// Extract tools array
	toolsRaw, ok := result["tools"]
	if !ok {
		return nil, fmt.Errorf("no tools in result")
	}

	toolsSlice, ok := toolsRaw.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid tools format")
	}

	tools := make([]McpToolInfo, 0, len(toolsSlice))
	for i, t := range toolsSlice {
		toolMap, ok := t.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid tool format at index %d", i)
		}

		tool := McpToolInfo{
			Name:        getString(toolMap, "name"),
			Description: getString(toolMap, "description"),
		}

		// Serialize inputSchema to JSON
		if schema, ok := toolMap["inputSchema"]; ok {
			schemaBytes, err := json.Marshal(schema)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal inputSchema: %w", err)
			}
			tool.InputSchema = schemaBytes
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

// CallTool sends a tools/call request and returns the result text.
func (c *McpClient) CallTool(name string, args map[string]any) (string, error) {
	params := map[string]any{
		"name":      name,
		"arguments": args,
	}

	id := c.nextID
	c.nextID++

	if err := c.transport.SendRequest("tools/call", params, id); err != nil {
		return "", fmt.Errorf("failed to send tools/call request: %w", err)
	}

	// Read response
	resp, err := c.transport.ReadResponse()
	if err != nil {
		return "", fmt.Errorf("failed to read tools/call response: %w", err)
	}

	// Check for JSON-RPC error
	if errMsg, ok := resp["error"]; ok {
		return "", fmt.Errorf("tools/call error: %v", errMsg)
	}

	// Extract result
	result, ok := resp["result"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid result format")
	}

	// Extract content array
	contentRaw, ok := result["content"]
	if !ok {
		return "", fmt.Errorf("no content in result")
	}

	contentSlice, ok := contentRaw.([]any)
	if !ok || len(contentSlice) == 0 {
		return "", fmt.Errorf("invalid content format")
	}

	// Get first content item
	firstContent, ok := contentSlice[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid content item format")
	}

	// Extract text
	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text in content")
	}

	// Check isError and return error if set
	if isError, ok := result["isError"].(bool); ok && isError {
		return text, fmt.Errorf("tool returned error: %s", text)
	}

	return text, nil
}

// getString helper to extract string from map
func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
