# Python Harness

## 什么是 Harness？

Harness 是 claude-code-Go 的**测试和评估基础设施**，用于在开发和 CI 环境中验证 Agent 的行为，无需消耗真实 API 费用。

## 为什么用 Python？

| 原因 | 说明 |
|------|------|
| **Mock Server** | FastAPI 写 SSE Mock 比 Go 快 3 倍，代码更简洁 |
| **测试框架** | pytest 的参数化测试、fixture 系统比 Go testing 灵活 |
| **数据分析** | 评估器需要 NLP 能力，Python 生态无可替代 |
| **开发速度** | 测试代码不需要编译，改完直接跑 |

**这不是生产代码** — Harness 只在 CI/开发时使用，不影响用户的 `go-code` 二进制。Claw Code 也采用同样的架构（Rust 主程序 + Python 测试）。

## 目录结构

```
harness/
├── mock_server/     ← Mock Anthropic API（模拟真实 API 响应）
├── evaluators/      ← 输出质量评估（文本完整性、工具正确性）
├── replay/          ← 会话回放调试（JSONL 会话重放）
├── test_*.py        ← Parity 测试（验证行为与预期一致）
└── requirements.txt ← Python 依赖
```

## 快速开始

### 安装依赖

```bash
cd harness
pip install -r requirements.txt
```

### 运行 Mock Server

```bash
python -m harness.mock_server
# Mock API 启动在 http://127.0.0.1:8765
```

### 运行测试

```bash
cd harness
python -m pytest -v
```

### 会话回放

```python
from harness.replay import SessionReplayer

replayer = SessionReplayer("~/.go-code/sessions/")
replayer.replay_latest()
```

## 在 CI 中的使用

GitHub Actions 会自动运行 Harness 测试：

```yaml
- name: Run Python tests
  working-directory: harness
  run: python -m pytest -v --tb=short
```
