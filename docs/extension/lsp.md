---
title: LSP Integration
description: LSP health checks, capability gating, diagnostics, and fallback behavior
---

# LSP Integration

> **Status: Experimental (v0.3 productization)** — The LSP client and health gate are implemented, but LSP-powered code intelligence is exposed only after a configured server passes a health check. `go-code doctor --offline` reports local configuration without starting network probes.

go-code uses Language Server Protocol (LSP) support as an optional extension surface for code intelligence. The core prompt, TUI, built-in tools, MCP tools, hooks, and skills continue to work when no LSP server is configured.

## Configuration

LSP health checks are opt-in through `GO_CODE_LSP_URL`:

```bash
export GO_CODE_LSP_URL="http://127.0.0.1:8080/lsp"
```

When the variable is not set, LSP is treated as unavailable rather than failed.

## Health Checks

The LSP gate sends an `initialize` request before advertising code-intelligence operations. A server is considered healthy only when initialization succeeds.

```bash
go-code doctor
```

Offline diagnostics do not contact the server:

```bash
go-code doctor --offline
```

Expected doctor states:

| State | Meaning |
| --- | --- |
| `[SKIP] lsp: not configured` | `GO_CODE_LSP_URL` is unset. |
| `[SKIP] lsp: configured ... skipped by --offline` | LSP URL is present, but network health was intentionally skipped. |
| `[PASS] lsp` | Initialization succeeded. |
| `[FAIL] lsp` | A configured server failed initialization or returned an invalid response. |

## Capability Gating

go-code advertises LSP operations only after the health check passes. The gate currently recognizes:

| Operation | LSP capability |
| --- | --- |
| Diagnostics | `publishDiagnosticsProvider` |
| Symbols | `workspaceSymbolProvider` or `documentSymbolProvider` |
| Definitions | `definitionProvider` |
| References | `referencesProvider` |
| Hover | `hoverProvider` |

If the server is not configured, cannot initialize, or does not advertise a capability, the corresponding operation remains unavailable.

## Trace Evidence

LSP health outcomes are recorded as non-fatal extension events in the session trace:

```json
{"type":"extension","name":"lsp","event":"health_check","status":"unavailable"}
```

Successful checks include advertised operations and server identity when provided. Failed checks record the initialization error without stopping unrelated workflows.

## Fallback Behavior

- No LSP config: continue without LSP operations.
- Offline doctor: report configuration state without network access.
- Initialization failure: keep core prompt/tool workflows available and report the LSP failure through diagnostics and trace evidence.
- Missing server capability: do not advertise that operation.

## Verification

Current coverage lives in `internal/lsp/gate_test.go` and proves unavailable, initialization success, initialization failure, trace recording, and operation gating behavior.

## Related

- [MCP Integration](./mcp.md)
- [Hooks System](./hooks.md)
- [Skills System](./skills.md)
