## Context

claudecode-go is intended to become a polished Go rewrite of a Claude Code-style coding agent, with Python used for deterministic harness engineering rather than production runtime. The repository already contains valuable pieces: a Go CLI, TUI/TTY surfaces, agent loop, provider abstraction, tool registry, permission policy, session replay, trace logging, docs, OpenSpec artifacts, and a Python harness. The main problem is alignment: current docs and specs describe a broad future, while the user-facing product needs a smaller, sharper path that runs, tests, and explains itself.

The reference projects are Claw Code and Claude Code's observable workflows. Claw Code is especially useful for the product shape: doctor-first onboarding, explicit usage docs, parity tracking, a roadmap, and a harness-centered implementation story. Claude Code should be treated as a workflow and UX reference unless legally usable source code is present in this repository; the project must not depend on private implementation details or claim private-source compatibility.

External model support must be updated from vague "OpenAI-compatible" language to concrete, tested families. DeepSeek's official API docs currently list `deepseek-v4-flash` and `deepseek-v4-pro`, with OpenAI and Anthropic-compatible base URLs. The older `deepseek-chat` and `deepseek-reasoner` names are compatibility aliases with a published deprecation date. MiMo's official pages identify `mimo-v2.5-pro` as the model tag and describe the V2.5 series as long-context, agentic-coding oriented, and intended for harnesses such as Claude Code/OpenCode/Kilo. Those facts should drive config, tests, and docs.

## Goals / Non-Goals

**Goals:**

- Restore a clear product identity: Go-first coding agent runtime, Python-first harness and evaluation infrastructure.
- Make DeepSeek and MiMo-V2.5 explicit first-class model families in provider configuration, docs, and tests.
- Define harness engineering as a release gate, not an optional demo layer.
- Reframe Claude Code / Claw Code parity as observable workflows with test evidence: prompt mode, streaming, tool calls, edit/bash safety, permissions, sessions, replay, recovery, and model switching.
- Preserve the good existing work: provider registry, doctor, permission policy, session replay, Python mock server, trace logs, docs site, and OpenSpec task discipline.
- Reduce drift by requiring docs/spec claims to map to runnable commands, passing tests, or an explicit "planned" status.
- Produce an implementation roadmap that is small enough to apply in phases and strict enough to improve code quality.

**Non-Goals:**

- Copying private Claude Code source, branding, or proprietary internals.
- Building every Claude Code feature before the core prompt/tool/session/provider path is reliable.
- Treating MCP, LSP, hooks, skills, or plugin ecosystems as launch blockers for the recentered roadmap.
- Adding new runtime dependencies for polish before the existing Go architecture is made coherent.
- Hiding provider differences behind a misleading "all OpenAI-compatible APIs work" claim.
- Rewriting the whole repository in one pass.

## Decisions

### Decision 1: Keep Go as the shipped runtime and Python as the harness

The binary users run remains Go. Go owns CLI/TUI, agent loop, provider calls, tools, permissions, config, sessions, logs, replay command, and release packaging. Python owns mock provider scenarios, deterministic parity tests, trace analysis, and evaluators.

Alternatives considered:
- Move more runtime behavior into Python. Rejected because the open-source value proposition is a Go-native agent binary.
- Rewrite the harness in Go immediately. Rejected because Python's FastAPI/pytest ecosystem is already productive for mock servers and scenario evaluation.

### Decision 2: Treat providers as named compatibility profiles, not only transports

DeepSeek and MiMo must be represented as explicit profiles even if they use an OpenAI-compatible transport. A profile includes provider name, base URL, model IDs, model capabilities, environment variable conventions, known limitations, and harness scenarios.

Alternatives considered:
- Keep only `openai` and ask users to configure any compatible base URL manually. Rejected because it makes DeepSeek/MiMo support invisible, poorly tested, and easy to misconfigure.
- Fork separate provider implementations prematurely. Rejected unless the API shape requires it; start with explicit profiles over the existing transport where possible.

### Decision 3: Update DeepSeek model truth before adding features

The registry and docs must prefer `deepseek-v4-flash` and `deepseek-v4-pro`. `deepseek-chat` and `deepseek-reasoner` may remain as legacy aliases only if marked deprecated with migration guidance.

Alternatives considered:
- Leave old model names because they still work. Rejected because public docs should not guide new users toward deprecated names.
- Remove old names immediately. Rejected until migration notes and tests prevent avoidable breakage.

### Decision 4: Add MiMo-V2.5 as a first-class target with verified API assumptions

MiMo support should begin with `mimo-v2.5-pro`, config examples, docs, and harness coverage. If the public API platform exposes an OpenAI-compatible endpoint, use the shared transport with a MiMo profile. If the API differs, add a narrow adapter after verifying official documentation or sample responses.

Alternatives considered:
- Document MiMo only as "custom OpenAI-compatible". Rejected because the user explicitly wants MiMo-V2.5 series support and the model is positioned for agentic coding harnesses.
- Hardcode undocumented endpoint details. Rejected because model integration must be based on official docs or user-provided credentials/examples.

### Decision 5: Make harness scenarios the gate for public claims

Every user-facing "supported" claim must be backed by at least one of: Go unit test, Python harness scenario, docs build check, or manual verification note in `PARITY.md`. The release gate should run `go test ./...` and `./scripts/run-harness.sh`.

Alternatives considered:
- Treat harness as optional examples. Rejected because agent failures usually occur across provider, tool, permission, session, and UI boundaries.
- Require live API tests in CI. Rejected because they are costly, flaky, and unsuitable for community contributors.

### Decision 6: Recenter docs around executable truth

README and docs should start with a working install/build path, `go-code doctor`, a prompt run, provider config for Anthropic/DeepSeek/MiMo, harness verification, parity status, and known gaps. Aspirational architecture pages should be marked roadmap until backed by code.

Alternatives considered:
- Keep broad docs and gradually catch implementation up. Rejected because it damages trust and makes the project harder to evaluate.
- Delete most docs. Rejected because the project needs strong bilingual onboarding and architecture clarity.

### Decision 7: Defer extension ecosystem emphasis until the core is smooth

MCP, LSP, hooks, skills, and update surfaces remain valuable, but the immediate roadmap should put them behind core reliability: prompt mode, tool use, permissions, sessions, model support, and harness coverage.

Alternatives considered:
- Lead with ecosystem extensibility to attract contributors. Rejected because a plugin surface without a stable agent core creates more support burden.
- Remove extension modules. Rejected because some pieces already exist and are useful once documented honestly.

## Risks / Trade-offs

- [Risk] The reset becomes another broad roadmap without implementation traction. -> Mitigation: tasks are ordered around runnable checkpoints and every phase must report verification evidence.
- [Risk] MiMo API details are less publicly documented than DeepSeek's. -> Mitigation: support only documented model tags and verified endpoint behavior; keep unknowns explicit.
- [Risk] Explicit provider profiles duplicate configuration logic. -> Mitigation: keep transport reuse internal while making user-facing profiles clear and testable.
- [Risk] Safer documentation may look less impressive in the short term. -> Mitigation: honest supported/planned status increases contributor trust and reduces failed onboarding.
- [Risk] Existing active OpenSpec changes overlap with this reset. -> Mitigation: this change acts as a recentering contract; subsequent apply work should reconcile or supersede overlapping tasks rather than stacking unreviewed scope.
- [Risk] Harness requirements can slow feature velocity. -> Mitigation: use small, deterministic scenarios and keep one-command local execution.

## Migration Plan

1. Run baseline verification and record current truth: `go test ./...`, `./scripts/run-harness.sh`, docs build if available, `go-code doctor --offline`, and current CLI prompt path.
2. Update provider registry and config semantics for explicit `anthropic`, `openai`, `deepseek`, and `mimo` profiles while preserving shared transports where compatible.
3. Add DeepSeek model metadata for `deepseek-v4-flash` and `deepseek-v4-pro`, legacy alias warnings, env examples, and harness tests.
4. Add MiMo-V2.5 metadata for `mimo-v2.5-pro`, config examples, compatibility notes, and mock/harness coverage based on verified API behavior.
5. Rework `PARITY.md`, README, and Chinese docs to mark each capability as supported, partial, planned, or removed, with test evidence links.
6. Tighten the product roadmap so each next feature has a spec, code task, harness case, and doc task.
7. Re-run the full verification gate and update OpenSpec task status only after evidence passes.

Rollback is phase-based: provider profile changes should preserve existing `openai` and `anthropic` configs; docs can be reverted independently; harness scenarios can be skipped only with a clear `Not-tested` note until fixed.

## Open Questions

- What is the final public binary/project name before launch: `go-code`, `claudecode-go`, or a distinct non-infringing name?
- Does the user have legally usable Claude Code source in the local workspace, or should all parity remain based on observable behavior and public docs?
- Which MiMo API endpoint, auth header, and streaming format should be considered canonical if the public API platform requires login-only docs?
- Should Tencent Coding Plan remain a first-class profile, a custom Anthropic-compatible example, or out of scope for the recentered roadmap?
- What is the minimum v0.1 public launch gate: passing harness only, docs parity, or signed release artifacts as well?
