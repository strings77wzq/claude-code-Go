package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sync"
)

// StdioTransport implements stdio-based transport for MCP servers.
type StdioTransport struct {
	cmd     *exec.Cmd
	stdin   *os.File
	stdout  *bufio.Reader
	process *os.Process
	mu      sync.Mutex
	closed  bool
}

// NewStdioTransport creates a new stdio transport with the given command and arguments.
func NewStdioTransport(command string, args []string, env map[string]string) *StdioTransport {
	cmd := exec.Command(command, args...)
	if env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	return &StdioTransport{
		cmd: cmd,
	}
}

// Start starts the subprocess and creates stdin/stdout pipes.
func (t *StdioTransport) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("transport already closed")
	}

	// Create pipes for stdin and stdout
	stdinPipe, err := t.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	t.stdin = stdinPipe.(*os.File)

	stdoutPipe, err := t.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	t.stdout = bufio.NewReader(stdoutPipe)

	// Start the process
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	t.process = t.cmd.Process
	return nil
}

// SendRequest writes a JSON-RPC request to stdin.
func (t *StdioTransport) SendRequest(method string, params map[string]any, id int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed || t.stdin == nil {
		return fmt.Errorf("transport not started or closed")
	}

	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	data = append(data, '\n')
	if _, err := t.stdin.Write(data); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}

	if err := t.stdin.Sync(); err != nil {
		return fmt.Errorf("failed to sync stdin: %w", err)
	}

	return nil
}

// ReadResponse reads a JSON-RPC response from stdout (line-based).
func (t *StdioTransport) ReadResponse() (map[string]any, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed || t.stdout == nil {
		return nil, fmt.Errorf("transport not started or closed")
	}

	line, err := t.stdout.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Remove trailing newline
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	var resp map[string]any
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

// Close kills the subprocess and closes the transport.
func (t *StdioTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	// Close stdin to signal EOF to the process
	if t.stdin != nil {
		t.stdin.Close()
	}

	// Kill the process if still running
	if t.process != nil {
		if err := t.process.Kill(); err != nil {
			slog.Error("failed to kill process", "error", err)
		}
		if _, err := t.process.Wait(); err != nil {
			slog.Error("failed to wait for process", "error", err)
		}
	}

	return nil
}
