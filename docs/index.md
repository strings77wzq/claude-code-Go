# go-code — Claude Code in Go

go-code is a Go implementation of Anthropic's Claude Code (Cline) agent system. It provides a command-line interface for interacting with Large Language Models using the Claude API, featuring a full agent loop with tool execution, permission management, and Model Context Protocol (MCP) support.

## Overview

go-code implements the core functionality of Claude Code in pure Go, offering:

- **Full Agent Loop**: State management, message history, and tool execution
- **6 Built-in Tools**: Read, Write, Edit, Glob, Grep, and Bash commands
- **Permission System**: User approval for dangerous operations
- **MCP Support**: Extensible tool ecosystem via Model Context Protocol
- **Streaming**: Real-time token-by-token response streaming

## Quick Links

- [Installation Guide](guide/install.md)
- [Quick Start](guide/quick-start.md)
- [Configuration](guide/config.md)
- [Architecture Overview](architecture/overview.md)
- [Agent Loop](architecture/agent-loop.md)
- [Built-in Tools](architecture/tools.md)
- [Python Harness](harness/overview.md)

## Why Go?

This implementation leverages Go's strengths:

- **Performance**: Fast execution and low memory footprint
- **Concurrency**: Built-in goroutines for parallel tool execution
- **Cross-Platform**: Single binary deployment across all major platforms
- **Standard Library**: Rich built-in packages for HTTP, JSON, and file operations

## Status

This is an educational implementation demonstrating how Claude Code works under the hood. It is not affiliated with or endorsed by Anthropic.

## License

MIT License - see [LICENSE](https://github.com/user/go-code/blob/main/LICENSE) for details.