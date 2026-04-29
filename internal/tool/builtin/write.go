package builtin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type WriteTool struct {
	workingDir string
}

func NewWriteTool(workingDir string) tool.Tool {
	return &WriteTool{
		workingDir: workingDir,
	}
}

func (w *WriteTool) Name() string {
	return "Write"
}

func (w *WriteTool) Description() string {
	return "Writes content to a file, creating parent directories if needed."
}

func (w *WriteTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Path to the file to write",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "Content to write to the file",
			},
		},
		"required": []string{"file_path", "content"},
	}
}

func (w *WriteTool) RequiresPermission() bool {
	return true
}

func (w *WriteTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelWorkspaceWrite
}

func (w *WriteTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return tool.Error("file_path is required")
	}

	resolvedPath, err := ResolvePath(filePath, w.workingDir)
	if err != nil {
		return tool.Error(err.Error())
	}

	content, ok := input["content"].(string)
	if !ok || content == "" {
		return tool.Error("content is required and cannot be empty")
	}
	if int64(len(content)) > permission.MaxFileSize {
		return tool.Error(fmt.Sprintf("content exceeds maximum allowed size (%d bytes)", permission.MaxFileSize))
	}

	if _, err := os.Stat(resolvedPath); err == nil {
		if _, err := permission.CheckFileSize(resolvedPath); err != nil {
			return tool.Error(fmt.Sprintf("file too large: %s", filePath))
		}
		if permission.IsBinaryFile(resolvedPath) {
			return tool.Error(fmt.Sprintf("refusing to overwrite binary file: %s", filePath))
		}
	} else if !os.IsNotExist(err) {
		return tool.Error(fmt.Sprintf("failed to stat file: %v", err))
	}

	dir := filepath.Dir(resolvedPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return tool.Error(fmt.Sprintf("failed to create directory: %v", err))
		}
	}

	if err := os.WriteFile(resolvedPath, []byte(content), 0644); err != nil {
		return tool.Error(fmt.Sprintf("failed to write file: %v", err))
	}

	return tool.Success(fmt.Sprintf("Successfully wrote to %s", filePath))
}
