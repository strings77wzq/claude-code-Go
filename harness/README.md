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
# From project root; builds bin/go-code first
./scripts/run-harness.sh

# Or run pytest directly after the binary exists
python -m pytest harness/ -v
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
  run: pytest harness/ -v --tb=short
```

失败时 CI 会上传 `harness/logs/pytest.log`，用于定位具体 scenario。

## 添加新 Scenario

1. 优先在 `harness/manifests/*.json` 添加一个 `agent-quality.v1` manifest，写清 `prompt`、`allowed_tools`、`assertions`、`trace.required_events` 和 `budgets`。
2. 如需 mock provider 行为，在 `harness/mock_server/scenarios.py` 添加一个 `Scenario`。
3. 如果需要自动路由，在 `harness/mock_server/app.py` 的 `select_scenario_name` 中加入触发条件。
4. 在 `harness/test_scenarios.py` 或 `harness/test_quality_gates.py` 增加断言，优先使用 `go-code -p ... -q -f json` 的非交互入口。
5. 本地运行 `./scripts/run-harness.sh`，确认不依赖真实 API key。

## Manifest Quality Gates

`harness/quality` 提供 manifest-driven release evidence：

- `manifest.py` 负责解析和校验 `agent-quality.v1` 场景。
- `runner.py` 将 stdout/stderr、return code、latency、trace events 汇总成脱敏 evidence。
- `comparison.py` 生成规范化对比报告，明确区分 `source=measured` 和 `source=manual`。

当前内置 manifest 覆盖：

- `repository-inspection`
- `safe-edit-and-test`
- `permission-denial-recovery`
- `provider-tool-failure-recovery`
- `user-facing-explanation`

---

## Setup for Development

### Prerequisites

- Python 3.10+
- pip

### Installation (Editable Mode)

For development, install the harness package in editable mode so imports work correctly:

```bash
# From project root
pip install -e .

# Or install with dev dependencies
pip install -e ".[dev]"
```

This allows you to:
- Run `pytest harness/` from the project root
- Import `harness` from any Python script
- Get IDE support for imports (VSCode/GoLand/PyCharm)

### Alternative: Using PYTHONPATH

If you don't want to install, set PYTHONPATH:

```bash
export PYTHONPATH="${PWD}:$PYTHONPATH"
pytest harness/
```

### Running Tests

After installation, run tests from project root:

```bash
# Run all harness tests
pytest harness/ -v

# Run with coverage
pytest harness/ -v --cov=harness

# Run specific test file
pytest harness/test_evaluators.py -v
```

### VSCode Configuration

The project includes `.vscode/settings.json` for optimal Python support:

- Import resolution is configured
- pytest is enabled as the test runner
- Test discovery is set to the `harness/` directory

### Troubleshooting

#### Import Error: No module named 'harness'

**Solution**: Install in editable mode:
```bash
pip install -e .
```

#### IDE shows "Unresolved reference"

**Solution**: 
1. Install the package: `pip install -e .`
2. Restart your IDE
3. For VSCode: Python extension should auto-detect the settings

#### pytest cannot find tests

**Solution**: Run from project root with the harness path:
```bash
pytest harness/ -v
```

Or from harness directory:
```bash
cd harness
python -m pytest -v
```
