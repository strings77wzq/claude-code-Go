---
title: 性能问题
description: 优化 claude-code-Go 性能
---

# 性能问题

优化 claude-code-Go 性能。

## 启动缓慢

### 症状
- 启动需要 5 秒以上
- 长时间显示加载动画

### 解决方案

**检查二进制文件位置**:
```bash
# 应为本地编译的二进制文件，而非 'go run'
file $(which go-code)
# 预期输出: ELF 64-bit executable
```

**清理旧的会话文件**:
```bash
rm ~/.go-code/sessions/*.jsonl
```

**禁用未使用的功能**:
```json
{
  "enableMCP": false,
  "enableLSP": false
}
```

## 内存占用过高

### 症状
- 内存使用超过 500MB
- 系统变慢
- OOM 被杀

### 解决方案

**检查内存使用情况**:
```bash
ps aux | grep go-code
```

**压缩上下文**:
```
> /compact
```

**重新开始**:
```
> /clear
```

**减少历史记录大小**:
```json
{
  "maxHistoryMessages": 20
}
```

## 响应缓慢

### 症状
- 消息之间长时间延迟
- 超时

### 解决方案

**检查网络**:
```bash
ping api.anthropic.com
```

**增加超时时间**:
```json
{
  "timeout": "60s"
}
```

**切换模型**:
```
> /model claude-haiku-4-6-20251001  # 更快的模型
```

**简化请求**:
- 将复杂任务拆分为更小的步骤
- 使用 `/clear` 清理上下文
- 使用具体的文件路径

## Token 使用

### 监控用量

```
> /tokens

上下文: 45,234 / 100,000 tokens (45%)
```

### 降低用量

1. **定期压缩**: `/compact`
2. **重新开始**: `/clear`
3. **更具体**: "读取 main.go" 而非 "读取所有文件"

### 自动压缩

```json
{
  "autoCompactThreshold": 0.7  // 达到 70% 时自动压缩
}
```

## 基准测试

测量性能:

```bash
# 启动时间
time go-code -p "Hello"

# 内存使用
/usr/bin/time -v go-code -p "Hello"
```

> 英文版本: [Performance Issues](/troubleshooting/performance-issues)
