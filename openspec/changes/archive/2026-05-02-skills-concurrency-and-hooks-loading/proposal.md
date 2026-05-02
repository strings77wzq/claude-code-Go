## Why

Skills Registry uses an unprotected `map[string]Skill` — `Register` runs at startup, `Execute`/`List` run at runtime, creating a race condition window. The Hooks system checks `~/.go-code/hooks/` in doctor diagnostics but never loads external hooks from that directory. PostExecute hook errors are silently discarded (`_ = err`), making debugging impossible.

## What Changes

- Add `sync.RWMutex` to `skills.Registry` protecting all map access
- Implement external hook loading from `~/.go-code/hooks/` directory (JSON files matching the `Hook` interface)
- Wire hook loading into `main.go` startup, after permission policy creation
- Replace `_ = err` with structured logging in `RunPostHooks`

## Capabilities

### New Capabilities
- `skills-concurrency`: Thread-safe skill registry with `sync.RWMutex` guarding Register/Get/List/Execute
- `hooks-external-loading`: Load user-defined hooks from `~/.go-code/hooks/` directory as JSON files

### Modified Capabilities
<!-- No existing specs directly cover the skills or hooks registration internals -->

## Impact

- `internal/skills/registry.go` — add `sync.RWMutex`, lock in all public methods
- `internal/hooks/hooks.go` — add external hook loading function, fix PostExecute error handling
- `cmd/go-code/main.go` — call hook loading from `~/.go-code/hooks/` on startup
- `cmd/go-code/doctor.go` — update hook check to verify loaded hooks (minor follow-up, not a blocking task)
