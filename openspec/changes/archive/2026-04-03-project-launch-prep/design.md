## Context

项目已完成全部 50 个开发任务（Go 运行时 + Python Harness + 文档 + CI/CD），但代码库中还包含：
- 分析参考项目时的临时目录（owncode-analysis/, claw-code-parity/）
- 缺少 .gitignore 和 LICENSE
- MkDocs 配置需要修复（mkdocs.yml 在 docs/ 子目录而非根目录）
- docs.yml 工作流的 publish_dir 路径需要修正
- 缺少一键上线准备脚本

## Goals / Non-Goals

**Goals:**
- 代码库干净，可直接 git push 到 GitHub
- 文档可部署到 GitHub Pages
- CI/CD 工作流正确运行
- 用户有一键上线脚本和清晰的操作指南

**Non-Goals:**
- 不修改任何 Go 业务逻辑代码
- 不修改 Python Harness 逻辑
- 不添加新功能

## Decisions

### 1. mkdocs.yml 移到根目录

**Decision**: mkdocs.yml 从 docs/ 移到项目根目录。

**Rationale**: MkDocs 标准做法，site 输出在根目录 site/，GitHub Pages 工作流更简单。

### 2. 清理参考项目目录

**Decision**: 删除 owncode-analysis/ 和 claw-code-parity/。

**Rationale**: 这些是调研时的克隆副本，不属于项目代码。.gitignore 也会排除它们。

### 3. launch.sh 一键脚本

**Decision**: 创建 launch.sh 自动化清理、构建、测试、git init。

**Rationale**: 用户只需运行一个脚本 + 手动 push 到 GitHub 即可上线。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 删除目录可能丢失有用信息 | 已归档到 OpenSpec archive，可随时恢复 |
| mkdocs.yml 移动可能破坏本地 docs 服务 | Makefile 已同步更新 |
