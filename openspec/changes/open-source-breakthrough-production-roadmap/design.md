## Context

claudecode-go is trying to occupy a valuable position: an open-source, Go-native coding-agent runtime inspired by Claude Code and Claw Code, with Python used as a deterministic verification harness. That combination is already distinctive. Most projects in this space are either thin provider wrappers, UI demos, framework-heavy agent stacks, or language ports without a serious reliability story. This repository has the raw materials for a stronger story: a shippable Go binary, local tool execution, permissions, session/replay, provider abstraction, docs, and harness tests.

The current gap is packaging and production discipline. A strong open-source project needs two things at the same time: visible "why this matters" moments that make users want to try it, and boring production gates that make maintainers trust it. The project should therefore evolve on two tracks: a breakout track that makes the value obvious, and a production track that makes the value reliable.

This design intentionally complements `recenter-claudecodego-agent-roadmap`. That change corrects the technical direction and provider/model alignment. This change defines how the corrected direction becomes an excellent open-source product.

## Goals / Non-Goals

**Goals:**

- Name the project's current open-source "explosive points" and turn them into public product surfaces.
- Define a maturity model that prevents premature "production-grade" claims.
- Create a five-minute wow path that demonstrates real coding-agent value without live API spend.
- Establish production release gates for tests, harness, docs, security, performance, packaging, and known risks.
- Build a community flywheel around clear docs, demos, benchmark evidence, contributor tasks, issue hygiene, and release cadence.
- Keep positioning honest: strong enough to attract attention, grounded enough to retain trust.

**Non-Goals:**

- Implementing the full production roadmap in this proposal.
- Claiming feature parity with private Claude Code internals.
- Turning the project into a hosted SaaS, IDE product, or marketplace before the CLI agent is reliable.
- Adding vanity metrics, fabricated benchmarks, or unsupported testimonials.
- Replacing the technical recentering work required for DeepSeek/MiMo/provider correctness.

## Decisions

### Decision 1: Productize five standout points

The project should present five concrete breakout points:

1. Go-native single-binary coding agent.
2. Harness-first reliability with Python mock/evaluation infrastructure.
3. Real local-agent safety through permission, sandbox, session, and replay design.
4. China-friendly and global model support, including DeepSeek and MiMo alongside Anthropic/OpenAI-compatible APIs.
5. Bilingual, spec-driven open-source process that makes the project easy to audit and contribute to.

Alternatives considered:
- Market the project only as "Claude Code in Go". Rejected because it is too derivative and legally/technically vague.
- Lead with every feature in the README. Rejected because feature sprawl hides the project's strongest differentiation.

### Decision 2: Define maturity levels before claiming production readiness

The project should label maturity explicitly:

- Demo: runs in a narrow happy path.
- Usable: quickstart works and core workflows have basic tests.
- Reliable: harness verifies core trajectories and docs match behavior.
- Production-grade: release gates cover tests, harness, docs, security, permissions, packaging, and rollback notes.
- Community-leading: benchmarks, examples, contributor flow, roadmap discipline, and repeatable releases are visible.

Alternatives considered:
- Use a binary production/not-production label. Rejected because the project needs staged credibility.
- Avoid maturity language. Rejected because users need to understand risk before installing a local coding agent.

### Decision 3: Make the five-minute wow path offline-capable

The first impressive experience should not require paid tokens. A user should be able to build, run doctor, execute a mock-provider prompt, inspect a tool call, and replay the session locally. Live provider setup becomes the second path.

Alternatives considered:
- Start with real provider setup. Rejected because API keys, quota, and network issues create avoidable onboarding failures.
- Start with docs only. Rejected because an agent project needs executable proof quickly.

### Decision 4: Treat harness evidence as a public asset

Harness results, parity status, and demo scenarios should become part of the public project narrative. Instead of hiding tests as internal plumbing, the project should show how harness engineering makes a local coding agent safe and reliable.

Alternatives considered:
- Keep harness docs developer-only. Rejected because harness-first reliability is one of the project's strongest differentiators.
- Publish benchmark claims without reproducibility. Rejected because weak benchmarks damage trust.

### Decision 5: Build community around trustworthy tasks, not hype

The project should provide contributor lanes: first issue, docs sync, provider profile, harness scenario, tool improvement, permission policy, and TUI polish. Each lane should include acceptance criteria and verification commands.

Alternatives considered:
- Ask contributors to explore freely. Rejected because agent-runtime projects have safety and architecture boundaries.
- Open broad roadmap issues without tasks. Rejected because it creates noise and abandoned PRs.

### Decision 6: Keep public claims evidence-linked

README, docs, showcase, and release notes should link important claims to commands, tests, specs, parity rows, benchmark scripts, or known limitation notes.

Alternatives considered:
- Use aspirational marketing copy first. Rejected because users of a local code-editing agent need trust more than enthusiasm.
- Avoid marketing entirely. Rejected because a high-quality open-source project still needs clear positioning.

## Risks / Trade-offs

- [Risk] The project over-focuses on presentation and under-delivers runtime quality. -> Mitigation: every breakout surface must have a verification task and cannot bypass production gates.
- [Risk] "Production-grade" becomes too strict and delays release indefinitely. -> Mitigation: use maturity levels and publish honest partial status for early milestones.
- [Risk] Offline demo path diverges from real provider behavior. -> Mitigation: mock scenarios must mirror provider protocol contracts and be supplemented by optional live smoke notes.
- [Risk] Community growth increases support load. -> Mitigation: issue templates, doctor output, troubleshooting docs, and known-gap labels route problems clearly.
- [Risk] Benchmarks become misleading. -> Mitigation: publish methodology, fixtures, environment, pass/fail criteria, and limitations before headline numbers.
- [Risk] Bilingual docs drift. -> Mitigation: docs parity checks and release tasks require English/Chinese support status alignment.

## Migration Plan

1. Add a project audit page that lists current breakout points, missing production gates, and maturity level.
2. Create the offline five-minute wow path using existing build, doctor/offline checks, mock provider harness, and session replay.
3. Update README/docs to lead with the strongest evidence-backed differentiators and move uncertain claims into roadmap/known gaps.
4. Define release gates and CI checks for production-grade milestones.
5. Add community contribution lanes with issue/PR templates and verification commands.
6. Publish a reproducible benchmark/showcase format before promoting benchmark numbers.
7. Review maturity level at each release and record known risks in changelog/release notes.

Rollback is documentation and process based: if a gate or claim proves too strict or inaccurate, downgrade the maturity label and move the claim back to planned/partial until evidence exists.

## Open Questions

- Which maturity label should the next public release target: usable, reliable, or production-grade?
- Should the offline wow path use a dedicated `go-code demo` command or document the existing harness and replay commands?
- Which demo scenario best shows the project's value: safe file edit, repo analysis, permission denial, session replay, or provider switching?
- What benchmark should become the first public quality signal: harness pass rate, tool-call correctness, edit success, startup latency, or token/cost accounting?
- Should the project keep the current name for community launch, or adopt a clearer non-infringing product name first?
