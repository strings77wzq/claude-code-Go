## Why

The permission mode is hardcoded to `WorkspaceWrite` with no CLI flag or environment variable to override it. In non-interactive mode (scripting/CI), permission prompts block forever on stdin, hanging the process. Additionally, the permission hierarchy uses equality comparison instead of proper level-based comparison, so a `DangerFullAccess` session incorrectly rejects tools requiring `WorkspaceWrite`. These three issues prevent scripted automation, break CI pipelines, and violate the principle of least privilege.

## What Changes

- Add `--permission-mode` CLI flag accepting `read-only`, `workspace-write`, `danger-full-access` with `GO_CODE_PERMISSION_MODE` env var fallback
- Fix `meetsModeRequirement` to use hierarchical comparison (ReadOnly < WorkspaceWrite < DangerFullAccess) instead of strict equality
- Non-interactive mode detection: when stdin is not a terminal, permission prompts return denied immediately instead of blocking indefinitely
- Clear error messages explaining which permission was required and how to grant it

## Capabilities

### New Capabilities
- `permission-mode-cli-config`: CLI flag and env var to set the agent's permission mode at startup, supporting least-privilege workflows
- `noninteractive-permission-fastfail`: Detection of non-interactive stdin and immediate denial of permission prompts with actionable error messages

### Modified Capabilities
- `permission-and-sandbox-flow`: Permission hierarchy comparison changes from equality to level-based (ReadOnly < WorkspaceWrite < DangerFullAccess), affecting which tools are allowed under each mode

## Impact

- `cmd/go-code/main.go` — new CLI flag parsing, flag→permission mode mapping, non-interactive stdin detection
- `internal/permission/policy.go` — `meetsModeRequirement` comparison logic rewrite
- `internal/config/loader.go` — `GO_CODE_PERMISSION_MODE` env var loading
- `internal/permission/prompter.go` — non-interactive fast-fail behavior
- Existing permission tests remain passing; new tests added for hierarchy, flag parsing, and non-interactive path
