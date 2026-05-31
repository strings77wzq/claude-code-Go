## Why

The Bash semantic validator is a safety-critical gate before shell execution, but it currently has little direct regression coverage despite containing path, redirect, subshell, and destructive-command logic. Hardening it now reduces the risk that future permission or tool changes accidentally allow workspace escapes or unsafe shell constructs.

## What Changes

- Add direct, table-driven regression coverage for Bash semantic validation behavior.
- Cover read-only command recognition, destructive command detection, redirect parsing, sed/awk write-path extraction, workspace boundary checks, subshell analysis, and full-command validation.
- Fix narrow validator defects exposed by the tests without changing the public CLI or permission-mode contract.
- Keep existing Bash tool behavior compatible except where current behavior is demonstrably unsafe or inconsistent with the spec.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `bash-semantic-validation`: Strengthen the requirement so semantic validation has deterministic test coverage and enforces workspace-safe handling for redirects, sed/awk writes, subshells, and destructive commands.

## Impact

- Affected code: `internal/permission/bash_semantic.go`, focused tests under `internal/permission`, and possibly narrow Bash tool integration tests under `internal/tool/builtin`.
- Affected specs: `openspec/specs/bash-semantic-validation/spec.md` via a change delta.
- No new runtime dependencies.
- No intended breaking changes to CLI flags, provider configuration, or tool schema.
