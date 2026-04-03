---
title: Design Philosophy
description: Deep dive into the intelligence/reliability separation — what the model does vs. what the harness guarantees
---

# Design Philosophy

## The Core Principle

**Model provides intelligence, harness provides reliability** — this is not a marketing slogan. It is an architectural constraint that shapes every line of code in go-code.

The principle emerges from a fundamental observation: large language models are remarkably capable at reasoning, planning, and generating appropriate responses, but they lack the properties required for safe, reliable system operation. They have no concept of timeouts, no understanding of resource constraints, no built-in safety guarantees, and no memory across sessions. Meanwhile, traditional software systems excel at these reliability properties but cannot reason about user intent or choose appropriate actions.

go-code sits at this boundary. The model handles the cognitive work — understanding what the user wants, deciding which tools to call, interpreting results, planning the next step. The harness handles the reliability work — ensuring operations execute within bounds, protecting against resource exhaustion, maintaining session state, and enforcing security boundaries.

This separation is not merely convenient. It is necessary because **the model can be wrong, but the harness must never fail**.

## What the Model Provides: Intelligence

When we say the model provides "intelligence," we mean a specific set of capabilities that go-code delegates to the LLM:

### Intent Understanding

The model interprets natural language input and extracts what the user is actually trying to accomplish. "Make this function faster" is not a precise instruction — the model must understand what "faster" means in context, which function is being referenced, and what performance characteristics currently exist. This is an intelligence task that no amount of rule-based parsing can replace.

### Tool Selection

Given an available toolset (Read, Write, Edit, Glob, Grep, Bash, plus any MCP tools), the model decides which tools to call and in what order. This requires understanding tool capabilities, matching them to the current task, and forming a plan. The harness provides the list of available tools, but the model determines which ones to use.

### Result Interpretation

Tool outputs are not self-explanatory. A grep search returns lines of text; the model must interpret whether those results indicate success, failure, or partial progress. A file write operation returns a simple confirmation; the model must decide whether to verify the write, proceed to the next step, or report completion to the user.

### Next-Step Planning

The agent loop is iterative. After each tool execution, the model must decide: continue with the next step, try a different approach, ask for clarification, or declare the task complete. This is a planning and reasoning task that requires understanding the overall goal and current progress.

These four capabilities — intent understanding, tool selection, result interpretation, and next-step planning — constitute what we mean by "intelligence" in this architectural context. They are the cognitive work that the LLM performs, and they are impossible to implement with deterministic code.

## What the Harness Provides: Reliability

Reliability is not a single property but a collection of guarantees that traditional software systems have long provided:

### Permission Control

The model may request any operation, but the harness decides which operations actually execute. Permission control is not advisory — it is a hard gate that blocks operations the user has not approved. The model might ask to delete all files in the project, but the harness intercepts this request and only proceeds if the user's permission policy allows it.

This is the most critical reliability property because the model operates with the user's full privileges. Without permission control, an AI agent would be indistinguishable from a malicious script.

### Timeout Protection

The model has no concept of time. It might call a grep search over a million files and wait indefinitely for results. It might execute a shell command that never terminates. The harness enforces timeout bounds — if a tool doesn't complete within the configured duration, it is terminated and an error is returned to the model.

Without timeout protection, a single misbehaving tool could hang the entire session indefinitely.

### Output Truncation

Tool outputs can be arbitrarily large. Reading a binary file, running a command that produces megabytes of output, or grep searching a massive codebase could produce results larger than the context window. The harness truncates tool outputs to prevent them from consuming the entire available context.

This is not merely a performance optimization — it is a correctness guarantee. Without truncation, the system could fail to provide the model with the information it needs to complete the task.

### Session Persistence

The model is stateless between API calls. The harness maintains conversation history across turns and can persist sessions to disk for later resumption. This allows the user to disconnect and reconnect without losing progress on a long-running task.

### Error Recovery

Tools can fail for reasons unrelated to the model's reasoning — network timeouts, file permission issues, invalid tool inputs. The harness catches these errors, formats them appropriately, and presents them to the model in a way that enables recovery. The model can decide to retry, try a different approach, or report the error to the user.

## Why This Separation Matters

The separation of intelligence and reliability is not an implementation detail — it is a fundamental architectural requirement that stems from the nature of current AI systems.

### The Model Can Be Wrong

LLMs are probabilistic systems. They can hallucinate, misinterpret instructions, choose suboptimal tools, and make mistakes. A model might:
- Call the wrong tool for the task
- Misinterpret a tool's output
- Enter an infinite loop of tool calls
- Misunderstand the user's intent

These failures are expected and should be handled gracefully. The harness cannot prevent model errors, but it can contain their impact.

### The Harness Must Never Fail

Unlike the model, the harness is deterministic software. It should never:
- Crash during permission checks
- Fail to enforce timeouts
- Lose user data due to session failures
- Allow unauthorized operations
- Produce unpredictable behavior

These are the same reliability guarantees that any production system must provide. The harness's job is to maintain these guarantees regardless of what the model does.

### The Interface Between Them

The model and harness communicate through a well-defined interface. The model receives tool definitions (name, description, input schema) and decides which tools to call. The harness executes tools and returns structured results (text output, errors, success/failure status).

This interface is intentionally narrow. The model cannot bypass the harness, cannot inspect harness internals, and cannot modify harness behavior. The harness is a hard boundary, not a soft layer.

## Design Decisions

| Decision | Rationale | Trade-off |
|----------|-----------|-----------|
| **Three-tier permission model** | Binary allow/deny is insufficient for varied development workflows | More complex configuration, but enables practical use |
| **Glob rule matching for paths** | Allows fine-grained control over which files can be modified | Requires users to understand glob patterns |
| **250-character tool descriptions** | Forces discipline in capability communication; prevents context bloat | May omit useful tool details |
| **MAX_TURNS = 50** | Prevents resource exhaustion from infinite loops | Complex tasks may need more iterations |
| **Three-tier message compaction** | Preserves semantic content while managing token limits | Compression is lossy; some context may be lost |
| **Pure Go implementation** | No runtime dependencies; consistent performance | Less flexible than interpreted languages |
| **MCP as external adapter** | Protocol-agnostic tool discovery; ecosystem integration | MCP servers add external dependencies |

## Comparison with Other Approaches

go-code's architecture can be contrasted with several alternative approaches to AI-assisted software development:

### Chatbot Paradigm

In a pure chatbot, the LLM generates text responses without taking direct actions. The user reads the response and executes commands manually. This provides maximum safety but minimum automation.

**go-code's position**: The model takes actions directly, with the harness ensuring safety through permission checks and timeout protection.

### Script Runner Paradigm

In a script runner, predefined scripts execute without LLM involvement. The LLM generates script arguments but has no control over the execution flow.

**go-code's position**: The model controls the execution flow, choosing which tools to call and in what order. The harness ensures execution stays within bounds.

### Full IDE Paradigm

In a full IDE, the AI has complete access to the development environment, including file system, build tools, and debugging capabilities. Safety is managed through user supervision rather than automated checks.

**go-code's position**: Permission control is enforced automatically. Human-in-the-loop approval is available but not required for routine operations.

### Agent Framework vs. Agent Product

Many "agent frameworks" provide just the agent loop and leave safety, persistence, and permission to the user. This is appropriate for experimentation but unsuitable for production use.

**go-code's position**: go-code is an agent *product* that ships with all reliability features included. Users can configure but not remove safety guarantees.

## The Philosophy in Practice

When building features in go-code, the intelligence/reliability separation serves as a decision-making framework:

- If a capability requires **understanding user intent** → model responsibility
- If a capability requires **choosing between options** → model responsibility
- If a capability requires **guaranteeing safety** → harness responsibility
- If a capability requires **resource management** → harness responsibility
- If a capability requires **persistence across sessions** → harness responsibility

This framework isn't perfect — some capabilities straddle the boundary. But it provides a consistent heuristic for architectural decisions.

## Related Documentation

- [Architecture Overview](overview.md) — System components and data flow
- [Agent Loop](agent-loop.md) — State machine mechanics and execution cycle
- [Permission System](tools.md#permission-system) — Three-tier model and policy configuration