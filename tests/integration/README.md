# 集成测试指南

本文档介绍如何为 claude-code-Go 编写集成测试。

## 概述

集成测试验证系统的多个组件协同工作。claude-code-Go 的主要集成点包括：

1. **Agent Loop** - 思考-行动-观察循环
2. **Tool Registry** - 工具注册和执行
3. **Permission System** - 权限检查
4. **API Client** - 与 LLM 提供商的通信
5. **Session Management** - 会话持久化

## 测试策略

### 1. 多轮对话测试

测试代理能够处理需要多轮交互的复杂任务：

```go
func TestMultiTurnConversation(t *testing.T) {
    // 创建测试代理
    ag := agent.New(agent.Config{
        MaxTurns: 10,
        Timeout:  30 * time.Second,
    })

    // 注册模拟工具
    ag.RegisterTool(tool.Definition{
        Name:        "MockRead",
        Description: "Read a file",
        Handler: func(args map[string]interface{}) (string, error) {
            return "file content", nil
        },
    })

    // 测试多轮对话流程
    messages := []string{
        "Read the main.go file",
        "Now edit it to add a comment",
        "Show me the diff",
    }

    for i, msg := range messages {
        response, err := ag.Process(msg)
        if err != nil {
            t.Fatalf("Turn %d failed: %v", i, err)
        }
        if response == "" {
            t.Fatal("Empty response")
        }
    }
}
```

### 2. 工具链测试

测试工具可以链式调用：

```go
func TestToolChain(t *testing.T) {
    ag := agent.New(agent.Config{
        MaxTurns: 5,
    })

    // 测试工具链
    result, err := ag.Process("Read config.json, then edit it to add a new key")
    if err != nil {
        t.Fatalf("Tool chain failed: %v", err)
    }

    // 验证代理使用了多个工具
    if len(ag.Session.ToolCalls) < 2 {
        t.Error("Expected at least 2 tool calls in chain")
    }
}
```

### 3. 错误恢复测试

测试代理能够从临时错误中恢复：

```go
func TestErrorRecovery(t *testing.T) {
    ag := agent.New(agent.Config{
        MaxRetries: 3,
    })

    // 注册一个起初失败然后成功的工具
    callCount := 0
    ag.RegisterTool(tool.Definition{
        Name: "FlakyTool",
        Handler: func(args map[string]interface{}) (string, error) {
            callCount++
            if callCount < 3 {
                return "", fmt.Errorf("temporary error")
            }
            return "success", nil
        },
    })

    result, err := ag.Process("Use FlakyTool")
    if err != nil {
        t.Fatalf("Error recovery failed: %v", err)
    }

    if result != "success" {
        t.Error("Expected success after retries")
    }

    if callCount != 3 {
        t.Errorf("Expected 3 calls, got %d", callCount)
    }
}
```

## 性能测试

### 代理处理性能

```go
func BenchmarkAgentProcessing(b *testing.B) {
    ag := setupTestAgent()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := ag.Process("Simple query")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 工具执行性能

```go
func BenchmarkToolExecution(b *testing.B) {
    ag := setupTestAgent()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := ag.ExecuteTool("Read", map[string]interface{}{
            "file_path": "test.txt",
        })
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 会话保存性能

```go
func BenchmarkSessionSave(b *testing.B) {
    ag := setupTestAgent()
    // 创建大会话
    for i := 0; i < 100; i++ {
        ag.Process(fmt.Sprintf("Message %d", i))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := ag.Session.Save()
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 测试辅助函数

```go
// setupTestAgent 创建一个配置好的测试代理
func setupTestAgent() *agent.Agent {
    return agent.New(agent.Config{
        MaxTurns:   10,
        Timeout:    30 * time.Second,
        MaxRetries: 3,
    })
}

// setupMockTools 注册一组模拟工具用于测试
func setupMockTools(ag *agent.Agent) {
    ag.RegisterTool(tool.Definition{
        Name:        "MockRead",
        Description: "Read file",
        Handler:     mockReadHandler,
    })
    
    ag.RegisterTool(tool.Definition{
        Name:        "MockWrite",
        Description: "Write file",
        Handler:     mockWriteHandler,
    })
}
```

## 运行集成测试

```bash
# 运行所有集成测试
go test ./tests/integration/... -v

# 运行特定测试
go test ./tests/integration/... -run TestMultiTurnConversation -v

# 运行性能测试
go test ./tests/integration/... -bench=. -benchmem

# 带覆盖率报告
go test ./tests/integration/... -cover -coverprofile=coverage.out
```

## 最佳实践

1. **使用模拟工具** - 不要依赖真实的 API 调用
2. **清理测试数据** - 每个测试结束后清理临时文件
3. **并行测试** - 使用 `t.Parallel()` 加速测试套件
4. **超时控制** - 始终设置合理的超时避免测试挂起
5. **确定性** - 确保测试可重复，不依赖外部状态