---
title: 配置参考
description: go-code 完整配置参考，包括环境变量、settings.json 模式和优先级链
---

# 配置参考

本文档提供 go-code 所有配置选项的完整参考，包括环境变量、配置文件设置和配置源的优先级链。

## 配置优先级

go-code 从多个来源加载配置，优先级顺序如下（从高到低）：

```
1. CLI 参数（最高优先级）
   ↓
2. 环境变量
   ↓
3. 项目配置文件: ./.go-code/settings.json
   ↓
4. 用户配置文件: ~/.go-code/settings.json
   ↓
5. 内置默认值（最低优先级）
```

这意味着您可以在用户配置中设置默认值，然后按项目、环境变量或 CLI 参数逐个覆盖。

---

## 环境变量

### API 配置

| 变量 | 必填 | 默认值 | 描述 |
|----------|----------|---------|-------------|
| `ANTHROPIC_API_KEY` | 是 | - | 用于认证的 API 密钥 |
| `ANTHROPIC_BASE_URL` | 否 | `https://api.anthropic.com` | 覆盖默认 API 端点 |
| `ANTHROPIC_MODEL` | 否 | `claude-sonnet-4-20250514` | 使用的默认模型 |

### MCP 配置

| 变量 | 必填 | 描述 |
|----------|----------|-------------|
| `MCP_SERVER_*` | 否 | 服务器特定的环境变量 |

### 会话配置

| 变量 | 必填 | 默认值 | 描述 |
|----------|----------|---------|-------------|
| `GO_CODE_SESSIONS_DIR` | 否 | `~/.go-code/sessions/` | 会话存储目录 |

### 更新配置

| 变量 | 必填 | 默认值 | 描述 |
|----------|----------|---------|-------------|
| `GO_CODE_UPDATE_URL` | 否 | GitHub releases | 检查更新的 URL |

### 调试

| 变量 | 必填 | 描述 |
|----------|----------|-------------|
| `GO_CODE_TRACE` | 否 | 启用跟踪日志 |
| `GO_CODE_DEBUG` | 否 | 启用调试模式 |

---

## settings.json 模式

配置文件使用 JSON 格式。以下是完整的模式：

### 根对象

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "apiKey": {
      "type": "string",
      "description": "API key for authentication"
    },
    "baseUrl": {
      "type": "string",
      "description": "API endpoint URL",
      "default": "https://api.anthropic.com"
    },
    "model": {
      "type": "string",
      "description": "Default model to use",
      "default": "claude-sonnet-4-20250514"
    },
    "maxTokens": {
      "type": "integer",
      "description": "Maximum tokens per response",
      "default": 4096
    },
    "temperature": {
      "type": "number",
      "description": "Sampling temperature (0-1)",
      "default": 0.7
    },
    "timeout": {
      "type": "integer",
      "description": "Request timeout in seconds",
      "default": 120
    },
    "sessionsDir": {
      "type": "string",
      "description": "Directory for session storage",
      "default": "~/.go-code/sessions"
    },
    "autoSave": {
      "type": "boolean",
      "description": "Auto-save sessions",
      "default": true
    },
    "maxHistorySize": {
      "type": "integer",
      "description": "Maximum history messages to keep",
      "default": 100
    }
  }
}
```

### 示例配置

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514",
  "maxTokens": 8192,
  "temperature": 0.7,
  "timeout": 180,
  "sessionsDir": "~/.go-code/sessions",
  "autoSave": true,
  "maxHistorySize": 50
}
```

---

## 配置文件

### 用户配置

位置：`~/.go-code/settings.json`

这是适用于当前用户所有会话的用户级配置。

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "model": "claude-opus-4-20250514"
}
```

### 项目配置

位置：`./.go-code/settings.json`（在项目根目录）

这是覆盖当前项目用户设置的项目级配置。

```json
{
  "model": "claude-haiku-4-20250514"
}
```

### 优先级示例

给定：
- 用户配置：`{ "apiKey": "user-key", "model": "sonnet" }`
- 项目配置：`{ "model": "opus" }`

结果：
- API 密钥：`user-key`（来自用户配置，未被覆盖）
- 模型：`opus`（来自项目配置，覆盖用户配置）

---

## MCP 配置

### 位置

`~/.go-code/mcp.json`

### 模式

```json
{
  "type": "object",
  "additionalProperties": {
    "type": "object",
    "properties": {
      "command": {
        "type": "string",
        "description": "The executable to run"
      },
      "args": {
        "type": "array",
        "items": {
          "type": "string"
        },
        "description": "Command-line arguments"
      },
      "env": {
        "type": "object",
        "additionalProperties": {
          "type": "string"
        },
        "description": "Environment variables (supports ${VAR} interpolation)"
      }
    },
    "required": ["command"]
  }
}
```

### 示例

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/directory"],
    "env": {
      "HOME": "${HOME}"
    }
  },
  "github": {
    "command": "uvx",
    "args": ["mcp-server-github"],
    "env": {
      "GITHUB_TOKEN": "${GITHUB_TOKEN}"
    }
  }
}
```

### 环境变量插值

MCP 配置支持 `${VAR}` 语法来插入环境变量：

```json
{
  "server": {
    "command": "my-server",
    "env": {
      "API_KEY": "${ANTHROPIC_API_KEY}",
      "HOME": "${HOME}"
    }
  }
}
```

这允许将敏感凭据从主机环境传入，而无需将它们存储在配置文件中。

---

## CLI 参数

| 标志 | 类型 | 描述 |
|------|------|-------------|
| `-p` | string | 要执行的单次提示 |
| `-f` | string | 输出格式：text 或 json |
| `-q` | bool | 安静模式（无旋转动画）|
| `-m` | string | 使用的模型 |
| `-c` | string | 配置文件路径 |

### 示例

```bash
# 单次提示
go-code -p "Explain the code"

# JSON 输出
go-code -p "List files" -f json

# 安静模式
go-code -p "What is 2+2?" -q

# 指定模型
go-code -m claude-opus-4-20250514
```

---

## 权限系统配置

go-code 使用三级权限系统。配置由内部处理，但可以通过以下方式控制：

1. **会话记忆**：权限决策在会话期间被记住
2. **Glob 规则**：自动授予权限的文件路径模式
3. **交互式提示**：需要时向用户请求权限

详情请参阅 [权限系统](../architecture/tools.md#permission-system)。

---

## 相关文档

- [配置指南](../guide/configuration.md) — 用户友好的配置指南
- [工具系统](../tools/overview.md) — 工具执行和权限
- [MCP 集成](../extension/mcp.md) — MCP 服务器配置