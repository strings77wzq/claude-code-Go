---
title: Roadmap
description: claude-code-Go development roadmap

---

# Roadmap

Here's the development roadmap for claude-code-Go, organized into three phases.

## Overview

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ Complete | Core foundation - Agent loop, tools, permission, MCP, SSE, sessions |
| Phase 2 | ✅ Complete | Enhanced capabilities - Skills, multi-provider, session resume |
| Phase 3 | 🔮 Planned | Advanced features - IDE integration, collaboration, cloud |

---

## Phase 1: Foundation ✅ Complete

The first phase established the core infrastructure for the AI coding assistant.

### Completed Features

- **Agent Loop** — Autonomous "think → act → observe" cycle with stop_reason dispatch
- **9 Built-in Tools** — 6 core (Read, Write, Edit, Glob, Grep, Bash) + 3 enhanced (Diff view, Tree, WebFetch)
- **Permission System** — Three-tier model with rule-based matching and session memory
- **MCP Integration** — Model Context Protocol with stdio transport and JSON-RPC
- **SSE Streaming** — Real-time token-by-token response with custom parser
- **Session Persistence** — Save and restore conversation state
- **Hooks System** — Pre/post execution callbacks for extensibility

---

## Phase 2: Enhanced Capabilities ✅ Complete

Phase 2 adds more powerful features to improve usability and flexibility.

### Completed Features

- **Skills System** — Custom commands and reusable workflows (e.g., `/review-pr`, `/deploy`)
- **Multi-Provider Support** — Anthropic, OpenAI, and any OpenAI-compatible API
- **Session Resume** — Load previous conversations and continue seamlessly
- **Enhanced Tools** — Diff view, tree visualization, web fetching
- **Manual Compaction** — `/compact` command to reduce context size
- **Auto-Update** — `/update` command for version checking and updates

---

## Phase 3: Advanced Features 🔮 Planned

Future plans to bring more advanced capabilities and better integration.

### Planned Features

- **VS Code Extension** — Native IDE integration with rich UI
- **Desktop Application** — Standalone GUI with full feature set
- **Team Collaboration** — Shared workflows, team dashboards, permissions
- **Cloud Agent** — Remote execution, API endpoints, monitoring
- **Plugin Marketplace** — Community-created skills and extensions

### Vision

Phase 3 aims to make claude-code-Go a full-fledged development environment, bridging the gap between CLI and IDE while maintaining the simplicity and performance of Go.

---

## Contributing

Want to influence the roadmap? See our [Feedback](/feedback) page for how to report issues and submit feature requests.

---

## Feature Comparison

| Feature | Phase 1 | Phase 2 | Claude Code |
|---------|---------|---------|-------------|
| Agent Loop | ✅ | ✅ | ✅ |
| Built-in Tools | 6 core | 9 total | 20+ |
| Permission System | ✅ | ✅ | ✅ |
| MCP Integration | ✅ | ✅ | ✅ |
| SSE Streaming | ✅ | ✅ | ✅ |
| Session Persistence | ✅ | ✅ | ✅ |
| Session Resume | ❌ | ✅ | ✅ |
| Skills System | ❌ | ✅ | ✅ |
| Multi-Provider | ❌ | ✅ | ❌ |
| Auto-Update | ❌ | ✅ | ✅ |
| Hooks System | ❌ | ✅ | ✅ |
| IDE Integration | ❌ | ❌ | ✅ |

---

*Last updated: April 2026*