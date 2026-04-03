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

创建 `~/.config/go-code/config.yaml`：

```yaml
api_key: "sk-ant-your-api-key-here"
```

## 交互式 REPL 模式运行

启动交互式会话：

```bash
./bin/go-code
```

您将看到一个提示符，可以在其中输入请求：

```
go-code> 当前目录是什么？
```

智能体将根据需要调用工具来回答您的问题。

## 单次命令模式运行

用于执行一次性任务：

```bash
./bin/go-code "创建一个 Go 的 hello world 程序"
```

此命令将执行完成后退出。

## 会话示例

```
$ ./bin/go-code
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

## 下一步

- [配置说明](configuration.md) — 自定义行为配置
- [架构概览](../architecture/overview.md) — 了解系统工作原理