---
title: Runtime Package Map
description: Package-level map for contributors working on claude-code-Go runtime behavior.
---

# Runtime Package Map

This page maps product workflows to the Go packages that own them. Use it before changing runtime behavior so fixes land in the correct layer.

## Runtime Flow

```text
cmd/go-code  (entrypoint, setup, doctor, replay)
  -> config loader          (internal/config)
  -> provider resolution    (internal/provider/registry)
  -> provider adapter       (internal/provider/anthropic | internal/provider/openai)
  -> tool registry          (internal/tool, internal/tool/init)
  -> MCP tools              (internal/tool/mcp)
  -> permission policy      (internal/permission)
  -> agent loop             (internal/agent)
  -> TUI or REPL            (pkg/tui | pkg/tty, internal/command)
```

## Package Ownership

| Package | Responsibility | Notes |
| --- | --- | --- |
| `cmd/go-code` | Process entrypoint, flags, subcommands (doctor, replay, setup), runtime wiring. | Keep orchestration here; subcommands live in the same package but must not import each other. |
| `pkg/tui` | Default Bubble Tea interactive TUI. | Should delegate slash commands to `internal/command`. |
| `pkg/tty` | Legacy line-based REPL. | Reference for slash command coverage until TUI reaches parity. |
| `internal/agent` | Agent loop, history, compaction, recovery, session trace hooks. | Owns stop-reason dispatch and tool-result feedback. |
| `internal/api` | Anthropic request/response types, HTTP client, SSE parsing, API error classification. | Keep provider-neutral behavior out of this package when possible. |
| `internal/provider` | Provider interface definition and normalized error types (`ClassifiedError`). | The `Provider` interface lives here; implementations live in sub-packages. |
| `internal/provider/registry` | Provider resolution, model-to-provider detection, base URL defaults. | Single place for routing `"anthropic"` vs `"openai"` provider strings at startup. |
| `internal/provider/anthropic` | Anthropic Messages API adapter. | Implements `Provider` from `internal/provider`. |
| `internal/provider/openai` | OpenAI Chat Completions API adapter. | Implements `Provider` from `internal/provider`. |
| `internal/config` | Configuration loading, defaults, env vars, CLI overrides. | Load order: defaults < user config < project config < env vars < CLI flags. |
| `internal/tool` | Tool interface, registry, and tool definition types. | Registry dispatches tool calls by name. |
| `internal/tool/init` | Registers all builtin tools (Bash, Read, Write, Edit, Glob, Grep, Diff, Tree, WebFetch, TodoWrite, Notebook) into the registry. | Single point for tool registration. |
| `internal/tool/builtin` | Built-in tool implementations. | Side effects must be permission-gated by the caller. Individual files: `bash.go`, `read.go`, `write.go`, `edit.go`, `glob.go`, `grep.go`, `diff.go`, `tree.go`, `webfetch.go`, `todo.go`, `notebook.go`. |
| `internal/tool/mcp` | MCP client, transport, and adapter. | MCP tools enter the same registry and permission flow as builtins. |
| `internal/permission` | Permission modes, policy evaluation, rules, prompter, file boundaries, bash semantic checks. | Evaluation order: deny rules > allow rules > session memory > tool requirement > mode default. |
| `internal/session` | JSONL session persistence, list, replay. | Needs a normalized schema for cross-UI replay. |
| `internal/command` | Shared slash command handler (`/help`, `/model`, `/clear`, `/sessions`, `/compact`, etc.). | Both TUI and legacy REPL route through the same handlers. |
| `internal/hooks` | Pre/post tool execution hooks. | Pre-hook failures should block tool execution. |
| `internal/skills` | User-defined reusable prompts/workflows: loader, registry, types. | Invalid skill files must not break startup. |
| `internal/lsp` | LSP client features (hover, definition, references, symbols, diagnostics, gate). | Optional surface; should fail closed when no server is configured. |
| `internal/update` | Release checking and binary download. | User-triggered network behavior only. |
| `internal/cost` | Token and cost tracking. | Provider pricing assumptions must be documented. |
| `internal/telemetry` | Optional local telemetry primitives and consent. | Keep disabled until product/privacy story is explicit. |
| `internal/logger` | Structured logging (slog) and trace-level helpers. | Used across all packages for consistent log output. |

## Contributor Rules

- UI packages (`pkg/tui`, `pkg/tty`) must not duplicate command semantics; use `internal/command` instead.
- Tool implementations perform local validation, but permission approval belongs above tool execution (in `internal/permission`).
- Provider adapters normalize errors via `ClassifyError` / `ClassifyHTTPStatus` before they reach UI code.
- Tests that require sockets, home-directory writes, or built binaries must document those requirements in their test functions.
- Public docs should only claim workflows that are verified by tests or listed in `PARITY.md`.
