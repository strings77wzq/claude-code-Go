## 1. README 修复

- [x] 1.1 修复项目名称（go-code → claude-code-Go）
- [x] 1.2 修复 clone URL（github.com/user/go-code → github.com/strings77wzq/claude-code-Go）
- [x] 1.3 更新文档命令（mkdocs → vitepress）
- [x] 1.4 更新项目结构描述（体现 en/zh 双语）
- [x] 1.5 添加 Star 引导和使用场景

## 2. 官网优化

- [x] 2.1 英文首页 hero 去冗余（去掉 "Claude Code in Go" 副标题）
- [x] 2.2 中文首页 hero 优化（突出核心设计理念）
- [x] 2.3 英文首页增加 CTA 底部区域
- [x] 2.4 中文首页增加 CTA 底部区域
- [x] 2.5 英文首页增加使用场景展示
- [x] 2.6 中文首页增加使用场景展示

## 3. Session Persistence

- [x] 3.1 创建 internal/session/session.go（JSONL 格式）
- [x] 3.2 创建 internal/session/session_test.go
- [x] 3.3 集成到 Agent Loop（会话结束后自动保存）

## 4. Hooks 系统

- [x] 4.1 创建 internal/hooks/hooks.go（Hook 接口 + Registry）
- [x] 4.2 创建 internal/hooks/builtin.go（内置 logging hook）
- [x] 4.3 集成到 Agent Loop（executeTools 中调用 pre/post hooks）
