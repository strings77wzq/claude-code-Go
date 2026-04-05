---
title: 贡献者指南
description: 如何为 claude-code-Go 贡献力量
---

# 贡献者指南

感谢您有兴趣为 claude-code-Go 贡献力量！本指南将帮助您设置开发环境、运行测试并提交拉取请求。

---

## 开发环境设置

### 前提条件

- **Go 1.24 或更高版本** — [安装 Go](https://go.dev/doc/install)
- **Git** — 版本控制
- **Make** — 构建自动化（可选，但推荐）

### 克隆仓库

```bash
git clone https://github.com/strings77wzq/claude-code-Go.git
cd claude-code-Go
```

### 安装依赖

```bash
go mod download
```

### 构建项目

```bash
# 使用 make（推荐）
make build

# 或直接使用 go build
go build -o bin/go-code ./cmd/go-code
```

### 验证安装

```bash
./bin/go-code --help
```

---

## 运行测试

### 运行所有测试

```bash
make test

# 或直接
go test -v ./...
```

### 运行特定测试

```bash
# 测试特定包
go test -v ./internal/tool/builtin/

# 带覆盖率测试
go test -cover ./...
```

### 监视模式运行测试

```bash
# 安装 gotest
go install github.com/gotestyourself/gotest@latest

# 监视模式
gotest -w ./...
```

---

## 代码风格指南

### Go 代码规范

1. **提交前格式化代码：**
   ```bash
   go fmt ./...
   ```

2. **运行静态分析：**
   ```bash
   go vet ./...
   ```

3. **为变量和函数使用有意义的名称**

4. **为导出的函数和类型添加注释：**
   ```go
   // Execute runs the tool with the given input and returns a result.
   func (t *Tool) Execute(ctx context.Context, input map[string]any) Result {
       // ...
   }
   ```

5. **保持函数专注和简洁**

6. **显式处理错误，不要忽略它们：**
   ```go
   if err != nil {
       return tool.Error(err.Error())
   }
   ```

### 文件组织

```
claude-code-Go/
├── cmd/go-code/          # 主入口
├── internal/
│   ├── agent/            # Agent 循环 + 上下文管理
│   ├── api/              # API 客户端 + SSE 流
│   ├── config/           # 配置加载
│   ├── permission/      # 权限系统
│   └── tool/             # 工具接口 + 内置工具
│       └── builtin/     # 内置工具实现
├── pkg/                  # 公共包
└── docs/                 # 文档
```

---

## 提交信息规范

我们遵循 [Conventional Commits](https://www.conventionalcommits.org/) 以获得清晰且结构化的提交信息。

### 格式

```
<type>(<scope>): <description>

[可选的正文]

[可选的页脚]
```

### 类型

| 类型 | 描述 |
|------|-------------|
| `feat` | 新功能 |
| `fix` | 错误修复 |
| `docs` | 文档更改 |
| `style` | 代码风格（格式化，无逻辑更改）|
| `refactor` | 代码重构 |
| `test` | 添加或更新测试 |
| `chore` | 维护、依赖项、构建更改 |

### 示例

```
feat(agent): add context compression for long sessions

fix(tool): correct glob pattern matching for hidden files

docs(readme): update installation instructions

refactor(api): simplify SSE streaming parser

chore: update go.mod dependencies
```

---

## 提交拉取请求

### 提交前

1. **确保测试通过：**
   ```bash
   go test -v ./...
   ```

2. **运行静态分析：**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **如果您的更改影响使用，请更新文档**

4. **保持更改专注** — 每个 PR 一个功能或修复

### PR 流程

1. **创建功能分支：**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **按照代码指南进行更改**

3. **提交您的更改：**
   ```bash
   git add .
   git commit -m "feat(tool): add new tool"
   ```

4. **推送到您的分支：**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **在 GitHub 上打开拉取请求**

### PR 描述模板

```markdown
## 摘要
简要描述此 PR 的内容。

## 相关问题
链接任何相关问题（例如，"Fixes #123"）

## 更改类型
- [ ] 错误修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 重构

## 测试
您如何测试您的更改？

## 检查清单
- [ ] 测试通过
- [ ] 代码已格式化
- [ ] 文档已更新（如适用）
```

---

## 其他开发命令

| 命令 | 描述 |
|---------|-------------|
| `make build` | 构建二进制文件 |
| `make install` | 安装到 `$GOPATH/bin` |
| `make test` | 运行所有测试 |
| `make vet` | 运行 go vet |
| `make build-all` | 为所有平台构建 |
| `make clean` | 清除构建产物 |
| `make docs` | 本地提供文档 |

---

## 行为准则

请注意，此项目是根据 [贡献者行为准则](https://github.com/strings77wzq/claude-code-Go/blob/main/CODE_OF_CONDUCT.md) 发布的。通过参与此项目，您同意遵守其条款。

---

## 获取帮助

- **问题**： [GitHub Issues](https://github.com/strings77wzq/claude-code-Go/issues)
- **讨论**： [GitHub Discussions](https://github.com/strings77wzq/claude-code-Go/discussions)
- **Discord**： 加入我们的社区

---

## 感谢您

您的贡献让开源变得更好。感谢您抽出时间为 claude-code-Go 贡献力量！

---

## 相关文档

- [CONTRIBUTING.md](https://github.com/strings77wzq/claude-code-Go/blob/main/CONTRIBUTING.md) — 项目贡献指南
- [架构概览](../architecture/overview.md) — 系统架构
- [工具系统](../tools/overview.md) — 工具开发指南