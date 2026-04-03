---
title: Skills 系统
description: Skills 系统技术深度解析 — JSON 格式、加载机制、自定义 Skills 创建、REPL 集成和最佳实践
---

# Skills 系统

go-code 实现了一个强大的 Skills 系统，允许用户通过命名提示词来扩展和自定义 agent 的行为。本文提供全面的技术概述。

## 什么是 Skills？

**Skills** 是命名提示词，当被调用时会注入到 agent 的系统提示词中。它们提供了以下能力：

- 为常见任务定义可重用的提示词
- 自定义 agent 在特定工作流中的行为
- 创建领域特定的指令
- 通过专门知识增强 agent

当你在 REPL 中使用 `/<skill-name>` 调用 Skills 时，Skills 的提示词会被前置到你的消息中，为任务提供专门的上下文。

## Skills 工作原理

### 加载机制

Skills 在启动时从 `.go-code/skills/` 目录加载：

```
.go-code/skills/
├── review-pr.json
├── explain-code.json
├── write-tests.json
└── custom-skill.json
```

加载器读取此目录中的所有 `.json` 文件并将它们解析为 Skills 结构体：

```go
// Skill represents a custom command skill
type Skill struct {
    Name        string   `json:"name" yaml:"name"`
    Description string   `json:"description" yaml:"description"`
    Prompt      string   `json:"prompt" yaml:"prompt"`
    Examples    []string `json:"examples" yaml:"examples"`
}
```

加载过程：
1. 读取 Skills 目录中的所有条目
2. 仅筛选 `.json` 文件
3. 将每个 JSON 文件解析为 Skills 结构体
4. 验证必填字段是否存在
5. 在注册表中注册有效的 Skills

### JSON 格式

每个 Skills 都定义为 JSON 文件，结构如下：

```json
{
  "name": "skill-name",
  "description": "这个 Skills 的用途",
  "prompt": "注入到 agent 上下文的指令内容",
  "examples": ["/skill-name"]
}
```

| 字段 | 必填 | 说明 |
|-------|----------|-------------|
| `name` | 是 | Skills 的唯一标识符（用于 REPL 命令） |
| `description` | 是 | 列出 Skills 时显示的简要描述 |
| `prompt` | 是 | 调用 Skills 时注入的实际提示词内容 |
| `examples` | 否 | 用法示例（用于文档） |

### REPL 集成

Skills 通过斜杠命令与 REPL 集成：

```
> /skills              # 列出所有可用 Skills
> /review-pr          # 调用 review-pr Skills
> /explain-code       # 调用 explain-code Skills
> /write-tests        # 调用 write-tests Skills
```

调用 Skills 时：
1. 从注册表中获取 Skills 的提示词
2. 将提示词前置到用户消息
3. 将组合消息发送给 agent
4. agent 根据 Skills 的专门上下文做出响应

## 创建自定义 Skills

### 步骤指南

1. **创建 Skills 目录**（如果不存在）：
   ```bash
   mkdir -p ~/.go-code/skills
   ```

2. **为 Skills 创建 JSON 文件**：
   ```bash
   touch ~/.go-code/skills/my-skill.json
   ```

3. **使用 JSON 格式定义 Skills**：
   ```json
   {
     "name": "my-skill",
     "description": "这个 Skills 的用途描述",
     "prompt": "你的自定义提示词内容...",
     "examples": ["/my-skill"]
   }
   ```

4. **重启 go-code** 以加载新 Skills

### 示例：代码审查 Skills

创建 `~/.go-code/skills/review-pr.json`：

```json
{
  "name": "review-pr",
  "description": "审查拉取请求的代码质量和问题",
  "prompt": "你正在进行代码审查。仔细分析提供的代码变更，就以下方面提供建设性反馈：\n\n1. 代码质量和可读性\n2. 潜在的 bug 或边界情况\n3. 性能考虑\n4. 安全漏洞\n5. 测试覆盖率\n\n请具体指出行号并提出改进建议。",
  "examples": ["/review-pr"]
}
```

### 示例：代码解释 Skills

创建 `~/.go-code/skills/explain-code.json`:

```json
{
  "name": "explain-code",
  "description": "详细解释代码工作原理",
  "prompt": "详细解释以下代码。涵盖：\n\n1. 代码的用途（总体目的）\n2. 代码如何工作（逐步逻辑）\n3. 关键函数及其作用\n4. 使用的有趣模式或习惯用法\n5. 潜在的改进或替代方案\n\n使用清晰的语言并在有帮助的地方提供示例。",
  "examples": ["/explain-code"]
}
```

### 示例：测试生成 Skills

创建 `~/.go-code/skills/write-tests.json`:

```json
{
  "name": "write-tests",
  "description": "为给定代码编写全面的测试",
  "prompt": "为提供的代码编写全面的测试。涵盖：\n\n1. 单个函数/方法的单元测试\n2. 边界情况和错误条件\n3. 适当的集成测试\n4. 使用该语言的适当测试框架\n5. 包含清晰的测试名称和文档\n\n确保测试是可维护的并遵循最佳实践。",
  "examples": ["/write-tests"]
}
```

### 示例：重构 Skills

创建 `~/.go-code/skills/refactor.json`:

```json
{
  "name": "refactor",
  "description": "重构代码以提高质量",
  "prompt": "重构以下代码以改进：\n\n1. 可读性 - 清晰的变量名、良好的格式\n2. 性能 - 优化昂贵操作\n3. 可维护性 - 清晰的结构、降低复杂性\n4. 可测试性 - 更易于单元测试\n5. DRY 原则 - 消除代码重复\n\n保留原始功能并确保所有现有测试继续通过。",
  "examples": ["/refactor"]
}
```

## 最佳实践

### 编写有效的 Skills

1. **具体明确**：定义针对特定任务的清晰、专注的提示词
2. **使用上下文**：包含关于领域或任务类型的相关上下文
3. **提供结构**：使用编号列表或章节来组织期望
4. **设定期望**：明确定义期望的输出格式或质量
5. **保持更新**：如果进行重大更改，请对 Skills 进行版本控制

### Skills 组织

- **将相关 Skills 分组**在同一目录
- **使用一致的命名约定**（例如，`动词-名词` 模式）
- **为复杂 Skills 添加文档**，包含详细描述
- **用实际用例测试** Skills

### 常见模式

```json
{
  "name": "security-audit",
  "description": "对代码进行安全审计",
  "prompt": "进行全面的安全审计，重点关注：\n- 输入验证\n- 身份验证/授权\n- 数据保护\n- 常见漏洞（OWASP Top 10）\n\n提供具体的发现和严重程度级别。",
  "examples": ["/security-audit"]
}
```

## 架构概述

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Skills 系统架构                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   REPL      │─────────▶│  Skills 注册表     │                     │
│   │  (/skills)  │          │                 │                     │
│   └─────────────┘          └────────┬────────┘                     │
│                                      │                               │
│                                      ▼                               │
│                              ┌─────────────────┐                     │
│                              │  Skills 加载器     │                     │
│                              │  (JSON 文件)   │                     │
│                              └────────┬────────┘                     │
│                                       │                              │
│                                       ▼                              │
│   ┌─────────────┐          ┌─────────────────┐                     │
│   │   Agent     │◄─────────│  Skills 目录       │                     │
│   │ (注入的)    │          │ ~/.config/...   │                     │
│   └─────────────┘          └─────────────────┘                     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 相关文档

- [配置指南](../guide/configuration.md) — 配置文件位置
- [工具系统概述](../tools/overview.md) — 工具接口和注册表
- [Agent 循环实现](../core-code/agent-loop-impl.md) — 工具执行流程

---

<div class="nav-prev-next">

- [扩展概述](./overview.md) ←
- → [钩子系统](./hooks.md)

</div>