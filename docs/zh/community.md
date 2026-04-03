---
title: 社区
description: 加入 claude-code-Go 社区

---

# 社区

欢迎来到 claude-code-Go 社区！本页面介绍如何参与贡献，以及如何充分利用这个项目。

---

## 加入社区

有多种方式可以与其他 claude-code-Go 用户和贡献者建立联系：

| 渠道 | 用途 |
|---------|---------|
| [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues) | Bug 报告、功能请求 |
| [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions) | 问答、创意展示 |
| [Discord](#) | 实时聊天、技术支持 |

---

## 如何贡献

### 报告问题

发现了 Bug 或有功能建议？

1. 搜索现有 issues 避免重复
2. 创建新 issue，清晰描述问题
3. 对于 Bug，提供复现步骤
4. 添加合适的标签（bug、enhancement、question）

### 提交 Pull Request

想贡献代码？步骤如下：

1. Fork 仓库
2. 创建功能分支（`git checkout -b feature/my-feature`）
3. 提交清晰的更改
4. 推送并提交 PR
5. 完整填写 PR 模板

### 贡献文档

文档改进同样欢迎！直接在 `docs/` 文件夹编辑并提交 PR。

---

## 创建自定义 Skills

扩展 claude-code-Go 最强大的方式之一是 **Skills 系统**。Skills 允许你创建自定义命令和可重用的工作流。

### 开始使用

1. 在配置目录创建 skill 定义文件
2. 定义命令名称、描述和操作
3. 在对话中使用 `/your-skill-name`

### 示例

```yaml
skills:
  - name: review-pr
    description: 审查 GitHub Pull Request
    action: |
      1. 获取 PR diff
      2. 分析代码更改
      3. 检查常见问题
      4. 提供反馈
```

### 分享你的 Skills

创造了有用的东西？在 [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions) 上分享！

---

## 贡献者指南

### 行为准则

我们遵循 [Contributor Covenant](https://www.contributor-covenant.org/)。请保持尊重和包容。

### 开发环境设置

```bash
# 克隆并设置
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go

# 构建
make build

# 运行测试
make test
```

### 代码规范

- 遵循 Go 惯用法（参考 [Effective Go](https://golang.org/doc/effective_go.html)）
- 新功能添加测试
- 相应更新文档
- 使用有意义的提交信息

### 提交更改

1. 确保所有测试通过（`go test ./...`）
2. 运行 linter（`make vet`）
3. 如有需要更新 CHANGELOG
4. 提交 PR 并清晰描述

---

## 致谢

感谢所有贡献者！查看 [GitHub Contributors](https://github.com/strings77wzq/claude-code-Go/graphs/contributors) 页面，了解谁在帮助构建 claude-code-Go。

---

*想更深入参与？在 Discord 上联系或在 GitHub 上发起讨论！*