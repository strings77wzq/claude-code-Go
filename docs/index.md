---
layout: home
title: claude-code-Go
titleTemplate: Claude Code in Go — AI-powered coding assistant

hero:
  name: claude-code-Go
  text: Claude Code in Go
  tagline: AI-powered coding assistant with a full agent loop, tool execution, and permission management — built in Go.
  image:
    src: /logo.svg
    alt: claude-code-Go Logo
  actions:
    - theme: brand
      text: Get Started
      link: /en/guide/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/strings77wzq/claude-code-Go

features:
  - icon: 🔄
    title: Agent Loop
    details: Autonomous "think → act → observe" cycle driven by stop_reason dispatch. Handles tool_use, end_turn, and max_tokens seamlessly.
  - icon: 🛠️
    title: 6 Built-in Tools
    details: Read, Write, Edit, Glob, Grep, and Bash — a complete toolset for software engineering tasks out of the box.
  - icon: 🔒
    title: Permission System
    details: Three-tier permission model (ReadOnly / WorkspaceWrite / DangerFullAccess) with rule-based matching and session memory.
  - icon: 🔌
    title: MCP Integration
    details: Model Context Protocol support with stdio transport, JSON-RPC client, and automatic tool discovery from external servers.
  - icon: 🌊
    title: SSE Streaming
    details: Real-time token-by-token response streaming with custom SSE parser. No external dependencies needed.
  - icon: 🧠
    title: Context Management
    details: Intelligent token estimation and automatic conversation compaction. Preserves context while staying within model limits.
---