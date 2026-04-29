# Parity Matrix

This matrix tracks Claude Code-style workflow parity for `claude-code-Go`. Status labels are intentionally conservative:

- `verified`: implemented and covered by automated tests or reproducible smoke checks.
- `partial`: implemented in some path, but missing coverage, consistency, or product polish.
- `planned`: accepted scope for the current roadmap, not yet implemented.
- `unsupported`: explicitly out of scope for now.

## Current Status

| Workflow | Status | Evidence | Known Gaps |
| --- | --- | --- | --- |
| Single prompt mode | verified | `go-code -p` exists in `cmd/go-code/main.go`; `harness/test_scenarios.py` runs it with `-q -f json`. | Broader JSON schema can still be improved. |
| Streaming text | verified | `harness/test_scenarios.py::TestStreamingText` passes against the mock Anthropic SSE server. | Real provider smoke tests remain manual. |
| Agent tool loop | verified | Harness read/bash scenarios cover `tool_use` -> tool result -> final answer loops. | Edit loop coverage is defined but still light. |
| Built-in file tools | verified | Unit tests cover read/edit/write-adjacent behavior in `internal/tool/builtin`. | Workspace policy integration is incomplete. |
| Bash tool | verified | Unit tests cover timeout/truncation/destructive blocking; harness covers bash tool roundtrip and denial scenario. | Interactive approval UX still needs TUI polish. |
| Permission prompts | verified | Agent/unit tests cover deny/remembered approvals; harness covers permission denial without executing the command. | Approval mode flags are still not exposed as CLI options. |
| Default TUI commands | partial | TUI supports basic input, `/help`, `/clear`, `/model`, and exit. | TUI advertises commands it does not implement yet. |
| Legacy REPL commands | partial | REPL supports sessions, resume, compact, update, skills, models. | Command behavior is not shared with TUI. |
| Sessions | partial | JSONL save/load tests pass. | Resume is not unified across UI surfaces; trace schema is not complete. |
| Context compaction | partial | Compaction unit tests pass. | User-facing status and replay/debug story are thin. |
| Multi-provider support | partial | Anthropic and OpenAI-compatible adapters exist. | Model validation and compatibility levels are not explicit enough. |
| MCP extension | planned | MCP manager and adapter code exists. | Not productized in default config/docs/permission flow. |
| LSP integration | planned | LSP client package exists. | Not exposed as a stable user workflow. |
| Doctor command | planned | OpenSpec requirement exists. | No CLI command yet. |
| Deterministic parity harness | verified | `./scripts/run-harness.sh` builds `bin/go-code` and runs `pytest harness/`; CI uploads harness logs on failure. | More scenarios should be added as features stabilize. |
| IDE extension | unsupported | Roadmap mentions future IDE integration. | Not part of current implementation scope. |
| Cloud/team collaboration | unsupported | Future concept only. | Not part of current implementation scope. |

## v0.2 Mandatory Workflows

The next credible public milestone should verify these workflows before promotion:

1. `go-code doctor` reports actionable local readiness.
2. `go-code -p "..."` returns visible text and JSON output in deterministic harness tests.
3. Streaming mock-provider scenarios pass.
4. Read and bash tool roundtrips pass through the harness.
5. Safe default permission mode asks before writes and shell execution.
6. TUI and legacy REPL share slash-command behavior for help, model, sessions, resume, compact, update, and permissions.
7. Session save, list, resume, and replay have automated coverage.
8. README and Chinese quick start match tested commands.

## Roadmap-Only Workflows

These remain important, but should not block v0.2:

- Stable MCP extension marketplace story.
- LSP-powered code intelligence commands.
- IDE extension.
- Desktop app.
- Remote/cloud agent.
- Team collaboration dashboard.
