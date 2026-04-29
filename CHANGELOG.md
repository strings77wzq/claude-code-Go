# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0] - 2026-04-30

### Added
- Baseline test coverage for 9 previously untested packages (61 new tests): LSP, Tool, Tool Init, MCP, Telemetry, Update, Anthropic/OpenAI Providers, TUI
- Provider model registry updated: DeepSeek v4 series (deepseek-v4-pro, deepseek-v4-flash), MiMo-V2.5 (mimo-v2.5-pro)
- Unknown-model passthrough with provider inference from model name prefixes
- Legacy model name deprecation warnings (deepseek-chat → deepseek-v4-pro, deepseek-reasoner → deepseek-v4-flash)
- Provider profile architecture reference docs in openspec/specs/

### Changed
- README: feature claims now labeled with PARITY.md verification status (verified/experimental/planned)
- MCP and LSP correctly labeled as "Planned (v0.3)" across all documentation
- Placeholder demo GIF, testimonials, and stale benchmark metrics removed
- Chinese docs (docs/zh/) synced with English docs — model names, provider config, feature status
- PARITY.md: verified rows now include evidence links to test files; unsupported claims downgraded to partial
- 3 overlapping roadmap changes archived; strategic insights extracted to openspec/specs/
- Doctor command now passes all checks offline (binary, tools, session dir, docs)

### Fixed
- All 24 Go packages pass `go test ./...` with zero failures
- Documentation build passes with zero errors
- Honest feature labeling prevents users from trying unsupported MCP/LSP features

## [v0.1.0] - 2026-04-05

## [v0.1.0] - 2026-04-05

### Added
- Full Agent Loop with stop_reason-driven state machine
- 10 built-in tools: Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch, TodoWrite
- 3-tier permission system (ReadOnly / WorkspaceWrite / DangerFullAccess)
- MCP (Model Context Protocol) support with stdio transport
- SSE streaming with custom parser
- Session persistence and resume (JSONL format)
- Skills system for custom commands and workflows
- Multi-Provider support (Anthropic, OpenAI-compatible: DeepSeek, Qwen, GLM)
- Runtime model switching with `/model` command
- LSP integration (symbols, references, diagnostics, definition, hover)
- Auto-recovery mechanism (API timeout, rate limit, tool error, context full)
- Bubbletea TUI with dark/light theme
- Bash semantic validation (937 LOC,对标 Claw Code 1004 LOC)
- File boundary guards (binary detection, size limit, symlink escape)
- Cost tracking and estimation
- Auto-update checker
- VitePress documentation site (English + Chinese)
- Python Harness (Mock API, evaluators, replay, parity tests)
- GoReleaser configuration for multi-platform releases
- GitHub Actions CI/CD
- Open source community files (CONTRIBUTING, SECURITY, CODE_OF_CONDUCT)

---

## Migration Notes

### Upgrading from v0.0.x
- Configuration files changed from YAML to JSON format
- MCP server configuration moved to `~/.go-code/mcp.json`
- Session files use new JSONL format

---

## Footnotes

[Unreleased]: https://github.com/strings77wzq/claude-code-Go/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/strings77wzq/claude-code-Go/releases/tag/v0.1.0