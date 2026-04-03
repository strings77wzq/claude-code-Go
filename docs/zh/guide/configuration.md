---
title: 配置说明
description: go-code 的配置文件格式、环境变量和命令行选项详解
---

# 配置说明

go-code 支持通过配置文件或环境变量进行配置。

## 配置文件位置

- Linux/macOS：`~/.config/go-code/config.yaml`
- Windows：`%APPDATA%\go-code\config.yaml`

## 配置文件格式

```yaml
# API 配置
api_key: "sk-ant-your-api-key-here"

# 模型设置
model: "claude-sonnet-4-20250514"
max_tokens: 4096

# 权限设置
auto_approve_read: true
auto_approve_write: false
auto_approve_bash: false

# MCP 服务器
mcp_servers:
  filesystem:
    command: "npx"
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
  github:
    command: "python"
    args: ["-m", "mcp.server.github", "--token", "your-token"]

# 日志配置
log_level: "info"
log_file: "go-code.log"
```

## 环境变量

| 变量名 | 描述 | 是否必需 |
|--------|------|----------|
| `ANTHROPIC_API_KEY` | 您的 Anthropic API 密钥 | 是 |
| `ANTHROPIC_MODEL` | 使用的模型 | 否（默认为 claude-sonnet-4-20250514）|
| `ANTHROPIC_MAX_TOKENS` | 最大响应 token 数 | 否（默认为 4096）|
| `GO_CODE_LOG_LEVEL` | 日志级别（debug, info, warn, error）| 否（默认为 info）|
| `GO_CODE_CONFIG` | 配置文件路径 | 否 |

## 配置优先级

配置按以下顺序加载，后续配置会覆盖之前的配置：

1. 默认值
2. 配置文件
3. 环境变量
4. 命令行参数

## 命令行参数

```
--help      显示帮助信息
--version   显示版本号
--model     指定模型（覆盖配置文件设置）
--config    指定配置文件路径
```

## 模型选项

可用的模型（取决于 API 可用性）：

- `claude-sonnet-4-20250514`（默认）
- `claude-opus-4-20250514`
- `claude-3-5-sonnet-20241022`

## MCP 配置

Model Context Protocol 服务器可在配置文件中进行配置：

```yaml
mcp_servers:
  server_name:
    command: "path/to/executable"
    args: ["arg1", "arg2"]
    env:
      KEY: "value"
```

更多详细信息，请参阅[工具系统](../architecture/tools.md)中的 MCP 集成部分。