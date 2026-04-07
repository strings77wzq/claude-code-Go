## Context

The `harness/` directory contains Python test infrastructure for the claude-code-Go project:
- Mock Anthropic API server (`mock_server/`)
- Output quality evaluators (`evaluators/`)
- Session replay utilities (`replay/`)
- Integration tests (`test_*.py`)

Currently, multiple issues prevent proper functioning:

1. **Import Error**: Running `pytest` from project root fails:
   ```
   ImportError while loading conftest '.../harness/conftest.py'.
   ModuleNotFoundError: No module named 'harness'
   ```

2. **IDE Errors**: VSCode/GoLand show "Unresolved reference 'Counter'" and similar errors because the Python path isn't properly configured.

3. **Missing Package Structure**: `harness/` subdirectories have `__init__.py`, but the root lacks it.

## Goals / Non-Goals

**Goals:**
- Enable `pytest` to run from project root without import errors
- Make `harness` importable as a Python package
- Configure IDE to resolve imports correctly
- Document the setup process for developers
- Ensure CI/CD can run harness tests

**Non-Goals:**
- Restructuring the harness directory layout
- Modifying test logic or assertions
- Adding new test cases
- Supporting Python versions below 3.10
- Changing the actual test implementations

## Decisions

### Decision 1: Add root-level `__init__.py`
**Choice**: Add `harness/__init__.py` to make `harness` a proper package.
**Rationale**: Python requires `__init__.py` to recognize a directory as a package. This is the standard approach.

### Decision 2: Create `pyproject.toml` with setuptools
**Choice**: Add `pyproject.toml` at project root with setuptools configuration.
**Rationale**: Modern Python packaging standard. Enables `pip install -e .` for development.
**Content**:
- Package name: `claude-code-harness`
- Source directory: root (harness/ as top-level package)
- Dependencies from requirements.txt

### Decision 3: Fix `conftest.py` with sys.path fallback
**Choice**: Modify `harness/conftest.py` to handle import resolution gracefully.
**Rationale**: When running pytest from project root without editable install, Python can't find `harness`.
**Approach**: Add project root to `sys.path` as a fallback before imports.

### Decision 4: Add VSCode settings
**Choice**: Create `.vscode/settings.json` with Python analysis configuration.
**Rationale**: Helps VSCode resolve imports correctly, especially for editable installs.
**Settings**:
- `python.analysis.extraPaths`: include project root
- `python.analysis.include`: include harness directory

### Decision 5: Keep requirements.txt as primary dependency source
**Choice**: Maintain `harness/requirements.txt` as the canonical dependency list.
**Rationale**: Simple, works without pip install, CI/CD already uses it.
**Note**: `pyproject.toml` will reference these dependencies.

## Risks / Trade-offs

- **[Risk]** Adding `pyproject.toml` may conflict with future packaging if we add a proper Python package later → **Mitigation**: Keep it minimal, only for development/editable installs
- **[Risk]** `sys.path` manipulation in `conftest.py` could have side effects → **Mitigation**: Only modify if import fails, use try/except
- **[Trade-off]** Editable install requires pip; developers using direct pytest may still have issues → **Mitigation**: Document both approaches in README
