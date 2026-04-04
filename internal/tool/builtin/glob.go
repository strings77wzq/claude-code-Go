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

type GlobTool struct{}

func NewGlobTool() tool.Tool {
	return &GlobTool{}
}

func (g *GlobTool) Name() string {
	return "Glob"
}

func (g *GlobTool) Description() string {
	return "Finds files matching a pattern, supports ** for recursive matching."
}

func (g *GlobTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern": map[string]any{
				"type":        "string",
				"description": "The glob pattern to match (e.g., **/*.go)",
			},
		},
		"required": []string{"pattern"},
	}
}

func (g *GlobTool) RequiresPermission() bool {
	return false
}

func (g *GlobTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}

func (g *GlobTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	pattern, ok := input["pattern"].(string)
	if !ok || pattern == "" {
		return tool.Error("pattern is required")
	}

	if strings.Contains(pattern, "**") {
		return g.executeRecursive(pattern)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return tool.Error(fmt.Sprintf("invalid pattern: %v", err))
	}

	if len(matches) == 0 {
		return tool.Success("No matches found")
	}

	return tool.Success(strings.Join(matches, "\n"))
}

func (g *GlobTool) executeRecursive(pattern string) tool.Result {
	parts := strings.SplitN(pattern, "**", 2)
	var baseDir string
	var suffixPattern string

	if parts[0] == "" {
		baseDir = "."
		suffixPattern = parts[1]
	} else {
		baseDir = parts[0]
		suffixPattern = parts[1]
	}

	suffixPattern = strings.TrimPrefix(suffixPattern, "/")

	var matches []string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return nil
		}

		if suffixPattern == "" {
			matches = append(matches, path)
			return nil
		}

		if match, err := filepath.Match(suffixPattern, relPath); err == nil && match {
			matches = append(matches, path)
		} else if match, err := filepath.Match(suffixPattern, filepath.Base(path)); err == nil && match {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return tool.Error(fmt.Sprintf("error walking directory: %v", err))
	}

	if len(matches) == 0 {
		return tool.Success("No matches found")
	}

	return tool.Success(strings.Join(matches, "\n"))
}
