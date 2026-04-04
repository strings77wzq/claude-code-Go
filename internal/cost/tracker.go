package cost

import "fmt"

type ModelPricing struct {
	Name               string
	InputPricePerMTok  float64
	OutputPricePerMTok float64
}

var DefaultPricing = map[string]ModelPricing{
	"claude-sonnet-4": {
		Name:               "claude-sonnet-4",
		InputPricePerMTok:  3.0,
		OutputPricePerMTok: 15.0,
	},
	"claude-opus-4": {
		Name:               "claude-opus-4",
		InputPricePerMTok:  15.0,
		OutputPricePerMTok: 75.0,
	},
	"claude-haiku-4": {
		Name:               "claude-haiku-4",
		InputPricePerMTok:  0.25,
		OutputPricePerMTok: 1.25,
	},
	"gpt-4o": {
		Name:               "gpt-4o",
		InputPricePerMTok:  2.5,
		OutputPricePerMTok: 10.0,
	},
	"default": {
		Name:               "default",
		InputPricePerMTok:  3.0,
		OutputPricePerMTok: 15.0,
	},
}

type CostTracker struct {
	model             string
	totalInputTokens  float64
	totalOutputTokens float64
	totalCostUSD      float64
	pricing           ModelPricing
}

func NewCostTracker(model string) *CostTracker {
	pricing, ok := DefaultPricing[model]
	if !ok {
		pricing = DefaultPricing["default"]
	}
	return &CostTracker{
		model:   model,
		pricing: pricing,
	}
}

func (ct *CostTracker) RecordUsage(inputTokens, outputTokens int) {
	inputCost := (float64(inputTokens) / 1_000_000) * ct.pricing.InputPricePerMTok
	outputCost := (float64(outputTokens) / 1_000_000) * ct.pricing.OutputPricePerMTok

	ct.totalInputTokens += float64(inputTokens)
	ct.totalOutputTokens += float64(outputTokens)
	ct.totalCostUSD += inputCost + outputCost
}

func (ct *CostTracker) GetTotalCost() (inputTokens, outputTokens int, costUSD float64) {
	return int(ct.totalInputTokens), int(ct.totalOutputTokens), ct.totalCostUSD
}

func (ct *CostTracker) GetSummary() string {
	return fmt.Sprintf(
		"Model: %s | Input: %.0f tokens | Output: %.0f tokens | Total Cost: $%.4f",
		ct.model,
		ct.totalInputTokens,
		ct.totalOutputTokens,
		ct.totalCostUSD,
	)
}
