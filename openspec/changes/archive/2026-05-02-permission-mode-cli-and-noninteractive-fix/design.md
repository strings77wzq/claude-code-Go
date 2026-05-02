## Context

Currently `permission.NewPolicy(permission.WorkspaceWrite)` is hardcoded at `cmd/go-code/main.go:167`. There is no CLI flag, environment variable, or config key to override this. The help text at line 268 references "permission mode" but the flag doesn't exist — it's aspirational copy.

When `-p "prompt"` is used in non-interactive mode (script/CI), and the agent needs to execute a tool requiring user approval, the `TerminalPrompter` or `StdinPrompter` calls `reader.ReadString('\n')` which blocks indefinitely on a closed or absent stdin. The process hangs.

The `meetsModeRequirement` function in `policy.go:156-167` already uses hierarchical comparison (`modeHierarchy` map with integer levels and `>=` check). This was fixed in a prior commit. No additional hierarchy work is needed.

## Goals / Non-Goals

**Goals:**
- Add `--permission-mode` CLI flag accepting `read-only`, `workspace-write`, `danger-full-access`
- Add `GO_CODE_PERMISSION_MODE` environment variable as fallback in config loader
- Non-interactive mode: when stdin is not a terminal, permission prompts immediately return `Deny` with a clear error message
- Invalid mode values produce clear, actionable error messages at startup
- Backward compatible: default remains `WorkspaceWrite`

**Non-Goals:**
- Changing the permission hierarchy comparison (already hierarchical)
- Adding new permission tiers beyond the existing 3
- Modifying the session memory or rule evaluation engine
- TUI permission dialog improvements (deferred)

## Decisions

### Decision 1: CLI flag design

**Chosen:** `--permission-mode <mode>` with values `read-only`, `workspace-write`, `danger-full-access` (kebab-case CLI, PascalCase internal).

**Why:** kebab-case is the Go CLI convention. PascalCase is already the `permission.Mode` type value. We normalize at the flag parsing boundary.

**Alternatives considered:**
- Short flag `-m`: conflicts with `/model` command convention, less discoverable
- Numeric levels (1/2/3): less readable, harder to document

### Decision 2: env var naming

**Chosen:** `GO_CODE_PERMISSION_MODE` in `config.Load()`.

**Why:** Follows the project's established `GO_CODE_*` prefix convention (see `GO_CODE_API_KEY`, `GO_CODE_BASE_URL`, etc.).

### Decision 3: non-interactive detection

**Chosen:** Check `term.IsTerminal(int(os.Stdin.Fd()))` at CLI startup, store as a field on `cliOptions`. Pass a `NonInteractivePrompter` that always returns `Deny` with a descriptive error message.

**Why:** Simple. Uses `github.com/charmbracelet/x/term` already in go.mod for TUI — no new dependency. The check happens once at startup, avoiding repeated syscalls per tool execution.

**Alternatives considered:**
- Context-based detection (cancel/timeout): adds complexity, doesn't give the user a clear error
- `os.Stdin` read attempt: racy, could consume input meant for the agent

### Decision 4: error message format

**Chosen:** When non-interactive prompt is denied, the error message includes:
- The tool name that required permission
- The required permission level
- A suggestion: "Re-run with --permission-mode danger-full-access to allow all operations"

## Risks / Trade-offs

- **Risk:** Users with scripts that currently work without triggering permission prompts may break if their prompts start requiring permissions → **Mitigation:** Default `WorkspaceWrite` preserved; only tools requiring `DangerFullAccess` (bash, MCP) hit prompts
- **Risk:** `term.IsTerminal` may behave unexpectedly in exotic terminal emulators or CI runners → **Mitigation:** Degrade safely — if detection fails, treat as non-interactive (fail closed)
