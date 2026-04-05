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
- Full agent loop implementation with think → act → observe cycle
- 10 built-in tools: Read, Write, Edit, Glob, Grep, Bash, Diff, Tree, WebFetch, TodoWrite
- 3-tier permission system (ReadOnly, WorkspaceWrite, DangerFullAccess) with glob rules
- MCP (Model Context Protocol) support with stdio transport
- SSE streaming for real-time token-by-token output
- Session persistence with auto-save and resume (JSONL format)
- Skills system for custom commands and reusable workflows
- Multi-provider support: Anthropic, Tencent Cloud Coding Plan, OpenAI-compatible APIs
- Runtime model switching with `/model` command
- Auto-update checker with `/update` command
- Bubbletea TUI with interactive prompts
- Complete VitePress documentation site

### Changed
- Project structure refactored for better organization
- Tool interface standardized across all implementations

### Removed
- Legacy CLI-only mode (replaced with TUI)

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