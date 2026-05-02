## Context

Unit tests and build checks prove implementation correctness, but they do not answer whether the agent behaves well on realistic coding tasks. The roadmap calls for a quality gate that captures task success, latency, tool behavior, permission decisions, and replayable evidence. This change turns the harness into an implementation-independent product gate.

## Goals / Non-Goals

**Goals:**
- Make scenario definitions explicit and reviewable.
- Capture enough evidence to debug failures and compare agent behavior over time.
- Track latency and tool-use budgets for critical workflows.
- Provide a release gate that can run in CI without live provider credentials for deterministic scenarios.

**Non-Goals:**
- Benchmark proprietary agents by calling their services automatically.
- Replace unit, integration, or security tests.
- Optimize provider response quality in this change.
- Build a public leaderboard.

## Decisions

1. Use manifest-driven scenarios.
   - Each scenario declares workspace setup, task prompt, allowed tools, expected file/test outcomes, trace assertions, and budgets.
   - Alternative rejected: embedding scenarios only in Go test code, because reviewers need to inspect task intent and release evidence.

2. Separate deterministic gates from live-agent trials.
   - CI uses deterministic fixtures and mock provider/tool paths where possible.
   - Live trials are documented and produce evidence artifacts but are not required for every PR.
   - Alternative rejected: requiring live model credentials in CI, because it is flaky and expensive.

3. Treat latency as a product signal.
   - Scenario output records total duration, provider wait when available, tool count, and timeout cause.
   - Alternative rejected: pass/fail only, because slow success can still make the agent unusable.

4. Make comparison workflows adapter-based.
   - The harness can record normalized evidence from this agent and manually supplied external-agent runs.
   - Alternative rejected: building direct competitor integrations now, because it creates maintenance and credential risks.

## Risks / Trade-offs

- [Risk] Scenario assertions can become brittle. → Mitigation: assert user-visible outcomes and trace invariants rather than exact model prose.
- [Risk] Latency budgets vary by machine and provider. → Mitigation: keep deterministic CI budgets focused on local harness overhead and mark live budgets as advisory.
- [Risk] Evidence artifacts can leak secrets. → Mitigation: reuse trace redaction from runtime safety work.

## Migration Plan

1. Add manifest schema and parser tests.
2. Convert existing shallow harness cases into manifests.
3. Add trace and latency collection to the runner.
4. Add release-quality scenarios for inspect/edit/test/recover/explain workflows.
5. Document live trial procedure for comparing with other coding agents.
