## Context

The repository now contains a runnable Go agent released as v0.2.0. The current architecture is intentionally local-first:

- `cmd/go-code` owns process entrypoints, configuration loading, provider selection, built-in/MCP tool registration, permission policy construction, agent construction, and TUI/legacy REPL selection.
- `internal/agent` owns the model loop, history, tool execution, hooks, permission checks, compaction, recovery, and JSONL trace/session writes.
- `internal/provider` owns Anthropic and OpenAI-compatible transports, model registry, error classification, and runtime model switching.
- `internal/tool` owns the tool interface, registry, built-ins, and MCP adapter/manager.
- `internal/permission` owns mode/rule/session-memory decisions and terminal prompts.
- `internal/session` owns JSONL trace, replay, redaction, and session listing/loading.
- `internal/lsp`, `internal/hooks`, and `internal/skills` are extension surfaces guarded by tests and doctor diagnostics.
- `pkg/tui` and `pkg/tty` are user-facing shells over the shared command service in `internal/command`.
- `harness` provides deterministic Python scenarios against local mock providers.
- `openspec/specs` is the canonical requirement surface, but many specs still carry archived placeholder Purpose text and old archived task files contain unchecked items.

The project is in a useful but fragile state: the release is usable, CI is green, and the parity matrix is conservative, but the next phase needs to prevent drift between architecture, specs, docs, tasks, and code.

### Current strengths to preserve

- Single Go binary with no runtime language dependency for users.
- Deterministic mock-provider harness and CI gates.
- Explicit permission model with file boundary and bash validation tests.
- Session trace/replay and redaction foundation.
- Honest `PARITY.md` labels for verified, partial, planned, and unsupported workflows.
- Bilingual documentation and OpenSpec-driven change history.
- Release automation through GitHub Actions and GoReleaser.

### Potential problems found

1. **Entrypoint composition is too dense.** `cmd/go-code/main.go` performs config, provider resolution, tool registration, MCP bootstrap, permission setup, agent setup, and UI startup in one function. This makes policy injection and integration tests harder than necessary.
2. **TUI request lifecycle has a cancellation hazard.** `pkg/tui.runAgent` creates a timeout context and defers `cancel()` before returning a `tea.Cmd`; the returned command may observe an already-cancelled context or a race with request startup.
3. **Non-interactive permission prompts can block or degrade.** Prompt mode uses stdin prompter paths even when stdin may be unavailable; permission behavior needs explicit non-interactive policy semantics.
4. **Permission mode is not user-configurable at the CLI.** The code hardcodes `WorkspaceWrite`, and harness comments document that permission scenarios cannot fully exercise prompt denial without a mode flag or injected policy.
5. **Permission hierarchy is inconsistent.** `meetsModeRequirement` exists but the active tool requirement logic compares equality, which can deny a higher privilege mode when a lower requirement is configured.
6. **Provider profile and transport are conflated.** User-facing providers such as DeepSeek, Qwen, GLM, Tencent, and MiMo are mostly represented through `openai`/`anthropic` transport choices and hard-coded model heuristics.
7. **Model registry is static and drift-prone.** The project needs a documented model/provider compatibility level and update process instead of treating all registry entries as equally verified.
8. **Runtime model switching is narrow.** `/model` rejects unknown models, while config resolution still permits unknown models with inferred provider; the UX needs a coherent distinction between verified registry models and explicit custom models.
9. **Agent sessions are per `Run` call.** TUI multi-turn conversations keep history, but session IDs and trace files are recreated per run, making one interactive session hard to reason about in replay.
10. **Trace schema is still too loose.** JSONL uses multiple internal line structs and map-based extension lines; it needs versioning, stable event names, and stricter replay compatibility tests.
11. **Tool panic recovery is ineffective.** `Registry.Execute` defers a recovery closure but does not assign the recovered error to the returned result, so a panic path may not produce the intended tool error result.
12. **There is a stale duplicate built-in registration stub.** `internal/tool/registry.go` still contains `RegisterBuiltinTools` with a TODO, while real registration lives under `internal/tool/init`.
13. **MCP lifecycle is basic.** MCP startup is eager, server failures are logged but not surfaced to the UI, there is no timeout per server startup/list/call, and external server trust boundaries are not strongly documented in runtime UX.
14. **LSP is mostly diagnostic, not productized.** Health gates and client methods exist, but user-facing commands/tools for diagnostics, hover, definition, references, and symbols are not yet complete.
15. **Hooks and skills are not first-class in the TUI.** They exist and have tests, but discovery, execution feedback, and failure reporting are still mostly diagnostic.
16. **TUI usability is minimal.** It lacks robust cancellation, transcript navigation, permission dialogs, model/provider selector, session picker, token/cost/status panels, and actionable error affordances.
17. **Legacy REPL and TUI parity is incomplete.** They share command handling, but the default experience and command affordances still differ.
18. **Harness scenarios are useful but shallow.** Many current tests validate "some output exists"; core agent quality needs scenario manifests, richer assertions, edit correctness, multi-step planning, recovery, and approval flows.
19. **OpenSpec hygiene is inconsistent.** Many archived specs still have `TBD` Purpose text, and at least one archived task file remains unchecked despite later work completing similar scope.
20. **Task evidence discipline is inconsistent.** Completed task checkboxes often do not include exact verification commands, logs, or not-tested notes.
21. **Docs and generated docs policy is unclear.** `docs/.vitepress/dist` frequently changes as generated output; the repository needs a clear source-vs-generated publication rule.
22. **Release maturity is not yet v1.0.** CI is green, but release gates need stronger end-to-end smoke checks, install-script checks, checksum validation, and platform-specific startup evidence.
23. **Security review is not separated as a gate.** Permission, MCP, shell, WebFetch, file boundaries, and secrets redaction should have an explicit security review checklist and tests.
24. **Performance baselines are not captured.** Context compaction, large tool output, long sessions, MCP startup, and TUI responsiveness need measurable baselines before optimization claims.
25. **Comparison against Codex/Claude-style agents lacks structured evaluation.** The project needs user-journey scenarios that compare actual task success, speed, safety, and recoverability rather than only feature parity.

## Goals / Non-Goals

**Goals:**

- Produce a detailed implementation roadmap that is ready to execute with `/opsx:apply`.
- Convert the current audit into staged, testable work with clear acceptance criteria.
- Prioritize defects that directly affect whether the author can use the agent for real coding work.
- Preserve the working v0.2 release while enabling v0.3/v0.4 development without destabilizing core flows.
- Strengthen OpenSpec and task discipline so future changes are reviewable, scoped, and evidence-backed.
- Define verification gates for each milestone, not just final release.

**Non-Goals:**

- Do not immediately implement the roadmap in this change.
- Do not claim parity with Codex, Claude Code, or any other agent until deterministic scenarios prove it.
- Do not add new third-party dependencies unless a later implementation task explicitly justifies them.
- Do not rewrite the entire agent loop before isolating lifecycle, permission, provider, and trace defects with tests.
- Do not treat marketing pages, showcase content, or enterprise ideas as blockers for the next runnable agent milestone.

## Decisions

### Decision 1: Use staged hardening instead of a broad rewrite

The next phase SHALL be split into stabilization, extension productization, agent quality, and release maturity stages. This avoids repeating previous large archived changes where many surfaces were touched at once.

Alternative considered: start a "v1 rewrite" that restructures the full runtime immediately. Rejected because the current system is already usable and tested; a rewrite would destroy useful evidence and make regressions harder to isolate.

### Decision 2: Treat user-trial workflows as the primary product metric

The roadmap SHALL start from workflows a developer actually tries:

1. install/build
2. doctor/configure provider
3. run one prompt
4. open TUI
5. read/search/edit files
6. approve/deny risky tools
7. resume/replay/debug a session
8. compare usefulness against other agents

Alternative considered: continue feature parity lists first. Rejected because feature inventory does not prove that the agent feels useful or safe in real coding work.

### Decision 3: Split bootstrap wiring from runtime capabilities

`cmd/go-code` SHOULD move toward a small command dispatcher plus app builder. Provider, permission, tools, MCP, hooks, skills, session, and UI construction SHOULD have testable constructors with explicit options.

Alternative considered: leave `main.go` as the only composition root. Rejected because it keeps mode flags, smoke tests, and policy injection tightly coupled to process startup.

### Decision 4: Make permission mode explicit before expanding autonomy

CLI/TUI config SHALL expose safe permission modes and non-interactive behavior before adding more autonomous agent behaviors. Permission prompts, denials, remembered approvals, and audit trace must be deterministic.

Alternative considered: default to broad access to make demos look stronger. Rejected because local coding agents fail trust tests quickly when safety boundaries are unclear.

### Decision 5: Define provider profiles separately from transports

The project SHOULD introduce user-facing provider profiles that map to reusable transports. For example, DeepSeek/Qwen/GLM/MiMo can use OpenAI-compatible transport while retaining profile-specific defaults, docs, errors, and compatibility status.

Alternative considered: keep all compatible providers as `LLM_PROVIDER=openai` only. Rejected because user setup, docs, and model compatibility become confusing.

### Decision 6: Version the session trace schema

Trace lines SHALL include schema/version discipline and stable event names before replay becomes a debugging contract. Unknown event handling should remain forward-compatible.

Alternative considered: continue extending ad hoc JSONL maps. Rejected because replay and release evidence will become brittle.

### Decision 7: Make OpenSpec hygiene a first-class task group

OpenSpec specs and tasks SHALL be cleaned as project artifacts, not just scaffolding. Placeholder Purpose text, archived incomplete task state, and missing evidence notes should be repaired.

Alternative considered: ignore archived specs because current validation passes. Rejected because the user explicitly asked to review spec/task design, and future agents will rely on these files.

### Decision 8: Use harness scenarios as release contracts

Each supported agent workflow SHALL map to at least one deterministic scenario with meaningful assertions. Coverage percentages remain useful, but scenario success is the stronger product signal.

Alternative considered: rely on Go unit tests and manual testing only. Rejected because provider/tool/permission/session behavior is cross-process and needs executable end-to-end evidence.

## Risks / Trade-offs

- **Scope creep** → Mitigation: keep milestone gates explicit and do not start v0.4 quality work until v0.2.x stabilization tasks pass.
- **Overfitting to mock harness** → Mitigation: add optional manual/live smoke notes while keeping CI deterministic and credential-free.
- **Provider docs drift** → Mitigation: mark model/profile entries with compatibility levels and require source/evidence notes when changing registry defaults.
- **Permission UX friction** → Mitigation: provide profiles such as read-only, workspace-write, ask-all, and danger-full-access, with clear TUI indicators.
- **Trace schema migration breaks old sessions** → Mitigation: add versioned replay tests for old and new trace lines.
- **TUI fixes become a rewrite** → Mitigation: first fix lifecycle/cancel/session correctness, then add visible UX affordances.
- **Generated docs churn pollutes reviews** → Mitigation: decide whether generated `dist` is committed; if yes, isolate it in release/docs commits.
- **Security-sensitive MCP expansion** → Mitigation: require explicit server trust status, startup/call timeouts, namespacing, and permission trace coverage.

## Migration Plan

1. Keep v0.2.0 as the stable release baseline.
2. Create and implement v0.2.x stabilization tasks first: TUI cancellation, panic recovery, CLI permission mode, session continuity, and stale stub cleanup.
3. Run the full release verification gate after each stabilization batch.
4. Productize v0.3 extension surfaces only after core prompt/tool/session flows remain green.
5. Expand harness scenarios and comparison workflows before claiming stronger agent quality.
6. Archive this change only after tasks are complete, evidence is recorded, and relevant specs are updated.

Rollback strategy: each task group should be independently revertible. Avoid global rewrites; prefer small commits with tests that preserve existing public commands.

## Open Questions

- Should generated `docs/.vitepress/dist` remain committed, or should releases publish docs from source only?
- Should v0.3 use explicit provider profiles such as `LLM_PROVIDER=deepseek`, or keep compatible providers under `openai` with profile metadata?
- Should non-interactive mode default to deny all prompts, fail fast, or support `--permission-mode` plus `--yes/--no` automation flags?
- Should TUI become the only interactive surface after parity is achieved, or should legacy REPL remain supported long-term?
- Which real-world coding benchmark should become the first "usefulness against Codex/Claude-style agents" scenario?
