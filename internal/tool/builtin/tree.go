package builtin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

const (
	defaultTreeDepth = 3
)

type TreeTool struct{}

func NewTreeTool() tool.Tool {
	return &TreeTool{}
}

func (t *TreeTool) Name() string {
	return "Tree"
}

func (t *TreeTool) Description() string {
	return "Displays directory tree structure as text."
}

func (t *TreeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Root directory path (default: current directory)",
			},
			"max_depth": map[string]any{
				"type":        "number",
				"description": "Maximum depth to traverse (default: 3)",
			},
		},
		"required": []string{},
	}
}

func (t *TreeTool) RequiresPermission() bool {
	return false
}

func (t *TreeTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}

func (t *TreeTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	path := "."
	if p, ok := input["path"].(string); ok && p != "" {
		path = p
	}

	maxDepth := defaultTreeDepth
	if d, ok := input["max_depth"]; ok {
		switch v := d.(type) {
		case float64:
			maxDepth = int(v)
		case int64:
			maxDepth = int(v)
		case int:
			maxDepth = v
		}
	}

	if maxDepth < 1 {
		maxDepth = 1
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tool.Error(fmt.Sprintf("path does not exist: %s", path))
		}
		return tool.Error(fmt.Sprintf("failed to stat path: %v", err))
	}

	if !info.IsDir() {
		return tool.Error(fmt.Sprintf("not a directory: %s", path))
	}

	result, err := walkTree(path, 0, maxDepth)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to walk directory: %v", err))
	}

	if result == "" {
		return tool.Success(".")
	}

	return tool.Success(result)
}

func walkTree(root string, currentDepth, maxDepth int) (string, error) {
	var result strings.Builder

	entries, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1

		prefix := getPrefixes(currentDepth, isLast)
		result.WriteString(prefix)
		result.WriteString(entry.Name())
		result.WriteString("\n")

		if entry.IsDir() && currentDepth < maxDepth-1 {
			childPath := filepath.Join(root, entry.Name())
			childResult, err := walkTree(childPath, currentDepth+1, maxDepth)
			if err != nil {
				continue
			}
			result.WriteString(childResult)
		}
	}

	return result.String(), nil
}

func getPrefixes(depth int, isLast bool) string {
	if depth == 0 {
		return ""
	}

	var prefix strings.Builder

	for i := 0; i < depth-1; i++ {
		prefix.WriteString("│   ")
	}

	if isLast {
		prefix.WriteString("└── ")
	} else {
		prefix.WriteString("├── ")
	}

	return prefix.String()
}
