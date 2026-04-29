# Project Audit

Last checked: 2026-04-28

This audit lists places where the repository's public story is ahead of the verified implementation. The goal is not to reduce ambition; it is to make the project trustworthy while the implementation catches up.

## High Priority

| Area | Current Issue | Recommended Action |
| --- | --- | --- |
| Default TUI | `/help` advertises commands such as `/models`, `/sessions`, `/resume`, `/compact`, and `/update`, but the default TUI handles only a subset. | Extract a shared command service and wire both TUI and legacy REPL to it. |
| Permission flow | Policy and prompter exist, but agent execution treats `Ask` as allowed. | Enforce human approval before side-effecting operations. |
| Harness | CLI scenarios fail with empty stdout after `bin/go-code` is built. | Make prompt mode and mock-provider scenarios pass before claiming parity. |
| README/demo | Demo GIF is marked as placeholder. | Replace with a real recorded demo or label the public demo section as pending. |
| Benchmark claims | Benchmark page contains fixed comparisons without reproducible command output. | Replace with methodology-first benchmarks and dated local measurements. |
| Showcase/testimonials | Some examples read like placeholders rather than verifiable user stories. | Remove or relabel until real user stories exist. |

## Medium Priority

| Area | Current Issue | Recommended Action |
| --- | --- | --- |
| Feature counts | Docs differ between 9, 10, and 11 built-in tools. | Generate or centralize tool count from registration. |
| Provider support | OpenAI-compatible routing exists but compatibility limits are not obvious. | Add provider matrix and validation behavior. |
| Sessions | Save/load exists, but trace schema and replay story are incomplete. | Normalize JSONL schema and provide replay tooling. |
| MCP/LSP | Code exists but default product path is unclear. | Mark experimental until configured, tested, and documented. |
| Telemetry | Module exists without a clear product value or privacy story. | Keep disabled or remove from public feature claims until decided. |

## Documentation Standard

Every public page should satisfy one of these:

- It documents a verified command or tested workflow.
- It clearly marks the feature as experimental, planned, or unsupported.
- It explains architecture using real package names and current code paths.

