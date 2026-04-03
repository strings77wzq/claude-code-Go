---
title: Agent Loop Deep Dive
description: Detailed explanation of the think → act → observe execution cycle, stop reasons, and history management
---

# Agent Loop Deep Dive

The agent loop is the core execution engine of go-code. It implements a think → act → observe cycle that allows the LLM to iteratively perform tasks by calling tools and processing their results.

## The Think → Act → Observe Cycle

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Agent Loop                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌────────────────────┐   │
│  │    THINK     │───▶│     ACT      │───▶│     OBSERVE        │   │
│  │              │    │              │    │                    │   │
│  │ 1. Build     │    │ 3. Call API  │    │ 5. Process result │   │
│  │    request   │    │    with      │    │    and decide     │   │
│  │    (history+ │    │    tools     │    │    next action    │   │
│  │    tools)    │    │              │    │                    │   │
│  └──────────────┘    └──────────────┘    └────────────────────┘   │
│                                                                     │
│         ▲                                                            │
│         │                                                            │
│         └────────────────────────────────────────────────────────────┤
│                            ITERATE                                  │
│         ┌───────────────────────────────────────────────────────────┤
│         │                                                            │
│         │  6. Add tool results to history                          │
│         │  7. Loop back to step 1                                  │
│         │                                                            │
│         ▼                                                            │
└─────────────────────────────────────────────────────────────────────┘
```

### Step-by-Step Breakdown

1. **Think**: Build the API request with conversation history, system prompt, and available tools
2. **Act**: Send the request to the Anthropic API and receive a streaming response
3. **Observe**: Process the response based on the `stop_reason` to decide the next action

## Stop Reason Dispatch

The API response includes a `stop_reason` field that indicates why the model stopped generating. go-code handles each case appropriately:

### end_turn / stop_sequence

The model believes the task is complete and has provided a final response.

```go
case "end_turn", "stop_sequence":
    // LLM believes task is complete
    return extractTextContent(resp.Content), nil
```

**Action**: Return the final text response to the user.

### tool_use

The model wants to call one or more tools to gather more information or perform actions.

```go
case "tool_use":
    // LLM wants to call tools - execute them and continue loop
    toolResults := a.executeTools(ctx, resp.Content)
    if err := a.history.AddToolResults(toolResults); err != nil {
        return "", fmt.Errorf("failed to add tool results: %w", err)
    }
    turns++
    continue
```

**Action**: Execute all requested tools, add results to history, and continue the loop.

### max_tokens

The response was truncated because the maximum token limit was reached.

```go
case "max_tokens":
    // Output was truncated
    return extractTextContent(resp.Content) + 
           "\n[Warning] Response was truncated (max_tokens reached).", nil
```

**Action**: Return the partial response with a warning.

### unknown / default

Fallback handling for unexpected stop reasons.

```go
default:
    // Unknown stop_reason - safe fallback
    return extractTextContent(resp.Content), nil
```

## MAX_TURNS Safety Limit

To prevent infinite loops, go-code enforces a maximum number of iterations:

```go
// MaxTurns is the maximum number of agent loop iterations to prevent infinite loops.
const MaxTurns = 50

turns := 0
for turns < MaxTurns {
    // ... agent loop ...
}
```

If the agent reaches 50 turns without completing, it returns:

```
[Agent loop stopped] Reached maximum turns (50).
```

This protects against:
- Tool execution loops (model calling same tool repeatedly)
- Unproductive cycles (model oscillating between actions)
- Resource exhaustion from long-running sessions

## History Management

The conversation history maintains the user/assistant message alternation required by the Anthropic API.

### Adding Messages

```go
// Add user message
if err := a.history.AddUserMessage(userInput); err != nil {
    return "", fmt.Errorf("failed to add user message: %w", err)
}

// Add assistant response
if err := a.history.AddAssistantMessage(resp.Content); err != nil {
    return "", fmt.Errorf("failed to add assistant message: %w", err)
}

// Add tool results
if err := a.history.AddToolResults(toolResults); err != nil {
    return "", fmt.Errorf("failed to add tool results: %w", err)
}
```

### Message Format

The API expects messages in this format:

```json
{
  "messages": [
    { "role": "user", "content": "..." },
    { "role": "assistant", "content": "..." },
    { "role": "tool", "tool_use_id": "...", "content": "..." }
  ]
}
```

### History Compaction

As the conversation grows, go-code can compact history to stay within token limits. This is handled by the compact module (`internal/agent/compact.go`).

## Streaming Output Callback

The agent supports real-time streaming output via a callback function:

```go
// Signature
func (a *Agent) Run(ctx context.Context, userInput string, 
                   outputCallback func(string)) (string, error)

// Usage in API client
resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
```

The callback is invoked for each text delta received from the API, allowing:

- Real-time display of model output
- Progress indicators during tool execution
- Interactive user experience in the REPL

## Tool Execution Flow

When the model requests tool execution:

```go
func (a *Agent) executeTools(ctx context.Context, 
                             content []api.ContentBlock) []api.ContentBlock {
    var toolResults []api.ContentBlock

    for _, block := range content {
        if block.Type != "tool_use" {
            continue
        }

        toolName := block.Name
        toolInput := block.Input
        toolUseID := block.ID

        // 1. Check permission
        if !a.checkPermission(toolName, toolInput) {
            toolResults = append(toolResults, api.ContentBlock{
                Type:      "tool_result",
                ToolUseID: toolUseID,
                IsError:   true,
            })
            continue
        }

        // 2. Execute tool
        result := a.toolRegistry.Execute(ctx, toolName, toolInput)
        
        // 3. Collect result
        toolResults = append(toolResults, api.ContentBlock{
            Type:      "tool_result",
            ToolUseID: toolUseID,
            Text:      result.Content,
            IsError:   result.IsError,
        })
    }

    return toolResults
}
```

## Request Building

Each iteration builds an API request with:

```go
func (a *Agent) buildRequest() *api.ApiRequest {
    toolDefs := make([]api.ToolDefinition, 0)
    for _, td := range a.toolRegistry.GetAllDefinitions() {
        toolDefs = append(toolDefs, api.ToolDefinition{
            Name:        td.Name,
            Description: td.Description,
            InputSchema: td.InputSchema,
        })
    }

    return &api.ApiRequest{
        Model:     a.model,
        MaxTokens: a.maxTokens,
        System:    a.systemPrompt,
        Stream:    true,
        Tools:     toolDefs,
        Messages:  a.history.GetMessages(),
    }
}
```

## Related Documentation

- [Architecture Overview](overview.md) — System components and data flow
- [Tool System](tools.md) — Tool implementations and the registry
- [Configuration Guide](../guide/configuration.md) — Model and token settings