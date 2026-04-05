# Contributing to claude-code-Go

Thank you for your interest in contributing to claude-code-Go! We welcome contributions from the community to make this AI coding assistant even better.

## Welcome Message

Whether you're reporting a bug, proposing a feature, or submitting a pull request, your contributions are valued. Please read this guide to understand how to contribute effectively.

## How to Contribute

### Bug Reports

If you find a bug, please help us by reporting it:

1. Check if the issue already exists
2. Create a new issue with a clear title and description
3. Include:
   - Steps to reproduce the issue
   - Expected behavior vs actual behavior
   - Go version (`go version`)
   - OS and environment details
   - Any relevant logs or error messages

### Feature Requests

We'd love to hear your ideas for new features:

1. Describe the feature and its use case
2. Explain why this feature would be valuable
3. Include any mockups or examples if applicable

### Pull Requests

For code contributions:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes following our code guidelines
4. Run tests to ensure everything passes
5. Commit your changes with clear commit messages
6. Push to your fork and submit a pull request

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git

### Clone and Setup

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
```

### Build

```bash
make build
# or
go build -o bin/go-code ./cmd/go-code
```

### Test

```bash
make test
# or
go test -v ./...
```

### Other Development Commands

| Command | Description |
|---------|-------------|
| `make install` | Install to `$GOPATH/bin` |
| `make vet` | Run go vet for static analysis |
| `make build-all` | Build for all platforms |
| `make clean` | Remove build artifacts |

## Code Style Guidelines

### Go Code

- Follow standard Go conventions (run `go fmt` before committing)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and concise
- Handle errors explicitly, don't ignore them

### Formatting

```bash
# Format code before committing
go fmt ./...
go vet ./...
```

### Testing

- Add tests for new features and bug fixes
- Ensure all tests pass before submitting PR
- Keep test coverage meaningful, not just for coverage metrics

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

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
| `chore` | Maintenance, deps, build changes |

### Examples

```
feat(agent): add context compression for long sessions
fix(tool): correct glob pattern matching for hidden files
docs(readme): update installation instructions
refactor(api): simplify SSE streaming parser
```

## Pull Request Process

### Before Submitting

1. **Ensure tests pass**: Run `go test -v ./...`
2. **Run static analysis**: Run `go vet ./...` and `go fmt ./...`
3. **Update documentation**: If your change affects usage, update relevant docs
4. **Keep changes focused**: One PR per feature or fix

### PR Description

Include in your PR description:

- **Summary**: What does this PR do?
- **Related issues**: Link to any related issues (e.g., "Fixes #123")
- **Type of change**: Bug fix, feature, docs, refactor, etc.
- **Testing**: How did you test your changes?

### Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, your PR will be merged

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
- **Discussions**: Start a GitHub Discussion

## Thank You

Your contributions make open source a better place. Thank you for taking the time to contribute!
