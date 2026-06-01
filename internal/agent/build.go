package agent

import (
	"github.com/strings77wzq/claude-code-Go/internal/api"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// buildRequest assembles the API request with progressive tool disclosure and prompt caching.
func (a *Agent) buildRequest() *api.ApiRequest {
	// Progressive tool disclosure: turn 1 sends all tools; subsequent turns send
	// core + any extension tools the model previously requested.
	var toolDefs []api.ToolDefinition
	if a.turnCount == 0 || a.cacheInvalidated {
		defs := a.toolRegistry.GetAllDefinitions()
		toolDefs = make([]api.ToolDefinition, 0, len(defs))
		for _, td := range defs {
			toolDefs = append(toolDefs, api.ToolDefinition{
				Name:        td.Name,
				Description: td.Description,
				InputSchema: td.InputSchema,
			})
			a.disclosedTools[td.Name] = true
		}
	} else {
		coreDefs := a.toolRegistry.GetDefinitionsByTier(tool.TierCore)
		toolDefs = make([]api.ToolDefinition, 0, len(coreDefs)+len(a.disclosedTools))
		for _, td := range coreDefs {
			toolDefs = append(toolDefs, api.ToolDefinition{
				Name:        td.Name,
				Description: td.Description,
				InputSchema: td.InputSchema,
			})
		}
		// Include any extension/MCP tools already disclosed this session
		for _, td := range a.toolRegistry.GetAllDefinitions() {
			if td.Tier != tool.TierCore && a.disclosedTools[td.Name] {
				toolDefs = append(toolDefs, api.ToolDefinition{
					Name:        td.Name,
					Description: td.Description,
					InputSchema: td.InputSchema,
				})
			}
		}
	}

	req := &api.ApiRequest{
		Model:     a.model,
		MaxTokens: a.maxTokens,
		Stream:    true,
		Tools:     toolDefs,
		Messages:  a.history.GetMessages(),
	}

	// Prompt caching: cache system prompt and tool definitions
	if a.cacheEnabled {
		req.System = api.CachedSystemPrompt{
			Type:         "text",
			Text:         a.systemPrompt,
			CacheControl: &api.CacheControl{Type: "ephemeral"},
		}
		if len(toolDefs) > 0 {
			req.Tools = a.wrapToolsWithCache(toolDefs)
		}
	} else {
		req.System = a.systemPrompt
	}

	if a.thinkingBudget > 0 {
		req.Thinking = &api.Thinking{
			Type:         "enabled",
			BudgetTokens: a.thinkingBudget,
		}
	}

	a.cacheInvalidated = false
	return req
}

// wrapToolsWithCache wraps the last tool definition with a cache breakpoint.
func (a *Agent) wrapToolsWithCache(defs []api.ToolDefinition) []api.ToolDefinition {
	if len(defs) == 0 {
		return defs
	}
	cached := make([]api.ToolDefinition, len(defs))
	copy(cached, defs)
	cached[len(cached)-1].CacheControl = &api.CacheControl{Type: "ephemeral"}
	return cached
}
