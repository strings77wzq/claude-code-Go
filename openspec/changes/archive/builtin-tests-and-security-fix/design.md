## Context

`internal/tool/builtin` contains 11 tool implementations (bash, read, write, edit, grep, glob, diff, tree, notebook, todo, webfetch) plus a `validate.go` utility file with `ResolvePath` and `ValidatePath` helpers. Currently only read, write, edit, and bash have test files. The remaining 8 tools lack any automated test coverage. All tools follow a common pattern: implement the `tool.Tool` interface (Name, Description, InputSchema, Execute, RequiresPermission, RequiredPermissionLevel).

`internal/permission/bash_semantic.go` has an active refactor (218 lines uncommitted) that converts `DetectDestructive` from map iteration (non-deterministic) to a priority-ordered slice, and adds `RemoteScriptPattern`/`RedirectPattern` regexps. Two bugs exist in the current code:
1. `WriteIndicators` includes `|`, `;`, `&&`, `||` — causing piped/chained read-only commands to be misclassified
2. `hasWriteArguments` contains dead code `strings.HasPrefix(field, wc+" ")` — `fields` are whitespace-split, so no field contains a space

## Goals / Non-Goals

**Goals:**
- Achieve 80%+ statement coverage for `internal/tool/builtin` (currently 20.5%)
- Every builtin tool has at least happy-path and critical-error test cases; path validation utilities also tested
- Fix `WriteIndicators` false-positive for pipes and command chains
- Remove unreachable dead code in `hasWriteArguments` (both `strings.HasPrefix(field, wc+" ")` and `strings.HasPrefix(field, "-"+wc)` are dead; only `field == wc` remains)
- Commit the pending bash_semantic.go refactor with deterministic `DetectDestructive`
- Update existing test `bash_semantic_test.go:23` (pipe false-positive) to expect correct read-only classification

**Non-Goals:**
- Integration/harness tests (these are Go unit tests only)
- New tools or tool behavior changes
- Permission policy changes beyond the classification fix
- TUI or REPL testing

## Decisions

### 1. Test strategy: table-driven tests with mock filesystem

Each tool test uses Go's standard table-driven test pattern. Tools that interact with the filesystem (read, write, edit, glob, grep, tree, diff) use `t.TempDir()` for isolated test fixtures. The bash tool tests use `t.Setenv("PATH", ...)` with a mock script to avoid real shell execution. WebFetch tests use `httptest.NewServer` for mock HTTP endpoints.

**Rationale**: Table-driven tests are the Go standard. Real temp dirs avoid mocking the filesystem which would add complexity without benefit. Mock HTTP servers are built into `net/http/httptest`.

### 2. bash_semantic.go: remove `|`, `;`, `&&`, `||` from WriteIndicators

These operators indicate command composition, not write operations. A command like `ls | grep foo | wc -l` is purely read-only. The semantic analyzer already checks individual pipeline stages and redirect targets, so removing these from `WriteIndicators` does not weaken security — destructive commands are still caught by `DetectDestructive`.

**Alternative considered**: Keep `|` but add special-case logic. Rejected as overly complex — the correct fix is to not treat composition operators as write indicators.

### 3. hasWriteArguments: simplify to exact field match only

Remove the dead `strings.HasPrefix(field, wc+" ")` and `strings.HasPrefix(field, "-"+wc)` checks. Keep only `field == wc` exact match, which is the only branch that can actually fire.

### 4. Commit the pending refactor as-is

The uncommitted `bash_semantic.go` changes (map→slice for deterministic ordering, RemoteScriptPattern, RedirectPattern) are improvements. They will be committed as part of this change after the two bug fixes above are applied on top.

## Risks / Trade-offs

- [Breaking existing tests] → Run full test suite before and after each change. The existing tests all pass, so any failure indicates a regression.
- [Bash test isolation] → Mock scripts in temp PATH ensure no real shell commands execute during testing.
- [WebFetch network dependency] → All WebFetch tests use `httptest.NewServer`, zero real network calls.
- [Classification change might allow unsafe commands] → Mitigated by `DetectDestructive` independently catching all destructive patterns. Edge case: `cat file | python -c 'write_stuff'` where the interpreter performs writes — this bypasses detection after the fix because `python` is not in the write-commands set. This is acceptable because (a) the Bash tool requires `LevelDangerFullAccess` permission anyway, (b) the semantic validator's role is catching accidental destructive patterns, not sandboxing arbitrary interpreters.
