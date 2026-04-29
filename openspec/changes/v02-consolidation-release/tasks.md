## 1. Change Consolidation

- [x] 1.1 Archive `recenter-claudecodego-agent-roadmap` via `openspec archive recenter-claudecodego-agent-roadmap` and extract provider profile, harness gate, and DeepSeek/MiMo decisions into `openspec/specs/` reference docs.
- [x] 1.2 Archive `open-source-breakthrough-production-roadmap` via `openspec archive open-source-breakthrough-production-roadmap` and extract maturity model and community flywheel insights into `openspec/specs/` reference docs.
- [x] 1.3 Archive `premium-project-upgrade` via `openspec archive premium-project-upgrade` (71/80 tasks complete, remaining 9 demo tasks deferred to v0.3).
- [x] 1.4 Verify `openspec list --json` shows only `world-class-go-agent-rewrite` and `v02-consolidation-release` as active changes.

## 2. Test Coverage — LSP Package

- [x] 2.1 Add `internal/lsp/client_test.go` with happy-path test (create LSPClient with valid URL, verify initialization) and error-path test (nil http client, unreachable server).
- [x] 2.2 Add `internal/lsp/hover_test.go` with happy-path test (mock LSP response for hover info) and error-path test (server returns error).

## 3. Test Coverage — Tool Package

- [x] 3.1 Add `internal/tool/tool_test.go` with happy-path test (create Tool implementation, verify Success/Error Result constructors) and error-path test (tool not found in registry).

## 4. Test Coverage — Tool Init Package

- [x] 4.1 Add `internal/tool/init/register_test.go` with happy-path test (register builtin tools, verify all 10 tools are present) and error-path test (double register same tool).

## 5. Test Coverage — Tool MCP Package

- [x] 5.1 Add `internal/tool/mcp/manager_test.go` with happy-path test (create manager, verify empty state) and error-path test (initialize with invalid config).
- [x] 5.2 Add `internal/tool/mcp/config_test.go` with happy-path test (parse valid MCP server config) and error-path test (parse invalid config).

## 6. Test Coverage — Telemetry Package

- [x] 6.1 Add `internal/telemetry/client_test.go` with happy-path test (create client, verify opt-in consent blocks sending) and error-path test (send with invalid endpoint).

## 7. Test Coverage — Update Package

- [x] 7.1 Add `internal/update/checker_test.go` with happy-path test (check update against mock version endpoint) and error-path test (network error during check).

## 8. Test Coverage — Anthropic Provider Package

- [x] 8.1 Add `internal/provider/anthropic/provider_test.go` with happy-path test (create provider, verify headers set correctly) and error-path test (missing API key).

## 9. Test Coverage — OpenAI Provider Package

- [x] 9.1 Add `internal/provider/openai/provider_test.go` with happy-path test (create provider, verify request format) and error-path test (missing API key).

## 10. Test Coverage — TUI Package

- [x] 10.1 Add `pkg/tui/tui_test.go` with happy-path test (create TUI model, verify initial state) and error-path test (nil agent interface).

## 11. Test Verification

- [x] 11.1 Run `go test ./...` and verify all packages report `ok` (zero `?` for packages with test files).
- [x] 11.2 Fix any failing tests or import errors discovered during the run.

## 12. Provider Registry Alignment

- [x] 12.1 Update model registry in `internal/provider/registry/registry.go`: replace `deepseek-chat` → `deepseek-v4-pro`, `deepseek-reasoner` → `deepseek-v4-flash`. Keep legacy names as deprecated aliases with warning logs.
- [x] 12.2 Add `mimo-v2.5-pro` model metadata with provider `mimo`. Add MiMo provider profile routing through the OpenAI-compatible transport (or a narrow adapter if API differs).
- [x] 12.3 Implement unknown-model passthrough: when a model name is not in the registry, infer provider from name prefix heuristics and log a warning instead of rejecting.
- [x] 12.4 Add tests for provider profile resolution for DeepSeek v4, MiMo, and unknown model passthrough.
- [x] 12.5 Update provider docs in English and Chinese with current model names and compatibility status.

## 13. Docs Truth Alignment

- [x] 13.1 Audit README against PARITY.md: label each feature claim with its PARITY.md status (`verified`, `partial`, `planned`). Add status badges or parenthetical notes.
- [x] 13.2 Remove placeholder demo GIF note. Replace with a single-sentence text description until a real recording exists.
- [x] 13.3 Remove placeholder testimonials and stale metrics from README and docs pages.
- [x] 13.4 Label MCP and LSP as "Planned (v0.3)" in README, docs, and website.
- [x] 13.5 Sync Chinese docs (`docs/zh/`) model names, provider setup, and feature status with English docs.
- [x] 13.6 Update PARITY.md: add evidence links to all `verified` rows, downgrade claims without evidence to `partial`.

## 14. Final Verification

- [x] 14.1 Run `go test ./...` and confirm all packages pass.
- [x] 14.2 Run `./scripts/run-harness.sh` and confirm all scenarios pass.
- [x] 14.3 Run docs build and confirm zero errors.
- [x] 14.4 Verify `go-code doctor --offline` runs and produces pass/fail output.
- [x] 14.5 Add v0.2 entry to CHANGELOG.md with summary, date, and links to this change and PARITY.md evidence.
- [x] 14.6 Update PARITY.md v0.2 mandatory workflow rows to `verified` with evidence links.
- [x] 14.7 Run `openspec validate v02-consolidation-release --strict` and confirm all artifacts pass.
