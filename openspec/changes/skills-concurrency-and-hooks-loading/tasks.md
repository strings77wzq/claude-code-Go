## 1. Skills registry concurrency

- [x] 1.1 Add `sync.RWMutex` to `skills.Registry` struct
- [x] 1.2 Add `RLock`/`RUnlock` to `Get` method
- [x] 1.3 Add `RLock`/`RUnlock` to `List` method
- [x] 1.4 Add `RLock`/`RUnlock` to `Execute` method
- [x] 1.5 Add `Lock`/`Unlock` to `Register` method

## 2. Hooks external loading

- [x] 2.1 Add `LoadHooksFromDir(dir string) ([]Hook, error)` function in `internal/hooks/hooks.go`
- [x] 2.2 Parse hook JSON files with `name`, `type`, `command` fields (ShellHook wrapper)
- [x] 2.3 Wire hook loading in `cmd/go-code/main.go` via `agent.LoadExternalHooks()`
- [x] 2.4 Add structured logging (`slog.Warn`) to `RunPostHooks` error path

## 3. Tests

- [x] 3.1 Add `TestRegistryConcurrency` in `internal/skills/registry_test.go` (concurrent Register + Execute/List)
- [x] 3.2 Add `TestLoadHooksFromDir` test in `internal/hooks/hooks_test.go`
- [x] 3.3 Add `TestLoadHooksFromDirEmptyDir` test
- [x] 3.4 Add `TestLoadHooksFromDirInvalidJSON` test + missing fields test

## 4. Verification

- [x] 4.1 Run `go test -race ./...` — 26/26 packages pass, zero data races
- [x] 4.2 Run `go test ./...` — all packages pass
- [x] 4.3 Run `go vet ./...` — zero warnings
- [x] 4.4 Run `go build ./...` — all packages compile
