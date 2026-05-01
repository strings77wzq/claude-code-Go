## Context

官网 404、Logo 错误、导航混乱、Release 冲突。

## Goals / Non-Goals

**Goals:**
- 修复所有 404 链接
- 替换 Logo 为 Go 风格
- 精简导航栏
- 修复 Release 页面

**Non-Goals:**
- 不重写整个官网
- 不改 Go 源代码

## Decisions

### 1. Logo 设计

**Decision**: 使用 Go Gopher 风格的 `<Go>` 组合图标。

```svg
<svg viewBox="0 0 64 64">
  <!-- Go 蓝渐变背景 -->
  <!-- 代码括号 < > -->
  <!-- Go 字母 -->
</svg>
```

### 2. 导航栏精简

**Decision**: 5 个顶级菜单，中英文一致。

```
Guide | Architecture | API | Resources | GitHub
```

### 3. 中英文对齐

**Decision**: 每个英文侧边栏项对应一个中文项，不存在的中文页面暂时移除链接。

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| 中文页面缺失 | 暂时移除链接，后续补充 |
| Logo 设计不满意 | 后续可替换 |
