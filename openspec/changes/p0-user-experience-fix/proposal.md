## Why

三个 P0 问题严重影响用户体验：1) 用户等待 API 响应时无任何反馈，不知道是卡住了还是在思考；2) 错误信息不分类，401/429/timeout 都显示同一个错误；3) 不支持非交互模式，无法用于脚本调试。参考 opencode 的最佳实践，需要立即修复。

## What Changes

### 1. 连接超时反馈
- 等待 >3 秒无响应：显示 "Connecting to API..."
- 等待 >30 秒无响应：显示 "Still connecting... check network or API key"
- 超时 (5 分钟)：显示明确的超时错误

### 2. 详细错误分类
- 401 → "API Key 无效，请检查 ANTHROPIC_API_KEY"
- 403 → "API 访问被拒绝"
- 429 → "请求频率过高，正在重试..."
- 500+ → "服务器错误，请稍后重试"
- Timeout → "请求超时 (5 分钟)，请检查网络和 API Key"
- Network → "网络连接失败，请检查网络"

### 3. 非交互模式
- `go-code -p "prompt"` 直接输出结果后退出
- `go-code -p "prompt" -f json` JSON 格式输出
- `go-code -p "prompt" -q` 静默模式（无 spinner）
- 所有权限自动批准

## Capabilities

### New Capabilities
- `connection-feedback`: 连接超时反馈机制
- `error-classification`: 详细错误分类
- `non-interactive-mode`: 非交互模式（-p 参数）

### Modified Capabilities
- `api-client`: 超时反馈 + 错误分类
- `cli-entry`: 非交互模式支持

## Impact

- 修改文件: `internal/api/client.go`（超时反馈 + 错误分类）
- 修改文件: `pkg/tui/tui.go`（连接状态显示）
- 修改文件: `cmd/go-code/main.go`（-p 参数支持）
- 不影响: Agent Loop、工具系统、权限系统
