# Tutorial 10: Contributing

How to contribute to claude-code-Go.

## Getting Started

### Prerequisites

- Go 1.24+
- Git
- Make (optional)

### Clone and Build

```bash
git clone https://github.com/strings77wc/claude-code-Go.git
cd claude-code-Go
make install
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Follow the coding standards:
- Run `go vet ./...`
- Run `go test ./...`
- Format with `go fmt`

### 3. Run Tests

```bash
make test
```

### 4. Submit PR

Open a pull request with:
- Clear description
- Test coverage
- Documentation updates

## Code Style

- Follow Go conventions
- Add tests for new features
- Document public APIs
- Keep functions focused

## Areas to Contribute

- New tools
- Bug fixes
- Documentation
- Performance improvements
- MCP servers

## Getting Help

- GitHub Discussions
- GitHub Issues
- Code of Conduct applies
