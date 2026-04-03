package builtin

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/go-code/internal/tool"
)

type DiffTool struct{}

func NewDiffTool() tool.Tool {
	return &DiffTool{}
}

func (d *DiffTool) Name() string {
	return "Diff"
}

func (d *DiffTool) Description() string {
	return "Compares two content strings and returns unified diff output."
}

func (d *DiffTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Optional file path to include in diff header (for reference only)",
			},
			"old_content": map[string]any{
				"type":        "string",
				"description": "The original content to compare",
			},
			"new_content": map[string]any{
				"type":        "string",
				"description": "The new content to compare against",
			},
		},
		"required": []string{"old_content", "new_content"},
	}
}

func (d *DiffTool) RequiresPermission() bool {
	return false
}

func (d *DiffTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	oldContent, ok := input["old_content"].(string)
	if !ok || oldContent == "" {
		return tool.Error("old_content is required")
	}

	newContent, ok := input["new_content"].(string)
	if !ok || newContent == "" {
		return tool.Error("new_content is required")
	}

	result, err := diffCommand(oldContent, newContent)
	if err != nil {
		result = diffPureGo(oldContent, newContent)
	}

	if result == "" {
		return tool.Success("No differences found.")
	}

	return tool.Success(result)
}

func diffCommand(oldContent, newContent string) (string, error) {
	oldFile, err := os.CreateTemp("", "diff_old_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(oldFile.Name())

	newFile, err := os.CreateTemp("", "diff_new_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(newFile.Name())

	if _, err := oldFile.WriteString(oldContent); err != nil {
		return "", err
	}
	oldFile.Close()

	if _, err := newFile.WriteString(newContent); err != nil {
		return "", err
	}
	newFile.Close()

	cmd := exec.Command("diff", "-u", oldFile.Name(), newFile.Name())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()

	if out.Len() == 0 {
		return "", nil
	}

	return out.String(), nil
}

func diffPureGo(oldContent, newContent string) string {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	var result strings.Builder

	m := len(oldLines)
	n := len(newLines)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	i, j := m, n
	type diffOp struct {
		op      string
		line    string
		lineNum int
	}

	var ops []diffOp

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			ops = append(ops, diffOp{op: " ", line: oldLines[i-1], lineNum: i})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			ops = append(ops, diffOp{op: "+", line: newLines[j-1], lineNum: j})
			j--
		} else if i > 0 {
			ops = append(ops, diffOp{op: "-", line: oldLines[i-1], lineNum: i})
			i--
		}
	}

	for k := len(ops) - 1; k >= 0; k-- {
		op := ops[k]
		switch op.op {
		case "+":
			result.WriteString(fmt.Sprintf("+%s\n", op.line))
		case "-":
			result.WriteString(fmt.Sprintf("-%s\n", op.line))
		case " ":
			result.WriteString(fmt.Sprintf(" %s\n", op.line))
		}
	}

	return result.String()
}
