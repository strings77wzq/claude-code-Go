## Description

<!-- Brief description of what this PR does. Link to OpenSpec proposal if applicable. -->

## Design Traceability

<!-- REQUIRED for feature PRs. Link your design artifacts. -->

- **Proposal**: <!-- openspec/changes/<name>/proposal.md -->
- **Design**: <!-- openspec/changes/<name>/design.md -->
- **Tasks**: <!-- openspec/changes/<name>/tasks.md -->
- **Requirements**: <!-- openspec/changes/<name>/requirement.md -->

## Related Issues

<!-- Link issues: "Fixes #123", "Closes #456" -->

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Refactor
- [ ] Documentation
- [ ] Test
- [ ] Chore

## Quality Gates

<!-- ALL must be checked before requesting review. -->

- [ ] **Format**: `gofmt -l .` — zero diffs
- [ ] **Vet**: `go vet ./...` — clean
- [ ] **Lint**: `golangci-lint run ./...` — zero new issues
- [ ] **Build**: `go build ./...` — passes
- [ ] **Test**: `go test -cover ./...` — ≥ 80%
- [ ] **Race**: `go test -race ./...` — clean
- [ ] **Security**: `gosec ./...` — zero high/critical
- [ ] **Eval**: `pytest harness/test_workflow_quality.py` — **Scorecard 100% green**

## Scope Verification

<!-- Confirm the diff matches what tasks.md declared. -->

```bash
$ git diff --name-only main...HEAD
```

- [ ] All changed files are within declared task boundaries
- [ ] No incidental changes ("顺便改") outside scope

## Testing Evidence

<!-- Paste test output or link to CI run. -->

```
$ go test -v -cover ./...
<output>

$ go test -race ./...
<output>

$ pytest harness/test_workflow_quality.py -v
<output>
```

## Additional Notes

<!-- Screenshots, migration steps, breaking change notices, performance data. -->
