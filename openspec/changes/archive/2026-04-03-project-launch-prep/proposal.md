## Why

项目代码已完成（50/50 tasks），但缺少上线所需的最后一步：清理非项目文件、创建 .gitignore/LICENSE、修复 MkDocs 配置、准备上线脚本和文档。

## What Changes

- 添加 .gitignore 排除构建产物和 IDE 文件
- 添加 MIT LICENSE
- 修复 MkDocs 配置（mkdocs.yml 移到根目录）
- 修复 docs.yml 工作流（publish_dir 路径）
- 创建 launch.sh 一键上线准备脚本
- 创建 docs/guide/launch.md 上线指南
- 清理非项目目录（owncode-analysis, claw-code-parity, test/）

## Capabilities

### New Capabilities
- `launch-prep`: 上线准备（清理、gitignore、license、脚本）
- `docs-launch-guide`: 上线指南文档

### Modified Capabilities
<!-- 无 -->

## Impact

- 新增文件: .gitignore, LICENSE, launch.sh, mkdocs.yml (根目录), docs/guide/launch.md
- 修改文件: .github/workflows/docs.yml, Makefile
- 删除目录: owncode-analysis/, claw-code-parity/, test/
