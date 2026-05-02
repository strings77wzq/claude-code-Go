## Why

The project has reached a runnable v0.2 release, but its architecture, OpenSpec history, and product surface now need a consolidated hardening roadmap before adding more features. Without a single evidence-backed plan, future work will keep mixing runtime fixes, docs truth cleanup, extension productization, and aspirational roadmap items in ways that make the agent hard to compare against Codex, Claude Code, or similar tools.

## What Changes

- Add a project-wide architecture and roadmap hardening plan based on the current implementation, OpenSpec baseline, archived tasks, docs, CI, and release evidence.
- Define a prioritized issue register covering runtime architecture, provider/model routing, permissions, TUI/CLI UX, session trace/replay, MCP/LSP/hooks/skills, harness, docs, release process, and OpenSpec hygiene.
- Establish staged milestones for the next usable agent versions: v0.2.x stabilization, v0.3 extension productization, v0.4 agent quality, and v1.0 trust/release readiness.
- Require detailed, executable tasks with acceptance criteria, test evidence, and stop conditions before implementation starts.
- Preserve current strengths: single Go binary, deterministic harness, conservative parity matrix, local-first execution, permission gates, replay evidence, and bilingual docs.

## Capabilities

### New Capabilities

- `agent-roadmap-hardening`: Covers architecture audit, issue triage, roadmap governance, implementation sequencing, and verification discipline for the next development phase.

### Modified Capabilities

- None. Existing runtime, docs, harness, and release capabilities are not directly changed by this proposal; they are referenced as implementation surfaces for the new hardening roadmap.

## Impact

- Affected architecture surfaces: `cmd/go-code`, `internal/agent`, `internal/provider`, `internal/permission`, `internal/tool`, `internal/tool/mcp`, `internal/lsp`, `internal/session`, `internal/hooks`, `internal/skills`, `pkg/tui`, `pkg/tty`, `harness`, `docs`, `PARITY.md`, and CI/release workflows.
- Affected OpenSpec surfaces: `openspec/specs/*`, archived `tasks.md` files, this new change's `design.md`, `specs/agent-roadmap-hardening/spec.md`, and `tasks.md`.
- Affected verification surfaces: `go test ./...`, `go vet ./...`, `./scripts/run-harness.sh`, docs build, OpenSpec strict validation, doctor/offline smoke checks, release workflow, and new scenario-specific harness tests.
