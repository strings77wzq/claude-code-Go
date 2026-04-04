package init

import (
	"fmt"

	"github.com/strings77wzq/claude-code-Go/internal/tool"
	"github.com/strings77wzq/claude-code-Go/internal/tool/builtin"
)

func RegisterBuiltinTools(r *tool.Registry, workingDir string) error {
	bashTool := builtin.NewBashTool(workingDir)
	if err := r.Register(bashTool); err != nil {
		return fmt.Errorf("failed to register bash tool: %w", err)
	}

	readTool := builtin.NewReadTool(workingDir)
	if err := r.Register(readTool); err != nil {
		return fmt.Errorf("failed to register read tool: %w", err)
	}

	writeTool := builtin.NewWriteTool(workingDir)
	if err := r.Register(writeTool); err != nil {
		return fmt.Errorf("failed to register write tool: %w", err)
	}

	editTool := builtin.NewEditTool(workingDir)
	if err := r.Register(editTool); err != nil {
		return fmt.Errorf("failed to register edit tool: %w", err)
	}

	globTool := builtin.NewGlobTool()
	if err := r.Register(globTool); err != nil {
		return fmt.Errorf("failed to register glob tool: %w", err)
	}

	grepTool := builtin.NewGrepTool()
	if err := r.Register(grepTool); err != nil {
		return fmt.Errorf("failed to register grep tool: %w", err)
	}

	diffTool := builtin.NewDiffTool()
	if err := r.Register(diffTool); err != nil {
		return fmt.Errorf("failed to register diff tool: %w", err)
	}

	treeTool := builtin.NewTreeTool()
	if err := r.Register(treeTool); err != nil {
		return fmt.Errorf("failed to register tree tool: %w", err)
	}

	webFetchTool := builtin.NewWebFetchTool()
	if err := r.Register(webFetchTool); err != nil {
		return fmt.Errorf("failed to register webfetch tool: %w", err)
	}

	todoTool := builtin.NewTodoTool()
	if err := r.Register(todoTool); err != nil {
		return fmt.Errorf("failed to register todo tool: %w", err)
	}

	return nil
}
