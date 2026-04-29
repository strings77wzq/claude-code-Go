---
title: Runtime Package Map
description: Package-level map for contributors working on claude-code-Go runtime behavior.
---

# Runtime Package Map

This page maps product workflows to the Go packages that own them. Use it before changing runtime behavior so fixes land in the right layer.

## Runtime Flow

```text
cmd/go-code
  -> config loader
  -> provider registry
  -> tool registry
  -> permission policy
  -> agent loop
  -> TUI or REPL
```

## Package Ownership

| Package | Responsibility | Notes |
| --- | --- | --- |
| `cmd/go-code` | Process entrypoint, flags, setup, runtime wiring. | Keep orchestration here; avoid embedding business logic in `main`. |
| `pkg/tui` | Default Bubble Tea interactive UI. | Should delegate slash commands to a shared command layer. |
| `pkg/tty` | Legacy line-based REPL. | Useful reference for command coverage until TUI reaches parity. |
| `internal/agent` | Agent loop, history, compaction, recovery, session trace hooks. | Owns stop-reason dispatch and tool-result feedback. |
| `internal/api` | Anthropic request/response types, HTTP client, SSE parsing, API error classification. | Keep provider-neutral behavior out of this package when possible. |
| `internal/provider` | Provider interface and provider adapters. | Anthropic and OpenAI-compatible differences belong here. |
| `internal/config` | Configuration loading and defaults. | Should expose config source details for doctor output. |
| `internal/tool` | Tool interface, registry, and built-in tool registration. | Tool definitions must stay concise and schema-valid. |
| `internal/tool/builtin` | Built-in local tools such as read, edit, bash, tree, diff, web fetch, todo, notebook. | Side effects must be permission-gated by the caller. |
| `internal/permission` | Permission policy, rules, prompter, file boundaries, bash semantic checks. | This package decides; agent/tool layers must enforce. |
| `internal/session` | JSONL session persistence and trace line helpers. | Needs a normalized schema for replay and debugging. |
| `internal/hooks` | Pre/post tool execution hooks. | Pre-hook failures should block execution. |
| `internal/skills` | User-defined reusable prompts/workflows. | Invalid skill files should not break startup. |
| `internal/tool/mcp` | MCP client, transport, and adapter. | MCP tools must enter the same registry and permission flow. |
| `internal/lsp` | LSP client features. | Optional surface; should fail closed when no server is configured. |
| `internal/update` | Release checking and binary update. | User-triggered network behavior only. |
| `internal/cost` | Token and cost tracking helpers. | Provider pricing assumptions must be documented. |
| `internal/telemetry` | Optional local telemetry primitives. | Keep disabled until product/privacy story is explicit. |

## Contributor Rules

- UI packages should not duplicate command semantics.
- Tool implementations should perform local validation, but permission approval belongs above tool execution.
- Provider adapters should normalize errors before they reach UI code.
- Tests that require sockets, home-directory writes, or built binaries must document those requirements.
- Public docs should only claim workflows that are verified by tests or listed in `PARITY.md`.
