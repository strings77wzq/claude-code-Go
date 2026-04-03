---
title: Introduction
description: What is claude-code-Go? A Go-native implementation of the Claude Code agent architecture
---

# Introduction

Welcome to claude-code-Go—a Go implementation of the Claude Code agent system.

## What is claude-code-Go?

claude-code-Go is a terminal-based AI coding assistant that implements the complete agent loop with built-in tools, a permission system, MCP (Model Context Protocol) support, and SSE (Server-Sent Events) streaming.

Core design philosophy: **Model provides intelligence, harness provides reliability.**

The LLM handles understanding intent, deciding tool calls, and interpreting results. The runtime handles permission control, timeout protection, output truncation, and session persistence—ensuring production-grade reliability.

It operates as:
- An AI assistant that can read, write, and edit files across your project
- An automation tool that executes shell commands for builds, tests, and scripts
- A code navigator that searches your codebase with glob patterns and content search
- An autonomous agent that collaborates with you interactively in a terminal REPL

## How is it Different from Chatbots?

Unlike traditional chatbots that simply respond to messages, claude-code-Go operates as an **autonomous agent** with agency and tool-use capabilities.

| Feature | Chatbot | claude-code-Go |
|---------|---------|----------------|
| **Interaction Model** | Request-Response | Continuous agent loop with tool calls |
| **Capabilities** | Text generation only | Read, write, edit, execute, search files |
| **State Management** | Stateless | Message history + context compaction |
| **User Control** | Passive recipient | Active approval for dangerous operations |
| **Use Case** | Q&A, conversation | Code development, automation |

## Tech Stack

claude-code-Go is built with a minimal, dependency-free philosophy:

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.23+ |
| **Standard Library** | Go stdlib only (no external dependencies in main binary) |
| **Python Harness** | Optional (for testing, mocking, evaluation) |
| **API** | Anthropic Messages API with SSE streaming |

The Python harness (`harness/`) is optional and provides:
- Mock API server for testing without API costs
- Session replay for debugging
- Quality evaluators for integration testing

## Core Numbers

- **~31** Go source files (non-test)
- **~4,000** lines of Go code
- **~2,000** lines of Python (harness)
- **8** internal modules
- **6** built-in tools (Read, Write, Edit, Glob, Grep, Bash)

## What You'll Learn

This documentation covers the architecture and implementation details of claude-code-Go:

1. **AI Agent Architecture** — How the agent loop works, decision-making process, and tool selection
2. **Streaming API Integration** — SSE token streaming from Anthropic's API with real-time output
3. **Tool System** — Built-in tools + MCP integration for extending capabilities
4. **Security Design** — Three-tier permission system for controlling dangerous operations
5. **Engineering Practice** — Structured logging, graceful shutdown, configuration management

## Target Audience

This documentation is designed for:

- **Go developers** interested in building AI-powered tools
- **Systems programmers** curious about agent architecture
- **DevOps engineers** looking for CLI automation with AI capabilities
- **Researchers** exploring LLM application patterns

You should have familiarity with:
- Go programming language basics
- Terminal/command-line interfaces
- HTTP APIs and JSON data formats

## Next Steps

Ready to get started? Head to the [Quick Start](quick-start.md) guide to run your first command.

Or dive deeper into the [Architecture Overview](../architecture/overview.md) to understand how the pieces fit together.