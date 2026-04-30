## Why

The project goal is to become a Go-native, open-source Claude Code-style coding agent that is genuinely usable by the community, influential as a reference implementation, and smooth enough for daily work. The current repository has a strong direction and useful modules, but the runtime surface, tests, permission flow, parity story, and documentation need to be aligned into a coherent product-grade system before the project can credibly grow.

## What Changes

- Establish a release-quality core agent runtime with a reliable CLI/TUI, deterministic tests, session persistence, and graceful recovery.
- Align the default user experience with the advertised capabilities: doctor checks, setup, model/provider selection, permissions, sessions, resume, compact, update, and clear error states.
- Introduce a Claude Code / Claw Code-inspired parity program so core workflows are tracked, tested, documented, and intentionally scoped.
- Make the permission system practical: safe defaults, human approval, workspace boundaries, command validation, remembered decisions, and audit traces.
- Build a deterministic harness and CI gate that proves the project runs, not only compiles.
- Refine the architecture into a clean Go implementation with small interfaces, explicit boundaries, and fewer disconnected placeholder modules.
- Rework documentation into a product-first Chinese/English journey: quick success path, troubleshooting, architecture, parity status, contribution path, and roadmap.
- Improve open-source readiness with contribution workflows, issue templates, release notes, benchmark methodology, and community-facing project positioning.
- Remove or rewrite unsupported claims, placeholder testimonials, stale metrics, and documentation that is ahead of the implementation.
- No new runtime dependency should be added unless it directly improves the core experience and is justified in design.

## Capabilities

### New Capabilities

- `runtime-health-check`: Doctor command and startup diagnostics covering config, provider access, filesystem permissions, tool availability, session paths, and docs for remediation.
- `smooth-cli-tui-experience`: Unified CLI/TUI command surface for setup, prompt mode, model switching, sessions, resume, compact, update, help, and user-friendly failures.
- `permission-and-sandbox-flow`: Practical permission enforcement with safe defaults, approval prompts, workspace boundaries, semantic bash validation, audit traces, and session memory.
- `parity-harness`: Deterministic mock provider and parity test suite for agent loop, tool use, permission decisions, streaming, recovery, and session replay.
- `provider-model-system`: Clean provider/model abstraction for Anthropic and OpenAI-compatible APIs, with config validation, runtime switching, and documented compatibility limits.
- `session-trace-and-replay`: Persistent sessions, tool traces, recovery traces, replay tooling, and inspectable logs for debugging and demonstrations.
- `mcp-lsp-extension-surface`: Productized extension surface for MCP tools, hooks, skills, and LSP features with configuration and safety integration.
- `docs-product-experience`: Documentation restructure for Chinese and English users, covering quick start, installation, doctor, usage, concepts, architecture, troubleshooting, parity, roadmap, and contribution.
- `open-source-release-readiness`: Community and release infrastructure including CI gates, benchmark methodology, changelog discipline, templates, and honest positioning.

### Modified Capabilities

- None. The repository currently has no archived base specs in `openspec/specs/`, so this change introduces new capability contracts rather than modifying existing specs.

## Impact

- Affected code: `cmd/go-code`, `pkg/tui`, `pkg/tty`, `internal/agent`, `internal/api`, `internal/provider`, `internal/permission`, `internal/tool`, `internal/session`, `internal/config`, `internal/hooks`, `internal/skills`, `internal/lsp`, `internal/tool/mcp`, `internal/update`, `internal/cost`, `internal/telemetry`.
- Affected tests: Go unit tests, integration tests, Python harness tests, CI workflows, deterministic mock provider scenarios, and docs build checks.
- Affected docs: `README.md`, `docs/`, `docs/zh/`, ADRs, roadmap, benchmark page, troubleshooting, API/command reference, contribution guide, and website homepage.
- Affected project surfaces: release process, GitHub templates, community onboarding, issue triage, parity tracking, and public positioning.
- External references: Claw Code's public project shape, especially its `doctor`-first onboarding, usage/parity/roadmap/philosophy documentation split, and deterministic parity harness approach.
