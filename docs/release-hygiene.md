# Release Hygiene

This page is the source-of-truth checklist for keeping OpenSpec changes, docs, generated artifacts, and release claims aligned.

## OpenSpec Hygiene Checklist

- Active implementation changes are small enough to implement, verify, and archive as one coherent unit.
- Umbrella changes are labeled as planning artifacts and are not applied directly as large mixed-scope implementation work.
- Every active or release-bound spec has a concrete purpose and testable scenarios.
- `openspec validate <change> --strict` passes before implementation completion and before archive.
- Completed tasks include evidence or an explicit validation gap.

## Task Evidence Format

Use this format when checking off tasks:

```text
Evidence:
- command or test name
- relevant artifact path
- result summary

Not-tested:
- explicit gap, or "none"
```

## Archive Readiness

Before archiving a change:

- All tasks are checked or deliberately carried forward.
- Evidence is recorded in `tasks.md`, implementation notes, or release notes.
- `go test ./...` passes for code changes.
- `./scripts/run-harness.sh` passes when agent behavior or release evidence changes.
- Docs source and generated output policy is followed.
- `openspec validate <change> --strict` passes.

## Active Change Inventory

| Change | Type | Evidence expectation |
| --- | --- | --- |
| `agent-roadmap-audit-and-hardening` | Umbrella roadmap | Keep as planning context; do not apply directly as one implementation unit. |
| `fix-core-runtime-safety` | Implementation | Go tests, trace fixtures, TUI cancellation tests, OpenSpec validation. |
| `productize-extension-boundaries` | Implementation | MCP/LSP/hooks/skills/provider tests, doctor/replay diagnostics, OpenSpec validation. |
| `harness-agent-quality-gates` | Implementation | Pytest harness, deterministic harness command, Go tests, OpenSpec validation. |
| `openspec-docs-release-hygiene` | Implementation | Docs inventory, hygiene script, install smoke, OpenSpec validation. |

## Release State Matrix

| State | Required evidence |
| --- | --- |
| Local/dev | `go test ./...`, targeted tests for changed packages, active change validation. |
| Release candidate | Local/dev evidence, `./scripts/run-harness.sh`, docs source check, install smoke, known gaps. |
| Published | Release-candidate evidence, CI green, release notes, artifact checksums when applicable, docs deployment result. |

## Known Gaps

- External-agent comparison evidence is manual and must be labeled as such.
- Generated docs under `docs/.vitepress/dist` are release artifacts, not the review source of truth.
- Publishing may require remote credentials and CI permissions outside local verification.
