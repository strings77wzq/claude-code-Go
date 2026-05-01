## Context

README 存在 4 处过时信息，未创建 Release，缺少 Issue 模板。

## Goals / Non-Goals

**Goals:**
- README 反映最新项目状态
- 创建 v0.1.0 Release
- 添加 Issue 模板方便社区贡献

**Non-Goals:**
- 不修改 Go 源代码
- 不修改 CI/CD 配置

## Decisions

### 1. README 更新策略

**Decision**: 全面更新 README，保持现有结构，只更新过时内容。

### 2. Release 策略

**Decision**: 创建 v0.1.0 tag，通过 GoReleaser 自动发布多平台二进制。

### 3. Issue 模板

**Decision**: 使用 GitHub 标准模板格式（.github/ISSUE_TEMPLATE/）。

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| README 更新遗漏 | 逐项对比实际项目结构 |
| Release 失败 | 手动创建 tag 和 Release |
