// Package eval provides a framework for evaluating agent quality across models.
package eval

import (
	"context"
	"fmt"
	"time"
)

// Case defines a single evaluation test case.
type Case struct {
	Name        string            // Human-readable name
	Prompt      string            // Input to the agent
	ExpectTools []string          // Expected tool names used (empty = don't check)
	ExpectText  string            // Substring expected in output (empty = don't check)
	MaxTokens   int               // Token budget (0 = default)
	Tags        map[string]string // Metadata for filtering
}

// Result holds the outcome of running a single eval case.
type Result struct {
	Case      Case
	Passed    bool
	Output    string
	ToolsUsed []string
	InputTok  int
	OutputTok int
	Duration  time.Duration
	Error     string
}

// Report summarizes eval results for a model.
type Report struct {
	Model    string
	Results  []Result
	Passed   int
	Failed   int
	TotalTok int
	Duration time.Duration
}

// Runner executes eval cases against an agent.
type Runner struct {
	executeFunc ExecuteFunc
}

// ExecuteFunc is the function signature for running an agent on a prompt.
type ExecuteFunc func(ctx context.Context, prompt string) (output string, toolsUsed []string, inputTok, outputTok int, err error)

// NewRunner creates a new eval runner with the given execute function.
func NewRunner(execute ExecuteFunc) *Runner {
	return &Runner{executeFunc: execute}
}

// Run executes all eval cases and returns a report.
func (r *Runner) Run(ctx context.Context, model string, cases []Case) Report {
	report := Report{Model: model}
	start := time.Now()

	for _, c := range cases {
		result := r.runCase(ctx, c)
		report.Results = append(report.Results, result)
		if result.Passed {
			report.Passed++
		} else {
			report.Failed++
		}
		report.TotalTok += result.InputTok + result.OutputTok
	}

	report.Duration = time.Since(start)
	return report
}

func (r *Runner) runCase(ctx context.Context, c Case) Result {
	start := time.Now()
	output, toolsUsed, inTok, outTok, err := r.executeFunc(ctx, c.Prompt)
	dur := time.Since(start)

	result := Result{
		Case:      c,
		Output:    output,
		ToolsUsed: toolsUsed,
		InputTok:  inTok,
		OutputTok: outTok,
		Duration:  dur,
	}

	if err != nil {
		result.Error = err.Error()
		result.Passed = false
		return result
	}

	// Check expected tools
	if len(c.ExpectTools) > 0 {
		if !containsAll(toolsUsed, c.ExpectTools) {
			result.Passed = false
			result.Error = fmt.Sprintf("expected tools %v, got %v", c.ExpectTools, toolsUsed)
			return result
		}
	}

	// Check expected text
	if c.ExpectText != "" {
		if !containsSubstring(output, c.ExpectText) {
			result.Passed = false
			result.Error = fmt.Sprintf("output missing expected text: %q", c.ExpectText)
			return result
		}
	}

	result.Passed = true
	return result
}

func containsAll(haystack, needles []string) bool {
	set := make(map[string]bool, len(haystack))
	for _, s := range haystack {
		set[s] = true
	}
	for _, n := range needles {
		if !set[n] {
			return false
		}
	}
	return true
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsAt(s, sub))
}

func containsAt(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
