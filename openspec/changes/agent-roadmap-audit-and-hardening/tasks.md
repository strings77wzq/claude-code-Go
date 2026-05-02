## 1. Baseline Evidence And Issue Register

- [ ] 1.1 Run `git status --short`, `git log -5 --oneline --decorate`, and record the current branch, latest release tag, and dirty files in this change's implementation notes.
- [ ] 1.2 Run `openspec list --json`, `openspec list --specs --json`, and `openspec validate --all --strict --json --no-interactive`; record active changes, spec count, and validation result.
- [ ] 1.3 Run `go test -count=1 ./...`, `go vet ./...`, `./scripts/run-harness.sh`, `cd docs && npm run build`, and `go build -o bin/go-code ./cmd/go-code`; record exact pass/fail evidence.
- [ ] 1.4 Run `./bin/go-code --version`, `./bin/go-code --help`, `./bin/go-code doctor --offline`, and one mock-provider prompt smoke; record expected exit codes and notable output.
- [ ] 1.5 Create `openspec/changes/agent-roadmap-audit-and-hardening/implementation-notes.md` with an issue register table: ID, priority, area, evidence, risk, target milestone, owner files, and verification command.
- [ ] 1.6 Populate the issue register with every problem listed in `design.md`, preserving the P0/P1/P2/P3 priority order used in this task file.
- [ ] 1.7 Add a "not implementing yet" section for deferred ideas such as IDE extension, desktop app, cloud/team collaboration, and marketing/showcase work.
- [ ] 1.8 Define milestone gates for v0.2.x stabilization, v0.3 extension productization, v0.4 agent quality, and v1.0 release trust.
- [ ] 1.9 Update `PARITY.md` with a link or reference to this roadmap change and mark this change as planning evidence, not feature verification.

## 2. P0 Runtime Correctness Stabilization

- [ ] 2.1 Add a failing unit test in `pkg/tui/tui_test.go` proving `runAgent` does not cancel the agent context before the returned `tea.Cmd` completes.
- [ ] 2.2 Fix `pkg/tui.runAgent` so timeout cancellation is owned by the returned command lifecycle, not by a `defer cancel()` that runs before command execution completes.
- [ ] 2.3 Add a TUI test proving user quit/cancel can stop a running request without leaking goroutines or leaving `isLoading` stuck.
- [ ] 2.4 Add a regression test in `internal/tool/tool_test.go` with a tool that panics during `Execute`.
- [ ] 2.5 Fix `internal/tool.Registry.Execute` so panic recovery returns `tool.Error("panic recovered: ...")` instead of discarding the recovered value.
- [ ] 2.6 Remove or replace the stale `RegisterBuiltinTools` TODO stub in `internal/tool/registry.go`; ensure only `internal/tool/init.RegisterBuiltinTools` remains as the real registration path.
- [ ] 2.7 Add tests proving the stale stub cannot accidentally be called or that the new replacement delegates correctly.
- [ ] 2.8 Add tests for unknown tool permission behavior in `internal/agent` so unknown tools never bypass permission or panic.
- [ ] 2.9 Run `go test -count=1 ./pkg/tui ./internal/tool ./internal/agent` and record results.

## 3. P0 Permission Modes And Non-Interactive Semantics

- [ ] 3.1 Define CLI flags and config keys for permission mode: `read-only`, `workspace-write`, `ask-all`, and `danger-full-access`.
- [ ] 3.2 Add `config.Config` fields and loader tests for permission mode from user config, project config, environment, and CLI override precedence.
- [ ] 3.3 Refactor `cmd/go-code` startup so permission policy construction is driven by config instead of hardcoded `WorkspaceWrite`.
- [ ] 3.4 Decide and document non-interactive prompt behavior: fail-fast deny by default, with explicit flags for approval automation if added.
- [ ] 3.5 Add prompt-mode tests proving write/bash/MCP approval requests do not block forever when stdin is unavailable.
- [ ] 3.6 Fix permission hierarchy so higher modes satisfy lower tool requirements unless an explicit deny rule blocks the operation.
- [ ] 3.7 Add policy tests for ReadOnly, WorkspaceWrite, AskAll, DangerFullAccess, deny-overrides-allow, and higher-mode-satisfies-lower-requirement.
- [ ] 3.8 Add CLI smoke tests for `--permission-mode read-only`, `--permission-mode workspace-write`, and invalid mode errors.
- [ ] 3.9 Update `docs/guide/quick-start.md`, `docs/api/config.md`, `docs/troubleshooting/permission-denied.md`, and Chinese equivalents with permission mode behavior.
- [ ] 3.10 Add harness scenarios for denied write, approved write, denied bash, allow-for-session reuse, and non-interactive denial.

## 4. P0 Session Continuity, Trace Versioning, And Replay Reliability

- [ ] 4.1 Define a versioned trace schema document under `docs/architecture/` or `internal/session/README.md` covering metadata, request, response, tool, permission, extension, error, message, and status events.
- [ ] 4.2 Add `schema_version` or equivalent version metadata to new session trace files while preserving replay compatibility with existing v0.2 traces.
- [ ] 4.3 Refactor `internal/agent.Agent` so an interactive TUI session can retain one session ID and trace file across multiple user turns.
- [ ] 4.4 Add tests proving prompt mode still creates a complete single-run session trace.
- [ ] 4.5 Add tests proving TUI multi-turn flow writes one coherent session trace with multiple user and assistant messages.
- [ ] 4.6 Add replay tests for old trace lines, new versioned lines, unknown future event types, and redaction of provider secrets.
- [ ] 4.7 Add trace events for model switch, permission mode, provider profile, command execution, and cancellation.
- [ ] 4.8 Improve `go-code replay --evidence` to summarize provider, model, permission mode, tool calls, denied operations, errors, and final status in a compact report.
- [ ] 4.9 Add a harness scenario that runs a multi-step prompt, then verifies replay evidence includes the expected session events.
- [ ] 4.10 Update `docs/guide/session-management.md`, `docs/architecture/agent-loop.md`, and Chinese equivalents.

## 5. P1 Entrypoint And App Composition Refactor

- [ ] 5.1 Introduce an internal app/bootstrap package or `cmd/go-code` helper that builds config, provider, tools, extensions, policy, agent, and UI from explicit options.
- [ ] 5.2 Keep `main()` as a thin dispatcher for subcommands and flag parsing.
- [ ] 5.3 Add constructor tests for app bootstrap with fake provider, fake MCP config, fake hooks path, fake skills path, and temporary session directory.
- [ ] 5.4 Ensure doctor can reuse bootstrap validation logic without starting a live agent request.
- [ ] 5.5 Ensure prompt mode, TUI mode, legacy REPL mode, doctor, replay, setup, version, and help remain behavior-compatible after the refactor.
- [ ] 5.6 Add integration tests for command dispatch without `os.Exit` by moving exit-code logic into testable functions.
- [ ] 5.7 Update architecture docs to show the new composition boundary and dependency direction.
- [ ] 5.8 Run `go test -count=1 ./cmd/go-code ./internal/config ./internal/provider ./internal/tool/... ./pkg/tui ./pkg/tty`.

## 6. P1 Provider Profiles And Model Compatibility

- [ ] 6.1 Design provider profile metadata separate from transport metadata: profile name, transport, default base URL, API key env var, compatible model patterns, compatibility status, and docs URL.
- [ ] 6.2 Implement provider profiles for Anthropic, OpenAI, DeepSeek, Qwen, GLM, Tencent Cloud Coding Plan, and MiMo without duplicating transport code.
- [ ] 6.3 Define compatibility levels: `verified`, `compatible`, `experimental`, `custom`, and `unsupported`.
- [ ] 6.4 Update config resolution so explicit provider profile names are validated before transport selection.
- [ ] 6.5 Reconcile `/model` switching with config unknown-model passthrough: verified registry models switch directly; explicit custom models require compatible provider/profile context.
- [ ] 6.6 Add tests for provider profile resolution, invalid provider/profile, invalid base URL, custom model with explicit profile, unknown model without profile, and provider/model mismatch.
- [ ] 6.7 Ensure runtime logs and doctor output name both provider profile and underlying transport.
- [ ] 6.8 Update `docs/api/config.md`, `docs/architecture/providers.md`, README provider table, and Chinese equivalents.
- [ ] 6.9 Add harness scenarios for at least Anthropic mock, OpenAI-compatible mock, and one China-friendly profile mock.
- [ ] 6.10 Record model registry update policy in docs: who updates names, what evidence is needed, and how deprecated aliases are handled.

## 7. P1 TUI Product Experience

- [ ] 7.1 Add a TUI status bar showing provider profile, model, permission mode, session ID, request state, and elapsed time.
- [ ] 7.2 Add a clear connection and request lifecycle indicator for queued, streaming, tool-running, waiting-for-permission, error, and completed states.
- [ ] 7.3 Add structured permission prompt UI for allow once, deny, allow for session, and cancel request.
- [ ] 7.4 Add keyboard-accessible transcript scrolling and stable rendering for long outputs.
- [ ] 7.5 Add model selector UI or command flow that uses provider profile metadata and explains unsupported model choices.
- [ ] 7.6 Add session picker/resume UI for recent sessions.
- [ ] 7.7 Add command palette or help overlay for `/help`, `/models`, `/sessions`, `/resume`, `/compact`, `/permissions`, `/replay`, and `/doctor`.
- [ ] 7.8 Add TUI tests for state transitions, permission prompt rendering, model switch rendering, session picker empty/non-empty states, and long-output rendering.
- [ ] 7.9 Add a screenshot or terminal-recording manual smoke note for the TUI after implementation.
- [ ] 7.10 Update README and quick-start docs so new users know whether to start with TUI or prompt mode.

## 8. P1 Harness Scenario Manifest And Agent Quality Gates

- [ ] 8.1 Create a machine-readable harness manifest mapping public workflows to scenario files, expected capabilities, and evidence owners.
- [ ] 8.2 Replace weak assertions such as "stdout length > 0" with scenario-specific assertions for tool calls, final answer content, permission decisions, and trace events.
- [ ] 8.3 Add edit correctness scenarios: unique edit, non-unique edit refusal, multi-file read-before-edit, and diff/replay evidence.
- [ ] 8.4 Add multi-step coding scenarios: inspect project, create small file change, run tests, handle failure, and summarize result.
- [ ] 8.5 Add recovery scenarios for rate limit, timeout, malformed stream, provider server error, and network error.
- [ ] 8.6 Add session scenarios for resume, replay latest, replay evidence, and redaction.
- [ ] 8.7 Add permission scenarios for read-only, workspace-write, ask-all, danger-full-access, deny, allow-once, and allow-for-session.
- [ ] 8.8 Add extension scenarios for MCP registration, MCP tool timeout, MCP denied action, LSP unavailable, LSP healthy mock, hooks warning/block, and invalid skills.
- [ ] 8.9 Add a local command `make verify-agent` or equivalent that runs Go tests, harness manifest scenarios, docs build, OpenSpec validation, doctor smoke, and binary smoke.
- [ ] 8.10 Update `PARITY.md` so every verified row names a manifest scenario or explicit unit test.

## 9. P1 MCP, LSP, Hooks, And Skills Productization

- [ ] 9.1 Add startup and call timeouts for MCP servers and tools; test slow startup, slow list-tools, slow call, and closed transport.
- [ ] 9.2 Surface MCP server status in doctor, TUI status, and replay evidence without requiring live external credentials.
- [ ] 9.3 Add MCP trust documentation explaining namespacing, permissions, process execution, environment variables, and secret handling.
- [ ] 9.4 Add user-facing LSP commands or tools for diagnostics, hover, definition, references, and symbols only when health gate passes.
- [ ] 9.5 Add tests proving LSP commands are hidden or return actionable unavailable messages before health passes.
- [ ] 9.6 Add hook execution feedback to trace/replay and user-facing output when a hook blocks or warns.
- [ ] 9.7 Add skill list/execute/help commands to the TUI shared command surface if they are not already fully exposed.
- [ ] 9.8 Add extension docs that clearly label each surface as verified, partial, experimental, or planned.
- [ ] 9.9 Add extension security review checklist and run it before marking any extension workflow verified.

## 10. P2 OpenSpec Hygiene And Task Discipline

- [ ] 10.1 Audit every `openspec/specs/*/spec.md` file for `TBD - created by archiving` Purpose text.
- [ ] 10.2 Replace each placeholder Purpose with a concrete purpose statement that matches the current requirement content.
- [ ] 10.3 Audit archived `tasks.md` files for unchecked items and record whether they were completed later, intentionally deferred, or incorrectly archived.
- [ ] 10.4 Add an `implementation-notes.md` or archive note for any archived change whose tasks remain unchecked but whose scope was superseded.
- [ ] 10.5 Define a task evidence convention: each completed task group must record tests run, files changed, and not-tested gaps.
- [ ] 10.6 Update OpenSpec authoring docs or project contributing docs with proposal/design/spec/tasks expectations.
- [ ] 10.7 Add a CI or script check that fails on new baseline specs containing placeholder Purpose text.
- [ ] 10.8 Add a CI or script check that reports active changes with apply-required artifacts missing.
- [ ] 10.9 Run `openspec validate --all --strict --json --no-interactive` after cleanup and record summary.

## 11. P2 Docs Truth, Generated Docs, And Onboarding

- [ ] 11.1 Decide whether `docs/.vitepress/dist` is committed source-of-truth or generated release artifact; document the policy.
- [ ] 11.2 If generated docs stay committed, add task guidance requiring docs source changes and dist changes to be reviewed separately.
- [ ] 11.3 If generated docs are not source-of-truth, update `.gitignore`, CI, and deploy workflow accordingly.
- [ ] 11.4 Audit README, docs, Chinese docs, `PARITY.md`, and release notes for claims that exceed harness or unit-test evidence.
- [ ] 11.5 Add a "try this agent against your own project" guide with safe permission mode, doctor, prompt mode, TUI, replay, and rollback steps.
- [ ] 11.6 Add a "compare against Codex/Claude-style agents" guide using deterministic tasks and subjective notes without unsupported superiority claims.
- [ ] 11.7 Add troubleshooting pages for stuck permissions, provider mismatch, TUI no output, MCP server hang, LSP unavailable, and replay missing sessions.
- [ ] 11.8 Run docs build and link checks; update docs evidence in implementation notes.

## 12. P2 Security And Trust Review

- [ ] 12.1 Create a security review checklist covering file read/write/edit boundaries, symlink escapes, shell validation, WebFetch/network access, MCP external processes, LSP URLs, hooks, skills, trace redaction, and config secrets.
- [ ] 12.2 Add tests for WebFetch URL validation, timeout, and redaction if coverage is missing.
- [ ] 12.3 Add tests for hook and skill files that try to expose secrets or execute unexpected behavior.
- [ ] 12.4 Add tests for MCP environment interpolation and secret redaction in trace/replay.
- [ ] 12.5 Add tests for workspace boundary behavior across read, write, edit, notebook edit, diff, tree, grep, and glob tools.
- [ ] 12.6 Add a security note to release verification requiring review for any change touching permission, bash, MCP, hooks, WebFetch, file tools, or config.
- [ ] 12.7 Run `go test -race ./internal/agent ./internal/tool/... ./internal/session ./pkg/tui` if feasible; record failures or not-tested gaps.

## 13. P2 Performance And Reliability Baselines

- [ ] 13.1 Define benchmark scenarios for large session history, large file read, large grep output, long tool output, MCP startup, and TUI long transcript rendering.
- [ ] 13.2 Add Go benchmarks or harness timing checks for agent request construction, compaction, trace append, replay parse, and tool output truncation.
- [ ] 13.3 Add baseline docs recording command, environment, date, and result for each benchmark.
- [ ] 13.4 Add regression thresholds only where stable enough; otherwise mark measurements as informational.
- [ ] 13.5 Add memory and goroutine leak checks around TUI request cancellation, MCP process close, and agent max-turn loops.
- [ ] 13.6 Update benchmark docs to avoid performance claims without reproducible commands.

## 14. P3 Agent Quality And Real-World Trial Workflows

- [ ] 14.1 Define three local trial projects: tiny fixture, medium Go package, and this repository itself.
- [ ] 14.2 Create workflow A: "explain architecture and identify risks" with expected evidence references and output quality rubric.
- [ ] 14.3 Create workflow B: "make a small safe code change and run tests" with expected file diff and verification command.
- [ ] 14.4 Create workflow C: "debug a failing test" with seeded failure, expected diagnosis, patch, and green test.
- [ ] 14.5 Create workflow D: "permission denial recovery" with denied shell/write operation and expected alternative behavior.
- [ ] 14.6 Create workflow E: "resume and replay" with interrupted session and expected replay summary.
- [ ] 14.7 Run the workflows manually through prompt mode and TUI; record usability notes, friction, and comparison notes against other agents without making unsupported public claims.
- [ ] 14.8 Convert stable workflows into deterministic harness scenarios where possible.
- [ ] 14.9 Update `PARITY.md` and docs with workflow status only after evidence exists.

## 15. Release Maturity And Install Verification

- [ ] 15.1 Add release candidate checklist covering version alignment, changelog, tags, GoReleaser config, checksums, docs build, OpenSpec validation, CI, CodeQL, and release assets.
- [ ] 15.2 Add local or CI checks for install scripts on Linux/macOS shell and PowerShell syntax where feasible.
- [ ] 15.3 Add smoke checks for release archives: download artifact, verify checksum, extract, run `go-code --version`, run `go-code doctor --offline`.
- [ ] 15.4 Add platform notes for Linux amd64/arm64, macOS amd64/arm64, Windows amd64/arm64.
- [ ] 15.5 Update `CHANGELOG.md` discipline so every release includes user-visible changes, verification commands, known gaps, and upgrade notes.
- [ ] 15.6 Add release rollback instructions for deleting a bad tag/release and publishing a patch release.
- [ ] 15.7 Run the release workflow on a test tag or patch tag only after main CI passes.

## 16. Final Verification And Archive Readiness

- [ ] 16.1 Run `gofmt -l .` and confirm it prints nothing.
- [ ] 16.2 Run `git diff --check` and confirm no whitespace errors.
- [ ] 16.3 Run `go vet ./...` and confirm success.
- [ ] 16.4 Run `go test -count=1 ./...` and confirm all packages pass.
- [ ] 16.5 Run `./scripts/run-harness.sh` and confirm all scenarios pass.
- [ ] 16.6 Run `cd docs && npm run build` and confirm success.
- [ ] 16.7 Run `openspec validate --all --strict --json --no-interactive` and confirm every spec/change passes.
- [ ] 16.8 Run `go build -o bin/go-code ./cmd/go-code` and confirm binary builds.
- [ ] 16.9 Run `./bin/go-code --version`, `./bin/go-code --help`, and `./bin/go-code doctor --offline`; record expected results.
- [ ] 16.10 Run all new milestone-specific harness scenarios and record scenario count and pass/fail summary.
- [ ] 16.11 Update `implementation-notes.md` with files changed, tests run, remaining risks, and not-tested gaps.
- [ ] 16.12 Update `PARITY.md`, roadmap docs, and release notes so public status matches completed evidence.
- [ ] 16.13 Run `openspec status --change agent-roadmap-audit-and-hardening` and confirm apply-required artifacts are complete.
- [ ] 16.14 Archive the change only after all tasks are complete and evidence is recorded.
