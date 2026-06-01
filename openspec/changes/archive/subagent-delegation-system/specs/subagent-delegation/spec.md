## ADDED Requirements

### Requirement: Sub-agent types have filtered tool sets

The system SHALL provide three specialized sub-agent types, each with a restricted set of tools appropriate for its role. Explore SHALL have read-only file exploration tools. Review SHALL have read and diff tools. TestGen SHALL have read, write, edit, and bash tools.

#### Scenario: Explore sub-agent tool set
- **WHEN** ExploreToolSet() is called
- **THEN** it returns Read, Glob, Grep, and Tree tools

#### Scenario: Review sub-agent tool set
- **WHEN** ReviewToolSet() is called
- **THEN** it returns Read and Diff tools

#### Scenario: TestGen sub-agent tool set
- **WHEN** TestGenToolSet() is called
- **THEN** it returns Read, Write, Edit, and Bash tools

### Requirement: Task tool enables sub-agent invocation

The system SHALL provide a Task builtin tool that accepts a subagent_type and prompt, spawns the requested sub-agent, and returns structured results.

#### Scenario: Task tool spawns Explore sub-agent
- **WHEN** the Task tool is called with subagent_type "Explore" and a prompt
- **THEN** an Explore sub-agent runs with isolated context and returns its findings

#### Scenario: Task tool rejects invalid sub-agent type
- **WHEN** the Task tool is called with an unknown subagent_type
- **THEN** it returns an error indicating the valid types

### Requirement: Sub-agent execution returns structured results

Sub-agent results SHALL include the sub-agent type label and the trimmed output for clear synthesis by the main agent.

#### Scenario: Structured result format
- **WHEN** a sub-agent completes successfully
- **THEN** the result includes "[{Type} Sub-Agent Result]" prefix followed by the output
