## 1. Skills 系统

- [x] 1.1 创建 `internal/skills/types.go`（Skill 结构体定义）
- [x] 1.2 创建 `internal/skills/loader.go`（JSON 加载 + 注册）
- [x] 1.3 创建 `internal/skills/registry.go`（Skills 注册表）
- [x] 1.4 集成到 REPL（`/skills` 列出、`/<skill-name>` 执行）
- [x] 1.5 创建示例 Skills（review-pr、explain-code、write-tests）
- [x] 1.6 创建 Skills 文档（docs/guide/skills.md + 中文版）

## 2. 多 Provider 支持

- [x] 2.1 创建 `internal/provider/interface.go`（Provider 接口）
- [x] 2.2 创建 `internal/provider/anthropic/provider.go`（Anthropic 适配器）
- [x] 2.3 创建 `internal/provider/openai/provider.go`（OpenAI 适配器）
- [x] 2.4 更新 config 支持 provider 配置
- [x] 2.5 更新 REPL 显示当前 provider
- [x] 2.6 创建 Provider 文档（docs/architecture/providers.md + 中文版）

## 3. Session 恢复

- [x] 3.1 创建 `internal/session/list.go`（列出可用 sessions）
- [x] 3.2 LoadSession 已存在于 session.go（无需新建 load.go）
- [x] 3.3 集成到 REPL（`/resume <id>`、`/sessions` 命令）
- [x] 3.4 创建 Session 管理文档

## 4. 更多内置工具

- [x] 4.1 创建 `internal/tool/builtin/diff.go`（文件差异对比）
- [x] 4.2 创建 `internal/tool/builtin/tree.go`（目录树展示）
- [x] 4.3 创建 `internal/tool/builtin/webfetch.go`（网页内容抓取）
- [x] 4.4 注册新工具到 builtin registry
- [x] 4.5 创建工具文档

## 5. 手动压缩 + 自动更新

- [x] 5.1 创建 `/compact` REPL 命令
- [x] 5.2 创建 `internal/update/checker.go`（GitHub Releases 检查）
- [x] 5.3 创建 `internal/update/downloader.go`（二进制下载 + 替换）
- [x] 5.4 集成 `/update` REPL 命令

## 6. 官网更新

- [x] 6.1 首页增加 Phase 2 新功能卡片（Skills、多 Provider、Session 恢复）
- [x] 6.2 创建 Roadmap 页面（docs/roadmap.md + 中文版）
- [x] 6.3 创建 Community 页面（docs/community.md + 中文版）
- [ ] 6.4 更新功能对比表（Phase 1 vs Phase 2 vs Claude Code 官方）
- [x] 6.5 更新导航栏（新增 Roadmap、Community 链接）
- [x] 6.6 更新侧边栏（新增 Skills 指南、Provider 文档）
