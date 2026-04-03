---
title: Community
description: Join the claude-code-Go community

---

# Community

Welcome to the claude-code-Go community! This page explains how you can get involved, contribute, and make the most of your experience with the project.

---

## Join the Community

There are several ways to connect with other claude-code-Go users and contributors:

| Channel | Purpose |
|---------|---------|
| [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues) | Bug reports, feature requests |
| [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions) | Q&A, ideas, showcase |
| [Discord](#) | Real-time chat, support |

---

## How to Contribute

### Reporting Issues

Found a bug or have a feature request?

1. Search existing issues to avoid duplicates
2. Create a new issue with a clear description
3. Include reproduction steps for bugs
4. Tag appropriately (bug, enhancement, question)

### Submitting Pull Requests

Want to contribute code? Here's how:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes with clear commits
4. Push to your fork and submit a PR
5. Fill out the PR template completely

### Contributing Docs

Documentation improvements are always welcome! Edit directly in the `docs/` folder and submit a PR.

---

## Creating Custom Skills

One of the most powerful ways to extend claude-code-Go is through the **Skills System**. Skills allow you to create custom commands and reusable workflows.

### Getting Started

1. Create a skill definition file in your config directory
2. Define the command name, description, and action
3. Use it in your conversations with `/your-skill-name`

### Example

```yaml
skills:
  - name: review-pr
    description: Review a GitHub Pull Request
    action: |
      1. Fetch the PR diff
      2. Analyze code changes
      3. Check for common issues
      4. Provide feedback
```

### Share Your Skills

Created something useful? Share it on [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)!

---

## Contributor Guide

### Code of Conduct

We follow the [Contributor Covenant](https://www.contributor-covenant.org/). Please be respectful and inclusive.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go

# Build
make build

# Run tests
make test
```

### Coding Standards

- Follow Go idioms (refer to [Effective Go](https://golang.org/doc/effective_go.html))
- Add tests for new features
- Update documentation accordingly
- Use meaningful commit messages

### Submitting Changes

1. Ensure all tests pass (`go test ./...`)
2. Run linter (`make vet`)
3. Update CHANGELOG if applicable
4. Submit PR with clear description

---

## Recognition

Thank you to all our contributors! Check the [GitHub Contributors](https://github.com/strings77wzq/claude-code-Go/graphs/contributors) page to see who's helping build claude-code-Go.

---

*Want to get more involved? Reach out on Discord or start a discussion on GitHub!*