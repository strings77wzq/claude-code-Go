# 上线指南

本文档指导你如何将 go-code 项目部署到 GitHub 并上线官网。

## 前置条件

- [Git](https://git-scm.com/) 已安装
- [Go 1.23+](https://go.dev/dl/) 已安装
- [GitHub](https://github.com/) 账号

## 步骤 1：运行一键准备脚本

```bash
chmod +x launch.sh
./launch.sh
```

脚本会自动完成：
- 清理非项目目录（owncode-analysis/, claw-code-parity/, test/）
- 生成 go.sum
- 构建二进制文件
- 运行所有测试
- 初始化 Git 仓库

## 步骤 2：创建 GitHub 仓库

1. 打开 https://github.com/new
2. 填写以下信息：
   - **Repository name**: `go-code`
   - **Description**: `Claude Code in Go — AI-powered coding assistant`
   - **Visibility**: Public（开源项目）
   - ⚠️ **不要勾选** "Add a README file"（我们已有 README.md）
   - ⚠️ **不要勾选** "Add .gitignore"（我们已有 .gitignore）
   - ⚠️ **不要勾选** "Choose a license"（我们已有 LICENSE）
3. 点击 **Create repository**

## 步骤 3：添加远程并推送

```bash
# 替换 YOUR_USERNAME 为你的 GitHub 用户名
git remote add origin git@github.com:YOUR_USERNAME/go-code.git

# 添加所有文件
git add .

# 首次提交
git commit -m "feat: initial release — Claude Code clone in Go

- Agent Loop with tool_use support
- 6 built-in tools (Bash, Read, Write, Edit, Glob, Grep)
- Permission system (ReadOnly/WorkspaceWrite/DangerFullAccess)
- MCP integration (stdio transport)
- Context management (token estimation, auto compaction)
- Python Harness (mock server, parity tests, evaluators)
- MkDocs documentation site
- GitHub Actions CI/CD"

# 重命名分支为 main
git branch -M main

# 推送到 GitHub
git push -u origin main
```

## 步骤 4：启用 GitHub Pages

1. 打开你的仓库页面：`https://github.com/YOUR_USERNAME/go-code`
2. 点击 **Settings**（顶部导航栏）
3. 左侧菜单点击 **Pages**
4. 在 **Build and deployment** 区域：
   - **Source**: 选择 `GitHub Actions`
5. 保存后，`docs.yml` 工作流会自动运行

等待约 1-2 分钟，你的文档网站将上线：
```
https://YOUR_USERNAME.github.io/go-code/
```

## 步骤 5：创建第一个 Release

```bash
# 打标签
git tag v0.1.0
git push origin v0.1.0
```

然后：
1. 打开仓库页面 → **Releases** → **Draft a new release**
2. **Tag version**: 选择 `v0.1.0`
3. **Release title**: `v0.1.0 - Initial Release`
4. **Description**: 填写发布说明
5. 点击 **Publish release**

## 步骤 6：验证

### 验证代码仓库
```
✅ README.md 显示在仓库首页
✅ .gitignore 已生效（bin/ 等目录不在 git 中）
✅ LICENSE 已显示
✅ GitHub Actions 工作流已运行（Actions 标签页查看）
```

### 验证文档网站
```
✅ https://YOUR_USERNAME.github.io/go-code/ 可以访问
✅ 导航栏有 Home, Guide, Architecture, Harness
✅ 页面有 Material 主题样式
```

### 验证 Release
```
✅ Releases 页面有 v0.1.0
✅ 用户可以下载源码
```

## 常见问题

### Q: GitHub Actions 运行失败？
检查 Actions 标签页的日志。常见原因：
- Go 版本不匹配（需要 1.23+）
- Python 依赖未安装（harness/requirements.txt）

### Q: 文档网站 404？
- 确认 Settings → Pages → Source 选择了 `GitHub Actions`
- 等待 2-3 分钟，部署需要时间
- 检查 Actions 标签页中 docs.yml 是否成功运行

### Q: 想自定义域名？
在 Settings → Pages → Custom domain 中添加你的域名，然后在仓库根目录创建 `docs/CNAME` 文件。

## 上线后

上线后，以下功能会自动运行：

| 功能 | 触发条件 | 说明 |
|------|---------|------|
| CI 测试 | 每次 push/PR | Go 测试 + Python Harness 测试 |
| 多平台构建 | 每次 push/PR | linux/darwin 二进制 |
| 文档部署 | docs/ 或 mkdocs.yml 变更 | 自动部署到 GitHub Pages |
| 依赖更新检查 | 每周 | Dependabot 自动检查更新 |
