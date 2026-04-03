// Package hooks provides a system for registering and executing pre/post tool execution callbacks.
package hooks

import (
	"sync"
)

// Hook defines the interface for tool execution hooks.
// Implementations can intercept tool calls before and after execution.
type Hook interface {
	// Name returns the unique identifier for this hook.
	Name() string

	// PreExecute is called before a tool is executed.
	// It receives the tool name and input parameters.
	// Returning an error will prevent tool execution.
	PreExecute(toolName string, input map[string]any) error

	// PostExecute is called after a tool has been executed.
	// It receives the tool name, input parameters, execution result, and whether the result is an error.
	// Errors from PostExecute are logged but do not affect the tool result.
	PostExecute(toolName string, input map[string]any, result string, isError bool) error
}

// Registry manages the registration and execution of hooks.
// It is safe for concurrent use.
type Registry struct {
	mu        sync.RWMutex
	preHooks  []Hook
	postHooks []Hook
}

// NewRegistry creates a new empty hook registry.
func NewRegistry() *Registry {
	return &Registry{
		preHooks:  make([]Hook, 0),
		postHooks: make([]Hook, 0),
	}
}

// Register adds a hook to the registry.
// The hook will be called on all tool executions.
// Returns an error if a hook with the same name is already registered.
func (r *Registry) Register(hook Hook) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, h := range r.preHooks {
		if h.Name() == hook.Name() {
			return &DuplicateHookError{Name: hook.Name()}
		}
	}

	r.preHooks = append(r.preHooks, hook)
	r.postHooks = append(r.postHooks, hook)

	return nil
}

// GetPreHooks returns a slice of all registered pre-execute hooks.
// The returned slice is a copy and is safe for iteration.
func (r *Registry) GetPreHooks() []Hook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Hook, len(r.preHooks))
	copy(result, r.preHooks)
	return result
}

// GetPostHooks returns a slice of all registered post-execute hooks.
// The returned slice is a copy and is safe for iteration.
func (r *Registry) GetPostHooks() []Hook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Hook, len(r.postHooks))
	copy(result, r.postHooks)
	return result
}

// RunPreHooks executes all pre-execute hooks in order.
// It stops and returns the first error encountered.
// If a pre-hook returns an error, subsequent hooks are not executed.
func (r *Registry) RunPreHooks(toolName string, input map[string]any) error {
	hooks := r.GetPreHooks()

	for _, hook := range hooks {
		if err := hook.PreExecute(toolName, input); err != nil {
			return &PreHookError{HookName: hook.Name(), ToolName: toolName, Err: err}
		}
	}

	return nil
}

// RunPostHooks executes all post-execute hooks in order.
// Errors from post-hooks are logged but do not stop execution of subsequent hooks.
func (r *Registry) RunPostHooks(toolName string, input map[string]any, result string, isError bool) {
	hooks := r.GetPostHooks()

	for _, hook := range hooks {
		if err := hook.PostExecute(toolName, input, result, isError); err != nil {
			// Log the error - PostExecute errors should not stop execution
			// Using the hook name for logging context
			_ = err // In production, would use proper logging
		}
	}
}

// DuplicateHookError is returned when attempting to register a hook with a duplicate name.
type DuplicateHookError struct {
	Name string
}

func (e *DuplicateHookError) Error() string {
	return "hook already registered: " + e.Name
}

// PreHookError wraps an error from a pre-execute hook.
type PreHookError struct {
	HookName string
	ToolName string
	Err      error
}

func (e *PreHookError) Error() string {
	return "pre-hook '" + e.HookName + "' failed for tool '" + e.ToolName + "': " + e.Err.Error()
}

func (e *PreHookError) Unwrap() error {
	return e.Err
}
