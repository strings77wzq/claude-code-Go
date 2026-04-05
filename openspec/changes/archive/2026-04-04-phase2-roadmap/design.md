## Context

Phase 1 完成了 Agent Loop、6 内置工具、权限系统、MCP、SSE、Session 持久化、Hooks。Phase 2 需要补齐 Skills、多 Provider、Session 恢复、更多工具、自动更新，同时更新官网展示。

## Goals / Non-Goals

**Goals:**
- Skills 系统：自定义命令框架，社区可贡献
- 多 Provider：OpenAI、Gemini 兼容
- Session 恢复：从 JSONL 加载历史
- 3 个新内置工具：Diff、Tree、WebFetch
- 官网新增 Roadmap、Community 页面

**Non-Goals:**
- 不做 VS Code / JetBrains 扩展
- 不做 Desktop App
- 不做云端 Agent
- 不改 Phase 1 核心代码

## Decisions

### 1. Skills 系统设计

**Decision**: 基于 `.go-code/skills/` 目录的 YAML 定义 + Markdown 模板。

```yaml
# .go-code/skills/review-pr.yaml
name: review-pr
description: Review a pull request and provide feedback
prompt: |
  Review the following PR and provide detailed feedback on:
  - Code quality and best practices
  - Potential bugs and edge cases
  - Security concerns
  - Performance implications
```

**Rationale**: 简单、可读、社区可贡献。不需要编译，用户直接编辑 YAML 即可。

### 2. 多 Provider 架构

**Decision**: Provider 接口 + 适配器模式。

```go
type Provider interface {
    Name() string
    SendMessage(ctx context.Context, req *Request) (*Response, error)
    SendMessageStream(ctx context.Context, req *Request, onTextDelta func(string)) (*Response, error)
}
```

每个 Provider 实现这个接口，Agent Loop 不感知具体 Provider。

**Rationale**: 与现有 API 客户端解耦，新增 Provider 只需实现接口。

### 3. Session 恢复

**Decision**: `/resume <session-id>` 命令，从 `~/.go-code/sessions/` 加载 JSONL。

**Rationale**: 利用已有的 Session 持久化，只需加恢复逻辑。

### 4. 官网 Roadmap 页面

**Decision**: 三阶段路线图（Phase 1 ✅、Phase 2 🔄、Phase 3 🔮），用表格展示。

**Rationale**: 让社区看到项目规划，增加信任。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| Skills 系统增加复杂度 | 先做最小版本，只支持 YAML 定义 |
| 多 Provider 测试成本高 | 用 Mock Server 模拟各 Provider |
| 官网内容过多 | 分页面展示，首页保持简洁 |
