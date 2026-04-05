package cost

import (
	"testing"
)

func TestNewCostTrackerWithKnownModel(t *testing.T) {
	tracker := NewCostTracker("claude-sonnet-4")
	if tracker == nil {
		t.Fatal("NewCostTracker should not return nil")
	}

	if tracker.model != "claude-sonnet-4" {
		t.Errorf("Expected model 'claude-sonnet-4', got '%s'", tracker.model)
	}
}

func TestNewCostTrackerWithUnknownModel(t *testing.T) {
	tracker := NewCostTracker("unknown-model")
	if tracker == nil {
		t.Fatal("NewCostTracker should not return nil")
	}

	if tracker.model != "unknown-model" {
		t.Errorf("Expected model 'unknown-model', got '%s'", tracker.model)
	}

	if tracker.pricing.Name != "default" {
		t.Errorf("Expected default pricing, got '%s'", tracker.pricing.Name)
	}
}

func TestRecordUsage(t *testing.T) {
	tracker := NewCostTracker("claude-sonnet-4")

	tracker.RecordUsage(1000, 500)

	input, output, cost := tracker.GetTotalCost()
	if input != 1000 {
		t.Errorf("Expected 1000 input tokens, got %d", input)
	}
	if output != 500 {
		t.Errorf("Expected 500 output tokens, got %d", output)
	}
	if cost <= 0 {
		t.Error("Cost should be positive")
	}
}

func TestRecordUsageAccumulates(t *testing.T) {
	tracker := NewCostTracker("claude-sonnet-4")

	tracker.RecordUsage(1000, 500)
	tracker.RecordUsage(2000, 1000)

	input, output, _ := tracker.GetTotalCost()
	if input != 3000 {
		t.Errorf("Expected 3000 input tokens, got %d", input)
	}
	if output != 1500 {
		t.Errorf("Expected 1500 output tokens, got %d", output)
	}
}

func TestGetTotalCost(t *testing.T) {
	tracker := NewCostTracker("claude-opus-4")

	tracker.RecordUsage(1000000, 500000)

	input, output, cost := tracker.GetTotalCost()
	if input != 1000000 {
		t.Errorf("Expected 1000000 input tokens, got %d", input)
	}
	if output != 500000 {
		t.Errorf("Expected 500000 output tokens, got %d", output)
	}

	expectedCost := (1000000.0/1_000_000)*15.0 + (500000.0/1_000_000)*75.0
	if cost != expectedCost {
		t.Errorf("Expected cost %.4f, got %.4f", expectedCost, cost)
	}
}

func TestGetSummary(t *testing.T) {
	tracker := NewCostTracker("claude-haiku-4")

	tracker.RecordUsage(1000, 500)

	summary := tracker.GetSummary()
	if summary == "" {
		t.Error("GetSummary should not return empty string")
	}
}

func TestDefaultPricing(t *testing.T) {
	tests := []struct {
		model   string
		wantIn  float64
		wantOut float64
	}{
		{"claude-sonnet-4", 3.0, 15.0},
		{"claude-opus-4", 15.0, 75.0},
		{"claude-haiku-4", 0.25, 1.25},
		{"gpt-4o", 2.5, 10.0},
		{"unknown", 3.0, 15.0},
	}

	for _, tt := range tests {
		tracker := NewCostTracker(tt.model)
		if tracker.pricing.InputPricePerMTok != tt.wantIn {
			t.Errorf("Model %s: expected input price %.2f, got %.2f", tt.model, tt.wantIn, tracker.pricing.InputPricePerMTok)
		}
		if tracker.pricing.OutputPricePerMTok != tt.wantOut {
			t.Errorf("Model %s: expected output price %.2f, got %.2f", tt.model, tt.wantOut, tracker.pricing.OutputPricePerMTok)
		}
	}
}

func TestZeroTokens(t *testing.T) {
	tracker := NewCostTracker("claude-sonnet-4")

	tracker.RecordUsage(0, 0)

	input, output, _ := tracker.GetTotalCost()
	if input != 0 {
		t.Errorf("Expected 0 input tokens, got %d", input)
	}
	if output != 0 {
		t.Errorf("Expected 0 output tokens, got %d", output)
	}
}

func TestLargeTokenCount(t *testing.T) {
	tracker := NewCostTracker("claude-opus-4")

	tracker.RecordUsage(10_000_000, 5_000_000)

	input, output, cost := tracker.GetTotalCost()
	if input != 10_000_000 {
		t.Errorf("Expected 10000000 input tokens, got %d", input)
	}
	if output != 5_000_000 {
		t.Errorf("Expected 5000000 output tokens, got %d", output)
	}
	if cost <= 0 {
		t.Error("Cost should be positive for large token count")
	}
}
