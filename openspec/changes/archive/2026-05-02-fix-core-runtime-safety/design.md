## Context

The current roadmap audit treats runtime safety as the first executable slice because it protects every subsequent feature. The most important issues are request lifecycle correctness, permission behavior under non-interactive execution, tool panic containment, and session continuity. These areas are cross-cutting and should be fixed before MCP/LSP expansion, harness quality gates, or release polish.

## Goals / Non-Goals

**Goals:**
- Ensure TUI and CLI runs cannot remain stuck after cancellation, provider errors, permission denial, or tool panics.
- Ensure non-interactive execution never blocks waiting for a user approval prompt.
- Make permission mode decisions explicit, testable, and reusable across built-in and future extension tools.
- Persist structured trace events that support replay and debugging without leaking secrets.

**Non-Goals:**
- Redesign the entire TUI product experience.
- Add new provider transports or model registry behavior.
- Productize MCP, LSP, hooks, or skills; those are handled by `productize-extension-boundaries`.
- Introduce distributed tracing or telemetry.

## Decisions

1. Use one request lifecycle contract for CLI and TUI.
   - A run has explicit `started`, `completed`, `cancelled`, and `failed` states.
   - The lifecycle owner creates and cancels the request context only after the async command finishes.
   - Alternative rejected: relying on `defer cancel()` in the UI command constructor, because it can cancel work before the returned command completes.

2. Make permission evaluation table-driven.
   - Define a single mode hierarchy and action classification matrix.
   - The policy returns structured `allow`, `deny`, or `prompt-required` decisions with reasons.
   - Alternative rejected: spreading equality checks across tools, because it makes higher modes accidentally fail lower requirements.

3. Fail closed outside interactive approval.
   - If a tool requires approval and no TTY approval channel is available, the operation is denied with an agent-visible result.
   - Alternative rejected: implicit approval in non-interactive contexts, because it is unsafe for CI, scripts, and piped input.

4. Convert tool panics to structured tool errors.
   - Registry execution wraps each tool call with recovery, trace emission, and a safe error result.
   - Alternative rejected: process-level panic propagation, because one faulty tool should not terminate the agent loop.

5. Version trace envelopes centrally.
   - Trace writes go through constructors that apply schema versioning and redaction.
   - Alternative rejected: ad hoc `map[string]any` events from each package, because replay and redaction drift quickly.

## Risks / Trade-offs

- [Risk] Stricter fail-closed behavior may reject existing scripted workflows that relied on implicit approval. → Mitigation: document the permission mode and provide explicit opt-in flags.
- [Risk] Centralized policy code can become a bottleneck for new tool types. → Mitigation: keep action classification small and table-driven.
- [Risk] Trace schema versioning can require fixture updates. → Mitigation: add golden tests and migration notes.

## Migration Plan

1. Add regression tests that reproduce the current cancellation, permission, and panic risks.
2. Introduce shared runtime/policy/trace helpers behind existing call sites.
3. Move CLI and TUI call sites to the shared helpers.
4. Run targeted tests, then full `go test ./...`.
5. Update docs only for user-visible permission or resume behavior.
