package builtin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

const (
	defaultLimit    = 2000
	defaultOffset   = 0
	maxFileReadSize = 200 * 1024
)

type ReadTool struct{}

func NewReadTool() tool.Tool {
	return &ReadTool{}
}

func (r *ReadTool) Name() string {
	return "Read"
}

func (r *ReadTool) Description() string {
	return "Reads a file and returns its contents with line numbers."
}

func (r *ReadTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Path to the file to read",
			},
			"offset": map[string]any{
				"type":        "number",
				"description": "Line number to start reading from (0-based, default: 0)",
			},
			"limit": map[string]any{
				"type":        "number",
				"description": "Maximum number of lines to read (default: 2000)",
			},
		},
		"required": []string{"file_path"},
	}
}

func (r *ReadTool) RequiresPermission() bool {
	return false
}

func (r *ReadTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return tool.Error("file_path is required")
	}

	offset := defaultOffset
	if o, ok := input["offset"]; ok {
		switch v := o.(type) {
		case float64:
			offset = int(v)
		case int64:
			offset = int(v)
		case int:
			offset = v
		}
	}

	limit := defaultLimit
	if l, ok := input["limit"]; ok {
		switch v := l.(type) {
		case float64:
			limit = int(v)
		case int64:
			limit = int(v)
		case int:
			limit = v
		}
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tool.Error(fmt.Sprintf("file not found: %s", filePath))
		}
		return tool.Error(fmt.Sprintf("failed to stat file: %v", err))
	}

	if info.IsDir() {
		return tool.Error(fmt.Sprintf("cannot read directory: %s", filePath))
	}

	if info.Size() > maxFileReadSize {
		return tool.Error(fmt.Sprintf("file too large (max 200KB): %s", filePath))
	}

	file, err := os.Open(filePath)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to open file: %v", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	lineNum := 0

	for scanner.Scan() {
		if lineNum >= offset {
			lines = append(lines, scanner.Text())
			if len(lines) >= limit {
				break
			}
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return tool.Error(fmt.Sprintf("error reading file: %v", err))
	}

	var sb strings.Builder
	for i, line := range lines {
		sb.WriteString(strconv.Itoa(offset + i + 1))
		sb.WriteString(": ")
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return tool.Success(sb.String())
}
