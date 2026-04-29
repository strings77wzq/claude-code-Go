## Why

claudecode-go already has several rare ingredients for an influential open-source coding-agent project, but those ingredients need to be turned into visible product breakthroughs and production-grade gates. This change defines what the current "explosive points" are, how to package them honestly, and what must evolve before the project can credibly be called complete, excellent, and production-ready.

## What Changes

- Identify and productize the project's strongest open-source breakthroughs:
  - Go-native Claude Code-style agent runtime that can ship as one small binary.
  - Python deterministic harness that proves behavior without paid API calls.
  - Provider abstraction ready for DeepSeek, MiMo, Anthropic, OpenAI-compatible APIs, and future profiles.
  - Permission/sandbox/session/replay architecture that matches real coding-agent reliability needs.
  - Bilingual docs and OpenSpec-driven engineering discipline that can attract contributors.
- Define a production-grade maturity model with explicit levels: demo, usable, reliable, production-grade, and community-leading.
- Add release gates for build/test/harness/docs/security/performance before public "production-grade" claims.
- Add an open-source growth flywheel: quickstart success, demos, benchmark evidence, issue templates, contributor paths, roadmap clarity, and showcase scenarios.
- Add a user-facing "wow path" that demonstrates meaningful coding-agent value in under five minutes without requiring live provider spend.
- Define how the project should evolve into a complete product: stable CLI/TUI, trusted permissions, provider profiles, session/replay, deterministic evaluation, docs parity, release packaging, and community operations.
- No runtime feature is marked complete by this proposal alone; this proposal creates the specs and tasks that make completeness measurable.

## Capabilities

### New Capabilities

- `oss-breakout-positioning`: Defines the project's current standout points, public narrative, demo path, and differentiation as an open-source Claude Code-style Go agent.
- `production-readiness-maturity`: Defines measurable production-readiness levels, release gates, reliability/security/performance requirements, and evidence required for production-grade claims.
- `developer-experience-wow-path`: Defines the first-run, quickstart, doctor, demo, prompt, replay, and harness path that makes the project feel smooth and impressive within minutes.
- `community-release-flywheel`: Defines contribution workflows, release process, issue/PR templates, roadmap hygiene, benchmark/showcase artifacts, and community trust mechanisms.

### Modified Capabilities

- None. This is a new productization and production-readiness roadmap layered on top of the existing active implementation changes.

## Impact

- Affected runtime surfaces: `cmd/go-code`, `pkg/tui`, `pkg/tty`, provider/profile configuration, permission flow, session/replay commands, and doctor/setup commands.
- Affected verification surfaces: `go test ./...`, `./scripts/run-harness.sh`, docs build, CI workflows, security checks, benchmark scripts, and release checks.
- Affected docs: `README.md`, `PARITY.md`, `docs/`, `docs/zh/`, `docs/showcase.md`, `docs/benchmark.md`, troubleshooting, architecture docs, and contributor docs.
- Affected community surfaces: GitHub topics/about text, issue templates, PR templates, release notes, changelog, roadmap, demo assets, and good-first-issue labels.
- Relationship to `recenter-claudecodego-agent-roadmap`: that change corrects project direction and model/harness alignment; this change turns the aligned direction into visible open-source impact and production-grade operating standards.
