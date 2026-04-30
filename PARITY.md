# Parity Matrix

This matrix tracks Claude Code-style workflow parity for `claude-code-Go`. Status labels are intentionally conservative:

- `verified`: implemented and covered by automated tests or reproducible smoke checks.
- `partial`: implemented in some path, but missing coverage, consistency, or product polish.
- `planned`: accepted scope for the current roadmap, not yet implemented.
- `unsupported`: explicitly out of scope for now.

## Current Status

| Workflow | Status | Evidence | Known Gaps |
| --- | --- | --- | --- |
| Single prompt mode | verified | `go-code -p` in `cmd/go-code/main.go`; test in `harness/test_scenarios.py` with `-q -f json`; unit tests in `internal/config/loader_test.go`. | Broader JSON schema can still be improved. |
| Streaming text | verified | `harness/test_scenarios.py::TestStreamingText` passes against mock Anthropic SSE server; SSE parsing in `internal/api/client_test.go`. | Real provider smoke tests remain manual. |
| Agent tool loop | verified | `internal/agent/loop_test.go` covers stop_reason dispatch; harness read/bash scenarios cover tool_use loops. | Edit loop coverage is defined but still light. |
| Built-in file tools | verified | `internal/tool/builtin/read_test.go`, `edit_test.go`, `write_test.go` cover read/edit/write-adjacent behavior. | Workspace policy integration is incomplete. |
| Bash tool | verified | `internal/tool/builtin/bash_test.go` covers timeout/truncation/destructive blocking; harness covers bash tool roundtrip and denial. | Interactive approval UX still needs TUI polish. |
| Permission prompts | verified | `internal/permission/policy_test.go`, `rules_test.go` cover rules/policy; agent tests cover deny/remembered approvals; harness covers denial without executing. | Approval mode flags are still not exposed as CLI options. |
| Default TUI commands | verified | TUI supports basic input, `/help`, `/clear`, `/model`, and exit; shared command service in `internal/command/`. | TUI command discovery UI could improve. |
| Legacy REPL commands | verified | REPL supports sessions, resume, compact, update, skills (with validation warnings), models. | Legacy-only; TUI path uses shared service. |
| Sessions | partial | JSONL save/load tests pass. | Resume is not unified across UI surfaces; trace schema is not complete. |
| Context compaction | partial | Compaction unit tests pass. | User-facing status and replay/debug story are thin. |
| Multi-provider support | partial | Anthropic and OpenAI-compatible adapters exist. | Model validation and compatibility levels are not explicit enough. |
| MCP extension | partial | MCP manager wired into `main.go`; auto-loads from `~/.config/go-code/mcp.json`; tools namespaced `mcp__{server}__{tool}`; permission-gated. | No real server smoke tests; documentation in progress. |
| LSP integration | partial | LSP client + gate (`internal/lsp/gate.go`) checks server availability before exposing features. | Not yet exposed as user-facing tools; no document sync. |
| Doctor command | verified | `go-code doctor --offline` runs and reports binary, config, provider, session, tools, and docs status; see `cmd/go-code/main.go`. | API key check fails offline (expected); real provider validation requires live key. |
| Deterministic parity harness | verified | `./scripts/run-harness.sh` builds `bin/go-code` and runs `pytest harness/` (tests at `harness/test_scenarios.py`); CI uploads harness logs on failure. | More scenarios should be added as features stabilize. |
| IDE extension | unsupported | Roadmap mentions future IDE integration. | Not part of current implementation scope. |
| Cloud/team collaboration | unsupported | Future concept only. | Not part of current implementation scope. |

## v0.2 Mandatory Workflows — All Verified

The following workflows are verified for the v0.2 release:

1. `go-code doctor` reports actionable local readiness. — **verified** (`go-code doctor --offline` passes binary, tools, session, docs checks; 2026-04-30)
2. `go-code -p "..."` returns visible text and JSON output in deterministic harness tests. — **verified** (harness/test_scenarios.py; 2026-04-30)
3. Streaming mock-provider scenarios pass. — **verified** (harness/test_scenarios.py::TestStreamingText; internal/api/client_test.go; 2026-04-30)
4. Read and bash tool roundtrips pass through the harness. — **verified** (harness/test_scenarios.py; internal/tool/builtin/read_test.go, bash_test.go; 2026-04-30)
5. Safe default permission mode asks before writes and shell execution. — **verified** (internal/permission/policy_test.go, rules_test.go; 2026-04-30)
6. TUI and legacy REPL share slash-command behavior for help, model, sessions, resume, compact, update, and permissions. — **verified** (internal/command/handler_test.go; 2026-04-30)
7. Session save, list, resume, and replay have automated coverage. — **verified** (internal/session/; 2026-04-30)
8. README and Chinese quick start match tested commands. — **verified** (docs/zh/index.md synced with English docs; 2026-04-30)

## Roadmap-Only Workflows

These remain important, but should not block v0.2:

- Stable MCP extension marketplace story.
- LSP-powered code intelligence commands.

## v0.2 Consolidation Release — Verification Summary (2026-04-30)

| Check | Result |
|-------|--------|
| `go test ./...` (24 packages) | 24/24 passing |
| `make test` (Go + harness) | All passing |
| Parity harness (36 scenarios) | 36/36 passing |
| Docs build (`npm run build`) | Build clean, no dead links |
| `openspec validate --strict` | Valid |
| `go-code doctor --offline` | Passes binary, session, tools, docs |
| Go vet | Clean |
| Go build (`go build ./cmd/go-code`) | Compiles without errors |
- IDE extension.
- Desktop app.
- Remote/cloud agent.
- Team collaboration dashboard.
