---
title: Interactive Architecture Diagram
description: Click through the system components to explore the architecture
---

# Interactive Architecture Diagram

Click on each module to navigate to its detailed documentation.

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         claude-code-Go Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│    ┌─────────────────────────────────────────────────────────────────┐     │
│    │                           TUI Layer                             │     │
│    │  ┌─────────────────┐            ┌─────────────────────────┐     │     │
│    │  │   Bubbletea    │────────────│    Legacy REPL          │     │     │
│    │  │    (Default)   │            │    (--legacy-repl)      │     │     │
│    │  └────────┬────────┘            └─────────────────────────┘     │     │
│    └───────────┼───────────────────────────────────────────────────────┘     │
│                │                                                                │
│                ▼                                                                │
│    ┌─────────────────────────────────────────────────────────────────┐     │
│    │                        Core Engine                              │     │
│    │  ┌─────────────────────────────────────────────────────────────┐│     │
│    │  │                    Agent Loop + Context                     ││     │
│    │  │  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌─────────┐││     │
│    │  │  │  THINK    │─▶│   ACT     │─▶│ OBSERVE   │─▶│ DECIDE  │││     │
│    │  │  └───────────┘  └───────────┘  └───────────┘  └─────────┘││     │
│    │  └─────────────────────────────────────────────────────────────┘│     │
│    └───────────┬───────────────────────────────────────────────────────┘     │
│                │                                                                │
│         ┌──────┴──────┬──────────────────────┐                                │
│         ▼             ▼                      ▼                                 │
│    ┌─────────┐  ┌────────────┐       ┌────────────┐                           │
│    │  Tool   │  │   API      │       │ Permission │                           │
│    │Registry │  │  Client    │       │  System    │                           │
│    └────┬────┘  └─────┬──────┘       └─────┬──────┘                           │
│         │             │                    │                                   │
│         ▼             ▼                    ▼                                   │
│    ┌─────────┐  ┌────────────┐       ┌────────────┐                           │
│    │Built-in │  │   SSE      │       │   3-Tier   │                           │
│    │ Tools   │  │  Streaming │       │  Controls  │                           │
│    └─────────┘  └────────────┘       └────────────┘                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Module Details

<div class="module-grid">

### [TUI Layer](/architecture/overview)
**Entry point** — User interface for interacting with the agent
- **Bubbletea** (default): Modern TUI with rich formatting
- **Legacy REPL**: Simple line-by-line interface

### [Agent Loop](/architecture/agent-loop)
**Brain** — State machine managing the think-act-observe cycle
- State transitions based on `stop_reason`
- Context management with message compaction
- MAX_TURNS limit (50) to prevent infinite loops

### [API Client](/architecture/providers)
**Communication** — Handles all LLM provider interactions
- Anthropic API with SSE streaming
- OpenAI-compatible endpoints
- Multi-provider support

### [Tool Registry](/architecture/tools)
**Capabilities** — Central registry for all available tools
- Built-in tools: Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch
- MCP integration for external tools
- Hook system for pre/post execution

### [Permission System](/architecture/design-philosophy)
**Safety** — Three-tier authorization with glob matching
- **ReadOnly**: File reading, globbing, grep
- **WorkspaceWrite**: File creation, modification
- **DangerFullAccess**: Shell execution, deletions

</div>

## Data Flow

```
User Input
    │
    ▼
┌─────────────┐
│    TUI      │ ◄── User interaction (keyboard input)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Agent Loop  │ ◄── Orchestrates: think → act → observe
└──────┬──────┘
       │
       ├──────────────────────┐
       ▼                      ▼
┌─────────────┐        ┌─────────────┐
│   API      │        │    Tool     │
│   Client   │        │  Registry   │
└──────┬──────┘        └──────┬──────┘
       │                      │
       ▼                      ▼
┌─────────────┐        ┌─────────────┐
│ LLM Model  │        │   System    │
│  (Remote)  │        │ (Files/Tools)│
└─────────────┘        └─────────────┘
       │                      │
       └──────────┬───────────┘
                  ▼
           ┌─────────────┐
           │   Output    │
           │  (Terminal) │
           └─────────────┘
```

## Key Components

| Component | Location | Purpose |
|-----------|----------|---------|
| **Entry Point** | `cmd/go-code/main.go` | Application startup |
| **Agent Core** | `internal/agent/agent.go` | Main agent loop logic |
| **Tool System** | `internal/tool/` | Built-in tool implementations |
| **Permissions** | `internal/permission/` | 3-tier permission system |
| **API Layer** | `internal/api/` | SSE streaming, API clients |
| **Config** | `internal/config/` | Multi-source configuration |
| **Session** | `internal/session/` | Persistence and resume |
| **Skills** | `internal/skills/` | Custom command system |

## Explore Further

- [Architecture Overview](/architecture/overview) — Full architectural documentation
- [Design Philosophy](/architecture/design-philosophy) — Intelligence vs Reliability
- [Agent Loop Deep Dive](/architecture/agent-loop) — State machine mechanics
- [Tool System](/architecture/tools) — Built-in tools and extensions
- [Permission Design](/architecture/design-philosophy) — Three-tier security model