---
title: Architecture Overview
description: Harness-First Engineering — where the model provides intelligence and the harness provides reliability
---

# Architecture Overview

## The Defining Principle

**Harness-First Engineering** — the harness isn't an afterthought; it's the foundation.

In the rapidly evolving landscape of AI-assisted software engineering, the distinction between what the model knows and what the system guarantees has become the defining architectural question. This project embraces a fundamental principle: **the model provides intelligence, the harness provides reliability**.

This is not merely a division of labor — it is an architectural constraint that shapes every design decision, from the permission system that guards dangerous operations to the context management that maximizes the utility of limited context windows.

## Five Architectural Advantages

### a. Harness-First Engineering

When building an AI agent system that operates with the same privileges as its user, reliability cannot be layered on top — it must be embedded from the first line of code.

**Permission control**, **timeout protection**, **output truncation**, and **session persistence** aren't features — they're the safety net that makes AI agents production-ready. The model can hallucinate, misinterpret instructions, or choose the wrong tool. But the harness must never permit a destructive operation without authorization, must never hang indefinitely, and must never lose the user's work.

| Property | Guarantee |
|----------|-----------|
| Permission Control | No destructive operation executes without user consent or policy approval |
| Timeout Protection | Misbehaving tools (runaway grep, infinite loop) cannot hang the session |
| Output Truncation | Single tool response cannot consume the entire context window |
| Session Persistence | Agent resumes after interruption without losing conversation history |

### b. Context Window as Scarce Resource

Modern LLMs offer generous context windows, but "generous" is a relative term when representing an entire codebase, multiple files, and ongoing conversation history. This architecture treats context as a scarce resource to be managed, not a luxury to be spent freely.

- **Lazy Tool Loading**: Tools are registered in a central registry; only tool definitions (name, description, input schema) are sent to the model — not their implementations
- **250-Character Skill Descriptions**: Force discipline in how capabilities are communicated to the model
- **Three-Tier Compaction**: Messages retained in full until threshold, then compressed into summaries that preserve key information while reducing token count

### c. Layered Permission Defense

Trust is not binary, and neither is authorization. This system implements a three-tier permission model that reflects real-world security thinking:

| Level | Capability | Use Case |
|-------|------------|----------|
| **ReadOnly** | File reading, globbing, grep searching | Safe exploration of codebases |
| **WorkspaceWrite** | File creation, modification, editing | Controlled development work |
| **DangerFullAccess** | Shell execution, deletions, network operations | Full automation with explicit approval |

The permission system doesn't merely check a level — it evaluates each operation against **glob rule matching** to determine whether specific paths or operations are permitted. A policy might allow editing files in `/src/` while blocking modifications to `/config/` or `/tests/`.

**Session memory** means the system remembers approval decisions within a session, preventing repetitive prompts for the same operation. A **human-in-the-loop approval** mode exists for operations requiring explicit confirmation before execution — critical for destructive operations or shell commands that could compromise the system.

### d. Local-First Execution

This system runs as a pure Go binary with zero runtime dependencies. There is no Python environment to configure, no Node.js runtime to maintain, no virtual machine to manage. The agent executes in the full environment of the host machine, with access to all system resources, environment variables, and installed tooling.

This **local-first** approach eliminates virtualization overhead and provides maximum performance. The agent can invoke the same build tools, linters, and test runners that developers use daily. There is no abstraction layer between the model's intent and the system's actual capabilities.

### e. MCP as Universal Bridge

The Model Context Protocol (MCP) is not a closed ecosystem — it is an interoperability hub. This system treats MCP as a bridge to external tool ecosystems, enabling the agent to discover and utilize capabilities beyond its built-in set.

MCP servers can be configured to provide specialized tooling: database access, API integrations, cloud service management, or domain-specific analyzers. The agent treats MCP tools the same as built-in tools — with full permission enforcement and the same execution guarantees.

## Agent Loop: State Machine

```
┌──────────┐  tool_use  ┌──────────────┐
│ REQUEST  │───────────▶│ TOOL_EXEC    │
│          │◀───────────│              │
└────┬─────┘  result    └──────┬───────┘
     │                         │
end_turn│                 error│
     ▼                         ▼
┌──────────┐           ┌──────────────┐
│ TERMINATE│           │ CONTINUE     │
└──────────┘           └──────────────┘
```

| State | Transition | Condition |
|-------|------------|-----------|
| REQUEST | → TOOL_EXEC | Model invokes tool (tool_use) |
| REQUEST | → TERMINATE | Task complete (end_turn) |
| TOOL_EXEC | → CONTINUE | Tool succeeded, continue loop |
| TOOL_EXEC | → CONTINUE | Tool failed with recoverable error |
| TOOL_EXEC | → TERMINATE | Non-recoverable failure |

## Design Decisions

| Decision | Rationale | Trade-off |
|----------|-----------|-----------|
| **Go Implementation** | Static linking, single-binary deployment, predictable memory usage | Less flexible than interpreted languages |
| **Three-Tier Permission** | Binary allow/deny insufficient for varied development workflows | More complex configuration, enables practical use |
| **Glob Rule Matching** | Fine-grained control over which files can be modified | Requires users to understand glob patterns |
| **250-Character Tool Descriptions** | Forces discipline in capability communication; prevents context bloat | May omit useful tool details |
| **MAX_TURNS = 50** | Prevents resource exhaustion from infinite loops | Complex tasks may need more iterations |
| **Three-Tier Message Compaction** | Preserves semantic content while managing token limits | Compression is lossy; some context may be lost |
| **MCP as External Adapter** | Protocol-agnostic tool discovery; ecosystem integration | MCP servers add external dependencies |

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