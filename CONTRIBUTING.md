# Contributing to claude-code-Go

Thank you for contributing! This guide covers everything you need to get started.

## Welcome Message

Whether you're reporting a bug, proposing a feature, or submitting a pull request, your contributions are valued.

---

## Development Setup

### Prerequisites

- **Go 1.24** or later (`go version`)
- **Git** (`git --version`)
- **Python 3.10+** (for harness tests)

### Quick Start

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go

# Build
make build
# → bin/go-code

# Verify everything works
make test
# → go test -v ./... + harness tests
```

### Full Dev Environment

```bash
# Install Go tooling
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Install Python harness deps
cd harness && pip install -r requirements.txt && cd ..

# Run all checks
make vet           # go vet ./...
make test          # go test + harness
golangci-lint run ./...
gosec ./...
```

### Project Structure

```
cmd/go-code/         → CLI entry point
internal/
  agent/             → Core agent loop (loop.go is the hot path)
  tool/              → Tool registry & execution
  permission/        → Permission modes
  hooks/             → Hook system (PreToolUse/PostToolUse/Stop)
  skills/            → Skill loading & execution
  config/            → Configuration loading
  session/           → Session management
  lsp/               → LSP integration
  provider/          → Model provider adapters
  telemetry/         → Telemetry
  logger/            → Logging
  command/           → Slash command support
pkg/
  tty/               → TTY utilities
  tui/               → Terminal UI
tests/integration/   → Integration tests
openspec/            → OpenSpec specs & changes
docs/                → VitePress documentation
harness/             → Python eval harness
```

---

## Pull Request

### Before You Start

1. **Check existing issues** — Is someone already working on this?
2. **Open an issue first** for features or large changes — get design alignment before coding
3. **Read ARCHITECTURE.md** — understand the codebase structure in < 10 minutes

### Workflow (the same process used by maintainers)

This project uses a structured workflow. Every change goes through:

```
OpenSpec explore → brainstorm → propose → design review → approve
      ↓
TDD (red → green → refactor) → quality gates → verify
      ↓
Parallel review (code + security + behavior) → atomic commits → ship → archive
      ↓
Eval scorecard (pytest harness/test_workflow_quality.py) → merge
```

### Step-by-Step

1. **Fork** the repository on GitHub
2. **Clone** your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/claude-code-Go.git
   cd claude-code-Go
   git remote add upstream https://github.com/strings77wzq/claude-code-Go.git
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Write tests first** (TDD — see Testing section below)
5. **Implement** until all tests pass
6. **Run quality gates**:
   ```bash
   make vet && make test && golangci-lint run ./...
   ```
7. **Commit** with conventional commit format (see Commit Convention below)
8. **Push** to your fork: `git push origin feature/your-feature-name`
9. **Submit a pull request** to the main repository

### PR Description Template

```markdown
## Summary
Brief description of what this PR does.

## Related Issues
Fixes #123

## Type of Change
- [ ] Bug fix
- [ ] Feature
- [ ] Refactor
- [ ] Documentation
- [ ] Test

## Design Artifacts
- Proposal: openspec/changes/<name>/proposal.md
- Design: openspec/changes/<name>/design.md
- Tasks: openspec/changes/<name>/tasks.md

## Testing
- [ ] Unit tests added/updated
- [ ] Coverage ≥ 80%
- [ ] `go test -race ./...` passes
- [ ] Harness tests pass (`make test`)

## Quality Gates
- [ ] Lint: `golangci-lint run ./...` — clean
- [ ] Build: `go build ./...` — passes
- [ ] Test: `go test -cover ./...` — ≥ 80%
- [ ] Race: `go test -race ./...` — clean
- [ ] Security: `gosec ./...` — zero high/critical
- [ ] Eval: `pytest harness/test_workflow_quality.py` — all green
```

### Review Process

1. **CI must pass** — lint, test, coverage, race, security scans
2. **Maintainers review** — expect feedback within 1-3 days
3. **Address all CRITICAL/HIGH findings** before merge
4. **Eval scorecard must be green** — `pytest harness/test_workflow_quality.py`
5. Once approved, a maintainer will merge your PR

---

## Testing

### Test-Driven Development (Required)

This project practices TDD. All code changes must:

1. **Write a failing test first** (RED)
2. **Write minimal code to pass** (GREEN)
3. **Refactor while tests stay green** (REFACTOR)

### Running Tests

```bash
# Go unit + integration tests
go test -v ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Race detection (required before merge)
go test -race ./...

# Python harness tests
cd harness && python -m pytest -v && cd ..

# Eval scorecard (required before merge)
pytest harness/test_workflow_quality.py -v

# Everything
make test
```

### Coverage Target

**≥ 80%** for new code. Check with:
```bash
go test -cover ./...
```

### Test Structure

- **Table-driven tests** for Go (`t.Run` with subtests)
- **AAA pattern**: Arrange → Act → Assert
- **Descriptive names**: `TestHandler_ReturnsError_WhenInputIsEmpty`
- **Harness scenarios** for agent behavior: add manifests in `harness/manifests/`

### Harness Scenarios

For changes to agent loop, tools, or permissions, add a harness scenario:

1. Create a manifest in `harness/manifests/`
2. Add mock provider behavior in `harness/mock_server/scenarios.py`
3. Add assertions in `harness/test_scenarios.py` or `harness/test_quality_gates.py`

---

## Code Style

### Go Code

- **Format**: `go fmt ./...` before every commit (automated by hooks)
- **Lint**: `go vet ./...` + `golangci-lint run ./...`
- **Naming**: idiomatic Go — `camelCase` for unexported, `PascalCase` for exported
- **Comments**: document all exported functions and types
- **Errors**: handle explicitly, wrap with `%w`, never silently swallow
- **Immutability**: prefer returning new values over mutating pointers

### File Size Limits

- **Files**: ≤ 800 lines (agent/loop.go is ~733, approaching the limit)
- **Functions**: ≤ 50 lines
- **Nesting**: ≤ 4 levels deep

### Quality Gates (Auto-Enforced)

| Gate | Command | Threshold |
|------|---------|-----------|
| Format | `gofmt -l .` | Zero diffs |
| Vet | `go vet ./...` | Zero errors |
| Lint | `golangci-lint run ./...` | Zero new issues |
| Build | `go build ./...` | Must pass |
| Test | `go test -cover ./...` | ≥ 80% coverage |
| Race | `go test -race ./...` | Zero races |
| Security | `gosec ./...` | Zero high/critical |
| Eval | `pytest harness/test_workflow_quality.py` | All green |

### Anti-Patterns (Blockers)

- ❌ Hardcoded secrets/API keys (use env vars)
- ❌ Silent error swallowing (handle or propagate)
- ❌ `fmt.Printf` / `print()` debug statements (remove before commit)
- ❌ Assuming code works without running `go test`
- ❌ Referencing APIs/libraries without grep-verifying they exist
- ❌ "顺手改" (incidental changes) without updating tasks.md

---

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `refactor` | Code refactoring (no behavior change) |
| `test` | Adding or updating tests |
| `chore` | Maintenance, dependencies, build |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |

### Examples

```
feat(agent): add context compression for long sessions
fix(tool): correct glob pattern matching for hidden files
docs(readme): update installation instructions
refactor(api): simplify SSE streaming parser
test(permission): add table-driven tests for deny flow
```

---

## Bug Reports

1. Check if the issue already exists
2. Create a new issue with:
   - Steps to reproduce
   - Expected vs actual behavior
   - `go version` output
   - OS and environment
   - Relevant logs or error messages

## Feature Requests

1. Describe the feature and use case
2. Explain why it's valuable
3. Include mockups or examples if applicable
4. Ideally, start with `Skill: openspec-explore` to understand the impact

---

## Changelog Discipline

User-visible changes must include a CHANGELOG entry in [CHANGELOG.md](CHANGELOG.md) under `[Unreleased]`, using [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
### Added
- New feature description.

### Changed
- Changed behavior description.

### Fixed
- Bug fix description.
```

---

## First Good Issues

Look for issues tagged `good first issue` in the [issue tracker](https://github.com/strings77wzq/claude-code-Go/issues).

Good starter contributions:
- Adding tests to packages with low coverage
- Updating documentation to match current behavior
- Fixing small bugs tagged `help wanted`
- Adding examples to existing docs

When working on a `good first issue`:
1. Comment on the issue to claim it
2. Ask questions if scope is unclear
3. Start with `make test` from a clean checkout

---

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)

---

## Code of Conduct

This project follows a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to abide by its terms.

---

Thank you for contributing! Your work makes open source better.
