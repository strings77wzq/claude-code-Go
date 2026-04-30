---
title: 常见问题
description: claude-code-Go 常见问题及解决方案
---

# 常见问题

解决经常遇到的使用问题。

## 安装问题

### "command not found: go-code"

**问题**: 二进制文件不在 PATH 环境变量中。

**解决方案**:
```bash
# 检查二进制文件位置
which go-code

# 如果未找到，添加到 PATH
export PATH="$HOME/go/bin:$PATH"

# 永久修复，添加到 ~/.bashrc 或 ~/.zshrc
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
```

### "go install" 因权限被拒绝而失败

**问题**: 对 GOPATH 没有写入权限。

**解决方案**:
```bash
# 检查 GOPATH
go env GOPATH

# 更改所有权
sudo chown -R $(whoami) $(go env GOPATH)

# 或安装到其他位置
go build -o ~/bin/go-code ./cmd/go-code
```

## 配置问题

### "ANTHROPIC_API_KEY not set"

**问题**: 未配置 API 密钥。

**解决方案**:
```bash
# 快速设置
export ANTHROPIC_API_KEY="sk-ant-..."

# 永久设置
echo 'export ANTHROPIC_API_KEY="sk-ant-..."' >> ~/.bashrc

# 或创建配置文件
mkdir -p ~/.go-code
cat > ~/.go-code/settings.json << EOF
{
  "apiKey": "sk-ant-..."
}
EOF
```

### 配置文件未加载

**问题**: JSON 语法错误或位置错误。

**解决方案**:
```bash
# 检查文件位置
ls -la ~/.go-code/settings.json

# 验证 JSON 格式
cat ~/.go-code/settings.json | python3 -m json.tool

# 检查权限
chmod 600 ~/.go-code/settings.json
```

## 连接问题

### "connection refused" 或超时

**问题**: 无法连接到 API。

**解决方案**:
```bash
# 检查网络连接
curl https://api.anthropic.com/v1/health

# 检查防火墙
# 部分公司网络会阻止 API 调用

# 尝试不同的网络
# 使用移动热点进行测试
```

### "invalid api key"

**问题**: API 密钥错误或已过期。

**解决方案**:
```bash
# 验证密钥格式（应以 sk-ant- 开头）
echo $ANTHROPIC_API_KEY

# 测试密钥
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/models

# 在以下地址生成新密钥:
# https://console.anthropic.com/settings/keys
```

## 运行时问题

### 会话无法保存

**问题**: 无法持久化对话记录。

**解决方案**:
```bash
# 检查目录是否存在
mkdir -p ~/.go-code/sessions

# 检查磁盘空间
df -h ~/.go-code

# 检查权限
ls -la ~/.go-code/sessions
```

### 工具执行失败

**问题**: 工具返回错误。

**解决方案**:
1. 检查当前模式: `/mode`
2. 验证权限: `/rules`
3. 检查工具是否存在: `/tools`
4. 查看错误消息了解具体问题

### 内存不足

**问题**: 进程被终止或无响应。

**解决方案**:
```bash
# 检查内存使用情况
ps aux | grep go-code

# 清理旧的会话文件
rm ~/.go-code/sessions/*.jsonl

# 压缩当前会话
> /compact

# 重新开始
> /clear
```

## 获取更多帮助

1. 查看 [API 错误](api-errors.md)
2. 查看 [权限被拒绝](permission-denied.md)
3. 查看 [性能问题](performance-issues.md)
4. 在 [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions) 中提问
5. 搜索 [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)

> 英文版本: [Common Issues](/troubleshooting/common-issues.md)
