## Context

`internal/tool/builtin.BashTool` calls `permission.NewSemanticValidator(...).ValidateFullCommand(...)` before running a shell command. The validator already contains logic for read-only classification, destructive-pattern detection, redirect parsing, sed/awk write-path extraction, subshell parsing, command chaining, workspace path validation, and severity reporting.

The current test baseline covers adjacent permission policy, rule matching, prompters, and Bash tool execution, but it does not directly exercise most of `bash_semantic.go`. Coverage evidence from the analysis phase showed `internal/permission` at 19.1% and many semantic-validator functions at 0%. Because this code is security-sensitive, the implementation should add focused regression tests before changing behavior.

## Goals / Non-Goals

**Goals:**

- Add deterministic tests for the semantic validator's public behavior and important helper outputs.
- Preserve safe existing behavior while fixing clear false negatives or inconsistent parsing found by tests.
- Keep command validation workspace-scoped for paths discovered in redirects, sed/awk writes, and command arguments.
- Improve maintainability by making path and command parsing behavior explicit through table-driven tests.

**Non-Goals:**

- Do not redesign the full permission policy model or CLI permission modes.
- Do not add a shell parser dependency.
- Do not attempt to perfectly parse every Bash grammar edge case.
- Do not change the `BashTool` input schema or public CLI command shape.

## Decisions

1. Test through `SemanticValidator` first, then use Bash tool tests only for integration-critical behavior.

   Rationale: most risk is in validator classification and path extraction. Direct tests make failures smaller and easier to diagnose than only testing through command execution.

   Alternative rejected: test only `BashTool.Execute`. That would miss helper-level regressions and require running shell commands for behavior that can be checked without process execution.

2. Treat the validator as a conservative safety gate.

   Rationale: false positives may inconvenience a user, but false negatives can allow dangerous shell behavior. When a command contains destructive patterns, unsafe subshell content, or workspace-escaping write targets, validation MUST fail.

   Alternative rejected: allow ambiguous commands and defer to permission prompts. The semantic validator exists specifically to block known unsafe constructs before shell execution.

3. Keep parsing improvements narrow and dependency-free.

   Rationale: the repo currently avoids additional runtime dependencies for this surface. Small helper fixes are appropriate if tests expose clear gaps, but a full Bash parser is outside this change.

   Alternative rejected: replace the validator with a third-party shell AST parser. That would widen dependency and migration risk beyond this hardening task.

## Risks / Trade-offs

- [Risk] New tests expose commands whose intended behavior is ambiguous. → Mitigation: preserve current behavior unless the command clearly violates workspace or destructive-command requirements.
- [Risk] Regex-based parsing remains imperfect. → Mitigation: document behavior through tests and keep the validator conservative for unsafe constructs.
- [Risk] Path validation may reject legitimate absolute paths outside the workspace. → Mitigation: this is consistent with the existing workspace boundary contract for write-capable tool execution.
- [Risk] Adding broad test cases could lock accidental behavior. → Mitigation: focus tests on spec-level behavior: read-only classification, destructive blocking, workspace boundaries, redirects, sed/awk writes, and dangerous subshells.
