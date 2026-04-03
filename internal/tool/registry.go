package tool

import (
	"context"
	"fmt"
	"sync"
)

type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool already registered: %s", name)
	}

	r.tools[name] = tool
	return nil
}

func (r *Registry) GetTool(name string) Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.tools[name]
}

func (r *Registry) Execute(ctx context.Context, name string, input map[string]any) Result {
	tool := r.GetTool(name)
	if tool == nil {
		return Error(fmt.Sprintf("tool not found: %s", name))
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			Error(fmt.Sprintf("panic recovered: %v", recovered))
		}
	}()

	return tool.Execute(ctx, input)
}

func (r *Registry) GetAllDefinitions() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	definitions := make([]ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		definitions = append(definitions, ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return definitions
}

func RegisterBuiltinTools(r *Registry, workingDir string) error {
	// TODO: Implement builtin tools
	return nil
}
