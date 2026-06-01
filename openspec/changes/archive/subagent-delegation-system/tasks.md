## 1. Sub-Agent Package

- [ ] 1.1 Create `internal/agent/subagent/` package with sub-agent types and tool sets
- [ ] 1.2 Implement `Run()` function with isolated history, filtered tool registry, and max turns
- [ ] 1.3 Add sub-agent type constants: TypeExplore, TypeReview, TypeTestGen
- [ ] 1.4 Implement ExploreToolSet, ReviewToolSet, TestGenToolSet
- [ ] 1.5 Add sub-agent unit tests (tool sets, execution, error handling)

## 2. Task Tool

- [ ] 2.1 Create `internal/tool/builtin/task.go` — Task tool for sub-agent invocation
- [ ] 2.2 Implement InputSchema with subagent_type and prompt fields
- [ ] 2.3 Wire sub-agent execution through Task tool
- [ ] 2.4 Add Task tool unit tests (happy path, invalid type, missing fields)

## 3. Registration

- [ ] 3.1 Register Task tool in `internal/tool/init/register.go`
- [ ] 3.2 Pass API client and policy to TaskTool constructor

## 4. Verification

- [ ] 4.1 Run `go test -count=1 ./internal/agent/subagent` and confirm all tests pass
- [ ] 4.2 Run `go test -count=1 ./internal/tool/builtin` and confirm all tests pass
- [ ] 4.3 Run `go test ./...` and confirm all 27 packages pass
- [ ] 4.4 Run `go vet ./...` and `gofmt -l .` confirm clean
