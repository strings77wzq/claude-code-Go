package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

func TestNewMcpManager(t *testing.T) {
	mgr := NewMcpManager()
	if mgr == nil {
		t.Fatal("NewMcpManager returned nil")
	}

	if len(mgr.clients) != 0 {
		t.Errorf("expected empty clients map, got %d entries", len(mgr.clients))
	}
	if len(mgr.adapters) != 0 {
		t.Errorf("expected empty adapters map, got %d entries", len(mgr.adapters))
	}
}

func TestMcpManagerInitializeInvalidConfig(t *testing.T) {
	mgr := NewMcpManager()
	registry := tool.NewRegistry()

	configs := map[string]McpServerConfig{
		"nonexistent": {
			Command: "/no/such/binary/ever",
			Args:    []string{},
			Env:     nil,
		},
	}

	// InitializeAndRegister should return nil (it logs errors and continues)
	err := mgr.InitializeAndRegister(configs, registry)
	if err != nil {
		t.Fatalf("InitializeAndRegister returned unexpected error: %v", err)
	}

	// The server should not have been registered (no clients or adapters)
	if len(mgr.clients) != 0 {
		t.Errorf("expected 0 clients after failed init, got %d", len(mgr.clients))
	}
	if len(mgr.adapters) != 0 {
		t.Errorf("expected 0 adapters after failed init, got %d", len(mgr.adapters))
	}
}

func TestMcpToolAdapterPermissionRequired(t *testing.T) {
	adapter := &McpToolAdapter{
		serverName:  "test-server",
		toolName:    "test-tool",
		description: "A test MCP tool",
	}

	if !adapter.RequiresPermission() {
		t.Error("MCP tool adapter should require permission")
	}
	if adapter.RequiredPermissionLevel() != 2 {
		t.Errorf("expected LevelDangerFullAccess (2), got %d", adapter.RequiredPermissionLevel())
	}
	if adapter.Name() != "mcp__test-server__test-tool" {
		t.Errorf("expected namespaced name 'mcp__test-server__test-tool', got %q", adapter.Name())
	}
}

func TestMcpToolAdapterNoSchema(t *testing.T) {
	adapter, err := newMcpToolAdapter("srv", McpToolInfo{
		Name:        "no-schema-tool",
		Description: "Tool without input schema",
	}, nil)
	if err != nil {
		t.Fatalf("newMcpToolAdapter failed: %v", err)
	}
	if adapter.InputSchema() != nil {
		t.Error("expected nil input schema when none provided")
	}
}

func TestMcpToolAdapterInvalidSchema(t *testing.T) {
	_, err := newMcpToolAdapter("srv", McpToolInfo{
		Name:        "bad-schema-tool",
		Description: "Tool with invalid input schema",
		InputSchema: json.RawMessage(`{"type":`),
	}, nil)
	if err == nil {
		t.Fatal("expected invalid schema to return an error")
	}
}

func TestMcpToolAdapterExecuteReturnsClientError(t *testing.T) {
	adapter := &McpToolAdapter{
		serverName:  "srv",
		toolName:    "tool",
		description: "Tool backed by an unstarted client",
		client:      NewMcpClient(NewStdioTransport(os.Args[0], nil, nil)),
	}

	result := adapter.Execute(context.Background(), map[string]any{"x": "y"})
	if !result.IsError {
		t.Fatal("expected Execute to return tool error")
	}
	if !strings.Contains(result.Content, "transport not started or closed") {
		t.Fatalf("unexpected error content: %s", result.Content)
	}
}

func TestMcpManagerCloseEmpty(t *testing.T) {
	mgr := NewMcpManager()
	err := mgr.Close()
	if err != nil {
		t.Fatalf("Close on empty manager returned error: %v", err)
	}
}

func TestStdioTransportNotStartedAndClosedErrors(t *testing.T) {
	transport := NewStdioTransport(os.Args[0], nil, nil)

	if err := transport.SendRequest("tools/list", map[string]any{}, 1); err == nil {
		t.Fatal("expected SendRequest before Start to fail")
	}
	if _, err := transport.ReadResponse(); err == nil {
		t.Fatal("expected ReadResponse before Start to fail")
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("Close before Start returned error: %v", err)
	}
	if err := transport.SendRequest("tools/list", map[string]any{}, 2); err == nil {
		t.Fatal("expected SendRequest after Close to fail")
	}
	if _, err := transport.ReadResponse(); err == nil {
		t.Fatal("expected ReadResponse after Close to fail")
	}
	if err := transport.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}
}

func TestStdioTransportReadResponseContextTimesOut(t *testing.T) {
	reader, writer := io.Pipe()
	defer writer.Close()

	transport := &StdioTransport{
		stdout: bufio.NewReader(reader),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := transport.ReadResponseContext(ctx)
	if err == nil {
		t.Fatal("expected context timeout")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMcpClientListToolsContextTimesOut(t *testing.T) {
	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer stdinReader.Close()
	defer stdinWriter.Close()

	stdoutReader, stdoutWriter := io.Pipe()
	defer stdoutWriter.Close()

	transport := &StdioTransport{
		stdin:  stdinWriter,
		stdout: bufio.NewReader(stdoutReader),
	}
	client := NewMcpClient(transport)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = client.ListToolsContext(ctx)
	if err == nil {
		t.Fatal("expected list tools timeout")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetStringMissingOrWrongType(t *testing.T) {
	values := map[string]any{"name": 123}
	if got := getString(values, "name"); got != "" {
		t.Fatalf("getString wrong type = %q, want empty", got)
	}
	if got := getString(values, "missing"); got != "" {
		t.Fatalf("getString missing = %q, want empty", got)
	}
}

func TestMcpManagerRegistersToolsFromFixture(t *testing.T) {
	mgr := NewMcpManager()
	registry := tool.NewRegistry()

	err := mgr.InitializeAndRegister(map[string]McpServerConfig{
		"fixture": {
			Command: os.Args[0],
			Args:    []string{"-test.run=TestMcpFixtureServer"},
			Env: map[string]string{
				"GO_WANT_MCP_FIXTURE": "1",
			},
		},
	}, registry)
	if err != nil {
		t.Fatalf("InitializeAndRegister returned error: %v", err)
	}
	defer mgr.Close()

	registered := registry.GetTool("mcp__fixture__echo")
	if registered == nil {
		t.Fatalf("expected namespaced MCP tool to be registered")
	}

	result := registered.Execute(context.Background(), map[string]any{"text": "hello"})
	if result.IsError {
		t.Fatalf("expected MCP fixture tool success, got error: %s", result.Content)
	}
	if result.Content != "fixture: hello" {
		t.Fatalf("unexpected fixture result: %q", result.Content)
	}
}

func TestMcpToolAdapterUsesPermissionPolicy(t *testing.T) {
	adapter := &McpToolAdapter{
		serverName:  "external",
		toolName:    "write-file",
		description: "External write-like action",
		inputSchema: map[string]any{"type": "object"},
	}
	policy := permission.NewPolicy(permission.WorkspaceWrite)

	decision := policy.Evaluate(adapter.Name(), map[string]any{"path": "x"}, adapter.RequiresPermission())
	if decision != permission.Ask {
		t.Fatalf("expected MCP tool to require explicit approval in workspace mode, got %s", decision)
	}
}

func TestMcpFixtureServer(t *testing.T) {
	if os.Getenv("GO_WANT_MCP_FIXTURE") != "1" {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for scanner.Scan() {
		var req map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			writeFixtureResponse(encoder, nil, nil, fmt.Sprintf("invalid request: %v", err))
			continue
		}
		method, _ := req["method"].(string)
		id := req["id"]

		switch method {
		case "initialize":
			writeFixtureResponse(encoder, id, map[string]any{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]any{
					"name":    "fixture",
					"version": "test",
				},
				"capabilities": map[string]any{
					"tools": map[string]any{},
				},
			}, "")
		case "notifications/initialized":
			continue
		case "tools/list":
			writeFixtureResponse(encoder, id, map[string]any{
				"tools": []map[string]any{
					{
						"name":        "echo",
						"description": "Echo text from fixture",
						"inputSchema": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"text": map[string]any{"type": "string"},
							},
						},
					},
				},
			}, "")
		case "tools/call":
			params, _ := req["params"].(map[string]any)
			args, _ := params["arguments"].(map[string]any)
			text, _ := args["text"].(string)
			writeFixtureResponse(encoder, id, map[string]any{
				"content": []map[string]any{
					{"type": "text", "text": "fixture: " + text},
				},
			}, "")
		default:
			writeFixtureResponse(encoder, id, nil, "unknown method: "+method)
		}
	}
	os.Exit(0)
}

func writeFixtureResponse(encoder *json.Encoder, id any, result map[string]any, errMsg string) {
	resp := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
	}
	if errMsg != "" {
		resp["error"] = map[string]any{"message": errMsg}
	} else {
		resp["result"] = result
	}
	_ = encoder.Encode(resp)
}
