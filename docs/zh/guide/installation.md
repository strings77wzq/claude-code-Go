---
title: 安装指南
description: 了解如何从源码或预编译二进制包安装 go-code
---

# 安装指南

本指南介绍如何安装 go-code，包括从源码编译和使用预编译二进制包两种方式。

## 环境要求

- Go 1.23 或更高版本
- Python 3.x（可选，仅在使用测试 harness 时需要）

## 从源码安装

### 克隆仓库

```bash
git clone https://github.com/user/go-code.git
cd go-code
```

### 编译构建

```bash
make build
```

执行上述命令后，二进制文件将生成在 `bin/go-code`。

### 安装到系统 PATH

```bash
go install ./cmd/go-code
```

此命令会将二进制文件安装到 `$GOPATH/bin`（通常为 `~/go/bin`）。

## 预编译二进制包

从 releases 页面下载对应平台的二进制文件：

| 平台 | 架构 | 文件名 |
|------|------|--------|
| Linux | amd64 | go-code-linux-amd64 |
| macOS | amd64 | go-code-darwin-amd64 |
| macOS | arm64 | go-code-darwin-arm64 |

```bash
# 示例：下载并安装 Linux amd64 版本
curl -L -o go-code https://github.com/user/go-code/releases/latest/download/go-code-linux-amd64
chmod +x go-code
sudo mv go-code /usr/local/bin/
```

## 验证安装

```bash
# 检查二进制文件是否存在
ls -la bin/go-code

# 运行帮助命令
./bin/go-code --help
```

## 下一步

- [快速开始](quick-start.md) — 运行您的第一个命令
- [配置说明](configuration.md) — 设置 API 密钥和偏好配置