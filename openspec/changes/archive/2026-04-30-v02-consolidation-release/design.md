## Context

The project has four active OpenSpec changes with overlapping scopes. The core agent engine (agent loop, permission, session, recovery, built-in tools) is implemented and tested. The gap is not in the engine but in the surrounding productization layer: test coverage for secondary packages, an up-to-date provider model registry, honest documentation, and a credible release gate.

Current state:
- `go test ./...` passes all existing tests (14 packages with tests)
- 9 packages have zero test coverage: `lsp`, `tool`, `tool/init`, `tool/mcp`, `telemetry`, `update`, `provider/anthropic`, `provider/openai`, `tui`
- Provider model registry hardcodes stale model names and rejects unknown models
- README and docs claim MCP/LSP support as if implemented
- PARITY.md correctly identifies partial vs verified status
- 3 overlapping roadmap changes with 0 implementation tasks completed

## Goals / Non-Goals

**Goals:**
- Reduce active changes from 4 to 1 by archiving overlapping roadmap changes
- Add at least one happy-path test and one error-path test to each of the 9 zero-test packages
- Update the model registry to current model names and support unknown-model passthrough
- Strip or relabel unsupported claims in README, docs, and website
- Deliver a v0.2 release with all tests passing, harness green, and docs built

**Non-Goals:**
- MCP productization (config integration, UX flow) — deferred to v0.3
- LSP productization (user-facing LSP commands) — deferred to v0.3
- Full test coverage (80%+) — targeting baseline coverage only
- LLM-based context compaction — deferred to v0.3
- Multi-agent delegation — no timeline
- Demo GIF recording — deferred until product is stable

## Decisions

### D1: Archive overlapping changes, not merge

The three roadmap changes (`recenter`, `breakthrough`, `premium`) contain valuable strategic insights but zero implementation and overlapping task lists. Instead of merging their 167 tasks into the rewrite change, we archive them and extract key decisions into `openspec/specs/` as project reference specs. This keeps the active change list clean and prevents task duplication.

**Alternatives considered**: Merging tasks into one mega-change. Rejected because half the tasks are strategic duplicates and the other half are already covered by the rewrite change.

### D2: Baseline test coverage per package, not percentage target

Each zero-test package gets at least 1 happy-path test + 1 error-path test. No percentage target. This ensures every package's public API is exercised without encouraging low-value coverage padding.

**Alternatives considered**: 80% coverage target. Rejected because packages like `tui` (Bubble Tea UI) and `telemetry` (opt-in client) have inherently low testability without significant refactoring. Baseline coverage catches the critical gap (package can't even instantiate) without blocking on marginal coverage.

### D3: Config-driven model discovery, not periodic code updates

Instead of maintaining a hardcoded `var modelRegistry` that rots, we add a fallback path: if a model name is not found in the registry, treat it as an unknown model and construct a best-effort provider configuration from available context (inferred provider from name heuristics, user-supplied base URL, default transport). Known models remain validated strictly.

**Alternatives considered**: Periodic manual updates to the registry. Rejected — model names change faster than release cycles. Dynamic API-driven discovery. Rejected — most provider APIs don't expose a `/models` endpoint with compatible schemas.

### D4: Truth-first docs: label, don't delete

Unsupported features (MCP, LSP) are relabeled as "Planned (v0.3)" rather than deleted. This is honest about current state while preserving the roadmap signal. Placeholder content (demo GIF, testimonials) is removed entirely since it has no truthful counterpart.

**Alternatives considered**: Delete all unsupported claims. Rejected — removes roadmap visibility that contributors need. Keep as-is until implemented. Rejected — erodes user trust when people try features that don't work.

### D5: PARITY.md as the single source of truth for v0.2 gate

The 8 mandatory v0.2 workflows defined in PARITY.md become the release gate checklist. A workflow is "verified" only when it has both passing automated tests AND a documented manual smoke check path. This prevents the gate from being a rubber stamp.

## Risks / Trade-offs

- **Test coverage is baseline, not thorough**: Future regressions in untested code paths are still possible. Mitigation: each test exercises the package's primary public API, which catches import errors, nil panics, and fundamental breakage.
- **Unknown-model passthrough could mask config errors**: If a user types a model name slightly wrong, they get a 401/404 from the provider instead of a clear "unknown model" error. Mitigation: log a warning when an unknown model is used, pointing to docs for verified models.
- **Archiving three changes loses task granularity**: Some good ideas in the roadmap changes might be forgotten. Mitigation: extract strategic decisions into `openspec/specs/` as project reference before archiving.
- **Docs relabeling may look like the project is going backwards**: Publicly changing "Supports MCP" to "MCP planned" could be perceived as removing features. Mitigation: pair with a changelog entry that explains the honesty-first approach and links to the v0.3 roadmap.

## Migration Plan

1. Archive `recenter-claudecodego-agent-roadmap`, `open-source-breakthrough-production-roadmap`, `premium-project-upgrade` using `openspec archive <name>`.
2. Extract key strategic insights from archived changes into `openspec/specs/` as project-level reference specs.
3. Implement test coverage, provider alignment, and docs changes on the single remaining change.
4. Run full verification suite (`go test ./...`, harness, docs build).
5. Tag v0.2.

No data migration. No breaking config changes. No API changes.

## Open Questions

- Should MiMo-V2.5 be added as a named provider profile or as an OpenAI-compatible passthrough? Depends on whether its API is strictly OpenAI-compatible — needs verification against public docs or user-provided endpoint samples.
- Should the `premium-project-upgrade` demo recording tasks (tasks 1–9) be moved to a separate v0.3 change or deleted? Leaning toward deleting since recordings should happen after the product stabilizes.
