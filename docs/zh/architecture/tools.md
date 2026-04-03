---
title: 工具系统
description: go-code 内置工具详解及 MCP 协议集成
---

# 工具系统

go-code 内置六个核心工具，提供文件系统和命令执行能力。

## 工具概览

| 工具 | 描述 | 风险等级 |
|------|------|----------|
| `Read` | 读取文件内容 | 低 |
| `Write` | 创建或覆盖文件 | 中 |
| `Edit` | 进行针对性代码修改 | 中 |
| `Glob` | 按模式查找文件 | 低 |
| `Grep` | 搜索文件内容 | 低 |
| `Bash` | 执行 shell 命令 | 高 |

## Read（读取）

从文件系统中读取文件内容。

**工具模式：**
```json
{
  "name": "Read",
  "description": "读取文件内容。当需要检查文件内容时使用此工具。",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "要读取的文件路径"
      }
    },
    "required": ["file_path"]
  }
}
```

**调用示例：**
```
工具：Read
参数：{"file_path": "src/main.go"}
```

## Write（写入）

创建新文件或覆盖现有文件。

**工具模式：**
```json
{
  "name": "Write",
  "description": "创建新文件或用新内容覆盖现有文件。",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "要写入的文件路径"
      },
      "content": {
        "type": "string",
        "description": "要写入文件的内容"
      }
    },
    "required": ["file_path", "content"]
  }
}
```

**注意：** 除非启用自动批准，否则需要用户确认。

## Edit（编辑）

对现有文件进行针对性修改。采用基于行号的修改方式。

**工具模式：**
```json
{
  "name": "Edit",
  "description": "对现有文件进行针对性编辑。用新内容替换指定行。",
  "parameters": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "要编辑的文件路径"
      },
      "old_string": {
        "type": "string",
        "description": "要替换的确切文本（必须与文件内容匹配）"
      },
      "new_string": {
        "type": "string",
        "description": "替换后的文本"
      }
    },
    "required": ["file_path", "old_string", "new_string"]
  }
}
```

**注意：** 除非启用自动批准，否则需要用户确认。

## Glob（全局匹配）

查找与 glob 模式匹配的文件。

**工具模式：**
```json
{
  "name": "Glob",
  "description": "查找与 glob 模式匹配的文件。",
  "parameters": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "Glob 模式（例如：*.go, **/*.ts）"
      }
    },
    "required": ["pattern"]
  }
}
```

**模式示例：**
- `*.go` — 当前目录下的所有 Go 文件
- `**/*.js` — 递归查找所有 JavaScript 文件
- `src/**/*.ts` — src 目录下的 TypeScript 文件

## Grep（搜索）

在文件中搜索文本内容。

**工具模式：**
```json
{
  "name": "Grep",
  "description": "在文件中搜索文本模式。",
  "parameters": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "要搜索的正则表达式"
      },
      "path": {
        "type": "string",
        "description": "搜索路径（文件或目录）"
      }
    },
    "required": ["pattern", "path"]
  }
}
```

## Bash（命令行）

执行 shell 命令。

**工具模式：**
```json
{
  "name": "Bash",
  "description": "执行 shell 命令。",
  "parameters": {
    "type": "object",
    "properties": {
      "command": {
        "type": "string",
        "description": "要执行的 shell 命令"
      }
    },
    "required": ["command"]
  }
}
```

**注意：** 高风险操作 — 始终需要用户确认。

## MCP 集成

除了内置工具外，go-code 还支持 Model Context Protocol（MCP）来扩展功能。

### 配置 MCP 服务器

在 `~/.go-code/mcp.json` 中配置：

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path"]
  },
  "github": {
    "command": "python",
    "args": ["-m", "mcp.server.github", "--token", "your-token"]
  }
}

### MCP 工具发现

配置 MCP 服务器后，智能体执行以下步骤：

1. 启动时连接每个 MCP 服务器
2. 请求可用工具列表
3. 在工具注册表中注册发现的工具
4. 通过 MCP 客户端路由 MCP 工具调用

## 工具执行流程

```
模型响应 (含 tool_use)
       │
       ▼
┌─────────────────┐
│ 解析工具调用     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│ 检查工具注册表   │────▶│ 权限系统        │
│ （工具是否存在） │     │                │
└────────┬────────┘     └────────┬────────┘
         │ 是                      │
         ▼                         ▼
┌─────────────────┐     ┌─────────────────┐
│ 执行工具         │◀────│ 用户批准        │
│                 │     │ （如需要）       │
└────────┬────────┘     └─────────────────┘
         │
         ▼
┌─────────────────┐
│ 返回结果给模型  │
└─────────────────┘
```

## 权限系统集成

每个工具都与权限系统集成：

- **低风险**（Read、Glob、Grep）：通常自动批准
- **中风险**（Write、Edit）：需要用户确认或启用自动批准
- **高风险**（Bash）：始终需要用户确认

权限配置示例：

```yaml
auto_approve_read: true   # 自动批准读取操作
auto_approve_write: false # 写入操作需确认
auto_approve_bash: false  # 命令执行需确认
```

## 相关文档

- [Agent Loop 核心](agent-loop.md) — 工具调用的触发机制
- [架构概览](overview.md) — 权限系统详情