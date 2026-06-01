## Context

The current agent processes all tasks in a single context. Complex multi-step tasks (e.g., "find all security issues and generate fixes") benefit from decomposition: explore the codebase → review findings → generate tests. Each phase has different information needs and can run in isolation.

## Goals / Non-Goals

**Goals:**
- Define SubAgent interface with isolated history and filtered tool sets
- Three specialized sub-agent types: Explore, Review, TestGen
- Task tool for the main agent to spawn sub-agents
- Concurrent execution with context cancellation
- Structured result format for synthesis

**Non-Goals:**
- Persistent sub-agent memory across sessions
- Dynamic sub-agent type registration (v1 is hardcoded)
- Hierarchical sub-agent nesting
- Real-time streaming from sub-agents

## Decisions

### 1. Sub-agent as lightweight Agent with filtered tool registry

Sub-agents use the same Agent core but with:
- Isolated History (fresh context, no shared messages)
- Filtered tool Registry (only relevant tools)
- Controlled max turns (smaller than main agent)
- Result returned as structured text

### 2. Task tool for invocation

A new `Task` builtin tool accepts `{subagent_type, prompt}` and returns structured output. This follows Anthropic's own pattern where Claude Code uses Task tool for delegation.

### 3. Tool sets per sub-agent type

| Type | Tools | Max Turns |
|------|-------|-----------|
| Explore | Read, Glob, Grep, Tree | 15 |
| Review | Read, Diff | 10 |
| TestGen | Read, Write, Edit, Bash | 20 |

## Risks / Trade-offs

- [API cost] → Each sub-agent makes independent API calls. Mitigation: sub-agents use smaller context, fewer tools, and lower max tokens
- [Latency] → Sequential sub-agents add latency. Mitigation: concurrent execution option
- [Complexity] → Sub-agent system adds code. Mitigation: reuse existing Agent struct
