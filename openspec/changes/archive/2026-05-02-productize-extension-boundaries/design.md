## Context

Extension systems increase the usefulness of a local coding agent but also widen its trust boundary. The existing roadmap calls for MCP, LSP, hooks, and skills productization; the architect review narrowed this into security-first lifecycle and diagnostics work. This change depends conceptually on `fix-core-runtime-safety` because extension tools must share permission, trace, and runtime error semantics.

## Goals / Non-Goals

**Goals:**
- Make extension startup, health, timeout, shutdown, and diagnostic behavior consistent.
- Harden MCP process launch and tool calls before advertising MCP as a reliable product feature.
- Give doctor, TUI, and replay a shared diagnostic model for extension status.
- Keep provider profile metadata reusable and independent from transport code.

**Non-Goals:**
- Build a marketplace or remote plugin distribution system.
- Implement every LSP feature beyond the already scoped diagnostic/symbol/navigation surface.
- Add new provider SDKs.
- Change the permission model itself beyond consuming the shared policy.

## Decisions

1. Use a shared extension diagnostic type.
   - Diagnostics include component, severity, code, summary, detail, retryability, and redacted metadata.
   - Alternative rejected: each extension package printing its own warnings, because doctor/TUI/replay need the same data.

2. Treat MCP server launch as a trust boundary.
   - Configured command, args, working directory, environment, server identity, and timeout behavior are validated before use.
   - Alternative rejected: launching arbitrary configured commands without normalization, because it is difficult to audit.

3. Apply permission and trace policy to extension tools.
   - MCP tools are classified and evaluated through the same policy used by built-ins.
   - Alternative rejected: allowing MCP tools to self-declare safety without host policy enforcement.

4. Make LSP/hooks/skills failures non-fatal by default.
   - Invalid or unavailable extensions produce diagnostics and disabled capabilities, not core agent failure.
   - Alternative rejected: failing startup for optional extension errors, because local environments vary widely.

5. Keep provider profiles pure.
   - Profiles describe provider/model identity, compatibility, defaults, and capabilities; transports handle HTTP/streaming details.
   - Alternative rejected: profile-specific transport branching, because it duplicates provider behavior and makes runtime switching brittle.

## Risks / Trade-offs

- [Risk] MCP allowlisting can surprise users with existing configs. → Mitigation: report clear diagnostics and document the opt-in path.
- [Risk] Shared diagnostics may require touching several packages at once. → Mitigation: introduce the type first and migrate one extension family at a time.
- [Risk] Timeouts can break slow local servers. → Mitigation: make defaults safe and configurable with bounded maximums.

## Migration Plan

1. Add diagnostic type and golden output fixtures.
2. Harden MCP lifecycle and permission evaluation.
3. Migrate LSP, hooks, and skills to shared diagnostics.
4. Separate provider profile and transport tests.
5. Surface extension diagnostics in doctor, replay, and TUI status views.
