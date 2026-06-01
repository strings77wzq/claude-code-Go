## Why

`internal/tool/builtin` (12 tools, the agent's core capability layer) has only 20.5% test coverage, leaving file operations, web requests, and shell execution without automated regression protection. Simultaneously, `bash_semantic.go` classifies piped commands (`ls | grep foo`) as write operations due to `|` being in `WriteIndicators`, and contains unreachable dead code in `hasWriteArguments`. Fixing both closes the largest quality gap and the most impactful correctness bug in a single change.

## What Changes

- Add comprehensive unit tests for 8 untested or under-tested builtin tools: bash, glob, grep, tree, webfetch, diff, notebook, todo, plus the path validation utilities in validate.go
- Tests cover: happy path, error conditions, timeout, truncation, permission checks, workspace boundaries
- Fix `WriteIndicators` to remove false-positive classification of `|`, `;`, `&&`, `||` as write indicators
- Remove dead code in `hasWriteArguments` (`strings.HasPrefix(field, wc+" ")` can never match on space-split fields)
- Commit the pending `bash_semantic.go` refactor that improves `DetectDestructive` from map to ordered slice and adds `RemoteScriptPattern`/`RedirectPattern`

## Capabilities

### New Capabilities

- `builtin-tool-test-coverage`: Comprehensive unit tests for all 11 builtin tools plus path validation utilities (ResolvePath, ValidatePath) covering happy path, error handling, timeout, output truncation, permission gating, and workspace boundary validation. Target: 80%+ statement coverage (from 20.5%).

### Modified Capabilities

- `bash-semantic-validation`: Fix `WriteIndicators` to not classify pipes and command chains as write indicators. Remove unreachable dead code in `hasWriteArguments`. The requirement "full command chain is analyzed for safety" remains but the implementation must correctly distinguish read-only pipes from write operations.

## Impact

- `internal/tool/builtin/*_test.go` — 8-10 new test files
- `internal/permission/bash_semantic.go` — bug fixes (WriteIndicators, hasWriteArguments)
- `internal/permission/bash_semantic_test.go` — new/expanded semantic validation tests
- No API changes, no breaking changes, no dependency additions
