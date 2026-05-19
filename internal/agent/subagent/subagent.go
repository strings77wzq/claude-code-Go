package subagent

import (
	"context"
	"fmt"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/agent"
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// Type identifies a sub-agent specialization.
type Type string

const (
	TypeExplore Type = "Explore"
	TypeReview  Type = "Review"
	TypeTestGen Type = "TestGen"
)

func (t Type) String() string { return string(t) }

// Max turns per sub-agent type (smaller than main agent's 50).
const (
	maxExploreTurns = 15
	maxReviewTurns  = 10
	maxTestGenTurns = 20
)

// API client interface (matching agent's interface).
type apiClient interface {
	SendMessageStream(ctx context.Context, req *api.ApiRequest, onTextDelta func(text string)) (*api.ApiResponse, error)
}

// Tool registry interface.
type toolRegistry interface {
	GetTool(name string) tool.Tool
	Execute(ctx context.Context, name string, input map[string]any) tool.Result
	GetAllDefinitions() []tool.ToolDefinition
	GetDefinitionsByTier(tiers ...tool.ToolTier) []tool.ToolDefinition
}

// Permission policy interface.
type permissionPolicy interface {
	Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision
}

// ExploreToolSet returns tools for code exploration: read-only file tools.
func ExploreToolSet() map[string]tool.Tool {
	return map[string]tool.Tool{
		"Read": &readOnlyStub{name: "Read", desc: "Read file contents"},
		"Glob": &readOnlyStub{name: "Glob", desc: "Find files by pattern"},
		"Grep": &readOnlyStub{name: "Grep", desc: "Search file contents with regex"},
		"Tree": &readOnlyStub{name: "Tree", desc: "Display directory tree"},
	}
}

// ReviewToolSet returns tools for code review: read + diff.
func ReviewToolSet() map[string]tool.Tool {
	return map[string]tool.Tool{
		"Read": &readOnlyStub{name: "Read", desc: "Read file contents"},
		"Diff": &readOnlyStub{name: "Diff", desc: "Compare file contents"},
	}
}

// TestGenToolSet returns tools for test generation: read + write + edit + bash.
func TestGenToolSet() map[string]tool.Tool {
	return map[string]tool.Tool{
		"Read":  &readOnlyStub{name: "Read", desc: "Read file contents"},
		"Write": &readOnlyStub{name: "Write", desc: "Write file contents"},
		"Edit":  &readOnlyStub{name: "Edit", desc: "Edit file with string replacement"},
		"Bash":  &readOnlyStub{name: "Bash", desc: "Execute shell commands"},
	}
}

// readOnlyStub implements tool.Tool with stubbed execution.
// In real usage, the tool registry from the main agent provides real implementations.
type readOnlyStub struct {
	name, desc string
}

func (s *readOnlyStub) Name() string                { return s.name }
func (s *readOnlyStub) Description() string         { return s.desc }
func (s *readOnlyStub) InputSchema() map[string]any { return map[string]any{"type": "object"} }
func (s *readOnlyStub) RequiresPermission() bool    { return false }
func (s *readOnlyStub) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelReadOnly
}
func (s *readOnlyStub) Execute(_ context.Context, _ map[string]any) tool.Result {
	return tool.Success("stub: " + s.name)
}

// toolSet returns the tool set for the given sub-agent type.
func toolSet(st Type) map[string]tool.Tool {
	switch st {
	case TypeExplore:
		return ExploreToolSet()
	case TypeReview:
		return ReviewToolSet()
	case TypeTestGen:
		return TestGenToolSet()
	default:
		return nil
	}
}

// maxTurns returns the max turns for the given sub-agent type.
func maxTurns(st Type) int {
	switch st {
	case TypeExplore:
		return maxExploreTurns
	case TypeReview:
		return maxReviewTurns
	case TypeTestGen:
		return maxTestGenTurns
	default:
		return 5
	}
}

type subRegistry struct {
	tools map[string]tool.Tool
}

func (r *subRegistry) GetTool(name string) tool.Tool { return r.tools[name] }
func (r *subRegistry) Execute(ctx context.Context, name string, input map[string]any) tool.Result {
	t := r.tools[name]
	if t == nil {
		return tool.Error("tool not found: " + name)
	}
	return t.Execute(ctx, input)
}
func (r *subRegistry) GetAllDefinitions() []tool.ToolDefinition {
	defs := make([]tool.ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		defs = append(defs, tool.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return defs
}
func (r *subRegistry) GetDefinitionsByTier(_ ...tool.ToolTier) []tool.ToolDefinition {
	return r.GetAllDefinitions()
}

// Run executes a sub-agent with an isolated context and returns its result.
func Run(
	ctx context.Context,
	subType Type,
	prompt string,
	apiClient apiClient,
	mainRegistry toolRegistry,
	policy permissionPolicy,
) (string, error) {
	tools := toolSet(subType)
	if tools == nil {
		return "", fmt.Errorf("unknown sub-agent type: %s", subType)
	}

	// Build an isolated registry: use main registry tools when available,
	// fall back to stubs.
	isolated := &subRegistry{tools: make(map[string]tool.Tool)}
	for name := range tools {
		if real := mainRegistry.GetTool(name); real != nil {
			isolated.tools[name] = real
		} else {
			isolated.tools[name] = tools[name]
		}
	}

	sysPrompt := fmt.Sprintf(
		"You are a specialized %s sub-agent. Focus only on the task below. Be concise and return only the requested information. Do not use tools not in your tool set.",
		subType,
	)

	a := agent.NewAgent(apiClient, isolated, &subPolicyAdapter{policy}, sysPrompt, "")
	a.SetThinkingBudget(0) // sub-agents don't need extended thinking

	result, err := a.Run(ctx, prompt, func(string) {})
	if err != nil {
		return "", fmt.Errorf("%s sub-agent failed: %w", subType, err)
	}

	return fmt.Sprintf("[%s Sub-Agent Result]\n%s", subType, strings.TrimSpace(result)), nil
}

// subPolicyAdapter adapts a simple permission policy to agent's interface.
type subPolicyAdapter struct {
	p permissionPolicy
}

func (a *subPolicyAdapter) Evaluate(toolName string, input map[string]any, requiresPermission bool) permission.Decision {
	return a.p.Evaluate(toolName, input, requiresPermission)
}
