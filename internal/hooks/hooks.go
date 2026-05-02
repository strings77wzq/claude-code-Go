package hooks

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/strings77wzq/claude-code-Go/internal/diagnostic"
)

type Hook interface {
	Name() string
	PreExecute(toolName string, input map[string]any) error
	PostExecute(toolName string, input map[string]any, result string, isError bool) error
}

type Registry struct {
	mu        sync.RWMutex
	preHooks  []Hook
	postHooks []Hook
	policies  map[string]HookPolicy
}

type HookFailureMode string

const (
	HookFailureWarn  HookFailureMode = "warn"
	HookFailureBlock HookFailureMode = "block"
)

type HookPolicy struct {
	PreFailure HookFailureMode
}

func NewRegistry() *Registry {
	return &Registry{
		preHooks:  make([]Hook, 0),
		postHooks: make([]Hook, 0),
		policies:  make(map[string]HookPolicy),
	}
}

func (r *Registry) Register(hook Hook) error {
	return r.RegisterWithPolicy(hook, HookPolicy{PreFailure: HookFailureBlock})
}

func (r *Registry) RegisterWithPolicy(hook Hook, policy HookPolicy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, h := range r.preHooks {
		if h.Name() == hook.Name() {
			return &DuplicateHookError{Name: hook.Name()}
		}
	}

	r.preHooks = append(r.preHooks, hook)
	r.postHooks = append(r.postHooks, hook)
	r.policies[hook.Name()] = normalizeHookPolicy(policy)

	return nil
}

func (r *Registry) GetPreHooks() []Hook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Hook, len(r.preHooks))
	copy(result, r.preHooks)
	return result
}

func (r *Registry) GetPostHooks() []Hook {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Hook, len(r.postHooks))
	copy(result, r.postHooks)
	return result
}

func (r *Registry) RunPreHooks(toolName string, input map[string]any) error {
	hooks := r.GetPreHooks()

	for _, hook := range hooks {
		if err := hook.PreExecute(toolName, input); err != nil {
			if r.policyFor(hook.Name()).PreFailure == HookFailureBlock {
				return &PreHookError{HookName: hook.Name(), ToolName: toolName, Err: err}
			}
		}
	}

	return nil
}

func (r *Registry) RunPostHooks(toolName string, input map[string]any, result string, isError bool) {
	hooks := r.GetPostHooks()

	for _, hook := range hooks {
		if err := hook.PostExecute(toolName, input, result, isError); err != nil {
			slog.Warn("post-hook error", "hook", hook.Name(), "tool", toolName, "error", err)
		}
	}
}

type DuplicateHookError struct {
	Name string
}

func (e *DuplicateHookError) Error() string {
	return "hook already registered: " + e.Name
}

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

func HookErrorDiagnostic(err *PreHookError) diagnostic.Diagnostic {
	if err == nil {
		return diagnostic.Diagnostic{}
	}
	return diagnostic.Diagnostic{
		Component: "hooks",
		Severity:  diagnostic.SeverityError,
		Code:      "hooks.pre_failed",
		Summary:   "Pre-tool hook failed",
		Detail:    err.Err.Error(),
		Metadata: map[string]any{
			"hook": err.HookName,
			"tool": err.ToolName,
		},
	}
}

func (r *Registry) policyFor(name string) HookPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return normalizeHookPolicy(r.policies[name])
}

func normalizeHookPolicy(policy HookPolicy) HookPolicy {
	if policy.PreFailure == "" {
		policy.PreFailure = HookFailureBlock
	}
	return policy
}

// ShellHook implements Hook by running a shell command.
type ShellHook struct {
	name    string
	command string
}

func (h *ShellHook) Name() string { return h.name }

func (h *ShellHook) PreExecute(toolName string, input map[string]any) error {
	return h.run(toolName, input, "")
}

func (h *ShellHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
	return h.run(toolName, input, result)
}

func (h *ShellHook) run(toolName string, input map[string]any, result string) error {
	inputJSON, _ := json.Marshal(map[string]any{
		"tool":   toolName,
		"input":  input,
		"result": result,
	})
	cmd := exec.Command("sh", "-c", h.command)
	cmd.Stdin = strings.NewReader(string(inputJSON))
	return cmd.Run()
}

// hookFileDef is the JSON schema for a hook file in ~/.go-code/hooks/.
type hookFileDef struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Command string `json:"command"`
}

// LoadHooksFromDir loads hook definitions from JSON files in a directory.
// Only .json files are read. Invalid files are skipped.
func LoadHooksFromDir(dir string) ([]Hook, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read hooks dir %s: %w", dir, err)
	}

	var hooks []Hook
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("failed to read hook file", "path", path, "error", err)
			continue
		}
		var def hookFileDef
		if err := json.Unmarshal(data, &def); err != nil {
			slog.Warn("invalid hook JSON", "path", path, "error", err)
			continue
		}
		if def.Name == "" || def.Command == "" {
			slog.Warn("hook missing required fields", "path", path)
			continue
		}
		hooks = append(hooks, &ShellHook{name: def.Name, command: def.Command})
	}
	return hooks, nil
}
