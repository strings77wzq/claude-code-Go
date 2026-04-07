# Tutorial 3: Understanding the Agent Loop

Deep dive into how claude-code-Go's agent loop works.

## What is the Agent Loop?

The agent loop is the core decision-making engine. It's a continuous cycle where the AI:

1. **Thinks** about what to do
2. **Acts** by calling tools
3. **Observes** the results
4. **Repeats** until done

This is inspired by the ReAct (Reasoning + Acting) pattern from AI research.

## The Loop States

```
┌─────────────┐     ┌──────────┐     ┌────────────┐
│   User      │────▶│  Think   │────▶│   Act      │
│   Input     │     │          │     │  (Tool)    │
└─────────────┘     └──────────┘     └────────────┘
                                             │
                                             ▼
┌─────────────┐     ┌──────────┐     ┌────────────┐
│   Respond   │◀────│  Think   │◀────│  Observe   │
│   to User   │     │  Again   │     │  (Result)  │
└─────────────┘     └──────────┘     └────────────┘
```

## State Breakdown

### 1. User Input

You type a message. This could be:
- A question: "What does this code do?"
- A command: "Refactor this function"
- A conversation: "Let's plan the architecture"

### 2. Think

The AI processes your input and decides what to do:

```
User: "Read main.go and explain it"

AI thinks:
- The user wants me to read a file
- I should use the Read tool
- File path is "main.go"
- After reading, I'll need to explain the content
```

### 3. Act

The AI calls one or more tools:

```json
{
  "tool": "Read",
  "arguments": {
    "file_path": "main.go"
  }
}
```

### 4. Observe

The tool executes and returns results:

```
Result: package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
```

### 5. Think Again

The AI processes the result:

```
AI thinks:
- I've successfully read the file
- It's a simple Go program
- It prints "Hello"
- I should explain this to the user
```

### 6. Respond

The AI generates a response:

```
This is a simple Go program that:
1. Declares itself as package "main"
2. Imports the "fmt" package for formatting
3. Has a main() function that prints "Hello"

When run, it will output: Hello
```

## Stop Reasons

The loop continues until one of these conditions:

| Stop Reason | Meaning | Example |
|-------------|---------|---------|
| `end_turn` | AI is done responding | After explaining code |
| `tool_use` | AI wants to use a tool | "Let me read that file" |
| `max_tokens` | Context limit reached | Long conversation |
| `stop_sequence` | Special stop token | End of generation |

## Multi-Turn Example

Complex tasks require multiple loop iterations:

```
User: "Find all TODO comments and create a task list"

Loop 1:
  Think: "I need to search for TODO comments"
  Act: Grep {"pattern": "TODO", "path": "."}
  Observe: Found 3 TODOs in 2 files

Loop 2:
  Think: "I found TODOs, now I need to read them"
  Act: Read {"file_path": "main.go"}
  Observe: "// TODO: Add error handling"

Loop 3:
  Think: "I have the details, now create the task list"
  Act: TodoWrite {"todos": ["Add error handling in main.go"]}
  Observe: Task list created

Loop 4:
  Think: "All tasks complete, respond to user"
  Respond: "Created a task list with 1 item..."
```

## Loop Limits

To prevent infinite loops:

- **Max turns**: 10 iterations per user message (configurable)
- **Timeout**: 30 seconds per tool execution
- **Token limit**: Automatic context compaction when approaching limits

You can adjust these in settings:

```json
{
  "maxTurns": 20,
  "timeout": "60s"
}
```

## Visualizing the Loop

Enable debug mode to see the loop in action:

```
> /debug on

# Now you'll see each state transition:
[THINK] Analyzing user request...
[ACT] Calling tool: Read
[OBSERVE] Tool returned 150 lines
[THINK] Summarizing content...
[RESPOND] Generating response...
```

## Exercise: Trace the Loop

For this request, trace through each state:

```
> Find all functions in main.go and count them
```

<details>
<summary>Click to see the answer</summary>

**Loop 1:**
- Think: User wants to find functions in main.go
- Act: Read main.go
- Observe: File content returned

**Loop 2:**
- Think: I have the file content, need to count functions
- Act: No tool needed, I can analyze directly
- Respond: "Found 3 functions: main(), init(), helper()"

</details>

## Next Steps

- [Tutorial 4: Permission System](04-permission-system.md) - Learn about security
- [Architecture Deep Dive](../../architecture/agent-loop-state-machine.md) - Technical details
- [ADR-001: Agent Loop Design](../../adr/001-agent-loop.md) - Design decisions
