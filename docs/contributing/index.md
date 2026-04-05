---
title: Contributor Guide
description: How to contribute to claude-code-Go
---

# Contributor Guide

Thank you for your interest in contributing to claude-code-Go! This guide will help you set up your development environment, run tests, and submit pull requests.

---

## Development Setup

### Prerequisites

- **Go 1.24 or later** — [Install Go](https://go.dev/doc/install)
- **Git** — Version control
- **Make** — Build automation (optional, but recommended)

### Clone the Repository

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
```

### Install Dependencies

```bash
go mod download
```

### Build the Project

```bash
# Using make (recommended)
make build

# Or using go build directly
go build -o bin/go-code ./cmd/go-code
```

### Verify Installation

```bash
./bin/go-code --help
```

---

## Running Tests

### Run All Tests

```bash
make test

# Or directly
go test -v ./...
```

### Run Specific Tests

```bash
# Test a specific package
go test -v ./internal/tool/builtin/

# Test with coverage
go test -cover ./...
```

### Run Tests in Watch Mode

```bash
# Install gotest
go install github.com/gotestyourself/gotest@latest

# Watch mode
gotest -w ./...
```

---

## Code Style Guidelines

### Go Code Conventions

1. **Format code** before committing:
   ```bash
   go fmt ./...
   ```

2. **Run static analysis**:
   ```bash
   go vet ./...
   ```

3. **Use meaningful names** for variables and functions

4. **Add comments** for exported functions and types:
   ```go
   // Execute runs the tool with the given input and returns a result.
   func (t *Tool) Execute(ctx context.Context, input map[string]any) Result {
       // ...
   }
   ```

5. **Keep functions focused** and concise

6. **Handle errors explicitly**, don't ignore them:
   ```go
   if err != nil {
       return tool.Error(err.Error())
   }
   ```

### File Organization

```
claude-code-Go/
├── cmd/go-code/          # Main entry point
├── internal/
│   ├── agent/            # Agent loop + context management
│   ├── api/              # API client + SSE streaming
│   ├── config/           # Configuration loading
│   ├── permission/      # Permission system
│   └── tool/             # Tool interface + builtins
│       └── builtin/     # Built-in tool implementations
├── pkg/                  # Public packages
└── docs/                 # Documentation
```

---

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear and structured commit messages.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code style (formatting, no logic change) |
| `refactor` | Code refactoring |
| `test` | Adding or updating tests |
| `chore` | Maintenance, dependencies, build changes |

### Examples

```
feat(agent): add context compression for long sessions

fix(tool): correct glob pattern matching for hidden files

docs(readme): update installation instructions

refactor(api): simplify SSE streaming parser

chore: update go.mod dependencies
```

---

## Submitting a Pull Request

### Before Submitting

1. **Ensure tests pass:**
   ```bash
   go test -v ./...
   ```

2. **Run static analysis:**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **Update documentation** if your change affects usage

4. **Keep changes focused** — one PR per feature or fix

### PR Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code guidelines

3. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat(tool): add new tool"
   ```

4. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Open a pull request** on GitHub

### PR Description Template

```markdown
## Summary
Brief description of what this PR does.

## Related Issues
Link to any related issues (e.g., "Fixes #123")

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring

## Testing
How did you test your changes?

## Checklist
- [ ] Tests pass
- [ ] Code is formatted
- [ ] Documentation updated (if applicable)
```

---

## Other Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make install` | Install to `$GOPATH/bin` |
| `make test` | Run all tests |
| `make vet` | Run go vet |
| `make build-all` | Build for all platforms |
| `make clean` | Remove build artifacts |
| `make docs` | Serve documentation locally |

---

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/strings77wzq/claude-code-Go/blob/main/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

---

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)
- **Discord**: Join our community

---

## Thank You

Your contributions make open source a better place. Thank you for taking the time to contribute to claude-code-Go!

---

## Related Documentation

- [CONTRIBUTING.md](https://github.com/strings77wzq/claude-code-Go/blob/main/CONTRIBUTING.md) — Project contribution guidelines
- [Architecture Overview](../architecture/overview.md) — System architecture
- [Tool System](../tools/overview.md) — Tool development guide