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
| Phase 2 | 🔄 In Progress | Enhanced capabilities - Skills, multi-provider, session resume |
| Phase 3 | 🔮 Planned | Advanced features - IDE integration, collaboration, cloud |

---

## Phase 1: Foundation ✅ Complete

The first phase established the core infrastructure for the AI coding assistant.

### Completed Features

- **Agent Loop** — Autonomous "think → act → observe" cycle with stop_reason dispatch
- **6 Built-in Tools** — Read, Write, Edit, Glob, Grep, Bash
- **Permission System** — Three-tier model with rule-based matching and session memory
- **MCP Integration** — Model Context Protocol with stdio transport and JSON-RPC
- **SSE Streaming** — Real-time token-by-token response with custom parser
- **Session Persistence** — Save and restore conversation state
- **Hooks System** — Pre/post execution callbacks for extensibility

---

## Phase 2: Enhanced Capabilities 🔄 In Progress

Phase 2 adds more powerful features to improve usability and flexibility.

### Current Work

- **Skills System** — Custom commands and reusable workflows (e.g., `/review-pr`, `/deploy`)
- **Multi-Provider Support** — Anthropic, OpenAI, and any OpenAI-compatible API
- **Session Resume** — Load previous conversations and continue seamlessly
- **Enhanced Tools** — Diff view, tree visualization, web fetching
- **Auto-Update** — Automatic version checking and updates

### Upcoming in Phase 2

- Improved documentation and examples
- Performance optimizations
- Additional tool enhancements

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

Want to influence the roadmap? See our [Community](/community) page for how to contribute ideas, report issues, and submit pull requests.

---

*Last updated: April 2026*