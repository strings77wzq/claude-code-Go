---
title: 快速开始
description: 运行 go-code 的入门指南，第一次使用前的准备和基本操作
---

# 快速开始

本指南将带您完成首次运行 go-code 的完整流程。

## 前提条件

在运行 go-code 之前，您需要：

1. 一个 Anthropic API 密钥
2. 已编译完成的二进制文件（参见[安装指南](installation.md)）

## 配置 API 密钥

### 方式一：环境变量

```bash
export ANTHROPIC_API_KEY=sk-ant-your-api-key-here
```

建议将此配置添加到您的 shell 配置文件（`~/.bashrc`、`~/.zshrc` 等）中以实现持久化。

### 方式二：配置文件

创建 `~/.config/go-code/settings.json`：

```json
{
  "apiKey": "sk-ant-your-api-key-here"
}
```

配置加载器按以下顺序搜索（后者覆盖前者）：
1. 用户配置：`~/.config/go-code/settings.json`
2. 项目配置：`.go-code/settings.json`
3. 环境变量：`ANTHROPIC_API_KEY`

## 运行 go-code

### 交互式 REPL 模式

启动交互式会话：

```bash
./bin/go-code
```

您将看到欢迎界面：

```
  ____   _    ____ ___ 
 |  _ \ / \  / ___|_ _|
 | |_) / _ \ \___ \| | 
 |  __/ ___ \ ___) | | 
 |_| /_/   \_\____/___|

Welcome to go-code 0.1.0
Type /help for available commands

go-code> 
```

尝试输入一个请求：

```
go-code> 当前目录有哪些文件？
```

智能体将根据需要调用工具来回答您的问题。

### 单次命令模式

用于执行一次性任务：

```bash
./bin/go-code "创建一个 Go 的 hello world 程序"
```

此命令将执行完成后退出。

## 可用命令

在交互模式下，您可以使用以下特殊命令：

| 命令 | 说明 |
|------|------|
| `/help` | 显示可用命令 |
| `/clear` | 清除对话历史 |
| `/exit` | 退出程序 |
| `/quit` | 退出程序（与 /exit 相同）|
| `/model` | 显示当前模型 |

## 启动参数

go-code 支持以下启动选项：

### 位置参数

```bash
go-code [prompt]
```

- `prompt`（可选）：如果提供，go-code 将执行此单次命令然后退出

### 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `ANTHROPIC_API_KEY` | 您的 Anthropic API 密钥 | 是 |
| `ANTHROPIC_BASE_URL` | 覆盖默认 API 端点（可选）| 否 |
| `ANTHROPIC_MODEL` | 指定模型（默认：claude-3-5-sonnet-20241022）| 否 |

### 配置文件

创建 `~/.config/go-code/config.yaml` 用于持久化配置：

```yaml
api_key: "sk-ant-your-api-key-here"
model: "claude-3-5-sonnet-20241022"
base_url: "https://api.anthropic.com"
```

## 会话示例

```
$ ./bin/go-code
  ____   _    ____ ___ 
 |  _ \ / \  / ___|_ _|
 | |_) / _ \ \___ \| | 
 |  __/ ___ \ ___) | | 
 |_| /_/   \_\____/___|

Welcome to go-code 0.1.0
Type /help for available commands

go-code> 列出当前目录的文件

[智能体思考中...]
[调用工具：Glob，模式：*]
[工具结果：找到 3 个文件：main.go, Makefile, README.md]

当前目录包含 3 个文件：
- main.go
- Makefile
- README.md
```

## 权限提示

当工具需要执行潜在危险操作时，go-code 会提示您确认权限：

```
go-code> 删除当前目录下的所有文件

⚠️ 此操作将删除 3 个文件。是否批准？(yes/no)：no
```

输入 `yes` 表示批准，输入 `no` 表示拒绝。

权限系统控制：
- 文件删除和覆盖
- Shell 命令执行
- 网络请求
- 其他潜在危险操作

## 下一步

- [配置说明](configuration.md) — 自定义行为配置
- [架构概览](../architecture/overview.md) — 了解系统工作原理
- [智能体循环详解](../architecture/agent-loop.md) — 学习核心执行周期