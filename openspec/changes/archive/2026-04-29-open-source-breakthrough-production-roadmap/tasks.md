## 1. Breakout Audit And Positioning

- [ ] 1.1 Audit README, docs, `PARITY.md`, code modules, and harness tests to list evidence-backed current breakout points.
- [ ] 1.2 Classify each claimed feature as supported, partial, planned, or experimental with evidence links.
- [ ] 1.3 Write a concise project positioning statement that leads with Go runtime, Python harness, local-agent safety, provider reach, and bilingual spec-driven engineering.
- [ ] 1.4 Update README opening sections to present the strongest verified differentiators before secondary feature lists.
- [ ] 1.5 Add or update a project audit/showcase page that explains what is already strong and what is not production-grade yet.

## 2. Maturity Model And Production Gates

- [ ] 2.1 Define demo, usable, reliable, production-grade, and community-leading maturity levels in docs.
- [ ] 2.2 Create a production-readiness checklist covering build, Go tests, Python harness, docs build, security-sensitive review, packaging, release notes, and known risks.
- [ ] 2.3 Add CI or script support for the production-readiness checks that can already be automated.
- [ ] 2.4 Add explicit not-tested and known-risk sections to release or changelog templates.
- [ ] 2.5 Update docs to prevent production-grade claims unless the checklist evidence exists.

## 3. Five-Minute Wow Path

- [ ] 3.1 Design the primary first-run flow from clean checkout to visible successful output.
- [ ] 3.2 Add or document an offline demo path that uses the mock provider or harness without live API keys.
- [ ] 3.3 Ensure the demo path includes visible streaming or agent output, at least one tool interaction, and a replay/trace inspection step.
- [ ] 3.4 Update quickstart docs in English and Chinese to separate offline demo, live provider setup, and troubleshooting.
- [ ] 3.5 Add regression tests or harness scenarios that prove the documented wow path stays runnable.

## 4. Doctor And Troubleshooting Experience

- [ ] 4.1 Audit current doctor/setup behavior against the documented first-run path.
- [ ] 4.2 Add doctor checks for config, provider profile, model, session path, tool availability, docs links, and harness/offline demo readiness.
- [ ] 4.3 Ensure doctor failures include actionable remediation text and links to the right docs.
- [ ] 4.4 Add tests for common doctor failure states.
- [ ] 4.5 Update troubleshooting docs with doctor-first workflows.

## 5. Showcase And Benchmark Evidence

- [ ] 5.1 Define at least one reproducible showcase scenario for safe file edit, repo analysis, permission denial, session replay, or provider switching.
- [ ] 5.2 Add scripts or docs that reproduce the showcase from a clean checkout.
- [ ] 5.3 Define benchmark methodology before publishing any headline benchmark numbers.
- [ ] 5.4 Add benchmark fixtures for startup time, mock prompt latency, harness runtime, and tool-call correctness where feasible.
- [ ] 5.5 Update `docs/showcase.md` and `docs/benchmark.md` so claims include commands, environment notes, scoring criteria, and limitations.

## 6. Community Contribution Flywheel

- [ ] 6.1 Define contributor lanes for docs, harness scenarios, provider profiles, built-in tools, permission policy, TUI polish, and release infrastructure.
- [ ] 6.2 Update contributing docs with lane-specific scope, acceptance criteria, and verification commands.
- [ ] 6.3 Add or update GitHub issue templates for bug reports, provider requests, harness scenarios, docs gaps, and feature proposals.
- [ ] 6.4 Add or update the PR template to ask for tests, harness runs, docs changes, security-sensitive surfaces, and known risks.
- [ ] 6.5 Update roadmap docs so major items link to specs, tasks, status, blockers, and maturity impact.
- [ ] 6.6 Add good-first-issue style task seeds where the repository supports it.

## 7. Release Packaging And Trust

- [ ] 7.1 Audit install paths for source build, `go install`, scripts, and release binaries.
- [ ] 7.2 Align release notes, changelog, GitHub about text, topics, and README badges with verified project status.
- [ ] 7.3 Add package/release verification notes for checksums or artifact validation where feasible.
- [ ] 7.4 Ensure security/privacy docs cover local tool execution, secrets, telemetry status, update behavior, and provider data flow.
- [ ] 7.5 Update English and Chinese docs to share the same support and maturity status.

## 8. Final Verification

- [ ] 8.1 Run `go test ./...` and record the result.
- [ ] 8.2 Run `./scripts/run-harness.sh` and record the result.
- [ ] 8.3 Run docs build or docs lint and record the result.
- [ ] 8.4 Verify README, docs, showcase, benchmark, roadmap, contributing docs, and templates no longer overclaim unsupported behavior.
- [ ] 8.5 Run `openspec status --change open-source-breakthrough-production-roadmap` and confirm all artifacts remain complete.
- [ ] 8.6 Summarize changed files, breakout points productized, production gates added, remaining risks, and not-tested gaps.
