## Why

After fixing the Python harness import issues, two remaining problems prevent a clean CI run:

1. **CI configuration not using editable install** - The `python-harness` job in `.github/workflows/ci.yml` doesn't use `pip install -e .`, so the harness import fix isn't tested in CI
2. **Two failing harness tests** - `test_load_from_lines` and `test_timeline` have pre-existing assertion bugs that need fixing

These must be resolved before pushing to GitHub, otherwise CI will fail.

## What Changes

- Update `.github/workflows/ci.yml` to use `pip install -e .` and run tests from project root
- Fix `test_load_from_lines` assertion in `harness/test_evaluators.py` (tool_calls count mismatch)
- Fix `test_timeline` assertion in `harness/test_evaluators.py` (missing 'name' key in timeline entry)

## Capabilities

### New Capabilities
- `ci-harness-config`: CI properly runs Python harness tests with editable install
- `harness-test-fixes`: All harness tests pass without assertion errors

### Modified Capabilities
- None

## Impact

- **Affected**: `.github/workflows/ci.yml`, `harness/test_evaluators.py`
- **Dependencies**: None
- **Systems**: GitHub Actions CI pipeline
