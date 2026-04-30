## Context

claude-code-Go is positioned as a Go rewrite of a Claude Code-style coding agent. The repository already contains important ingredients: an agent loop, tool registry, permission policy, provider abstraction, session persistence, TUI/REPL surfaces, MCP/LSP modules, docs, and a Python harness. The problem is not lack of ideas; the problem is that the runtime surface, documentation, tests, and public claims are not yet aligned into a product-grade experience.

The target user is a developer who wants a dependable local coding agent, a contributor who wants a clear architecture to improve, and a learner who wants a readable reference implementation. Claw Code is a useful reference because it makes the health check, usage guide, parity tracking, roadmap, and project philosophy explicit. This project should borrow that operating shape while staying Go-native and dependency-light.

## Goals / Non-Goals

**Goals:**

- Make the default binary run a coherent product experience: setup, doctor, prompt mode, TUI, permissions, sessions, provider/model selection, and clear errors.
- Convert the project from demo-oriented to verification-oriented by requiring deterministic harness coverage and CI gates for the main workflows.
- Establish parity tracking against Claude Code-style workflows without claiming unsupported equivalence.
- Make permissions safe by default while still allowing power users to opt into higher automation.
- Keep Go code simple, explicit, and modular: small interfaces, no hidden global state, no placeholder modules exposed as complete features.
- Make Chinese and English docs trustworthy, task-oriented, and synchronized with implementation.
- Improve open-source impact by making first contribution, issue triage, releases, benchmarks, and roadmap easy to understand.

**Non-Goals:**

- Cloning private Claude Code internals or claiming affiliation with Anthropic.
- Building an IDE extension, desktop app, cloud service, or team collaboration layer in this change.
- Adding new runtime dependencies merely for polish.
- Optimizing for marketing claims before the product can pass its own doctor, tests, and parity harness.
- Solving every provider-specific API feature in one pass.

## Decisions

### Decision 1: Treat `go-code doctor` as the first product contract

The first reliable experience should be a health check that tells users whether the binary, config, provider, model, tools, filesystem permissions, session directory, and docs links are usable.

Alternatives considered:
- Keep troubleshooting in docs only. Rejected because users need machine-specific diagnosis.
- Make setup wizard do all validation. Rejected because doctor must be rerunnable in CI, bug reports, and support threads.

### Decision 2: Unify CLI, TUI, and legacy REPL around one command/application service layer

The CLI and TUI should call the same command handlers for `/model`, `/models`, `/sessions`, `/resume`, `/compact`, `/update`, `/permissions`, and `/help`. The current split lets the legacy REPL have capabilities that the default TUI only advertises.

Alternatives considered:
- Keep adding commands directly to each UI. Rejected because drift will continue.
- Remove legacy REPL immediately. Rejected until the TUI reaches feature parity and tests cover the command layer.

### Decision 3: Move from permissive default to explicit permission modes

The default mode should be safe for a new user: read and inspect by default, ask before writes or shell execution, and remember explicit session choices. `DangerFullAccess` remains available but must be a deliberate opt-in.

Alternatives considered:
- Keep `DangerFullAccess` default for smoother demos. Rejected because it undermines the main reliability claim.
- Disable shell by default. Rejected because shell is central to coding-agent usefulness; it should be governed, not removed.

### Decision 4: Define parity as a tracked matrix, not a vague claim

The project should maintain `PARITY.md` with explicit workflows: prompt mode, streaming, tool use, edits, bash, permissions, sessions, resume, compact, provider switching, MCP, and recovery. Each row should have status, tests, docs, and known gaps.

Alternatives considered:
- Compare only in marketing docs. Rejected because it is not actionable.
- Chase full Claude Code parity immediately. Rejected because scoped, verified parity is more credible.

### Decision 5: Use deterministic mock-provider harness as the release gate

The harness should simulate provider responses, streaming, tool calls, recoverable failures, permission prompts, and context pressure. CI must run both Go unit tests and harness scenarios that prove end-to-end behavior.

Alternatives considered:
- Rely on real API smoke tests. Rejected because they are slow, costly, flaky, and hard for contributors.
- Keep only unit tests. Rejected because agent bugs commonly occur across API, history, tools, permissions, and UI boundaries.

### Decision 6: Productize extension surfaces after core stability

MCP, hooks, skills, and LSP are valuable differentiators, but they should be surfaced through configuration, permission integration, tests, and docs only after the core doctor/CLI/TUI path is reliable.

Alternatives considered:
- Lead with plugin marketplace ideas. Rejected because ecosystem work needs a stable host.
- Remove extension modules. Rejected because they are part of the long-term differentiation.

### Decision 7: Make documentation honest and executable

Docs must describe verified behavior, include commands users can run, and avoid unsupported testimonials, stale metrics, or screenshots that imply unavailable features. Chinese docs should be first-class, not a partial translation.

Alternatives considered:
- Keep aspirational docs and catch implementation up later. Rejected because trust is more important than breadth.
- Focus only on README. Rejected because the project needs onboarding, architecture, troubleshooting, parity, and contribution paths.

## Risks / Trade-offs

- [Risk] The change is broad and could become a rewrite without checkpoints. -> Mitigation: split implementation into phases: baseline verification, command surface, permissions, harness, docs, release.
- [Risk] Safer permissions may feel less smooth. -> Mitigation: provide remembered decisions, clear prompts, and explicit trusted-workspace modes.
- [Risk] Parity tracking may expose gaps publicly. -> Mitigation: frame gaps as a transparent roadmap and link each gap to tests/tasks.
- [Risk] Provider abstraction may hide API differences. -> Mitigation: document compatibility levels and keep provider-specific adapters explicit.
- [Risk] Docs work can outrun implementation again. -> Mitigation: require docs pages to link to tested commands or known status markers.
- [Risk] Harness complexity may slow contributors. -> Mitigation: provide one-command local test scripts and short scenario authoring docs.

## Migration Plan

1. Establish baseline: run current tests, docs build, and harness; document failures as tasks before feature work.
2. Add doctor and shared command service; wire CLI, TUI, and legacy REPL to it.
3. Tighten permission defaults and implement approval prompt flow with regression tests.
4. Build deterministic harness scenarios for core workflows and require them in CI.
5. Productize provider/model configuration and session trace/replay.
6. Bring MCP/LSP/hooks/skills into the documented extension surface where tested.
7. Rewrite docs and website claims to match verified behavior.
8. Prepare release readiness: changelog, benchmark methodology, templates, parity status, and contributor guide.

Rollback is straightforward per phase because the work should be split into small, reversible changes. If a phase destabilizes the runtime, keep the prior command surface and disable new behavior behind explicit config until tests pass.

## Open Questions

- Should the final binary name remain `go-code`, or should the project adopt a more distinct product name before wider launch?
- Which Claude Code-style workflows are mandatory for v0.2, and which remain roadmap items?
- Should the Python harness remain the primary parity harness, or should selected scenarios move to Go integration tests for contributor simplicity?
- Should telemetry remain in scope, be disabled indefinitely, or be removed until there is a clear privacy and value story?
- What level of Chinese/English documentation parity is required before public promotion?
