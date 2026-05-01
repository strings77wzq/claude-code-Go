## 1. Release State Cleanup

- [x] 1.1 Record current audit findings: missing `CLAUDE.md`/`task.md`, dirty generated docs, deleted `world-class-go-agent-rewrite`, completed active changes, and current green test evidence.
- [x] 1.2 Archive or explicitly park completed active changes: `v02-consolidation-release`, `short-mid-term-roadmap`, `polish-for-production`, `fix-ci-and-harness-tests`, `fix-python-harness-ide-errors`, `website-and-release-fix`, and `release-readiness-fix`.
- [x] 1.3 Decide whether proposal-only changes `enterprise-readiness` and `content-marketing` remain parked, get tasks, or move to archive.
- [x] 1.4 Separate source docs changes from generated `docs/.vitepress/dist` changes and document the release policy for generated artifacts.
- [x] 1.5 Run `openspec list` and `openspec validate --strict` for remaining active changes.

## 2. Extension Diagnostics Surface

- [x] 2.1 Audit current MCP, LSP, hooks, and skills initialization paths in `cmd/go-code`, `internal/tool/mcp`, `internal/lsp`, `internal/hooks`, and `internal/skills`.
- [x] 2.2 Extend `go-code doctor --offline` or an equivalent command to report MCP config, LSP config/health, hooks path, and skills path status.
- [x] 2.3 Add unit tests for extension diagnostics with no config, invalid config, and readable valid config fixtures.
- [x] 2.4 Ensure diagnostics never require provider API keys or real external network services.

## 3. MCP Productization

- [x] 3.1 Add a local MCP fixture or mock transport for deterministic tests.
- [x] 3.2 Add tests proving MCP tool registration creates namespaced registry entries.
- [x] 3.3 Add tests proving MCP write-like or external actions pass through the same permission policy as built-in tools.
- [x] 3.4 Record MCP invocation, permission decision, and tool result in the session trace.
- [x] 3.5 Update MCP docs with supported config format, unavailable states, permission behavior, and known limits.

## 4. LSP Productization

- [x] 4.1 Add fixture coverage for LSP unavailable, initialization success, and initialization failure states.
- [x] 4.2 Expose LSP availability through diagnostics or commands without breaking core prompt workflows.
- [x] 4.3 Add tests for diagnostics, symbols, definitions, references, and hover being advertised only after health checks pass.
- [x] 4.4 Record LSP unavailable/error outcomes in session trace as non-fatal events.
- [x] 4.5 Update LSP docs with setup, health checks, supported operations, and fallback behavior.

## 5. Hooks, Skills, Trace, and Replay

- [x] 5.1 Add tests proving invalid skills are reported as warnings while valid skills remain available.
- [x] 5.2 Add tests proving hook failures block only when configured to block.
- [x] 5.3 Extend replay output to include extension events, permission decisions, hook/skill warnings, and final status.
- [x] 5.4 Add redaction for API keys, authorization headers, and provider secrets in trace/replay output.
- [x] 5.5 Add a concise replay evidence mode suitable for release notes and issue reports.

## 6. Harness and Documentation Gate

- [x] 6.1 Add harness scenarios for MCP registration, MCP permission denial, LSP unavailable behavior, and replay evidence output.
- [x] 6.2 Update PARITY.md so MCP, LSP, hooks, skills, and replay statuses match the new evidence.
- [x] 6.3 Align README, roadmap, docs, and Chinese docs with verified v0.2 support and v0.3 planned/experimental scope.
- [x] 6.4 Remove or label benchmark, testimonial, enterprise, and content-marketing claims that lack implementation evidence.
- [x] 6.5 Add a v0.3 changelog draft with verification commands and known risks.

## 7. Final Verification

- [x] 7.1 Run `go test ./...` and confirm all packages pass.
- [x] 7.2 Run `./scripts/run-harness.sh` and confirm all scenarios pass.
- [x] 7.3 Run the docs build and confirm it exits successfully.
- [x] 7.4 Run `go-code doctor --offline` and confirm extension diagnostics appear.
- [x] 7.5 Run `openspec validate v03-extension-productization --strict`.
- [x] 7.6 Update the implementation notes with remaining risks and any manual smoke checks not covered by automation.
