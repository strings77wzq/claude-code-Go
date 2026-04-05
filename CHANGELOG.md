# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- API reference documentation for built-in tools (Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch, TodoWrite)
- REPL commands reference (/help, /clear, /model, /models, /sessions, /resume, /compact, /update, /exit, /skills)
- Configuration reference with environment variables and settings.json schema
- Troubleshooting guide with common issues and error codes
- Contributor guide with development setup and PR process
- VitePress navigation updates for API Reference, Troubleshooting, and Contributing sections

### Fixed
- Documentation structure improvements

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