---
title: 故障排除指南
description: go-code 常见问题及解决方案
---

# 故障排除指南

本指南涵盖您在使用 go-code 时可能遇到的常见问题，并提供解决方饭。

---

## 常见问题

### API 密钥错误

#### "API key is required"

**问题：** 应用程序找不到有效的 API 密钥。

**解决方案：**

1. **设置环境变量：**
   ```bash
   export ANTHROPIC_API_KEY=sk-ant-your-key-here
   ```

2. **创建配置文件：**
   ```bash
   mkdir -p ~/.go-code
   echo '{"apiKey": "sk-ant-your-key-here"}' > ~/.go-code/settings.json
   ```

3. **验证密钥已设置：**
   ```bash
   echo $ANTHROPIC_API_KEY
   ```

#### "Invalid API key"

**问题：** API 密钥格式不正确。

**解决方案：**

1. **检查密钥格式：** Anthropic 密钥以 `sk-ant-` 开头
2. **重新生成密钥：** 访问 [Anthropic 控制台](https://console.anthropic.com/)
3. **检查额外字符：** 确保密钥中没有引号或空格

---

### 网络错误

#### "Connection timeout"

**问题：** 对 API 的请求超时。

**解决方案：**

1. **检查互联网连接：**
   ```bash
   ping api.anthropic.com
   ```

2. **增加超时时间**（通过配置）：
   ```json
   { "timeout": 300 }
   ```

3. **检查防火墙/代理**设置
4. **尝试不同的网络**来诊断问题

#### "Network error: dial tcp"

**问题：** 无法建立与 API 的 TCP 连接。

**解决方案：**

1. **验证 base URL** 是否正确
2. **检查 DNS 解析：**
   ```bash
   nslookup api.anthropic.com
   ```
3. **禁用 VPN/代理**（如果导致问题）
4. **检查公司防火墙**限制

---

### 权限被拒绝

#### "Permission denied" for file operations

**问题：** 工具无法访问或修改文件。

**解决方案：**

1. **检查文件权限：**
   ```bash
   ls -la /path/to/file
   ```

2. **验证路径在工作目录内：**
   - go-code 限制文件访问到工作目录树
   - 确保文件路径是相对路径或在项目内

3. **授予 Bash 执行权限：**
   ```bash
   chmod +x /path/to/script.sh
   ```

---

### 模型未找到

#### "Model not found: xxx"

**问题：** 指定的模型不可用。

**解决方案：**

1. **列出可用模型：**
   ```
   /models
   ```

2. **使用其他模型：**
   ```
   /model claude-sonnet-4-20250514
   ```

3. **检查模型名称拼写**
4. **验证 API 订阅**包含该模型

---

### 会话错误

#### "Failed to save session"

**问题：** 无法保存当前会话。

**解决方案：**

1. **检查会话目录是否存在：**
   ```bash
   mkdir -p ~/.go-code/sessions
   ```

2. **检查写入权限：**
   ```bash
   ls -la ~/.go-code/
   ```

3. **检查磁盘空间**

#### "Session not found"

**问题：** 无法恢复指定的会话。

**解决方案：**

1. **列出可用会话：**
   ```
   /sessions
   ```

2. **检查会话文件是否存在：**
   ```bash
   ls ~/.go-code/sessions/
   ```

3. **会话可能已被删除**

---

## 错误代码参考

| 错误代码 | 描述 | 解决方案 |
|------------|-------------|----------|
| `E001` | 需要 API 密钥 | 设置 `ANTHROPIC_API_KEY` 环境变量 |
| `E002` | API 密钥无效 | 验证密钥格式，必要时重新生成 |
| `E003` | 网络超时 | 检查连接并增加超时时间 |
| `E004` | 权限被拒绝 | 检查文件权限和工作目录 |
| `E005` | 模型未找到 | 使用 `/models` 列出可用模型 |
| `E006` | 会话未找到 | 使用 `/sessions` 列出有效的会话 ID |
| `E007` | 文件太大 | 最大文件大小为 200KB |
| `E008` | JSON 无效 | 检查配置文件语法 |
| `E009` | MCP 服务器错误 | 检查 MCP 服务器配置 |
| `E010` | 上下文超出限制 | 使用 `/compact` 压缩上下文 |

---

## 常见问题解答

### 通用问题

**问：如何获取 API 密钥？**
> 访问 [Anthropic 控制台](https://console.anthropic.com/) 并创建 API 密钥。

**问：为什么我的响应很慢？**
> 检查您的网络连接。大型响应或复杂任务需要更长时间。

**问：我可以离线使用 go-code 吗？**
> 不行，go-code 需要 API 连接才能让模型生成响应。

**问：会话保存在哪里？**
> 会话默认保存在 `~/.go-code/sessions/` 中。

### 配置问题

**问：如何切换模型？**
> 在 REPL 中使用 `/model` 命令：
> ```
> /model claude-opus-4-20250514
> ```

**问：如何设置不同的 base URL？**
> 设置 `ANTHROPIC_BASE_URL` 环境变量或在配置中设置 `baseUrl`。

**问：可以使用多个 API 密钥吗？**
> 不能同时使用。您可以通过更新配置来切换密钥。

### 工具问题

**问：为什么无法读取文件？**
> 检查文件是否在工作目录内且小于 200KB。

**问：为什么我的写操作被阻止？**
> 该工具需要权限批准。确认提示时批准。

**问：如何使用 MCP 服务器？**
> 在 `~/.go-code/mcp.json` 中配置 MCP 服务器。请参阅 [MCP 集成](../extension/mcp.md)。

### 故障排除问题

**问：如何调试问题？**
> 设置 `GO_CODE_TRACE=true` 环境变量以获取详细日志。

**问：如何报告错误？**
> 在 [GitHub](https://github.com/strings77wzq/claude-code-Go/issues) 上提交问题，包含：Go 版本、操作系统、错误消息和复现步骤。

**问：日志存储在哪里？**
> 日志写入 REPL 的 stdout。使用跟踪模式获取更多详细信息。

---

## 获取帮助

如果遇到此处未涵盖的问题：

1. **查看文档** — 请参阅本页底部的相关文档
2. **搜索现有问题** — [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
3. **在讨论区提问** — [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)
4. **报告错误** — 包括：Go 版本、操作系统、错误消息和复现步骤

---

## 相关文档

- [配置指南](../guide/configuration.md) — 配置选项
- [工具系统](../tools/overview.md) — 内置工具
- [MCP 集成](../extension/mcp.md) — MCP 服务器
- [会话管理](../guide/session-management.md) — 会话持久化