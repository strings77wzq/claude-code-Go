---
title: API 错误
description: AI 提供商的 API 错误处理与解决方案
---

# API 错误

处理来自 AI 提供商的错误。

## 速率限制

### "rate_limit_exceeded"

**原因**: 请求过于频繁。

**解决方案**:
- 等待 60 秒后重试
- 升级您的套餐
- 降低请求频率

### "rate_limit_request_enqueued"

**原因**: 因负载过高，请求被排队。

**解决方案**:
- 等待自动重试
- 减少并发请求数

## 认证错误

### "invalid_api_key"

**原因**: API 密钥错误。

**解决方案**:
- 检查密钥格式 (sk-ant-...)
- 生成新密钥
- 验证配置中的密钥

### "permission_denied"

**原因**: 密钥缺少所需权限。

**解决方案**:
- 检查密钥作用域
- 使用正确的权限生成密钥
- 联系 Anthropic 支持

## 内容错误

### "context_length_exceeded"

**原因**: 上下文过长。

**解决方案**:
```
> /compact
> 让我们使用压缩后的上下文继续
```

### "invalid_request_error"

**原因**: 请求格式错误。

**解决方案**:
- 更新到最新版本
- 检查 settings.json 格式
- 如持续出现，请报告错误

## 服务器错误

### "api_error" 或 "server_error"

**原因**: 提供商服务器问题。

**解决方案**:
- 等待后重试
- 检查提供商状态页面
- 尝试其他模型

## 重试逻辑

claude-code-Go 会自动重试:
- 速率限制（使用指数退避）
- 超时（最多 3 次尝试）
- 服务器错误（5xx）

手动重试:
```
> 执行失败了，请重试一次
```

> 英文版本: [API Errors](/troubleshooting/api-errors.md)
