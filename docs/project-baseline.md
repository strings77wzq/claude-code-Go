# Project Baseline

Last checked: 2026-04-28

This page records the baseline for the `world-class-go-agent-rewrite` OpenSpec change. It separates real product gaps from sandbox or environment constraints so future fixes can be measured against the same evidence.

## Summary

| Check | Result | Notes |
| --- | --- | --- |
| `go test -count=1 ./...` | Pass | All Go packages pass outside cache. |
| `make test` | Pass with unrestricted local test permissions | The normal sandbox blocks `httptest` sockets and user-home log/session writes. With local sockets and writable user cache paths allowed, the wrapper passes. |
| `npm run build` in `docs/` | Pass | VitePress build succeeds. |
| `python -m pytest harness -v` before building `bin/go-code` | Pass with skips | CLI scenario tests skip when `bin/go-code` is missing. |
| `go build -o bin/go-code ./cmd/go-code` | Pass | The sandbox may warn about read-only Go module stat cache; the binary is produced. |
| `python -m pytest harness -v` after building `bin/go-code` | Fail | 29 passed, 5 failed. Failures are CLI scenario assertions with empty stdout. |
| `openspec validate world-class-go-agent-rewrite --strict` | Pass | Change artifacts validate. |

## Harness Failures

After building `bin/go-code`, the Python harness can run the CLI scenarios and exposes these failures:

| Test | Failure |
| --- | --- |
| `harness/test_permission_flow.py::TestPermissionFlow::test_permission_allowed` | Expected CLI output, got empty stdout. |
| `harness/test_permission_flow.py::TestPermissionFlow::test_permission_denied` | Expected CLI output, got empty stdout. |
| `harness/test_scenarios.py::TestStreamingText::test_streaming_text` | Expected streaming text, got empty stdout. |
| `harness/test_scenarios.py::TestToolRoundtrip::test_tool_roundtrip_read` | Expected file-read output, got empty stdout. |
| `harness/test_scenarios.py::TestToolRoundtrip::test_tool_roundtrip_bash` | Expected bash output, got empty stdout. |

These are product-level gaps, not just test noise: the deterministic CLI harness is not proving prompt mode, streaming, tool roundtrip, or permission behavior yet.

## Environment Notes

- Localhost socket tests require permission to bind loopback ports.
- Some Go tests write under the user's home directory for logs or sessions.
- The docs build writes generated files under `docs/.vitepress/dist`.
- Harness CLI scenarios depend on a built binary at `bin/go-code`.

