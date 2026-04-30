## 1. Baseline And Scope Control

- [x] 1.1 Run and record current baseline for `go test ./...`, `make test`, docs build, OpenSpec validation, and any failing harness commands.
- [x] 1.2 Create `PARITY.md` with initial workflow matrix, status labels, tests, docs links, and known gaps.
- [x] 1.3 Audit README, docs, and website claims for unsupported features, stale metrics, placeholder testimonials, and implementation drift.
- [x] 1.4 Decide and document v0.2 mandatory workflows versus roadmap-only workflows.
- [x] 1.5 Add a short contributor-facing architecture map that identifies runtime package boundaries and ownership.

## 2. Runtime Health Check

- [x] 2.1 Add a `doctor` command entrypoint to the CLI command parser.
- [x] 2.2 Implement local checks for binary version, config file, environment variables, session directory, working directory, writable paths, and tool availability.
- [x] 2.3 Implement provider/model validation with offline/no-network skip support.
- [x] 2.4 Format doctor output with pass/fail/skip statuses and actionable remediation.
- [x] 2.5 Add tests for healthy config, missing API key, invalid session directory, and offline doctor mode.
- [x] 2.6 Document doctor usage in README and Chinese/English quick start pages.

## 3. Shared Command Surface

- [x] 3.1 Extract slash-command handling into a shared command service used by TUI and legacy REPL.
- [x] 3.2 Implement consistent `/help`, `/clear`, `/model`, `/models`, `/sessions`, `/resume`, `/compact`, `/update`, and `/permissions` behavior.
- [x] 3.3 Wire default Bubble Tea TUI to the shared command service.
- [x] 3.4 Keep legacy REPL behavior working through the same service.
- [x] 3.5 Add tests for command parsing, unknown commands, model switching, session listing, resume, and compact.
- [x] 3.6 Update command reference docs to match the implemented command surface.

## 4. Permission And Sandbox Flow

- [x] 4.1 Change default permission mode from full access to safe read-first behavior.
- [x] 4.2 Implement approval prompt decisions for allow, deny, allow-once, and allow-for-session.
- [x] 4.3 Ensure denied operations return structured tool results to the agent without executing side effects.
- [x] 4.4 Integrate file boundary checks, file size checks, binary checks, and symlink resolution into write/edit/read paths where appropriate.
- [x] 4.5 Integrate semantic bash validation before shell execution.
- [x] 4.6 Add permission audit entries to session trace without storing secrets.
- [x] 4.7 Add regression tests for read allow, write prompt, denied bash, remembered approval, path escape, and destructive command handling.
- [x] 4.8 Update permission-system docs and troubleshooting pages.

## 5. Provider And Model System

- [x] 5.1 Define a provider registry contract that validates provider, base URL, model, and API key source before requests.
- [x] 5.2 Fix runtime model switching so unsupported models are rejected without changing the active model.
- [x] 5.3 Normalize provider errors into auth, rate limit, timeout, server, network, invalid request, and unexpected categories.
- [x] 5.4 Add compatibility notes for Anthropic, OpenAI, DeepSeek, Qwen, GLM, and Tencent Cloud paths.
- [x] 5.5 Add tests for provider detection, OpenAI-compatible routing, model switch success, model switch failure, and error classification.
- [x] 5.6 Update provider and configuration docs.

## 6. Session Trace And Replay

- [x] 6.1 Normalize session JSONL schema for metadata, messages, tool calls, tool results, permission decisions, recovery events, and final status.
- [x] 6.2 Ensure interrupted or failed sessions are saved with recoverable context and explicit status.
- [x] 6.3 Implement reliable session listing and resume by ID across CLI/TUI surfaces.
- [x] 6.4 Add replay tooling to inspect or summarize a saved session without a real provider call.
- [x] 6.5 Add tests for save, list, resume, interrupted save, tool trace, permission trace, and replay output.
- [x] 6.6 Update session management and debugging docs.

## 7. Deterministic Parity Harness

- [x] 7.1 Inventory current Python harness scenarios and remove stale examples that no longer match the runtime.
- [x] 7.2 Add mock-provider scenarios for streaming text, read tool loop, edit tool loop, bash tool loop, permission denial, retryable provider error, context pressure, and malformed stream event.
- [x] 7.3 Add a one-command local harness runner from the repository root.
- [x] 7.4 Wire harness execution into CI for runtime-affecting changes.
- [x] 7.5 Publish harness logs in CI when a scenario fails.
- [x] 7.6 Link each parity matrix row to at least one test or explicit known gap.
- [x] 7.7 Document how contributors add new parity scenarios.

## 8. Extension Surface Productization

- [x] 8.1 Add configuration loading for MCP servers and register namespaced MCP tools through the existing registry.
- [x] 8.2 Apply permission policy consistently to MCP tool execution.
- [x] 8.3 Add tests for MCP server load failure, tool registration, and permission-gated MCP execution.
- [x] 8.4 Document hook lifecycle guarantees and add tests for pre-hook block and post-hook error behavior.
- [x] 8.5 Gate LSP commands/tools behind configured healthy LSP servers and add unavailable-state behavior.
- [x] 8.6 Validate skill files, list valid skills, report invalid skills non-fatally, and document the skill format.
- [x] 8.7 Update extension docs for MCP, hooks, LSP, and skills.

## 9. Documentation Product Experience

- [x] 9.1 Rewrite README around verified quick start, doctor, first prompt, core features, and honest status.
- [x] 9.2 Restructure Chinese docs into install, configure, doctor, usage, troubleshooting, architecture, extension, parity, roadmap, and contribution paths.
- [x] 9.3 Mirror critical Chinese docs in English or mark translation gaps explicitly.
- [x] 9.4 Replace aspirational showcase/testimonials with real examples, reproducible demos, or clearly labeled placeholders removed from public pages.
- [x] 9.5 Rewrite benchmark docs with methodology, commands, environment, date, and reproduction steps.
- [x] 9.6 Add architecture docs that map runtime flows to actual packages and interfaces.
- [x] 9.7 Add troubleshooting pages for API auth, provider routing, permission denial, session resume, command not found, and docs build issues.
- [x] 9.8 Run docs build and fix broken navigation or dead internal links.

## 10. Open Source And Release Readiness

- [x] 10.1 Ensure CI runs formatting, vet/static checks, Go tests, harness tests, docs build, and OpenSpec validation.
- [x] 10.2 Update issue templates, pull request template, and contribution guide with required test evidence.
- [x] 10.3 Add first-good-issue guidance and package-boundary notes for contributors.
- [x] 10.4 Review release workflow for platform binaries, checksums, release notes, and install instructions.
- [x] 10.5 Add changelog discipline for user-visible changes.
- [x] 10.6 Update roadmap with now, next, later, and not-planned sections.
- [x] 10.7 Add ownership/affiliation disclaimer where the project references Claude Code or Anthropic.

## 11. Final Verification

- [x] 11.1 Run `go test ./...` and record results.
- [x] 11.2 Run full repository test command and record results.
- [x] 11.3 Run parity harness and record scenario summary.
- [x] 11.4 Run docs build and record results.
- [x] 11.5 Run `openspec validate world-class-go-agent-rewrite --strict`.
- [x] 11.6 Run `go-code doctor` against a documented local configuration.
- [x] 11.7 Update `PARITY.md`, roadmap, and release notes with final status and remaining gaps.
