## Context

当前安装体验：install.sh 只下载二进制，打印静态提示文字。用户通过 go install 或预编译二进制安装后，没有引导配置 API key。缺少 Windows 安装脚本。缺少给 AI Agent 看的安装指南。

现有代码库：`internal/config/loader.go` 从 `~/.go-code/settings.json` 加载配置，支持 JSON 格式。CLI 入口在 `cmd/go-code/main.go`。

## Goals / Non-Goals

**Goals:**
- 所有安装方式最终都引导用户进入 setup 向导
- setup 向导跨平台（Go 实现，支持 Linux/macOS/Windows）
- API key 只检查格式，不发真实请求
- For LLM Agents 文档简化版，4 步完成
- Windows 用户有 PowerShell 一键安装

**Non-Goals:**
- 不实现 API key 真实验证（不发送请求）
- 不实现复杂的模型选择 UI
- 不修改现有 config loader 逻辑

## Decisions

### Decision 1: setup 向导用 Go 实现而非 shell 脚本

**Why**: `go install` 和预编译二进制用户不会跑 install.sh。Go 实现确保所有安装方式都能用同一个 setup 向导。install.sh 只负责下载 + 调用 `go-code --setup`。

### Decision 2: API key 只检查格式

**Why**: 发真实请求会花用户 token，且需要网络。格式检查足够：Anthropic `sk-ant-`，OpenAI `sk-`，Custom 非空。

### Decision 3: For LLM Agents 文档简化版

**Why**: oh-my-openagent 的安装指南有几百行，包含大量订阅问卷、模型匹配等。我们只需要核心 4 步：问 provider → 安装 → 配置 → 验证。

### Decision 4: install.sh 调用 go-code --setup 而非内嵌逻辑

**Why**: 避免 shell 和 Go 两套 setup 逻辑。install.sh 只负责下载，setup 交给 Go 代码统一处理。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| `go-code --setup` 需要二进制已存在 | install.sh 先下载再调用；手动用户先安装再运行 setup |
| Windows 终端不支持 ANSI 颜色 | setup 向导检测终端能力，降级为无颜色输出 |
| 用户不想立即配置 API key | setup 向导允许跳过，后续可随时 `go-code --setup` 重新配置 |
| install.sh 调用 go-code --setup 失败 | 捕获错误，打印手动配置说明作为 fallback |
