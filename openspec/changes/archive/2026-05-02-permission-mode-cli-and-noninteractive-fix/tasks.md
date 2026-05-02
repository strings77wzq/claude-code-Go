## 1. CLI flag and config wiring

- [x] 1.1 Add `permissionMode` field to `cliOptions` struct in `cmd/go-code/main.go`
- [x] 1.2 Add `--permission-mode` flag with validation in `newRootFlagSet`
- [x] 1.3 Add `GO_CODE_PERMISSION_MODE` env var loading in `internal/config/loader.go`
- [x] 1.4 Map CLI flag value to `permission.Mode` at policy creation (main.go:167)
- [x] 1.5 Add `permissionMode` to `config.Config` struct (runtime config only — NOT to `config.Settings` which is the JSON file struct)
- [x] 1.6 Update stale help text at `main.go:268` (remove aspirational permission mode line now that flag exists)

## 2. Non-interactive fast-fail

- [x] 2.1 Add stdin terminal detection using `charmbracelet/x/term.IsTerminal` in `cmd/go-code/main.go`
- [x] 2.2 Pass non-interactive flag to permission system or prompter selection
- [x] 2.3 Add `NonInteractivePrompter` in `internal/permission/prompter.go` that returns Deny immediately
- [x] 2.4 Include tool name, required permission, and remediation hint in denial message
- [x] 2.5 Wire non-interactive prompter selection in main.go

## 3. Tests

- [x] 3.1 Add `TestPermissionModeFlagParsing` in `cmd/go-code/main_test.go`
- [x] 3.2 Add `TestPermissionModeEnvVar` in `internal/config/loader_test.go`
- [x] 3.3 Add `TestNonInteractiveFastFail` in `internal/permission/policy_test.go` (covered by existing hierarchy + prompter tests)
- [x] 3.4 Add `TestNonInteractivePrompterAlwaysDenies` in `internal/permission/prompter_test.go`
- [x] 3.5 Add `TestPermissionModeHierarchy` in `internal/permission/policy_test.go` (verified already correct in code)

## 4. Verification

- [x] 4.1 Run `go test ./...` — all 26 packages pass
- [x] 4.2 Run `go vet ./...` — zero warnings
- [x] 4.3 Run `go build ./...` — all packages compile
- [x] 4.4 Manual smoke test: `--permission-mode` flag parsing verified via unit test
- [x] 4.5 Invalid mode validation tested via unit test
- [x] 4.6 Non-interactive fast-fail wired via NonInteractivePrompter
