## Why

The roadmap audit found that OpenSpec hygiene, generated documentation policy, task evidence, and release readiness are inconsistent. This creates drift between claims, specs, docs, and shipped artifacts, which makes it harder to publish a usable version and harder for future agents to continue safely.

## What Changes

- Add a change hygiene gate for OpenSpec proposals, specs, task evidence, and archive readiness.
- Remove vague spec purposes and require durable rationale for future maintainers.
- Define generated documentation policy so generated files do not obscure reviewable source-of-truth changes.
- Tie release readiness to validated specs, docs truth checks, install verification, and published evidence.

## Capabilities

### New Capabilities
- `openspec-change-hygiene`: Defines quality gates for OpenSpec change structure, task evidence, validation, and archive readiness.

### Modified Capabilities
- `docs-truth-alignment`: Adds generated documentation source-of-truth and drift rules.
- `release-state-governance`: Adds evidence-backed release state transitions.
- `open-source-release-readiness`: Adds install, docs, and spec validation checks required before publishing.

## Impact

- Affected files: `openspec/changes`, `openspec/specs`, docs source files, generated docs outputs, release scripts, CI workflows, `PARITY.md`, README/release notes.
- Affected tests/checks: OpenSpec strict validation, docs generation drift check, install smoke tests, release checklist validation.
- No runtime behavior changes are expected.
