---
title: 配置说明
description: go-code 的配置文件格式、环境变量和 MCP 服务器配置详解
---

# 配置说明

本文档涵盖 go-code 的所有配置选项，包括配置文件、环境变量和 MCP 服务器设置。

## 配置文件位置

go-code 从多个位置加载配置，优先级如下（从高到低）：

1. **CLI 参数**（最高优先级）
2. **环境变量**
3. **项目配置文件**：`./.go-code/settings.json`
4. **用户配置文件**：`~/.go-code/settings.json`
5. **内置默认值**（最低优先级）

这意味着您可以在用户配置中设置默认值，然后通过项目配置、环境变量或 CLI 参数覆盖它们。

## 配置文件格式

配置文件使用 JSON 格式：

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514"
}
```

### 配置项说明

| 字段 | 类型 | 默认值 | 说明 |
|-------|------|---------|-------------|
| `apiKey` | string | (必需) | API 密钥 |
| `baseUrl` | string | `https://api.anthropic.com` | API 端点 URL |
| `model` | string | `claude-sonnet-4-20250514` | 使用的模型 |

## 环境变量

| 变量 | 说明 |
|----------|-------------|
| `ANTHROPIC_API_KEY` | 您的 API 密钥 |
| `ANTHROPIC_BASE_URL` | 覆盖默认 API 端点 |

### 示例：设置环境变量

```bash
# 设置 API 密钥
export ANTHROPIC_API_KEY=sk-ant-your-api-key-here

# 可选：覆盖 API 端点（用于测试或代理）
export ANTHROPIC_BASE_URL=https://custom-api.example.com
```

将这些添加到您的 shell 配置文件中以实现持久化：

```bash
# ~/.bashrc 或 ~/.zshrc
echo 'export ANTHROPIC_API_KEY=sk-ant-xxx' >> ~/.bashrc
source ~/.bashrc
```

## MCP 服务器配置

go-code 支持 Model Context Protocol (MCP) 来扩展外部工具和服务的能力。

### MCP 配置文件位置

MCP 服务器配置存储在：

```
~/.go-code/mcp.json
```

### MCP 配置格式

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

### MCP 配置结构

每个 MCP 服务器条目包含：

| 字段 | 类型 | 说明 |
|-------|------|-------------|
| `command` | string | 要执行的可执行文件 |
| `args` | array | 命令行参数 |
| `env` | object | 环境变量（支持 `${VAR}` 插值） |

### 环境变量插值

MCP 配置支持 `${VAR}` 语法来插值环境变量：

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

这允许从主机环境传递敏感凭据，而无需将它们存储在配置文件中。

## 完整配置示例

### 用户配置 (~/.go-code/settings.json)

```json
{
  "apiKey": "sk-ant-your-api-key-here",
  "baseUrl": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514"
}
```

### 项目配置 (./.go-code/settings.json)

```json
{
  "model": "claude-opus-4-20250514"
}
```

此项目特定配置覆盖模型，同时使用用户配置中的 API 密钥。

### MCP 配置 (~/.go-code/mcp.json)

```json
{
  "filesystem": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "./workspace"]
  },
  "git": {
    "command": "uvx",
    "args": ["mcp-server-git"],
    "env": {
      "GIT_TOKEN": "${GIT_TOKEN}"
    }
  }
}
```

## 运行时切换模型

你可以在不重启 go-code 的情况下运行时切换模型。

### 使用 `/model` 命令

```
go-code> /model
Current model: claude-sonnet-4-20250514

go-code> /model hunyuan-2.0-instruct
Model switched to: hunyuan-2.0-instruct
```

### 使用 `/models` 命令

列出所有可用模型：

```
go-code> /models
Available models:

  Anthropic:
    claude-sonnet-4-20250514 (default)
    claude-opus-4-20250514
    claude-haiku-4-20250514

  Tencent Coding Plan:
    tc-code-latest (Auto)
    hunyuan-2.0-instruct
    hunyuan-2.0-thinking
    minimax-m2.5
    kimi-k2.5
    glm-5
    hunyuan-t1
    hunyuan-turbos

Switch model: /model <model-name>
```

### 腾讯云 Coding Plan 配置

使用腾讯云 Coding Plan 时，配置如下：

```bash
export ANTHROPIC_API_KEY="sk-sp-你的密钥"
export ANTHROPIC_BASE_URL="https://api.lkeap.cloud.tencent.com/coding/anthropic"
export ANTHROPIC_MODEL="tc-code-latest"
```

或在 `~/.go-code/settings.json` 中：

```json
{
  "apiKey": "sk-sp-你的密钥",
  "baseUrl": "https://api.lkeap.cloud.tencent.com/coding/anthropic",
  "model": "tc-code-latest"
}
```

## 故障排除

### "API key is required" 错误

确保您已设置以下任一选项：
- `ANTHROPIC_API_KEY` 环境变量，或
- 配置文件中的 `apiKey`

### 配置文件未找到

验证配置文件是否存在且具有有效的 JSON：

```bash
# 验证 JSON 语法
cat ~/.go-code/settings.json | python -m json.tool
```

### MCP 服务器未加载

检查：
1. MCP 配置文件存在于 `~/.go-code/mcp.json`
2. 命令可执行文件在您的 PATH 中
3. 已安装所需的依赖项

## 相关文档

- [快速开始](quick-start.md) — 基本设置和首次运行
- [架构概览](../architecture/overview.md) — 系统组件
- [工具系统](../architecture/tools.md) — 内置和 MCP 工具
