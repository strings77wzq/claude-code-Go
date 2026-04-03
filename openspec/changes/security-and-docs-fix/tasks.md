## 1. 安全修复

- [x] 1.1 创建 `internal/tool/builtin/validate.go`（路径验证工具函数）
- [x] 1.2 修改 Read/Write/Edit 工具调用路径验证
- [x] 1.3 修复 Bash 工具竞态条件（sync.WaitGroup）
- [x] 1.4 修复 API 客户端 Response Body 泄漏
- [x] 1.5 修复 MCP Transport 进程资源泄漏
- [x] 1.6 修复信号处理为优雅关闭

## 2. 文档翻译修复

- [x] 2.1 全局替换 "驾驭系统" → "Harness"（中文文档）
- [x] 2.2 全局替换 "技能" → "Skills"（中文文档）
- [x] 2.3 统一 Agent/Provider/Token 术语
- [x] 2.4 修复 awkward 翻译（"人在循环审批" → "Human-in-the-loop 审批"等）

## 3. 文档一致性修复

- [x] 3.1 统一工具数量（6 → 9，所有文档）
- [x] 3.2 修复配置加载顺序矛盾
- [x] 3.3 补全缺失中文页面（installation-for-agents）
- [x] 3.4 添加缺失的 frontmatter（skills.md）

## 4. README 修复

- [x] 4.1 修复标题为英文
- [x] 4.2 修复目录结构描述
- [x] 4.3 移除混入的中文

## 5. 构建修复

- [x] 5.1 生成 go.sum
- [x] 5.2 修复 Makefile 语法错误
- [x] 5.3 CI 启用 Go 模块缓存
