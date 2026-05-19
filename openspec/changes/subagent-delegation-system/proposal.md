## Why

The current agent is a single linear loop — one model, one context, one thread of execution. Modern agent architectures (AutoGen, CrewAI, LangGraph) demonstrate that delegating specialized work to sub-agents with focused contexts and tool sets produces better results for complex multi-step tasks. A lightweight delegation system allows the main agent to spawn Explore, Review, and TestGen sub-agents that work in isolation and return structured results.

## What Changes

- Sub-agent spawner with isolated contexts, focused tool sets, and result synthesis
- Three built-in sub-agent types: Explore (code exploration), Review (code review), TestGen (test generation)
- Task tool integration: main agent invokes sub-agents via a builtin `Task` tool
- Concurrent execution via goroutines with timeout and cancellation

## Capabilities

### New Capabilities

- `subagent-delegation`: Lightweight sub-agent system with isolated contexts, specialized tool sets, concurrent execution, and structured result synthesis

### Modified Capabilities

None.

## Impact

- New package: `internal/agent/subagent/` — sub-agent types and runner
- `internal/tool/builtin/task.go` — new Task tool for invoking sub-agents
- `internal/agent/loop.go` — optional sub-agent dispatch integration
