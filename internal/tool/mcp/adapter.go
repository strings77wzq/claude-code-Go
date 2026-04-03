package mcp

import (
	"context"
	"encoding/json"

	"github.com/user/go-code/internal/tool"
)

// McpToolAdapter wraps an MCP tool and implements the tool.Tool interface.
type McpToolAdapter struct {
	serverName  string
	toolName    string
	description string
	inputSchema map[string]any
	client      *McpClient
}

// Name returns the unique name of the tool in format mcp__{serverName}__{toolName}.
func (a *McpToolAdapter) Name() string {
	return "mcp__" + a.serverName + "__" + a.toolName
}

// Description returns the human-readable description of the tool.
func (a *McpToolAdapter) Description() string {
	return a.description
}

// InputSchema returns the JSON schema for the tool's input parameters.
func (a *McpToolAdapter) InputSchema() map[string]any {
	return a.inputSchema
}

// RequiresPermission always returns true since external MCP tools require approval.
func (a *McpToolAdapter) RequiresPermission() bool {
	return true
}

// Execute delegates to McpClient.CallTool and returns a tool.Result.
func (a *McpToolAdapter) Execute(ctx context.Context, input map[string]any) tool.Result {
	result, err := a.client.CallTool(a.toolName, input)
	if err != nil {
		return tool.Error(err.Error())
	}
	return tool.Success(result)
}

// newMcpToolAdapter creates a new MCP tool adapter from tool info.
func newMcpToolAdapter(serverName string, toolInfo McpToolInfo, client *McpClient) (*McpToolAdapter, error) {
	var schema map[string]any
	if len(toolInfo.InputSchema) > 0 {
		if err := json.Unmarshal(toolInfo.InputSchema, &schema); err != nil {
			return nil, err
		}
	}

	return &McpToolAdapter{
		serverName:  serverName,
		toolName:    toolInfo.Name,
		description: toolInfo.Description,
		inputSchema: schema,
		client:      client,
	}, nil
}
