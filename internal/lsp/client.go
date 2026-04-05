package lsp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/strings77wzq/claude-code-Go/internal/logger"
)

const (
	JSONRPCVersion                       = "2.0"
	MethodInitialize                     = "initialize"
	MethodShutdown                       = "shutdown"
	MethodExit                           = "exit"
	MethodTextDocumentPublishDiagnostics = "textDocument/publishDiagnostics"
	MethodWorkspaceSymbol                = "workspace/symbol"
	MethodTextDocumentDocumentSymbol     = "textDocument/documentSymbol"
	MethodTextDocumentReferences         = "textDocument/references"
	MethodTextDocumentDefinition         = "textDocument/definition"
	MethodTextDocumentHover              = "textDocument/hover"
)

type LSPClient struct {
	serverURL    string
	httpClient   *http.Client
	mu           sync.Mutex
	nextID       int64
	initialized  bool
	capabilities ServerCapabilities
	log          *slog.Logger
}

type ServerCapabilities struct {
	TextDocumentSync           interface{} `json:"textDocumentSync,omitempty"`
	HoverProvider              interface{} `json:"hoverProvider,omitempty"`
	DefinitionProvider         interface{} `json:"definitionProvider,omitempty"`
	ReferencesProvider         interface{} `json:"referencesProvider,omitempty"`
	DocumentSymbolProvider     interface{} `json:"documentSymbolProvider,omitempty"`
	WorkspaceSymbolProvider    interface{} `json:"workspaceSymbolProvider,omitempty"`
	PublishDiagnosticsProvider interface{} `json:"publishDiagnosticsProvider,omitempty"`
}

type InitializeParams struct {
	ProcessID    interface{}        `json:"processId,omitempty"`
	RootURI      string             `json:"rootUri,omitempty"`
	RootPath     string             `json:"rootPath,omitempty"`
	ClientInfo   *ClientInfo        `json:"clientInfo,omitempty"`
	Capabilities ClientCapabilities `json:"capabilities,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type ClientCapabilities struct {
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
}

type TextDocumentClientCapabilities struct {
	Synchronization *TextDocumentSyncCapabilities `json:"synchronization,omitempty"`
}

type TextDocumentSyncCapabilities struct {
	WillSave          bool `json:"willSave,omitempty"`
	DidSave           bool `json:"didSave,omitempty"`
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
}

type WorkspaceClientCapabilities struct {
	ApplyEdit        bool `json:"applyEdit,omitempty"`
	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`
	WorkspaceEdit    bool `json:"workspaceEdit,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func NewLSPClient(serverURL string) *LSPClient {
	return &LSPClient{
		serverURL:  serverURL,
		httpClient: &http.Client{},
		nextID:     1,
		log:        logger.Log(),
	}
}

func (c *LSPClient) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("client already initialized")
	}

	params := InitializeParams{
		ProcessID: 1,
		ClientInfo: &ClientInfo{
			Name:    "claude-code-Go",
			Version: "1.0.0",
		},
		Capabilities: ClientCapabilities{
			TextDocument: &TextDocumentClientCapabilities{
				Synchronization: &TextDocumentSyncCapabilities{
					WillSave: true,
					DidSave:  true,
				},
			},
			Workspace: &WorkspaceClientCapabilities{
				WorkspaceFolders: true,
			},
		},
	}

	var result InitializeResult
	err := c.call(ctx, MethodInitialize, params, &result)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	c.capabilities = result.Capabilities
	c.initialized = true
	c.log.Info("LSP client initialized successfully")
	return nil
}

func (c *LSPClient) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return fmt.Errorf("client not initialized")
	}

	var result interface{}
	err := c.call(ctx, MethodShutdown, nil, &result)
	if err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	c.initialized = false
	c.log.Info("LSP client shutdown completed")
	return nil
}

func (c *LSPClient) IsInitialized() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.initialized
}

func (c *LSPClient) GetCapabilities() ServerCapabilities {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.capabilities
}

func (c *LSPClient) call(ctx context.Context, method string, params interface{}, result interface{}) error {
	id := c.nextID
	c.nextID++

	reqParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	request := JSONRPCRequest{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  method,
		Params:  reqParams,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.serverURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vscode-jsonrpc")
	req.Header.Set("Accept", "application/vscode-jsonrpc")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)

	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read headers: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			fmt.Sscanf(line, "Content-Length: %d", &contentLength)
		}
	}

	if contentLength == 0 {
		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		contentLength = len(bodyBytes)
		reader = bufio.NewReader(bytes.NewReader(bodyBytes))
	}

	respBody := make([]byte, contentLength)
	_, err = io.ReadFull(reader, respBody)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("JSON-RPC error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	if len(rpcResp.Result) > 0 {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return nil
}

func (c *LSPClient) callNotification(ctx context.Context, method string, params interface{}) error {
	reqParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	notification := Notification{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  reqParams,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.serverURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vscode-jsonrpc")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("notification failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
