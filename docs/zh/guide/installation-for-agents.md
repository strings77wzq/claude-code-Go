---
title: AI Agent 安装指南
description: AI Agent 帮助用户安装和配置 go-code 的分步说明
---

# AI Agent 安装指南

如果您是帮助用户设置 go-code 的 AI Agent，请按照以下步骤操作。

## 步骤 0：询问用户的 API 设置

询问用户他们想使用哪个 LLM 提供商：

1. **使用哪家提供商？**
   - Anthropic (Claude) → 需要以 `sk-ant-` 开头的 `ANTHROPIC_API_KEY`
   - OpenAI (GPT) → 需要以 `sk-` 开头的 `OPENAI_API_KEY`
   - 自定义（OpenAI 兼容）→ 需要 base URL + API key

2. **您已经有 API key 了吗？**
   - 如果有 → 请他们粘贴
   - 如果没有 → 引导他们从提供商官网获取

## 步骤 1：安装 go-code

使用一键安装脚本：

**Linux / macOS：**
```bash
curl -fsSL https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.sh | bash
```

**Windows (PowerShell)：**
```powershell
irm https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.ps1 | iex
```

脚本会：
1. 检测操作系统和架构
2. 下载对应的预编译二进制文件
3. 安装到 PATH
4. 自动启动设置向导

### 替代方式：go install

```bash
go install github.com/strings77wzq/claude-code-Go/cmd/go-code@latest
```

执行 `go install` 后，运行 `go-code --setup` 进行配置。

## 步骤 2：配置 API key

设置向导（`go-code --setup`）将引导用户完成：
1. 选择提供商
2. 输入 API key（格式验证）
3. 选择模型
4. 写入 `~/.go-code/settings.json`

如果用户跳过了向导，请手动配置：

```bash
mkdir -p ~/.go-code
cat > ~/.go-code/settings.json << 'EOF'
{
  "apiKey": "sk-ant-your-key-here",
  "provider": "anthropic",
  "model": "claude-sonnet-4-20250514"
}
EOF
```

## 步骤 3：验证并首次运行

```bash
# 验证安装
go-code --help

# 如果尚未完成，运行设置向导
go-code --setup

# 开始使用 go-code
go-code
```

用户应该能看到 REPL 欢迎界面，并能够与 agent 进行交互。
