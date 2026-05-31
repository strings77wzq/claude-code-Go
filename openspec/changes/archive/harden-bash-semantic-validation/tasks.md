## 1. Red: Lock Semantic Validator Behavior

- [x] 1.1 Add table-driven tests for `VerifyReadOnly`, including allowed read-only commands and rejected writes, pipes, redirects, and command substitution.
- [x] 1.2 Add tests for `DetectDestructive`, including recursive deletes, privilege escalation, remote script execution, device writes, and process/system termination commands.
- [x] 1.3 Add tests for parsing helpers: `ParsePipes`, `ParseRedirects`, `ParseSubshells`, and `ParseCommandChaining`.
- [x] 1.4 Add tests for sed/awk write-path extraction and workspace validation.
- [x] 1.5 Add tests for `ValidatePath`, `AnalyzeSemantics`, and `ValidateFullCommand` covering workspace-safe paths, traversal, blocked system paths, workspace escapes, and dangerous subshells.
- [x] 1.6 Run `go test ./internal/permission` and confirm the new tests fail only where behavior needs implementation.

## 2. Green: Fix Narrow Validator Defects

- [x] 2.1 Fix path, redirect, sed/awk, or subshell parsing defects exposed by the failing tests while preserving conservative safety behavior.
- [x] 2.2 Keep helper behavior deterministic and dependency-free; avoid broad shell-parser rewrites.
- [x] 2.3 Run `go test ./internal/permission` and confirm all semantic validator tests pass.

## 3. Refactor And Integration

- [x] 3.1 Refactor validator helpers only where tests show duplication or unclear behavior.
- [x] 3.2 Run `go test ./internal/tool/builtin` to confirm Bash tool integration still passes.
- [x] 3.3 Run `go test ./internal/permission ./internal/tool/builtin` with coverage and record the updated package coverage.

## 4. Verification

- [x] 4.1 Run `gofmt -l internal/permission internal/tool/builtin` and confirm no files are listed.
- [x] 4.2 Run `go test ./...` and confirm all Go packages pass.
- [x] 4.3 Run `go vet ./...` and confirm no diagnostics.
- [x] 4.4 Run `./scripts/run-harness.sh` and confirm all harness scenarios pass.
- [x] 4.5 Run `openspec validate harden-bash-semantic-validation --strict --json --no-interactive` and confirm the change validates.
- [x] 4.6 Update this task list to checked only after each task has fresh evidence.
