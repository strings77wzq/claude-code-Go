## Context

The repository has reached a credible v0.2 technical baseline: `go test ./...` passes, the Python harness reports 36/36 passing, and the v0.2 OpenSpec artifacts validate. The remaining risk is coordination debt: completed changes are still active, a removed `world-class-go-agent-rewrite` change is visible in the working tree, some docs generated output is dirty, and older audit documents conflict with newer parity evidence.

The next development increment should not start with broad enterprise or marketing scope. The highest-leverage path is v0.3 extension productization: MCP, LSP, hooks, skills, and replay already have code and specs, but they need a coherent product surface, diagnostics, docs, and harness evidence before being marketed as supported.

## Goals / Non-Goals

**Goals:**
- Close the v0.2 state loop by archiving or resolving completed OpenSpec changes and documenting the current release evidence.
- Make extension surfaces diagnosable through `doctor`, commands, logs, or replay without requiring real providers or paid API calls.
- Add deterministic tests and harness scenarios for MCP/LSP availability, permission handling, hooks/skills behavior, and replay traces.
- Align English and Chinese docs with the verified support matrix.
- Keep generated documentation artifacts and source documentation changes intentionally separated.

**Non-Goals:**
- Enterprise SSO, RBAC, admin dashboards, or multi-user organization features.
- Content marketing, testimonials, blog infrastructure, or competitive claims that outrun evidence.
- IDE extensions, plugin marketplace, cloud agent, or team collaboration UX.
- Introducing new dependencies unless local mocks cannot cover the required extension verification.

## Decisions

### 1. Treat release state as a first-class gate

Completed changes should be archived or explicitly parked before new implementation begins. This avoids the current situation where OpenSpec says several changes are complete but still active, while the working tree separately shows deleted change artifacts.

Alternative considered: continue implementing v0.3 directly. Rejected because it would compound uncertainty about which specs are authoritative.

### 2. Productize MCP/LSP through diagnostics before UX polish

The first v0.3 milestone should make configured MCP servers and LSP capabilities visible, healthy, unavailable, or failing in deterministic ways. `doctor --offline`, slash commands, trace output, and docs are better first surfaces than a large TUI redesign.

Alternative considered: start with a richer interactive extension UI. Rejected because the main risk is supportability and verification, not presentation.

### 3. Use local mocks for extension proof

MCP and LSP behavior should be validated with local mock servers and fixture workspaces. Real servers can be optional smoke checks, but release gates should not depend on external services, credentials, or network availability.

Alternative considered: require real MCP/LSP integrations in CI. Rejected because it would make the release gate flaky and harder for contributors.

### 4. Expand replay from session history to debugging evidence

Replay should show enough request, tool, permission, hook, MCP/LSP, and error information to explain a failed or unavailable extension run. This makes v0.3 extension support auditable instead of relying on terminal anecdotes.

Alternative considered: keep replay as a basic transcript printer. Rejected because extension failures need structured debugging context.

### 5. Defer enterprise and marketing proposals

`enterprise-readiness` and `content-marketing` remain parked until the core open-source product surface is stable. They should not consume implementation capacity before v0.3 extension productization is proven.

Alternative considered: build marketing and enterprise proposals in parallel. Rejected because the project still needs stronger product evidence.

## Risks / Trade-offs

- Active-change cleanup could conflict with uncommitted user work -> inspect `git status` and archive only completed OpenSpec changes intentionally.
- Extension mocks may miss real server behavior -> add manual smoke check notes after deterministic mocks pass.
- Diagnostics-first work may feel less visible than UI features -> tie each diagnostic to docs, commands, and replay evidence.
- Docs source and generated `docs/.vitepress/dist` changes may mix -> keep generated artifacts out of review unless release publishing requires them.
- Replay trace expansion can expose sensitive data -> redact secrets and provider keys in trace and replay output.

## Migration Plan

1. Snapshot current verification evidence and dirty working tree state.
2. Archive or park completed/stale OpenSpec changes, then validate the remaining active set.
3. Implement extension diagnostics and local mock coverage in small slices.
4. Expand replay/trace output for extension and permission decisions.
5. Update docs and PARITY.md after tests prove behavior.
6. Run Go tests, harness, docs build, `go-code doctor --offline`, and OpenSpec validation before claiming v0.3 readiness.

Rollback is straightforward because each slice is additive: revert the affected diagnostic, test, or docs slice and keep v0.2 behavior unchanged.
