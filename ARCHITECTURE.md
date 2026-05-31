# ARCHITECTURE.md — claude-code-Go

> 10-minute read for new contributors. Understand the codebase structure, core abstractions, data flow, and design decisions.

---

## Project Identity

**claude-code-Go** is a pure-Go implementation of the Claude Code CLI — an AI coding assistant that runs in your terminal. It provides an agent loop powered by LLMs, with tool execution, permission control, hooks, LSP integration, and a Python eval harness.

| Fact | Value |
|------|-------|
| Language | Go 1.24.2 |
| Dependencies | Zero external Go frameworks (stdlib only) |
| Testing | Go stdlib `testing` + Python harness |
| Documentation | VitePress static site |
| Process | OpenSpec (`explore → propose → apply → archive`) |
| License | MIT |

---

## Project Structure

```
claude_code_Go/
├── cmd/go-code/           # main.go — CLI entry point
├── internal/
│   ├── agent/             # 🔥 Core agent loop (hottest path)
│   │   ├── loop.go        # Main loop: read→think→act→respond (~733 lines)
│   │   └── ...
│   ├── api/               # Anthropic API client (SSE streaming)
│   ├── tool/              # Tool registration & execution framework
│   │   ├── registry.go    # Tool catalog
│   │   ├── bash.go        # Bash tool (sandboxed)
│   │   ├── edit.go        # File edit tool
│   │   ├── write.go       # File write tool
│   │   └── task.go        # Sub-agent delegation
│   ├── permission/        # Permission control
│   │   └── modes: default, acceptEdits, plan, bypassPermissions
│   ├── hooks/             # Hook system (PreToolUse, PostToolUse, Stop)
│   ├── skills/            # Skill loading, concurrent resolution, hook injection
│   ├── config/            # Configuration loading
│   ├── session/           # Session & conversation management
│   │   └── memory/        # Agent memory (KV store + semantic search)
│   ├── lsp/               # Language Server Protocol integration
│   ├── provider/          # Model provider adapters (Anthropic, Bedrock, Vertex)
│   ├── telemetry/         # Usage telemetry
│   ├── logger/            # Structured logging
│   └── command/           # Slash command support
├── pkg/
│   ├── tty/               # TTY utilities
│   └── tui/               # Terminal UI components
├── tests/integration/     # Integration tests
├── openspec/
│   ├── specs/             # 60+ specification files
│   │   ├── agent-quality-gates/
│   │   ├── agent-product-recenter/
│   │   ├── provider-model-system/
│   │   └── ...
│   └── changes/           # Active & archived change proposals
├── harness/               # Python eval harness
│   ├── evaluators/        # Output quality assessment
│   ├── quality/           # Manifest-driven quality gates
│   ├── replay/            # Session replay & trace analysis
│   ├── workflow/          # Workflow quality evaluation (NEW)
│   └── manifests/         # Test scenario manifests
├── docs/                  # VitePress documentation site
├── CLAUDE.md              # Project instructions (loaded by Claude)
├── ARCHITECTURE.md         # ← This file
└── CONTRIBUTING.md         # Contributor guide
```

---

## Core Abstractions

### 1. Agent Loop (`internal/agent/`)

The central loop that powers the CLI. This is the **hottest path** in the codebase:

```
 User Input
     ↓
 [Read]  ← context, tools, history
     ↓
 [Think] ← LLM inference (streaming SSE)
     ↓
 [Act]   ← tool execution / text response
     ↓
 [Loop]  ← back to Read until stop condition
```

**Key file**: `loop.go` (~733 lines, approaching 800-line limit). This is the most referenced module (11x cross-references). Future refactoring: split into `read.go`, `think.go`, `act.go`.

**State machine**: The loop uses a state machine pattern with handlers:
- `handleRead` — build context, manage prompt caching
- `handleThink` — stream LLM response, parse tool calls
- `handleAct` — execute tools, handle permissions, apply results

### 2. Tool System (`internal/tool/`)

Pluggable tool registry. Each tool implements:

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() JSONSchema
    Execute(ctx Context, params Params) (Result, error)
}
```

**Built-in tools**:
| Tool | File | Description |
|------|------|-------------|
| Bash | `bash.go` | Sandboxed shell execution |
| Read | `read.go` | File reading |
| Write | `write.go` | File writing |
| Edit | `edit.go` | Exact string replacement |
| Glob | `glob.go` | File pattern matching |
| Grep | `grep.go` | Content search |
| Task | `task.go` | Sub-agent delegation |

**Concurrent execution**: `task.go` supports parallel sub-agent spawning with worktree isolation.

### 3. Permission System (`internal/permission/`)

Four permission modes:
| Mode | Behavior |
|------|----------|
| `default` | Prompt user for each tool execution |
| `acceptEdits` | Auto-accept file edits, prompt for destructive ops |
| `plan` | Read-only, no mutations allowed |
| `bypassPermissions` | Skip all permission checks |

### 4. Hook System (`internal/hooks/`)

Lifecycle hooks in the style of Claude Code settings:
- **PreToolUse** — before tool execution (validation, modification)
- **PostToolUse** — after tool execution (formatting, linting)
- **Stop** — session end (build verification)

Hooks are configured via `settings.json` and executed as external processes.

### 5. Skills (`internal/skills/`)

Skill loading and execution. Skills are markdown files with YAML frontmatter that define agent capabilities. Loading is concurrent with resolution of cross-references.

### 6. LSP Integration (`internal/lsp/`)

Language Server Protocol client for IDE-like features:
- Diagnostics (errors, warnings)
- Go-to-definition, find-references
- Document symbols, hover information
- Code actions (refactoring, quick fixes)

### 7. Provider Adapters (`internal/provider/`)

Abstraction over multiple LLM providers:
- **Anthropic** (direct API)
- **AWS Bedrock**
- **Google Vertex AI**
- **DeepSeek** (via MIMO proxy)

Each provider adapts to a common interface for streaming and tool use.

---

## Data Flow

### Agent Loop Data Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Config  │ →  │  Session │ →  │  Agent   │
│  Loader  │    │  Manager │    │  Loop    │
└──────────┘    └──────────┘    └────┬─────┘
                                     │
                    ┌─────────────────┼─────────────────┐
                    ↓                 ↓                  ↓
              ┌──────────┐    ┌──────────┐     ┌──────────┐
              │  Tools   │    │   LLM    │     │  Hooks   │
              │ Registry │    │ Provider │     │  Engine  │
              └──────────┘    └──────────┘     └──────────┘
```

### Tool Execution Flow

```
User prompt → Agent Loop → LLM responds with tool call
     ↓
Permission check (mode-dependent)
     ↓
PreToolUse hooks (validate, modify)
     ↓
Tool.Execute(ctx, params)
     ↓
PostToolUse hooks (format, lint)
     ↓
Result → Agent Loop → next iteration
```

### OpenSpec Change Flow

```
explore (understand) → brainstorm (alternatives) → propose (artifacts)
     ↓
proposal.md + design.md + tasks.md + requirement.md
     ↓
architect review → user approve
     ↓
TDD implementation → quality gates → parallel review → ship → archive
     ↓
eval scorecard (harness/test_workflow_quality.py)
```

---

## Design Decisions

### Why Go?
- Zero-dependency binary distribution (single `go build`)
- Native concurrency (goroutines for sub-agent parallelism)
- Fast compilation (instant feedback loop)
- Matches Anthropic's own CLI performance profile

### Why Python Harness?
- FastAPI for mock server (SSE streaming)
- pytest for parameterized testing
- NLP capabilities for output quality evaluation
- Not shipped to users — dev/CI only

### Why OpenSpec?
- Structured change tracking (not ad-hoc git commits)
- Every design decision is traceable: tasks → design → proposal
- LESSONS.md captures failure patterns for continuous improvement
- Enables contributors to understand WHY, not just WHAT

### Why Pure Go (No Frameworks)?
- Minimizes dependency risk
- Easier for contributors (no framework to learn)
- Go stdlib is comprehensive enough for CLI tools
- Smaller binary, faster builds

---

## Concurrency Model

```
Main Agent Loop (single goroutine, event-driven)
     │
     ├── Sub-agents (goroutine pool, bounded concurrency)
     │   └── Each in isolated worktree (optional)
     │
     ├── Hook execution (external processes, timeout-guarded)
     │
     └── Skill loading (concurrent resolution with sync.Once)
```

**Safety**: All goroutine paths tested with `go test -race ./...`. No shared mutable state without synchronization.

---

## Quality Gates

Every change must pass:

| Gate | Tool | Threshold |
|------|------|-----------|
| Format | `go fmt` | Zero diffs |
| Vet | `go vet` | Zero errors |
| Lint | `golangci-lint` | Zero new issues |
| Build | `go build` | Must pass |
| Test | `go test -cover` | ≥ 80% |
| Race | `go test -race` | Zero races |
| Security | `gosec` | Zero high/critical |
| Eval | `pytest harness/test_workflow_quality.py` | Scorecard green |

---

## Key Constraints

1. **File size**: ≤ 800 lines per file (`agent/loop.go` at ~733, monitor closely)
2. **Function size**: ≤ 50 lines per function
3. **Error handling**: explicit, wrap with `%w`, never silently swallow
4. **Immutability**: prefer returning new values over pointer mutation
5. **Concurrency**: all goroutines verified with `go test -race`
6. **Secrets**: never hardcoded; env vars or config files only

---

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) for the full development workflow.

Quick start:
```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
make build   # → bin/go-code
make test    # verify everything works
```

First good issues are tagged in the [issue tracker](https://github.com/strings77wzq/claude-code-Go/issues).

---

## Further Reading

- [CLAUDE.md](CLAUDE.md) — project-level instructions for AI assistants
- [CONTRIBUTING.md](CONTRIBUTING.md) — contributor workflow and guidelines
- [CHANGELOG.md](CHANGELOG.md) — release history
- [openspec/specs/](openspec/specs/) — 60+ specification documents
- [harness/README.md](harness/README.md) — Python eval harness documentation
