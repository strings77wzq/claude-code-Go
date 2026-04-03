# Agent Loop

The agent loop is the core execution engine that orchestrates the interaction between the model, tools, and user.

## Overview

The agent loop implements the following cycle:

```
┌────────────────────────────────────────────┐
│            User Prompt                      │
└─────────────────┬──────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────────┐
│  1. Add User Message to History            │
└─────────────────┬──────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────────┐
│  2. Send History to API                     │
│     (with tool schemas)                     │
└─────────────────┬──────────────────────────┘
                  │
                  ▼
┌────────────────────────────────────────────┐
│  3. Receive Model Response                  │
│     (streaming)                             │
└─────────────────┬──────────────────────────┘
                  │
                  ▼
         ┌────────────────┐
         │ Contains Tool  │
         │    Calls?       │
         └────────┬───────┘
            YES   │   NO
            ┌─────▼─────┐   ┌─────────────────┐
            │           │   │ Display Final   │
            │ 4. Execute│   │    Response     │
            │   Tools   │   └────────┬────────┘
            │           │            │
            └─────┬─────┘            │
                  │                   │
                  ▼                   │
         ┌────────────────┐           │
         │ Add Tool Result│           │
         │    to History  │           │
         └────────┬───────┘           │
                  │                    │
                  └────────┬───────────┘
                           │
                           ▼
                    Continue to Step 2
```

## Implementation Details

### Message History

The agent maintains a list of messages that includes:
- User messages (prompts)
- Assistant messages (model responses)
- Tool result messages (tool outputs)

```go
type Message struct {
    Role    string    // "user", "assistant", "tool"
    Content string    // Text content
    ToolID  string    // Tool call ID (for tool results)
}
```

### Tool Selection

When the model responds with tool calls, the agent:
1. Parses the tool call from the model's response
2. Validates the tool exists in the registry
3. Checks permissions for the tool operation
4. Executes the tool with the provided arguments
5. Captures the result and adds it to history

### Continuation Criteria

The loop continues when:
- The model requests tool execution
- The response is incomplete (streaming)

The loop terminates when:
- The model provides a final text response (no tool calls)
- An error occurs
- The maximum iterations are reached (to prevent infinite loops)

### Tool Schema

The model is provided with a schema describing all available tools:

```json
{
  "name": "Read",
  "description": "Read the contents of a file",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to read"
      }
    },
    "required": ["file_path"]
  }
}
```

## Error Handling

The agent loop handles several error conditions:

- **API Errors**: Retry with exponential backoff
- **Permission Denied**: Skip tool execution, inform model
- **Tool Execution Failed**: Return error to model, allow recovery
- **Timeout**: Return timeout error to model

## Configuration

Agent loop behavior can be configured:

```yaml
# Maximum tool calls per conversation (prevents infinite loops)
max_iterations: 100

# Timeout for tool execution (seconds)
tool_timeout: 30

# Whether to show tool execution progress
verbose: true
```

## Related Concepts

- [Built-in Tools](tools.md) - Available tool implementations
- [Permission System](overview.md) - Permission handling
- [API Client](../api/overview.md) - Anthropic API communication