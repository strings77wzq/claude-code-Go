---
title: Roadmap
description: claude-code-Go 项目 Roadmap

---

# Roadmap

以下是 claude-code-Go 的 Roadmap，分为三个阶段。

## 概览

| 阶段 | 状态 | 描述 |
|-------|--------|-------------|
| 第一阶段 | ✅ 已完成 | 核心基础 — Agent loop、工具、权限、SSE、会话 |
| 第二阶段 | ✅ 已完成 | 增强功能 — Skills、多 Provider、会话恢复 |
| 第三阶段 | 🟡 进行中 | v0.3 扩展产品化 — MCP、LSP、Hooks、Skills、Replay evidence |

---

## 第一阶段：基础架构 ✅ 已完成

第一阶段建立了 AI 编程助手的核心基础设施。

### 已完成功能

- **Agent Loop** — 由 stop_reason 驱动的「思考→行动→观察」自主循环
- **9 大内置工具** — 6 个核心（Read、Write、Edit、Glob、Grep、Bash）+ 3 个增强（Diff view、Tree、WebFetch）
- **权限系统** — 三级权限模型，支持 glob 规则匹配与会话记忆
- **MCP 集成** — 当前为 v0.3 partial：配置、命名空间、权限门控和 harness 证据正在产品化
- **SSE 流式** — 逐 token 实时流式输出
- **会话持久化** — 保存和恢复对话状态
- **Hooks 系统** — 执行前后回调，扩展能力

---

## 第二阶段：增强功能 ✅ 已完成

第二阶段添加了更强大的功能，提升可用性和灵活性。

### 已完成功能

- **Skills 系统** — 自定义命令和可重用工作流（如 `/review-pr`、`/deploy`）
- **多 Provider 支持** — Anthropic、OpenAI 及任何 OpenAI 兼容 API
- **会话恢复** — 加载之前的对话并无缝衔接
- **增强工具** — Diff 视图、树形可视化、网页抓取
- **手动压缩** — `/compact` 命令减少上下文大小
- **自动更新** — `/update` 命令检查版本并更新

---

## 第三阶段：扩展产品化 🟡 进行中

当前重点是把已有扩展代码变成可诊断、可验证、可文档化的产品能力。

### 当前范围

- **MCP 集成 (Partial v0.3)** — 配置、工具命名空间、权限门控、文档和 harness 场景
- **LSP 集成 (Partial v0.3)** — 健康检查、能力门控、doctor 诊断、文档和 unavailable harness 场景
- **Replay evidence** — extension events、permission decisions、secret redaction 和 evidence mode
- **Docs/PARITY 对齐** — 中英文文档只声明已有证据支持的能力

### 愿景

第三阶段的目标是先稳定 CLI/TUI runtime 和扩展证据，再考虑 IDE、团队协作或云端能力。

---

## 参与贡献

想影响 Roadmap？查看 [反馈](/zh/feedback) 页面了解如何提交问题报告和功能请求。

---

## 功能对比

| 功能 | Phase 1 | Phase 2 | Claude Code |
|------|---------|---------|-------------|
| Agent Loop | ✅ | ✅ | ✅ |
| 内置工具 | 6 个核心 | 9 个 | 20+ |
| 权限系统 | ✅ | ✅ | ✅ |
| MCP 集成 (Partial v0.3) | 🟡 | 🟡 | ✅ |
| SSE 流式 | ✅ | ✅ | ✅ |
| 会话持久化 | ✅ | ✅ | ✅ |
| 会话恢复 | ❌ | ✅ | ✅ |
| Skills 系统 | ❌ | ✅ | ✅ |
| 多 Provider 支持 | ❌ | ✅ | ❌ |
| 自动更新 | ❌ | ✅ | ✅ |
| Hooks 系统 | ❌ | ✅ | ✅ |
| 运行时切换模型 | ❌ | ✅ | ✅ |
| IDE 集成 | ❌ | ❌ | ✅ |

---

*最后更新：2026 年 5 月*
