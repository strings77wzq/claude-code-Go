## 1. 后端：模型切换功能

- [x] 1.1 `api.Client` 增加 `SetModel(model string)` 方法（带 mutex 保护）
- [x] 1.2 `agent.Agent` 增加 `SetModel(model string)` 方法
- [x] 1.3 REPL `/model` 命令增强：无参数显示当前模型，带参数切换
- [x] 1.4 新增 `/models` 命令（列出 Anthropic + Tencent Coding Plan 模型）
- [x] 1.5 更新欢迎信息显示当前模型

## 2. 前端官网：文档更新

- [x] 2.1 更新英文配置文档（说明运行时切换模型）
- [x] 2.2 更新中文配置文档（说明运行时切换模型）
- [x] 2.3 更新英文 Roadmap 功能对比表
- [x] 2.4 更新中文 Roadmap 功能对比表
