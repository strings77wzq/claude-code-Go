# go-code: Claude Code Go Implementation

Go 实现的 Claude Code CLI 工具，包含 agent loop、工具系统、权限控制、hooks、LSP 集成。

## 项目结构

```
cmd/go-code/         → 主入口
internal/
  agent/             → Agent 循环 (loop.go ~733行，核心热路径)
  api/               → API 客户端
  tool/              → 工具注册/执行
  permission/        → 权限控制
  hooks/             → Hook 系统
  skills/            → Skill 加载与执行
  config/            → 配置加载
  session/           → 会话管理
  lsp/               → LSP 集成
  provider/          → 模型提供商适配
  telemetry/         → 遥测
  logger/            → 日志
  command/           → Slash command 支持
pkg/
  tty/               → TTY 工具
  tui/                → 终端 UI
tests/integration/   → 集成测试
openspec/            → OpenSpec specs + changes
docs/                → VitePress 文档
harness/             → Python 测试 harness
```

## 常用命令

```bash
make build          # go build -o bin/go-code ./cmd/go-code
make test           # go test -v ./... + harness tests
make vet            # go vet ./...
make build-all      # 全平台交叉编译
make docs           # 启动文档开发服务器
make docs-build     # 构建生产文档
```

## 单独命令

```bash
go test -v ./...                              # 全量测试
go test -v -coverprofile=coverage.out ./...   # 带覆盖率
go test -race ./...                           # 数据竞争检测
go vet ./...                                  # 静态分析
golangci-lint run ./...                       # Lint (如已安装)
go build -o bin/go-code ./cmd/go-code         # 构建
```

## 技术栈

- Go 1.24.2
- 纯 Go 实现，无外部框架依赖
- 测试: Go 标准 testing 包
- 文档: VitePress
- Python harness 用于评估测试

## 关键模块说明

- **agent/loop.go** — 核心 agent 循环，最热路径，11x 引用
- **tool/** — 工具注册与执行框架
- **permission/** — 权限模式 (default/acceptEdits/plan/bypassPermissions)
- **hooks/** — PreToolUse/PostToolUse/Stop 钩子系统
- **skills/** — Skill 加载、并发解析、hook 注入
- **openspec/** — 已有 specs + changes，变更管理使用 OpenSpec 流程

## 开发规则

遵循全局 Harness Constitution (~/.claude/CLAUDE.md) 定义的完整流程。
本项目额外规则：

1. **Go 惯用性**: 遵循 `Skill: golang-patterns`，使用 `Agent: go-reviewer` 审查
2. **并发安全**: 所有 goroutine 必须通过 `go test -race` 验证
3. **错误处理**: 绝不静默吞掉 error，使用 `%w` wrap
4. **测试**: 表驱动测试，AAA 模式，覆盖率 ≥ 80%
5. **文件大小**: 单文件 ≤ 800 行 (agent/loop.go 当前 733 行，接近上限，考虑拆分)
6. **不可变性**: 优先返回新值而非修改指针

## Skill routing

When the user's request matches an available skill, invoke it via the Skill tool. When in doubt, invoke the skill.

Key routing rules:
- Product ideas/brainstorming → invoke /office-hours
- Strategy/scope → invoke /plan-ceo-review
- Architecture → invoke /plan-eng-review
- Design system/plan review → invoke /design-consultation or /plan-design-review
- Full review pipeline → invoke /autoplan
- Bugs/errors → invoke /investigate
- QA/testing site behavior → invoke /qa or /qa-only
- Code review/diff check → invoke /review
- Visual polish → invoke /design-review
- Ship/deploy/PR → invoke /ship or /land-and-deploy
- Save progress → invoke /context-save
- Resume context → invoke /context-restore
