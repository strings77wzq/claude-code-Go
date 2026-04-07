## Context

Multiple issues identified during code review that need immediate fixing:

1. **Git repository contains build artifacts** - `claude_code_harness.egg-info/` was committed
2. **No security scanning** - Missing CodeQL and dependency vulnerability checks
3. **No visual demo** - README is text-only, no GIF showing usage
4. **Incomplete CI** - Python harness tests don't use editable install
5. **Missing documentation** - No pricing, benchmark, or showcase pages
6. **No Docker support** - Enterprise users need containerization
7. **No telemetry** - Cannot track adoption or errors

## Goals / Non-Goals

**Goals:**
- Clean repository with proper .gitignore
- Security scanning enabled
- README includes Demo GIF section
- CI runs tests correctly with editable install
- Website has pricing, benchmark, showcase pages
- Docker image available
- Telemetry framework in place (opt-in)

**Non-Goals:**
- No new core features (just polish)
- No breaking changes
- No paid features (just documentation)

## Decisions

### Decision 1: Comprehensive .gitignore
**Choice**: Update .gitignore with standard Python, Go, and IDE patterns.
**Rationale**: Prevent build artifacts from being committed.
**Patterns to add**:
- `*.egg-info/`, `__pycache__/`, `.pytest_cache/`
- `.vscode/`, `.idea/`, `*.swp`
- `bin/`, `dist/`, `*.log`

### Decision 2: GitHub Security Scanning
**Choice**: Add CodeQL workflow for Go and Python.
**Rationale**: Identify security vulnerabilities early.
**Workflow**: Run on PR and push to main.

### Decision 3: Demo GIF Placeholder
**Choice**: Add Demo GIF section to README with placeholder.
**Rationale**: Visual demonstration increases adoption.
**Approach**: Create placeholder, record actual GIF later.

### Decision 4: CI Fix
**Choice**: Update python-harness job to use `pip install -e .`.
**Rationale**: Properly test import resolution.
**Change**: Remove working-directory, add editable install step.

### Decision 5: Documentation Pages
**Choice**: Add pricing.md, benchmark.md, showcase.md to docs.
**Rationale**: Marketing and comparison need dedicated pages.
**Structure**: Follow existing docs pattern.

### Decision 6: Docker Support
**Choice**: Create multi-stage Dockerfile.
**Rationale**: Single binary deployment in container.
**Base**: Alpine Linux for minimal size.

### Decision 7: Telemetry Framework
**Choice**: Add optional anonymous usage tracking.
**Rationale**: Understand adoption without privacy invasion.
**Principle**: Opt-in, no PII, transparent.

## Risks / Trade-offs

- **[Risk]** Adding telemetry may concern privacy-conscious users → **Mitigation**: Strict opt-in, open source, no PII
- **[Risk]** Docker adds complexity → **Mitigation**: Optional, binary still works standalone
- **[Risk]** Security scanning may slow CI → **Mitigation**: Run in parallel, cache results
