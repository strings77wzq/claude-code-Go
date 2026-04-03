---
title: Architecture Overview
description: Harness-First Engineering — where the model provides intelligence and the harness provides reliability
---

# Architecture Overview

In the rapidly evolving landscape of AI-assisted software engineering, the distinction between what the model knows and what the system guarantees has become the defining architectural question. go-code embraces a fundamental principle: **the model provides intelligence, the harness provides reliability**.

This is not merely a division of labor — it is an architectural philosophy that shapes every design decision, from the permission system that guards dangerous operations to the context management that maximizes the utility of limited context windows.

## The Five Pillars of go-code Architecture

### Harness-First Engineering

The harness isn't an afterthought; it's the foundation. When building an AI agent system that operates with the same privileges as its user, reliability cannot be layered on top — it must be baked in from the first line of code.

**Permission control** ensures that destructive operations (file deletion, shell command execution, network requests) never execute without explicit user consent or appropriate policy configuration. **Timeout protection** guarantees that a misbehaving tool — whether a runaway grep search or an infinite loop — cannot hang the entire session. **Output truncation** prevents a single tool response from consuming the entire context window. **Session persistence** means the agent can resume after interruption without losing conversation history or context.

These aren't features bolted onto the agent loop. They are the safety net that makes AI agents production-ready. The model can hallucinate, misinterpret, or choose the wrong tool — but the harness must never permit a destructive operation without authorization, must never hang indefinitely, and must never lose the user's work.

### Context Window as Scarce Resource

Modern LLMs offer generous context windows, but "generous" is a relative term when representing an entire codebase, multiple files, and ongoing conversation history. go-code treats context as a scarce resource to be managed, not a luxury to be spent freely.

The architecture implements a **lazy tool loading** strategy — tools are registered in a central registry and only the tool definitions (name, description, input schema) are sent to the model, not their implementations. Tool descriptions themselves are capped at 250 characters, forcing discipline in how capabilities are communicated to the model.

For long conversations, go-code implements **three-tier compaction**: messages are retained in full until a threshold is reached, then compressed into summaries that preserve key information while reducing token count. This allows conversations to span dozens of turns without hitting API limits.

### Layered Permission Defense

Trust is not binary, and neither is authorization. go-code implements a three-tier permission model that reflects real-world security thinking:

| Level | Capability | Use Case |
|-------|------------|----------|
| **ReadOnly** | File reading, globbing, grep searching | Safe exploration of codebases |
| **WorkspaceWrite** | File creation, modification, editing | Controlled development work |
| **DangerFullAccess** | Shell execution, deletions, network operations | Full automation with explicit approval |

The permission system doesn't merely check a level — it evaluates each operation against **glob rule matching** to determine whether specific paths or operations are permitted. A policy might allow editing files in `/src/` while blocking modifications to `/config/` or `/tests/`.

**Session memory** means the system remembers approval decisions within a session, preventing repetitive prompts for the same operation. A human-in-the-loop approval mode exists for operations requiring explicit confirmation before execution — critical for destructive operations or shell commands that could compromise the system.

### Local-First Execution

go-code runs as a pure Go binary with zero runtime dependencies. There is no Python environment to configure, no Node.js runtime to maintain, no virtual machine to manage. The agent executes in the full environment of the host machine, with access to all system resources, environment variables, and installed tooling.

This **local-first** approach eliminates virtualization overhead and provides maximum performance. The agent can invoke the same build tools, linters, and test runners that developers use daily. There's no abstraction layer between the model's intent and the system's actual capabilities.

### MCP as Universal Bridge

The Model Context Protocol (MCP) is not a closed ecosystem — it is an interoperability hub. go-code treats MCP as a bridge to external tool ecosystems, enabling the agent to discover and utilize capabilities beyond its built-in set.

MCP servers can be configured to provide specialized tooling: database access, API integrations, cloud service management, or domain-specific analyzers. The agent treats MCP tools the same as built-in tools — with full permission enforcement and the same execution guarantees.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                 go-code                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                              User                                      │   │
│  │                          (REPL / CLI)                                  │   │
│  └─────────────────────────────┬────────────────────────────────────────┘   │
│                                │                                              │
│                                ▼                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                           Agent Loop                                   │   │
│  │                    State Machine: THINK → ACT → OBSERVE              │   │
│  │                                                                       │   │
│  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────────────┐    │   │
│  │  │   History    │───▶│   Request    │───▶│    Stop Reason       │    │   │
│  │  │  Management  │    │   Builder   │    │    Dispatch          │    │   │
│  │  └──────────────┘    └──────────────┘    └──────────────────────┘    │   │
│  └─────────────────────────────┬────────────────────────────────────────┘   │
│                                │                                              │
│                                ▼                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                          API Client                                   │   │
│  │                    SSE Streaming + Auth                              │   │
│  └─────────────────────────────┬────────────────────────────────────────┘   │
│                                │                                              │
│                                ▼                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                        Anthropic API                                  │   │
│  │                     (Claude Model + Tools)                            │   │
│  └─────────────────────────────┬────────────────────────────────────────┘   │
│                                │                                              │
│  ═══════════════════════════════╪═══════════════════════════════════════════   │
│                                │                                              │
│  ┌─────────────────────────────┴────────────────────────────────────────┐   │
│  │                         Tool Layer                                     │   │
│  │                                                                       │   │
│  │  ┌──────────────────┐        ┌──────────────────┐                   │   │
│  │  │  Tool Registry   │───────▶│  Permission Gate │                   │   │
│  │  │                  │        │                  │                   │   │
│  │  │  • Built-in      │        │  • Tier Check    │                   │   │
│  │  │  • MCP Adapter  │        │  • Glob Rules    │                   │   │
│  │  │  • Lazy Loading │        │  • Human Approval│                   │   │
│  │  └──────────────────┘        └──────────────────┘                   │   │
│  │                                    │                                   │   │
│  │                                    ▼                                   │   │
│  │  ┌────────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Execution Layer                             │  │   │
│  │  │                                                                  │  │   │
│  │  │  Read  │  Write  │  Edit  │  Glob  │  Grep  │  Bash  │  MCP... │  │   │
│  │  │                                                                  │  │   │
│  │  │  • Timeout Protection  • Output Truncation  • Error Recovery  │  │   │
│  │  └────────────────────────────────────────────────────────────────┘  │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
User Input
    │
    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Agent Loop                                          │
│  1. Add user message to conversation history                                 │
│  2. Apply compaction if context threshold exceeded                          │
│  3. Build API request: messages + tool definitions + system prompt          │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           API Client                                          │
│  4. Send request via SSE streaming                                           │
│  5. Receive incremental tokens                                              │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │
                              ▼
                       ┌──────────────┐
                       │ Anthropic    │
                       │     API      │
                       └──────┬───────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Agent Loop                                            │
│  6. Parse stop_reason for state transition:                                  │
│     • end_turn: Task complete → Return final response                        │
│     • tool_use: Need more info → Execute tools, continue                     │
│     • max_tokens: Truncated → Warn user, return partial                     │
│     • max_turns: Loop limit → Return with warning                            │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Permission Gate                                        │
│  7. Evaluate operation against permission tier                             │
│  8. Check glob rules for path-specific permissions                          │
│  9. Prompt user if human-in-the-loop mode enabled                           │
│  10. Deny or approve based on policy                                         │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Tool Registry                                           │
│  11. Route to appropriate tool handler (built-in or MCP)                   │
│  12. Execute with timeout protection                                         │
│  13. Truncate output if exceeds limits                                       │
│  14. Return structured result or error                                       │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Agent Loop                                            │
│  15. Add tool results to history                                            │
│  16. Continue loop (back to step 2) or terminate                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Design Rationale

### Why Go?

Go was chosen for its balance of performance, simplicity, and static linking. The agent loop handles real-time streaming, maintains conversation state, and coordinates concurrent tool executions — all requiring predictable memory usage and low latency. Go's single-binary deployment model aligns with the local-first philosophy: no environment configuration, no dependency hell.

### Why Three-Tier Permission?

A binary allow/deny model is insufficient for AI agents that need to perform varied operations across different contexts. A developer might want to allow automatic file editing in a workspace while requiring confirmation for shell commands. The tiered model with glob matching provides the granularity needed for practical, secure automation.

### Why Context Compaction?

Token limits aren't a theoretical concern — they directly impact what the model can reason about. Without compaction, a 50-turn conversation would consume the entire context window, leaving no room for the model to see the code it's editing. Compaction preserves the conversation's semantic essence while making room for new information.

### Why MAX_TURNS?

The model has no intrinsic understanding of time or resource constraints. An agent loop can theoretically continue indefinitely if the model keeps choosing to call tools. MAX_TURNS provides a hard boundary that prevents resource exhaustion while still allowing complex tasks to complete. It's a reliability guarantee that operates independently of model behavior.

## Target Audience

This documentation is written for:

- **AI Agent Architecture Researchers** — examining how harness design influences agent capabilities
- **Senior Engineers** evaluating Harness-First engineering patterns for production systems
- **Tech Leads** seeking reference implementations for production-grade AI tooling

The emphasis throughout is on *why* rather than *what* — the code reveals implementation details, but these pages explain the reasoning behind them.

## Related Documentation

- [Design Philosophy](design-philosophy.md) — Deep dive into the intelligence/reliability separation
- [Agent Loop](agent-loop.md) — State machine mechanics and history management
- [Tool System](tools.md) — Built-in tools and extension mechanisms
- [Configuration Guide](../guide/configuration.md) — Permission policies and model settings
