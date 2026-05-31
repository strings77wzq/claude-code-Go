## 1. Security Fix: bash_semantic.go Classification Bugs

- [ ] 1.1 Fix `WriteIndicators`: remove `|`, `;`, `&&`, `||` — these are command composition operators, not write indicators
- [ ] 1.2 Fix `hasWriteArguments`: remove dead code — both `strings.HasPrefix(field, wc+" ")` (field is whitespace-split, cannot contain space) and `strings.HasPrefix(field, "-"+wc)` (meaningless prefix check). Keep only `field == wc` exact match
- [ ] 1.3 Commit the pending `bash_semantic.go` refactor (map→ordered slice in `DetectDestructive`, `RemoteScriptPattern`, `RedirectPattern`)
- [ ] 1.4 Update existing test `bash_semantic_test.go:23`: `cat ./README.md | wc -l` should now be classified as read-only (was false-positive); change `want: false` to `want: true`
- [ ] 1.5 Run existing tests to confirm no regression from the fixes

## 2. bash_semantic.go Semantic Validation Tests

- [ ] 2.1 Add test: piped read-only commands (`ls | grep foo | wc -l`) are classified as read-only
- [ ] 2.2 Add test: `&&` chained read-only commands are not falsely classified as write
- [ ] 2.3 Add test: `;` chained read-only commands are not falsely classified as write
- [ ] 2.4 Add test: `$()` command substitution is still flagged as potential write indicator
- [ ] 2.5 Add test: `tee` command is still flagged as write indicator
- [ ] 2.6 Add test: `hasWriteArguments` correctly detects write commands in arguments (e.g., `echo cp`)

## 3. Bash Tool Tests

- [ ] 3.1 Add test: happy path — simple echo command returns output
- [ ] 3.2 Add test: empty command returns error
- [ ] 3.3 Add test: timeout kills long-running command
- [ ] 3.4 Add test: output truncation when exceeding maxOutputSize
- [ ] 3.5 Add test: exit code propagation for failing commands
- [ ] 3.6 Add test: destructive command (rm -rf) is blocked by semantic validation
- [ ] 3.7 Add test: command with workspace-escaped path is blocked

## 4. Glob Tool Tests

- [ ] 4.1 Add test: happy path — glob pattern matches files in temp directory
- [ ] 4.2 Add test: empty pattern returns error
- [ ] 4.3 Add test: pattern with no matches returns empty list
- [ ] 4.4 Add test: recursive glob (`**`) finds nested files

## 5. Grep Tool Tests

- [ ] 5.1 Add test: happy path — pattern found in temp file
- [ ] 5.2 Add test: no matches returns empty output
- [ ] 5.3 Add test: empty pattern returns error
- [ ] 5.4 Add test: file not found returns error

## 6. Tree Tool Tests

- [ ] 6.1 Add test: happy path — prints directory tree with nested structure
- [ ] 6.2 Add test: empty directory returns empty listing
- [ ] 6.3 Add test: non-existent directory returns error

## 7. WebFetch Tool Tests

- [ ] 7.1 Add test: happy path — fetches content from mock HTTP server
- [ ] 7.2 Add test: URL not found (404) returns error
- [ ] 7.3 Add test: invalid URL returns error
- [ ] 7.4 Add test: timeout handling for slow server

## 8. Diff Tool Tests

- [ ] 8.1 Add test: happy path — shows unified diff between two temp files
- [ ] 8.2 Add test: identical files produce empty diff
- [ ] 8.3 Add test: non-existent file returns error

## 9. NotebookEdit Tool Tests

- [ ] 9.1 Add test: happy path — reads and modifies a notebook cell
- [ ] 9.2 Add test: invalid notebook path returns error
- [ ] 9.3 Add test: cell not found returns error

## 10. TodoWrite Tool Tests

- [ ] 10.1 Reset `globalTodoList` before each test to ensure test isolation (package-level singleton; parallel tests would interfere)
- [ ] 10.2 Add test: happy path — creates todo list from input
- [ ] 10.3 Add test: empty todos returns formatted empty list
- [ ] 10.4 Add test: todos with various statuses (pending, in_progress, completed)

## 11. Path Validation Utility Tests (validate.go: ResolvePath, ValidatePath)

- [ ] 11.1 Add test: `ResolvePath` rejects blocked system paths (`/dev/sda`, `/proc/cpuinfo`, `/sys/kernel`)
- [ ] 11.2 Add test: `ResolvePath` accepts paths within workspace
- [ ] 11.3 Add test: `ResolvePath` detects path traversal (`../etc/passwd`)
- [ ] 11.4 Add test: `ValidatePath` returns error for empty path

## 12. Final Verification

- [ ] 12.1 Run `go test -cover ./internal/tool/builtin` and confirm 80%+ coverage
- [ ] 12.2 Run `go test -cover ./internal/permission` and confirm all tests pass
- [ ] 12.3 Run `go test ./...` and confirm all 26 packages pass
- [ ] 12.4 Run `go vet ./...` and confirm clean
- [ ] 12.5 Run `gofmt -l .` and confirm no formatting issues
