package builtin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

type EditTool struct {
	workingDir string
}

func NewEditTool(workingDir string) tool.Tool {
	return &EditTool{
		workingDir: workingDir,
	}
}

func (e *EditTool) Name() string {
	return "Edit"
}

func (e *EditTool) Description() string {
	return "Performs exact string replacement in a file."
}

func (e *EditTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Path to the file to edit",
			},
			"old_string": map[string]any{
				"type":        "string",
				"description": "The exact string to replace",
			},
			"new_string": map[string]any{
				"type":        "string",
				"description": "The replacement string",
			},
			"replace_all": map[string]any{
				"type":        "boolean",
				"description": "Replace all occurrences (default: false)",
			},
		},
		"required": []string{"file_path", "old_string", "new_string"},
	}
}

func (e *EditTool) RequiresPermission() bool {
	return true
}

func (e *EditTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelWorkspaceWrite
}

func (e *EditTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return tool.Error("file_path is required")
	}

	if err := ValidatePath(filePath, e.workingDir); err != nil {
		return tool.Error(err.Error())
	}

	oldString, ok := input["old_string"].(string)
	if !ok || oldString == "" {
		return tool.Error("old_string is required")
	}

	newString, ok := input["new_string"].(string)
	if !ok {
		return tool.Error("new_string is required")
	}

	replaceAll := false
	if ra, ok := input["replace_all"]; ok {
		if b, ok := ra.(bool); ok {
			replaceAll = b
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tool.Error(fmt.Sprintf("file not found: %s", filePath))
		}
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	fileContent := string(content)

	count := strings.Count(fileContent, oldString)
	if count == 0 {
		return tool.Error(fmt.Sprintf("old_string not found in file: %s", filePath))
	}

	if !replaceAll && count > 1 {
		return tool.Error(fmt.Sprintf("old_string appears %d times (use replace_all=true to replace all)", count))
	}

	newContent := strings.Replace(fileContent, oldString, newString, -1)

	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return tool.Error(fmt.Sprintf("failed to write file: %v", err))
	}

	return tool.Success(fmt.Sprintf("Successfully edited %s", filePath))
}
