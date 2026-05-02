# Documentation Inventory

## Source of truth

These files are reviewed as authored documentation:

- `README.md`
- `PARITY.md`
- `docs/**/*.md`
- `docs/.vitepress/config.ts`
- `harness/README.md`
- `openspec/specs/**/*.md`
- `openspec/changes/**/{proposal.md,design.md,tasks.md,spec.md}`

## Generated output

These paths are generated output and should be reviewed only for release or docs-publish work:

- `docs/.vitepress/dist`
- `docs/.vitepress/.temp`
- `harness/**/__pycache__`
- `.pytest_cache`
- `bin`

## Generation and drift checks

- Build docs from source with `cd docs && npm run build`.
- Run release hygiene checks with `./scripts/check-release-hygiene.sh`.
- Keep generated `docs/.vitepress/dist` changes separate from ordinary source documentation changes unless a release task explicitly requires publishing artifacts.
