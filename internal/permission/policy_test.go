package permission

import (
	"strings"
	"testing"
)

func TestDenyRulesOverrideAllowRules(t *testing.T) {
	policy := NewPolicy(WorkspaceWrite)

	denyRule, _ := ParseRule("bash(rm -rf:*)")
	policy.AddDenyRule(denyRule)

	allowRule, _ := ParseRule("bash(*)")
	policy.AddAllowRule(allowRule)

	input := map[string]any{"command": "rm -rf /"}
	decision := policy.Evaluate("Bash", input, true)

	if decision != Deny {
		t.Errorf("Expected Deny, got %v", decision)
	}
}

func TestSessionMemoryAlwaysChoice(t *testing.T) {
	policy := NewPolicy(WorkspaceWrite)

	policy.SetSessionMemory("bash:git commit", Allow)

	input := map[string]any{"command": "git commit"}
	decision := policy.Evaluate("Bash", input, true)

	if decision != Allow {
		t.Errorf("Expected Allow from session memory, got %v", decision)
	}
}

func TestReadOnlyMode(t *testing.T) {
	policy := NewPolicy(ReadOnly)

	inputRead := map[string]any{"file_path": "/home/file.txt"}
	decisionRead := policy.Evaluate("Read", inputRead, false)

	if decisionRead != Allow {
		t.Errorf("Expected Allow for Read tool in ReadOnly mode, got %v", decisionRead)
	}

	inputBash := map[string]any{"command": "ls"}
	decisionBash := policy.Evaluate("Bash", inputBash, true)

	if decisionBash != Ask {
		t.Errorf("Expected Ask for Bash tool in ReadOnly mode, got %v", decisionBash)
	}
}

func TestWorkspaceWriteMode(t *testing.T) {
	policy := NewPolicy(WorkspaceWrite)

	inputRead := map[string]any{"file_path": "/home/file.txt"}
	decisionRead := policy.Evaluate("Read", inputRead, false)

	if decisionRead != Allow {
		t.Errorf("Expected Allow for Read tool in WorkspaceWrite mode, got %v", decisionRead)
	}

	inputWrite := map[string]any{"file_path": "/home/file.txt", "content": "test"}
	decisionWrite := policy.Evaluate("Write", inputWrite, true)

	if decisionWrite != Ask {
		t.Errorf("Expected Ask for Write tool in WorkspaceWrite mode (requires explicit permission), got %v", decisionWrite)
	}
}

func TestDangerFullAccessMode(t *testing.T) {
	policy := NewPolicy(DangerFullAccess)

	inputBash := map[string]any{"command": "rm -rf /"}
	decision := policy.Evaluate("Bash", inputBash, true)

	if decision != Allow {
		t.Errorf("Expected Allow for any tool in DangerFullAccess mode, got %v", decision)
	}
}

func TestToolRequirementOverridesMode(t *testing.T) {
	policy := NewPolicy(DangerFullAccess)

	policy.SetToolRequirement("Bash", ReadOnly)

	inputBash := map[string]any{"command": "ls"}
	decision := policy.Evaluate("Bash", inputBash, true)

	if decision != Deny {
		t.Errorf("Expected Deny when tool requirement is ReadOnly but active mode is DangerFullAccess, got %v", decision)
	}
}

func TestRequiresPermissionFalseAllowsRead(t *testing.T) {
	policy := NewPolicy(ReadOnly)

	inputGlob := map[string]any{"pattern": "*.go"}
	decision := policy.Evaluate("Glob", inputGlob, false)

	if decision != Allow {
		t.Errorf("Expected Allow for Glob with requiresPermission=false, got %v", decision)
	}
}

func TestNoMatchingRulesReturnsAsk(t *testing.T) {
	policy := NewPolicy(ReadOnly)

	input := map[string]any{"command": "unknown-command"}
	decision := policy.Evaluate("Bash", input, true)

	if decision != Ask {
		t.Errorf("Expected Ask when no rules match and requiresPermission=true, got %v", decision)
	}
}

func TestPolicyGetActiveMode(t *testing.T) {
	policy := NewPolicy(WorkspaceWrite)

	if policy.GetActiveMode() != WorkspaceWrite {
		t.Errorf("Expected WorkspaceWrite, got %v", policy.GetActiveMode())
	}
}

func TestPolicySetAndGetSessionMemory(t *testing.T) {
	policy := NewPolicy(ReadOnly)

	policy.SetSessionMemory("test-key", Allow)

	decision, exists := policy.GetSessionMemory("test-key")
	if !exists {
		t.Error("Expected key to exist in session memory")
	}
	if decision != Allow {
		t.Errorf("Expected Allow, got %v", decision)
	}
}

type mockReader struct {
	inputs   []string
	idx      int
	callback func() string
}

func (r *mockReader) ReadString(delim byte) (string, error) {
	if r.callback != nil {
		return r.callback(), nil
	}
	if r.idx >= len(r.inputs) {
		return "n\n", nil
	}
	input := r.inputs[r.idx]
	r.idx++
	return input + "\n", nil
}

type mockWriter struct {
	output strings.Builder
}

func (w *mockWriter) Write(p []byte) (int, error) {
	w.output.Write(p)
	return len(p), nil
}

func TestTerminalPrompterAcceptsYes(t *testing.T) {
	reader := &mockReader{inputs: []string{"y"}}
	writer := &mockWriter{}
	prompter := NewTerminalPrompter(reader, writer)

	input := map[string]any{"command": "ls"}
	decision := prompter.Decide("Bash", input, "test")

	if decision != Allow {
		t.Errorf("Expected Allow, got %v", decision)
	}
}

func TestTerminalPrompterAcceptsNo(t *testing.T) {
	reader := &mockReader{inputs: []string{"n"}}
	writer := &mockWriter{}
	prompter := NewTerminalPrompter(reader, writer)

	input := map[string]any{"command": "ls"}
	decision := prompter.Decide("Bash", input, "test")

	if decision != Deny {
		t.Errorf("Expected Deny, got %v", decision)
	}
}

func TestTerminalPrompterAcceptsAlways(t *testing.T) {
	reader := &mockReader{inputs: []string{"a"}}
	writer := &mockWriter{}
	prompter := NewTerminalPrompter(reader, writer)

	input := map[string]any{"command": "ls"}
	decision := prompter.Decide("Bash", input, "test")

	if decision != Allow {
		t.Errorf("Expected Allow for always, got %v", decision)
	}
}

func TestTerminalPrompterRejectsInvalidInput(t *testing.T) {
	reader := &mockReader{inputs: []string{"invalid", "y"}}
	writer := &mockWriter{}
	prompter := NewTerminalPrompter(reader, writer)

	input := map[string]any{"command": "ls"}
	decision := prompter.Decide("Bash", input, "test")

	if decision != Allow {
		t.Errorf("Expected Allow after invalid input, got %v", decision)
	}
}
