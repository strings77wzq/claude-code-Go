# v03 Extension Productization Implementation Notes

## Current Audit Findings

Last checked: 2026-04-30.

- Steering files: no repository-level `CLAUDE.md`, `claude.md`, `TASK.md`, or `task.md` were found in the active project tree. Current planning source is OpenSpec plus `README.md`, `PARITY.md`, `CHANGELOG.md`, `docs/project-audit.md`, `docs/project-baseline.md`, and `docs/roadmap.md`.
- CI failure root cause: GitHub Actions run `25153758009`, job `docs-build`, failed in the `Build docs` step because VitePress found 34 dead links on clean commit `a49d83d65cf278bc91b174071e048c11c3582910`.
- CI evidence after local fixes: workflow YAML parses, `npm ci && npm run build` in `docs/` passes, `go test -count=1 ./...` passes, `go vet ./...` passes, `./scripts/run-harness.sh` passes with 36 scenarios, and `openspec validate v03-extension-productization --strict` passes.
- `go-code doctor --offline`: command runs and now reports MCP, LSP, hooks, and skills diagnostics. It still returns failure when no API key is configured, which is current provider/config behavior.
- Dirty generated docs: `docs/.vitepress/dist/` contains many modified generated files after docs builds. These are build outputs and should not be mixed into ordinary feature or docs-source reviews unless a release/publish task explicitly requires generated artifacts.
- Dirty OpenSpec state: `world-class-go-agent-rewrite` is already deleted in the working tree but is not in `openspec list`; treat it as pre-existing cleanup and do not restore it.
- Completed active changes before cleanup: `v02-consolidation-release`, `short-mid-term-roadmap`, `polish-for-production`, `fix-ci-and-harness-tests`, `fix-python-harness-ide-errors`, `website-and-release-fix`, and `release-readiness-fix`.
- Proposal-only active changes before cleanup: `enterprise-readiness` and `content-marketing` had no tasks. They were parked until the core v0.3 extension productization work is verified.

## Parking Decisions

- `enterprise-readiness`: archived with `--skip-specs` as a parked proposal-only idea. Do not implement SSO, RBAC, admin dashboard, organization/team management, pricing tier, or support workflows until the open-source runtime and extension surface are stable and verified.
- `content-marketing`: archived with `--skip-specs` as a parked proposal-only idea. Do not add testimonials, benchmark superiority claims, blog infrastructure, or comparison marketing until the verification matrix and v0.3 feature evidence are current.

## OpenSpec Cleanup Results

- Archived with spec updates: `v02-consolidation-release`, `short-mid-term-roadmap`, `polish-for-production`, `fix-ci-and-harness-tests`, and `fix-python-harness-ide-errors`.
- Archived with `--skip-specs`: `website-and-release-fix` and `release-readiness-fix`. Their old delta files attempted to modify specs that did not exist in the active spec tree; the implemented code/docs state is already represented by current source and verification evidence.
- Archived with `--skip-specs` as parked ideas: `enterprise-readiness` and `content-marketing`.
- Active OpenSpec lane after cleanup: `v03-extension-productization` only.

## Extension Diagnostics Audit

- MCP startup path: `cmd/go-code/main.go` loads `~/.config/go-code/mcp.json` through `internal/tool/mcp.LoadMcpConfigs`, initializes servers with `McpManager.InitializeAndRegister`, and registers namespaced MCP tools in the main tool registry.
- LSP startup path: `internal/lsp` has `LSPClient` and `LSPGate`, but no regular startup registration path yet. Doctor diagnostics use `GO_CODE_LSP_URL` as the explicit health-check opt-in.
- Hooks path: `internal/agent.Agent` owns a hook registry and pre/post hook execution in the tool loop. There is no file-backed startup hook loader yet; doctor reports the optional `~/.go-code/hooks` directory only as a filesystem readiness check.
- Skills startup path: `main.go` loads `~/.go-code/skills` through `skills.LoadSkillsWithWarnings` and registers valid skills for the legacy REPL. Doctor reports valid skill count and warning count without treating invalid skill files as fatal.
- Offline behavior: `go-code doctor --offline` does not start MCP servers and does not probe LSP network health. It only parses local config and checks local paths.

## MCP Productization Notes

- Added a local MCP fixture server inside `internal/tool/mcp/manager_test.go` using the Go test binary as a stdio MCP subprocess.
- Fixed `internal/tool/mcp/transport.go` so stdio requests no longer call `Sync()` on a pipe. On Linux this returned `invalid argument` and prevented MCP initialization.
- Added coverage proving MCP tools are registered as `mcp__{server}__{tool}` and can execute through the adapter.
- Added coverage proving MCP tools require approval through the normal permission policy in workspace mode.
- Added agent trace coverage proving MCP permission decisions and MCP tool results are written to the session JSONL trace.
- Updated English and Chinese MCP docs with config path, doctor diagnostics, unavailable states, permission behavior, and current limits.

## Generated Artifact Policy

- Source docs live under `docs/**/*.md`, `docs/.vitepress/config.ts`, and `docs/.vitepress/theme/**`.
- Generated docs output under `docs/.vitepress/dist/` is excluded from ordinary review scope.
- Include `docs/.vitepress/dist/` changes only for a release/publish task that explicitly requires generated site artifacts.
- When fixing docs CI, validate with `cd docs && npm ci && npm run build`; prefer committing source link fixes rather than generated HTML churn.

## LSP Productization Notes

- Added `LSPGate.HealthCheckWithTrace` so unavailable, healthy, and initialization-error outcomes are recorded as non-fatal `extension` events in the session trace.
- Added LSP gate fixture coverage for unconfigured, initialization success, initialization failure, and operation advertisement.
- Added `AdvertisedOperations` so diagnostics, symbols, definitions, references, and hover are exposed only after a successful health check and only when the server advertises the matching capability.
- Added English and Chinese LSP docs covering `GO_CODE_LSP_URL`, doctor states, health checks, supported operations, trace evidence, and fallback behavior.

## Replay, Hooks, Skills, and Redaction Notes

- Skills warning behavior is covered by `internal/skills/loader_test.go`: invalid skills produce warnings while valid skills remain available.
- Hook pre-execution failures now support explicit warning/blocking policy through `RegisterWithPolicy`; default `Register` preserves blocking pre-hook behavior.
- Replay now summarizes extension events, permission decisions, hook/skill-style warnings represented as extension events, errors, and final status.
- Trace writes and replay summaries redact API keys, bearer authorization values, provider tokens, secrets, and passwords while preserving token-count fields such as `input_tokens`.
- `go-code replay --evidence` provides concise release/issue evidence output.

## Harness and Docs Gate Notes

- Added `harness/test_extension_productization.py` for MCP registration, MCP permission denial, LSP unavailable behavior, and replay evidence output.
- `PARITY.md`, README, English/Chinese roadmaps, homepage copy, introduction pages, pricing, benchmark, and changelog now distinguish verified v0.2 behavior from partial v0.3 extension productization.
- Business/pricing/enterprise concepts are labeled as parked with no price, timeline, or support commitment.
- Benchmark docs now describe methodology only and avoid quantitative claims until reproducible raw results exist.

## Verification Evidence

Last checked: 2026-05-01.

- `go test ./...`: passed.
- `./scripts/run-harness.sh`: passed, 40 scenarios.
- `cd docs && npm run build`: passed.
- `ANTHROPIC_API_KEY=sk-ant-test ./bin/go-code doctor --offline`: passed and reported MCP, LSP, hooks, and skills diagnostics.
- `openspec validate v03-extension-productization --strict`: passed.

## Remaining Risks

- `npm audit` reports 4 moderate vulnerabilities in the docs toolchain (`vitepress`/`vite`/`esbuild`/`postcss`). This is not the current CI failure, but it remains a security maintenance item.
- GitHub Actions Node 20 deprecation has been addressed in workflow config by moving docs jobs to `actions/setup-node@v6` and Node 24, but the remote run still needs to be verified after push.
- Multiple old completed OpenSpec changes have been archived; the remaining OpenSpec risk is validating newly archived specs against current implementation as v0.3 proceeds.
- Real MCP and LSP server smoke checks remain manual; deterministic release gates use local fixtures and offline diagnostics.
- LSP operations are gated internals, not a complete user-facing IDE/code-intelligence command set.
- Generated `docs/.vitepress/dist/` output is dirty after docs builds and should remain outside ordinary review unless publishing requires it.
