package builtin

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

const maxGrepOutput = 100 * 1024

type GrepTool struct{}

func NewGrepTool() tool.Tool {
	return &GrepTool{}
}

func (g *GrepTool) Name() string {
	return "Grep"
}

func (g *GrepTool) Description() string {
	return "Searches for a pattern in files and returns matching lines."
}

func (g *GrepTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern": map[string]any{
				"type":        "string",
				"description": "The regex pattern to search for",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Directory to search in (default: .)",
			},
			"include": map[string]any{
				"type":        "string",
				"description": "File pattern to include (e.g., *.go)",
			},
			"exclude": map[string]any{
				"type":        "string",
				"description": "File pattern to exclude",
			},
		},
		"required": []string{"pattern"},
	}
}

func (g *GrepTool) RequiresPermission() bool {
	return false
}

func (g *GrepTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	pattern, ok := input["pattern"].(string)
	if !ok || pattern == "" {
		return tool.Error("pattern is required")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return tool.Error(fmt.Sprintf("invalid regex pattern: %v", err))
	}

	path := "."
	if p, ok := input["path"].(string); ok && p != "" {
		path = p
	}

	var include, exclude string
	if inc, ok := input["include"].(string); ok {
		include = inc
	}
	if exc, ok := input["exclude"].(string); ok {
		exclude = exc
	}

	var matches []string
	var totalSize int

	err = filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if exclude != "" {
			if matched, _ := filepath.Match(exclude, info.Name()); matched {
				return nil
			}
		}

		if include != "" {
			if matched, _ := filepath.Match(include, info.Name()); !matched {
				return nil
			}
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			if re.MatchString(line) {
				match := fmt.Sprintf("%s:%d:%s", filePath, lineNum+1, line)
				matches = append(matches, match)
				totalSize += len(match) + 1
				if totalSize > maxGrepOutput {
					return fs.SkipDir
				}
			}
		}

		return nil
	})

	if err != nil && err != fs.SkipDir {
		return tool.Error(fmt.Sprintf("error searching files: %v", err))
	}

	if len(matches) == 0 {
		return tool.Success("No matches found")
	}

	if totalSize > maxGrepOutput {
		truncated := strings.Join(matches[:len(matches)/2], "\n")
		return tool.Success(truncated + "\n... (output truncated)")
	}

	return tool.Success(strings.Join(matches, "\n"))
}
