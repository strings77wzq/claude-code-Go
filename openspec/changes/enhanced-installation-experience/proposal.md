## Why

当前安装体验只有基础的一键脚本和文档，缺少 oh-my-openagent 那样的三层安装体验（For Humans / For LLM Agents / 交互式引导）。用户使用 `go install` 或预编译二进制安装后，没有引导配置 API key，导致首次运行失败。

## What Changes

- 新增 `go-code --setup` 交互式配置向导（Go 代码实现，跨平台）
- 增强 `install.sh`：下载后自动调用 `go-code --setup`
- 新增 `install.ps1`：Windows PowerShell 一键安装脚本
- 新增 `docs/guide/installation-for-agents.md`：给 AI Agent 看的简化版安装指南
- 更新 README 和 installation.md：增加 For LLM Agents 链接

## Capabilities

### New Capabilities
- `setup-wizard`: 交互式配置向导，引导用户选择 provider、输入 API key、验证安装
- `installation-for-agents`: 给 AI Agent 阅读的结构化安装指南
- `windows-installer`: Windows PowerShell 一键安装脚本

### Modified Capabilities
- `installation-experience`: 增强现有安装流程，install.sh 调用 setup 向导

## Impact

- `cmd/go-code/main.go` — 增加 `--setup` 标志处理
- `cmd/go-code/setup.go` — 新增 setup 向导实现
- `install.sh` — 增加 setup 引导调用
- `install.ps1` — 新增 Windows 安装脚本
- `docs/guide/installation-for-agents.md` — 新增 LLM Agent 安装指南
- `README.md` — 安装指南增加 For LLM Agents 链接
- `docs/guide/installation.md` — 增加 For LLM Agents 章节
