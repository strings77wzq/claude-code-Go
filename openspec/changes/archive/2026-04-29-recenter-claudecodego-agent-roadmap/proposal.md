## Why

claudecode-go needs a product and engineering reset before more implementation: the goal is not a generic AI CLI, but an elegant Go-first Claude Code-style coding agent with a Python verification harness, first-class DeepSeek and MiMo support, and open-source quality that users can trust. The current repository has promising building blocks, but the design story, docs, specs, model plan, and runnable test gates are not yet aligned with that goal.

## What Changes

- Recenter the roadmap around a Go runtime plus Python harness architecture: Go owns the shipped CLI/TUI, agent loop, providers, tools, permissions, sessions, and traces; Python owns deterministic evaluation, mock providers, replay analysis, and CI harness scenarios.
- Define Claude Code / Claw Code inspiration as observable workflow parity and harness engineering, not private-source cloning or unsupported equivalence claims.
- Make DeepSeek and MiMo-V2.5 first-class model families with documented configuration, compatibility levels, model metadata, smoke tests, and harness coverage.
- Replace broad aspirational specs with a staged product contract: baseline runs, model/provider correctness, harness gates, core agent ergonomics, and honest documentation.
- Introduce an explicit “harness engineering” quality gate for release readiness: deterministic scenarios must prove tool calls, streaming, permissions, session replay, recovery, and provider adapters before docs claim support.
- Keep and strengthen the good existing direction: provider abstraction, doctor, permissions, session replay, Python harness, docs site, OpenSpec-driven work, and release-readiness mindset.
- Correct current drift: remove stale DeepSeek model assumptions, add MiMo series support, reduce premature extension-surface emphasis, and make specs/docs negotiate implementation truth rather than over-describing future intent.
- No breaking runtime API is required yet, but undocumented or unverified claims should be treated as not supported until backed by tests and docs.

## Capabilities

### New Capabilities

- `agent-product-recenter`: Defines the product identity, scope, non-goals, and staged roadmap for an elegant Claude Code-style Go rewrite with Python harness support.
- `deepseek-mimo-provider-support`: Defines first-class model/provider behavior for DeepSeek and MiMo-V2.5 series, including config, defaults, compatibility, validation, and tests.
- `harness-engineering-gate`: Defines the deterministic harness and CI release gate required before features are documented as supported.
- `docs-spec-alignment`: Defines how docs, OpenSpec artifacts, parity status, and public claims stay synchronized with verified implementation.

### Modified Capabilities

- None. There are active change-local specs, but no archived base specs under `openspec/specs/`; this realignment introduces new capability contracts that future changes can apply independently.

## Impact

- Affected code: `cmd/go-code`, `internal/agent`, `internal/provider`, `internal/api`, `internal/config`, `internal/permission`, `internal/session`, `internal/tool`, `internal/logger`, `pkg/tui`, `pkg/tty`.
- Affected harness: `harness/mock_server`, `harness/replay`, `harness/evaluators`, `harness/test_*.py`, and `scripts/run-harness.sh`.
- Affected docs/specs: `README.md`, `PARITY.md`, `docs/`, `docs/zh/`, `docs/architecture/`, `harness/README.md`, `docs/project-audit.md`, and OpenSpec change artifacts.
- External constraints: DeepSeek official API docs currently expose OpenAI/Anthropic-compatible access with `deepseek-v4-flash` and `deepseek-v4-pro`; MiMo official pages expose `mimo-v2.5-pro` as the model tag and emphasize 1M-context long-horizon agentic coding with a proper harness.
- Public positioning impact: the project should state it is inspired by Claude Code / Claw Code workflows and harness patterns, while avoiding claims of private Claude Code source compatibility unless source is explicitly provided and legally usable.
