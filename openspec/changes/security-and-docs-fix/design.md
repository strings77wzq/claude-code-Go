## Context

Code Review 发现 5 个严重安全问题和大量文档翻译问题。Phase 1 修复安全问题，Phase 2 修复文档。

## Goals / Non-Goals

**Goals:**
- 修复 5 个严重安全问题
- 修复所有中文文档翻译错误
- 统一术语（Harness、Skills、Agent、Provider）
- 修复工具数量不一致
- 补全缺失页面
- 修复 README

**Non-Goals:**
- 不改变 Agent Loop 核心逻辑
- 不新增功能
- 不改官网样式

## Decisions

### 1. 路径验证策略

**Decision**: 所有文件操作工具（Read、Write、Edit）校验路径在 workingDir 内。

```go
func validatePath(filePath, workingDir string) error {
    absPath, _ := filepath.Abs(filePath)
    absWd, _ := filepath.Abs(workingDir)
    if !strings.HasPrefix(absPath+string(filepath.Separator), absWd+string(filepath.Separator)) && absPath != absWd {
        return fmt.Errorf("path outside working directory")
    }
    return nil
}
```

**Rationale**: 防止读写系统文件，是生产级工具的基本要求。

### 2. Bash 竞态修复

**Decision**: 用 `sync.WaitGroup` 替代 `sync.Mutex` 保护 output。

**Rationale**: goroutine 和主线程的同步应该用 WaitGroup 确保 goroutine 完成后再读取 output，Mutex 只能保护写操作但不能保证顺序。

### 3. 术语统一

**Decision**: 以下术语保留英文不翻译：
- Harness（核心架构概念）
- Skills（功能名称）
- Agent（AI Agent 领域通用）
- Provider（API 领域通用）
- Token（LLM 领域通用）

**Rationale**: 这些是技术专有名词，翻译后反而降低可读性。中文技术文档惯例是保留英文术语。

### 4. 工具数量统一

**Decision**: 所有文档统一为 9 个工具（6 核心 + 3 增强）。

**Rationale**: 实际代码已有 9 个工具，文档应反映真实状态。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 路径验证可能影响 MCP 工具 | MCP 工具有自己的路径处理，不受影响 |
| WaitGroup 可能引入新 bug | 仔细测试 goroutine 生命周期 |
| 术语替换可能遗漏 | 用 grep 全局搜索确保一致性 |
