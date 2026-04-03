## Context

Code review 发现 4 类问题需要修复：
1. README 多处硬编码错误（名称、URL、过时命令）
2. 官网 hero 视觉冗余 + 缺少 CTA
3. Session Persistence 和 Hooks 未实现
4. 缺少 Star 引导和使用场景展示

## Goals / Non-Goals

**Goals:**
- README 完全匹配实际仓库（名称、URL、命令、结构）
- 官网 hero 简洁专业，无冗余
- 官网有 CTA 和使用场景
- Session Persistence 和 Hooks 基础实现

**Non-Goals:**
- 不改名（保持 claude-code-Go）
- 不实现完整 Hooks 生态（只做基础框架）
- 不实现 Session 恢复 UI（只做后端逻辑）

## Decisions

### 1. README 修复策略

**Decision**: 全面重写 README，不逐个修复。

**Rationale**: 问题太多（名称、URL、命令、结构都错），逐个改容易遗漏。重写确保一致性。

### 2. 官网 Hero 优化

**Decision**: 去掉 "Claude Code in Go" 副标题，保留 claude-code-Go 作为唯一品牌名。

**Rationale**: 仓库名已经包含 "Go"，重复强调显得不专业。

### 3. Session Persistence 实现

**Decision**: JSONL 格式，每行一个消息，支持 load/save。

**Rationale**: 与 claw-code-parity 一致，简单可靠。

### 4. Hooks 系统实现

**Decision**: 简单的 pre/post 回调接口，在 executeTools 中调用。

**Rationale**: 最小化实现，为后续扩展留接口。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| README 重写可能遗漏内容 | 保留原有架构图文档 |
| Hooks 实现增加代码复杂度 | 只做接口定义，不实现具体 hook |
