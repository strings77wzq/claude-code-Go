## Context

当前只支持 Anthropic 格式 API，根目录杂乱，Harness 缺文档。

## Goals / Non-Goals

**Goals:**
- 支持 Anthropic + OpenAI 兼容格式（覆盖 90% 主流模型）
- 更新默认模型到最新版本
- 清理根目录，提升源码阅读体验
- 补充 Harness 文档

**Non-Goals:**
- 不支持 Gemini 专属格式（Google API 差异太大）
- 不支持 MiniMax/Doubao 专属格式（后续按需添加）
- 不改 Harness 代码（Python 是设计选择）

## Decisions

### 1. Provider 架构

**Decision**: 基于接口抽象，每种格式一个 Provider 实现。

```go
type Provider interface {
    Name() string
    SendMessage(ctx, req) (*Response, error)
    SendMessageStream(ctx, req, callback) (*Response, error)
}
```

- `provider/anthropic/` — 现有代码迁移
- `provider/openai/` — 新实现，兼容所有 OpenAI 格式模型

### 2. 模型配置

**Decision**: 通过 `ANTHROPIC_MODEL` 环境变量指定模型名，Provider 自动识别格式。

```
模型名自动识别 Provider：
├── claude-* → Anthropic Provider
├── gpt-*    → OpenAI Provider
├── deepseek-* → OpenAI Provider
├── qwen-*   → OpenAI Provider
├── glm-*    → OpenAI Provider
└── 其他     → 根据 base_url 判断
```

### 3. 默认模型

**Decision**: 更新为 `claude-sonnet-4-6-20251001`。

### 4. 目录重组

**Decision**: 最小化改动，只移动脚本和删除残留。

```
scripts/
├── install.sh
├── install.ps1
└── launch.sh
```

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| OpenAI Provider 实现不完整 | 先支持基本聊天，后续加工具调用 |
| 目录重组破坏现有引用 | 只移动脚本，不影响 Go 代码 |
| 模型名识别错误 | fallback 到 Anthropic Provider |
