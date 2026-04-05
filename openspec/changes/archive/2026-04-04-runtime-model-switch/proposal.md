## Why

当前 `/model` 命令只能显示当前模型，不能切换。用户需要重启进程才能换模型，体验差。特别是腾讯云 Coding Plan 支持 8 个模型，运行时切换是刚需。

## What Changes

### 后端
- `api.Client` 增加 `SetModel(model string)` 方法
- `agent.Agent` 增加 `SetModel(model string)` 方法
- REPL `/model` 命令增强：无参数显示当前模型，带参数切换模型
- 新增 `/models` 命令列出支持的模型（含 Coding Plan 模型列表）

### 前端官网
- 更新配置文档，说明运行时切换模型的方法
- 更新 Roadmap 功能对比表，标记 "运行时切换模型" 为已完成

## Capabilities

### New Capabilities
- `runtime-model-switch`: 运行时切换模型（`/model <name>`）
- `model-list`: 列出支持的模型（`/models`）

### Modified Capabilities
- `docs-website`: 更新配置文档和 Roadmap

## Impact

- 修改文件: `internal/api/client.go`, `internal/agent/loop.go`, `pkg/tty/repl.go`
- 修改文件: `docs/guide/configuration.md`, `docs/zh/guide/configuration.md`, `docs/roadmap.md`, `docs/zh/roadmap.md`
- 不影响: 核心 Agent Loop 逻辑、工具系统、权限系统
