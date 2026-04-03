## Why

当前官网使用 MkDocs 默认主题，纯英文、纯文档风格，无法体现项目的专业性和产品感。参考 opencode.ai 和 schhaohao.github.io/docs/ 的水准，需要一个现代化的产品级官网——双语支持、专业视觉设计、功能展示卡片，让项目看起来像资深开发者研发的产品。

## What Changes

- **BREAKING**: 用 VitePress 替换 MkDocs 作为官网框架
- 首页重构：Hero 区域（项目名 + 标语 + 一键安装命令）+ 功能特性卡片网格
- 中英文双语支持，右上角语言切换按钮
- 专业视觉设计：深色主题、代码高亮、动画效果
- 文档内容从 MkDocs 迁移到 VitePress 格式
- 更新 GitHub Actions docs.yml 工作流适配 VitePress 构建

## Capabilities

### New Capabilities
- `vitepress-site`: VitePress 站点配置、主题定制、多语言支持
- `bilingual-content`: 中英文双语内容体系，首页 hero、功能卡片、文档页面
- `professional-design`: 专业视觉设计（深色主题、代码高亮、动画效果）

### Modified Capabilities
- `docs-website`: 从 MkDocs 迁移到 VitePress，构建流程、部署方式变更

## Impact

- 删除文件: `docs/mkdocs.yml`, 所有 MkDocs 格式文档
- 新增文件: `docs/.vitepress/config.ts`, VitePress 格式的文档（中英文）
- 修改文件: `.github/workflows/docs.yml`（构建命令从 `mkdocs build` 改为 `vitepress build docs`）
- 新增依赖: `vitepress`, `markdown-it`（devDependencies）
- 不影响: Go 源代码、Python Harness、CI 工作流
