package permission

import (
	"testing"
)

func TestNonInteractivePrompterAlwaysDenies(t *testing.T) {
	p := NewNonInteractivePrompter()

	tests := []struct {
		name     string
		toolName string
	}{
		{"bash prompt denied", "Bash"},
		{"write prompt denied", "Write"},
		{"edit prompt denied", "Edit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := p.Decide(tt.toolName, map[string]any{"command": "test"}, "requires_approval")

			if decision != Deny {
				t.Errorf("expected Deny, got %s", decision)
			}
		})
	}
}

func TestNonInteractivePrompterDenialMessage(t *testing.T) {
	p := NewNonInteractivePrompter()

	decision := p.Decide("Bash", map[string]any{"command": "rm -rf /"}, "requires_approval")

	if decision != Deny {
		t.Fatalf("expected Deny, got %s", decision)
	}
}
