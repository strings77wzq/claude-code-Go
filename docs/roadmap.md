---
title: Roadmap
description: claude-code-Go development roadmap — now, next, later, and not-planned
---

# Roadmap

This roadmap is kept current with the project. Status is verified against tests and parity harness results.

## Now (v0.2)

Shipped and verified. See [PARITY.md](https://github.com/strings77wzq/claude-code-Go/blob/main/PARITY.md) for evidence links.

| Feature | Status |
|---------|--------|
| Agent Loop (think → act → observe) | Verified |
| 11 Built-in Tools | Verified |
| Permission System (3-tier) | Verified |
| Doctor Health Check | Verified |
| Multi-Provider (Anthropic, OpenAI-compatible) | Verified |
| Session Persistence + Resume | Verified |
| Session Replay | Verified |
| Slash Commands | Verified |
| Skills System | Verified |
| Hooks System | Verified |
| Deterministic Parity Harness | Verified |
| Go Test Suite (20 packages, all passing) | Verified |

## Next (v0.3)

Active extension productization. Code paths exist; public support depends on diagnostics, tests, docs, and harness evidence.

| Feature | Status |
|---------|--------|
| MCP Integration | Partial — config, namespacing, permission gate, docs, and harness evidence in progress |
| LSP Integration | Partial — health gate, capability gating, doctor diagnostics, docs, and harness evidence in progress |
| Replay evidence mode | Partial — extension events, permission decisions, redaction, and concise evidence output covered by tests |
| Extension docs + configuration | In progress |
| Enhanced permission memory | Planned |
| Improved TUI command coverage | Planned |

## Later

Features recognized as valuable but not yet scheduled.

| Feature | Notes |
|---------|-------|
| IDE Extension (VS Code) | Requires stable core first |
| Plugin/Skill Marketplace | Depends on skill format stabilization |
| Team Collaboration | Requires multi-session architecture |
| Cloud Agent | Requires auth and tenant model |

## Not Planned

Features explicitly out of scope for the foreseeable future.

| Feature | Reason |
|---------|--------|
| Desktop Application | Focus is CLI/TUI + eventual IDE |
| Mobile App | Not a target platform |
| Proprietary/Closed Features | Project is MIT-licensed open source |

## Feature Comparison

Comparison with Claude Code (as of April 2026). Status reflects verified implementation, not planned work.

| Feature | go-code | Claude Code |
|---------|---------|-------------|
| Agent Loop | Verified | Yes |
| Built-in Tools | 11 | 20+ |
| Permission System | Verified | Yes |
| Session Persistence | Verified | Yes |
| Session Resume | Verified | Yes |
| Streaming (SSE) | Verified | Yes |
| Multi-Provider | Verified | No (Anthropic-only) |
| MCP Integration | Partial v0.3 | Yes |
| LSP Integration | Partial v0.3 | No |
| Skills System | Verified | Yes |
| Hooks System | Verified | Yes |
| IDE Integration | Later | Yes |
| Cloud/Team | Later | Yes |

## Contributing

See [CONTRIBUTING.md](https://github.com/strings77wzq/claude-code-Go/blob/main/CONTRIBUTING.md) to get started. First-good-issues are tagged in the issue tracker.

---

*Last updated: May 2026*
