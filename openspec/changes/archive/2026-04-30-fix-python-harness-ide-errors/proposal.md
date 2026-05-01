## Why

The Python harness test infrastructure (`harness/` directory) has multiple issues that prevent it from running correctly and cause IDE errors (VSCode/GoLand show "Unresolved reference" warnings). The specific problems are:

1. **Missing `__init__.py`** - The `harness/` root directory lacks `__init__.py`, preventing Python from recognizing it as a package
2. **Missing `pyproject.toml`** - No Python package configuration for editable installs
3. **Import path issues** - `conftest.py` imports fail when running `pytest` from project root
4. **IDE configuration missing** - No `.vscode/settings.json` to help IDEs resolve imports

These issues block the entire harness testing framework, which represents ~13% of the project's test coverage.

## What Changes

- Add `harness/__init__.py` to make harness a proper Python package
- Create `pyproject.toml` with proper package configuration for editable installs
- Fix `harness/conftest.py` import resolution
- Add `.vscode/settings.json` for VSCode Python path configuration
- Create `harness/setup.py` as alternative installation method
- Add `harness/README.md` with setup instructions
- Verify all harness tests can run with `pytest`

## Capabilities

### New Capabilities
- `python-package-structure`: Proper Python package configuration for the harness module
- `ide-import-resolution`: IDE can resolve harness imports without errors
- `harness-test-runner`: Working pytest execution for harness tests

### Modified Capabilities
- None (this is a bug fix, not a requirement change)

## Impact

- **Affected**: `harness/` directory, `pyproject.toml` (new), `.vscode/settings.json` (new)
- **Dependencies**: Python 3.10+, pytest, fastapi, uvicorn, httpx
- **Systems**: CI/CD pipeline can now run harness tests; developers can use IDE features
