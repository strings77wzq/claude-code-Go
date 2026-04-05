package hooks

import (
	"errors"
	"strings"
	"testing"
)

// mockHook is a mock implementation of the Hook interface for testing.
type mockHook struct {
	name        string
	preExecErr  error
	postExecErr error
}

func (h *mockHook) Name() string {
	return h.name
}

func (h *mockHook) PreExecute(toolName string, input map[string]any) error {
	return h.preExecErr
}

func (h *mockHook) PostExecute(toolName string, input map[string]any, result string, isError bool) error {
	return h.postExecErr
}

// TestHookInterfaceImplementation verifies that mockHook correctly implements the Hook interface.
func TestHookInterfaceImplementation(t *testing.T) {
	hook := &mockHook{
		name: "test-hook",
	}

	// Verify Name() works
	if hook.Name() != "test-hook" {
		t.Errorf("Expected name 'test-hook', got '%s'", hook.Name())
	}

	// Verify PreExecute() works
	input := map[string]any{"key": "value"}
	if err := hook.PreExecute("Read", input); err != nil {
		t.Errorf("PreExecute returned unexpected error: %v", err)
	}

	// Verify PostExecute() works
	if err := hook.PostExecute("Read", input, "result", false); err != nil {
		t.Errorf("PostExecute returned unexpected error: %v", err)
	}

	// Verify error handling in PreExecute
	hook.preExecErr = errors.New("pre-exec error")
	if err := hook.PreExecute("Read", input); err == nil {
		t.Error("PreExecute should have returned an error")
	}

	// Verify error handling in PostExecute
	hook.postExecErr = errors.New("post-exec error")
	if err := hook.PostExecute("Read", input, "result", false); err == nil {
		t.Error("PostExecute should have returned an error")
	}
}

// TestRegistryRegister verifies the Register method of the Registry.
func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	hook1 := &mockHook{name: "hook1"}
	hook2 := &mockHook{name: "hook2"}
	hook3 := &mockHook{name: "hook1"} // duplicate name

	// Register first hook - should succeed
	if err := registry.Register(hook1); err != nil {
		t.Errorf("Register hook1 failed: %v", err)
	}

	// Register second hook - should succeed
	if err := registry.Register(hook2); err != nil {
		t.Errorf("Register hook2 failed: %v", err)
	}

	// Register duplicate hook - should fail
	if err := registry.Register(hook3); err == nil {
		t.Error("Register duplicate hook should have failed")
	} else {
		errStr := err.Error()
		if !strings.Contains(errStr, "hook1") || !strings.Contains(errStr, "already registered") {
			t.Errorf("Expected error message about duplicate hook, got: %v", err)
		}
	}
}

// TestRegistryGetPreHooks verifies the GetPreHooks method.
func TestRegistryGetPreHooks(t *testing.T) {
	registry := NewRegistry()

	// Empty registry should return empty slice
	preHooks := registry.GetPreHooks()
	if len(preHooks) != 0 {
		t.Errorf("Expected 0 pre-hooks, got %d", len(preHooks))
	}

	// Register hooks and verify they are returned
	hook1 := &mockHook{name: "hook1"}
	hook2 := &mockHook{name: "hook2"}

	_ = registry.Register(hook1)
	_ = registry.Register(hook2)

	preHooks = registry.GetPreHooks()
	if len(preHooks) != 2 {
		t.Errorf("Expected 2 pre-hooks, got %d", len(preHooks))
	}

	// Verify the returned slice is a copy (modifying it doesn't affect the registry)
	preHooks[0] = nil
	if registry.GetPreHooks()[0] == nil {
		t.Error("GetPreHooks should return a copy, not the original slice")
	}
}

// TestRegistryGetPostHooks verifies the GetPostHooks method.
func TestRegistryGetPostHooks(t *testing.T) {
	registry := NewRegistry()

	// Empty registry should return empty slice
	postHooks := registry.GetPostHooks()
	if len(postHooks) != 0 {
		t.Errorf("Expected 0 post-hooks, got %d", len(postHooks))
	}

	// Register hooks and verify they are returned
	hook1 := &mockHook{name: "hook1"}
	hook2 := &mockHook{name: "hook2"}

	_ = registry.Register(hook1)
	_ = registry.Register(hook2)

	postHooks = registry.GetPostHooks()
	if len(postHooks) != 2 {
		t.Errorf("Expected 2 post-hooks, got %d", len(postHooks))
	}
}

// TestRunPreHooks verifies the RunPreHooks method with successful hooks.
func TestRunPreHooks(t *testing.T) {
	registry := NewRegistry()

	hook1 := &mockHook{name: "hook1"}
	hook2 := &mockHook{name: "hook2"}

	_ = registry.Register(hook1)
	_ = registry.Register(hook2)

	input := map[string]any{"file": "test.txt"}
	err := registry.RunPreHooks("Read", input)
	if err != nil {
		t.Errorf("RunPreHooks failed: %v", err)
	}
}

// TestRunPreHooksWithFailingHook verifies error handling when a pre-hook fails.
func TestRunPreHooksWithFailingHook(t *testing.T) {
	registry := NewRegistry()

	hook1 := &mockHook{name: "hook1"}
	hook2 := &mockHook{name: "hook2", preExecErr: errors.New("hook2 failed")}

	_ = registry.Register(hook1)
	_ = registry.Register(hook2)

	input := map[string]any{"file": "test.txt"}
	err := registry.RunPreHooks("Read", input)

	if err == nil {
		t.Error("RunPreHooks should have returned an error")
	}

	// Verify the error message contains expected info
	errStr := err.Error()
	if !strings.Contains(errStr, "hook2") || !strings.Contains(errStr, "Read") || !strings.Contains(errStr, "failed") {
		t.Errorf("Expected error message to contain 'hook2', 'Read', and 'failed', got: %s", errStr)
	}
}

// TestRunPostHooks verifies the RunPostHooks method.
func TestRunPostHooks(t *testing.T) {
	registry := NewRegistry()

	hook := &mockHook{name: "hook1"}
	_ = registry.Register(hook)

	input := map[string]any{"file": "test.txt"}
	result := "file content"

	// Should not panic or error even if hook returns error
	registry.RunPostHooks("Read", input, result, false)
}

// TestRunPostHooksWithErrorHook verifies post-hooks continue even if one returns error.
func TestRunPostHooksWithErrorHook(t *testing.T) {
	registry := NewRegistry()

	hook1 := &mockHook{name: "hook1", postExecErr: errors.New("hook1 error")}
	hook2 := &mockHook{name: "hook2"}

	_ = registry.Register(hook1)
	_ = registry.Register(hook2)

	// Should not panic - post-hook errors are logged but don't stop execution
	registry.RunPostHooks("Read", map[string]any{}, "result", false)
}

// TestDuplicateHookError verifies the DuplicateHookError type.
func TestDuplicateHookError(t *testing.T) {
	err := &DuplicateHookError{Name: "test-hook"}
	expectedMsg := "hook already registered: test-hook"
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestPreHookError verifies the PreHookError type.
func TestPreHookError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &PreHookError{
		HookName: "my-hook",
		ToolName: "Read",
		Err:      innerErr,
	}

	// Verify Error() method
	expectedMsg := "pre-hook 'my-hook' failed for tool 'Read': inner error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}

	// Verify Unwrap() method
	if !errors.Is(err, innerErr) {
		t.Error("Unwrap should return the inner error")
	}
}
