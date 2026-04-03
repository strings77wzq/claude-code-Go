---
title: Agent Loop Implementation
description: Deep dive into the Run() method — stop_reason dispatch, history management, tool execution, and session persistence
---

# Agent Loop Implementation

The agent loop is the core execution engine in go-code, implemented in `internal/agent/loop.go`. This document provides a detailed walkthrough of the `Run()` method and its key components.

## The Run() Method Overview

```go
func (a *Agent) Run(ctx context.Context, userInput string, 
                   outputCallback func(string)) (string, error)
```

The method takes:
- `ctx` — Context for cancellation
- `userInput` — The user's message
- `outputCallback` — Function to receive streaming text

Returns the final text response or an error.

## Execution Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Run() Method Flow                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Initialize session (generate ID, record start time)            │
│  2. Add user message to history                                    │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    FOR LOOP (MaxTurns)                      │   │
│  │                                                              │   │
│  │  3. Compact history if needed                               │   │
│  │  4. Build API request (system + tools + messages)          │   │
│  │  5. Send to API and receive streaming response             │   │
│  │  6. Add assistant message to history                       │   │
│  │  7. Dispatch on stop_reason:                               │   │
│  │     - end_turn / stop_sequence → return text               │   │
│  │     - max_tokens → return text + warning                   │   │
│  │     - tool_use → execute tools, add results, continue      │   │
│  │     - default → return text                                 │   │
│  │  8. Save session on exit                                   │   │
│  │                                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Session Initialization

```go
func (a *Agent) Run(ctx context.Context, userInput string, 
                   outputCallback func(string)) (string, error) {
    a.sessionID = generateSessionID()
    a.startTime = time.Now()

    var totalInputTokens, totalOutputTokens int

    if err := a.history.AddUserMessage(userInput); err != nil {
        return "", fmt.Errorf("failed to add user message: %w", err)
    }
```

Each session gets:
- Unique ID: `sess_<timestamp>`
- Start timestamp for session tracking
- Token counters for usage monitoring

## The Main Loop

```go
turns := 0

for turns < MaxTurns {
    // Step 1: Compact history if needed
    CompactIfNeeded(a.history, a.contextConfig)

    // Step 2: Build and send API request
    req := a.buildRequest()
    resp, err := a.apiClient.SendMessageStream(ctx, req, outputCallback)
```

### History Compaction

```go
CompactIfNeeded(a.history, a.contextConfig)
```

As the conversation grows, the history is periodically compacted to stay within token limits. This is handled by the compact module to prevent context overflow.

### Request Building

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

Each request includes:
- Model identifier
- Max tokens limit
- System prompt
- Tool definitions from registry
- Current conversation history

## Stop Reason Dispatch

The API response includes a `stop_reason` field that determines the next action:

### end_turn / stop_sequence

```go
switch resp.StopReason {
case "end_turn", "stop_sequence":
    result := extractTextContent(resp.Content)
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

The model believes the task is complete. Return the text content as the final response.

### max_tokens

```go
case "max_tokens":
    result := extractTextContent(resp.Content) + 
              "\n[Warning] Response was truncated (max_tokens reached)."
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

The response was truncated. Return partial content with a warning.

### tool_use

```go
case "tool_use":
    toolResults := a.executeTools(ctx, resp.Content)
    if err := a.history.AddToolResults(toolResults); err != nil {
        a.saveSession(turns, totalInputTokens, totalOutputTokens)
        return "", fmt.Errorf("failed to add tool results: %w", err)
    }
    turns++
    continue
```

The model wants to call tools. Execute all requested tools, add results to history, and continue the loop.

### unknown / default

```go
default:
    result := extractTextContent(resp.Content)
    a.saveSession(turns, totalInputTokens, totalOutputTokens)
    return result, nil
```

Fallback for unexpected stop reasons. Return whatever content was generated.

## MAX_TURNS Safety Limit

```go
// MaxTurns is the maximum number of agent loop iterations to prevent infinite loops.
const MaxTurns = 50
```

The loop enforces a maximum of 50 iterations to prevent:
- Infinite tool execution loops
- Unproductive model oscillation
- Resource exhaustion

If reached:

```go
result := "[Agent loop stopped] Reached maximum turns (" + 
          fmt.Sprintf("%d", MaxTurns) + ")."
a.saveSession(turns, totalInputTokens, totalOutputTokens)
return result, nil
```

## History Management

### Message Types

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

### API Message Format

The Anthropic API expects messages in this format:

```json
{
  "messages": [
    { "role": "user", "content": "..." },
    { "role": "assistant", "content": "..." },
    { "role": "tool", "tool_use_id": "...", "content": "..." }
  ]
}
```

## Tool Execution with Permission Check

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

        // Permission check
        if !a.checkPermission(toolName, toolInput) {
            toolResults = append(toolResults, api.ContentBlock{
                Type:      "tool_result",
                ToolUseID: toolUseID,
                IsError:   true,
            })
            continue
        }

        // Execute pre-hooks
        if a.hooksRegistry != nil {
            if err := a.hooksRegistry.RunPreHooks(toolName, toolInput); err != nil {
                toolResults = append(toolResults, api.ContentBlock{
                    Type:      "tool_result",
                    ToolUseID: toolUseID,
                    Text:      "pre-hook error: " + err.Error(),
                    IsError:   true,
                })
                continue
            }
        }

        // Execute tool
        result := a.toolRegistry.Execute(ctx, toolName, toolInput)

        // Execute post-hooks
        if a.hooksRegistry != nil {
            a.hooksRegistry.RunPostHooks(toolName, toolInput, 
                                         result.Content, result.IsError)
        }

        // Collect result
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

### Permission Check

```go
func (a *Agent) checkPermission(toolName string, 
                                input map[string]any) bool {
    t := a.toolRegistry.GetTool(toolName)
    requiresPermission := t != nil && t.RequiresPermission()

    decision := a.permissionPolicy.Evaluate(toolName, input, requiresPermission)
    return decision == permission.Allow || decision == permission.Ask
}
```

If the tool's `RequiresPermission()` returns `true`, the permission policy evaluates whether to allow, deny, or ask the user.

## Session Persistence

```go
func (a *Agent) saveSession(turnCount, inputTokens, outputTokens int) {
    if a.sessionID == "" {
        return
    }

    s := &session.Session{
        ID:           a.sessionID,
        Model:        a.model,
        StartTime:    a.startTime,
        EndTime:      time.Now(),
        TurnCount:    turnCount,
        InputTokens:  inputTokens,
        OutputTokens: outputTokens,
    }

    messages := a.convertHistoryToSessionMessages()
    dir := getSessionsDir()

    if err := session.SaveSession(s, messages, dir); err != nil {
        fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
    }
}
```

Sessions are saved to `~/.claude-code-go/sessions/` with:
- Session ID and timestamps
- Model used
- Turn count and token usage
- Full conversation history

## Related Documentation

- [Entry Point Walkthrough](entry-point.md) — main.go initialization sequence
- [Tool System Overview](../tools/overview.md) — Tool interface and registry
- [Architecture Overview](../architecture/overview.md) — System components

---

<div class="nav-prev-next">

- [Entry Point Walkthrough](entry-point.md) ←
- → [Tool System Overview](../tools/overview.md)

</div>