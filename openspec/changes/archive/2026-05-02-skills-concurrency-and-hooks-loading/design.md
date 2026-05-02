## Context

The Skills Registry (`internal/skills/registry.go`) stores skills in a plain `map[string]Skill`. Skills are loaded at startup via `LoadSkillsWithWarnings` and registered via `Registry.Register`. At runtime, `Execute` and `List` read the map. While Go's race detector might not catch this in single-goroutine REPL mode, the TUI mode uses goroutines for agent runs, creating a legitimate race window.

The Hooks system supports `PreExecute`/`PostExecute` hooks with a `Registry` that has proper `sync.RWMutex` protection. However, hooks are only registered in code — the `~/.go-code/hooks/` directory checked by `doctor.go:446-470` is never read for actual hook loading.

PostExecute hook errors are discarded with `_ = err` at `hooks.go:129`.

## Goals / Non-Goals

**Goals:**
- Protect Skills Registry map with `sync.RWMutex`
- Implement external hook loading from JSON files in `~/.go-code/hooks/`
- Replace silent error discard in `RunPostHooks` with structured logging
- Backward compatible — all existing tests pass unchanged

**Non-Goals:**
- Hot-reloading of skills or hooks (file watchers)
- Hook sandboxing or permission control
- Changing the Hook interface

## Decisions

### Decision 1: RWMutex over Mutex

**Chosen:** `sync.RWMutex` — `RLock` for `Get`/`List`/`Execute`, `Lock` for `Register`.

**Why:** `List` and `Execute` are read-heavy operations called frequently at runtime. `Register` is write-only called at startup. RWMutex allows concurrent reads.

### Decision 2: JSON-based hook files

**Chosen:** Load `.json` files from `~/.go-code/hooks/`. Each file defines a hook with `name`, `type` (`pre`/`post`), and `command` (shell command to execute with stdin JSON input).

**Why:** Consistent with the project's existing JSON-based config and skill formats. Shell commands as hooks match the project's extension philosophy (like MCP server processes).

**Alternatives considered:**
- Go plugins: Too complex, not cross-platform
- WASM: Overkill for hook execution
- Lua/embedded scripting: Adds runtime dependency

### Decision 3: Structured logging for PostExecute errors

**Chosen:** Use `WARN level (non-blocking PostExecute)` instead of `_ = err`. No injection needed — `slog.SetDefault(logger)` is already called at startup.

**Why:** `slog` is already the project's logging standard (see `main.go:82-88`). No new dependency.

## Risks / Trade-offs

- **Risk:** External hook commands could be malicious if user's `~/.go-code/hooks/` is compromised → **Mitigation:** This is a local tool — if attacker can write to the user's home directory, hooks are not the primary concern
- **Risk:** Hook JSON parsing errors at startup could block agent launch → **Mitigation:** Skip invalid hook files with warnings (same pattern as skills loader)
