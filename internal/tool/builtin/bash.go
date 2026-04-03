package builtin

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

const (
	defaultTimeout = 120 * time.Second
	maxOutputSize  = 100 * 1024
)

type BashTool struct {
	workingDir string
}

func NewBashTool(workingDir string) tool.Tool {
	return &BashTool{
		workingDir: workingDir,
	}
}

func (b *BashTool) Name() string {
	return "Bash"
}

func (b *BashTool) Description() string {
	return "Executes a bash command and returns its output."
}

func (b *BashTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The bash command to execute",
			},
			"timeout": map[string]any{
				"type":        "number",
				"description": "Timeout in seconds (default: 120)",
			},
		},
		"required": []string{"command"},
	}
}

func (b *BashTool) RequiresPermission() bool {
	return true
}

func (b *BashTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	command, ok := input["command"].(string)
	if !ok || command == "" {
		return tool.Error("command is required")
	}

	timeout := defaultTimeout
	if t, ok := input["timeout"]; ok {
		switch v := t.(type) {
		case float64:
			timeout = time.Duration(v) * time.Second
		case int64:
			timeout = time.Duration(v) * time.Second
		case int:
			timeout = time.Duration(v) * time.Second
		}
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = b.workingDir
	cmd.Stderr = cmd.Stdout

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to create stdout pipe: %v", err))
	}

	var mu sync.Mutex
	var output strings.Builder

	go func() {
		buf := make([]byte, 4096)
		for {
			n, readErr := stdout.Read(buf)
			if n > 0 {
				mu.Lock()
				if output.Len()+n > maxOutputSize {
					if output.Len() < maxOutputSize {
						remaining := maxOutputSize - output.Len()
						output.Write(buf[:remaining])
						output.WriteString("... (output truncated)")
					}
					mu.Unlock()
					return
				}
				output.Write(buf[:n])
				mu.Unlock()
			}
			if readErr != nil {
				break
			}
		}
	}()

	err = cmd.Run()

	mu.Lock()
	result := output.String()
	mu.Unlock()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return tool.Error(fmt.Sprintf("command timed out after %v", timeout))
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			if result != "" {
				return tool.Result{
					Content: result + fmt.Sprintf("\n[exit code: %d]", exitErr.ExitCode()),
					IsError: true,
				}
			}
			return tool.Error(fmt.Sprintf("command failed with exit code %d", exitErr.ExitCode()))
		}
		return tool.Error(fmt.Sprintf("command failed: %v", err))
	}

	return tool.Success(result)
}
