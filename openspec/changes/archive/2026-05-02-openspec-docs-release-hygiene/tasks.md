## 1. OpenSpec Hygiene Gate

- [x] 1.1 Define a checklist for active change size, concrete purpose, testable scenarios, task evidence, and strict validation.
- [x] 1.2 Identify active or release-relevant specs with placeholder or vague purpose text.
- [x] 1.3 Replace placeholder purpose text in priority specs with durable project context.
- [x] 1.4 Add a lightweight script or documented command sequence for hygiene checks.
- [x] 1.5 Produce an active-change inventory that identifies umbrella changes, implementation changes, and their evidence expectations.

## 2. Task Evidence Discipline

- [x] 2.1 Define the evidence format for completed tasks, including tested and not-tested cases.
- [x] 2.2 Update active implementation changes to use evidence-backed task completion notes.
- [x] 2.3 Add archive readiness steps requiring strict validation and evidence review.
- [x] 2.4 Document the exact evidence note format used when checking off tasks.

## 3. Documentation Truth Policy

- [x] 3.1 Inventory source-of-truth docs and generated documentation outputs.
- [x] 3.2 Document the generation command or source mapping for generated docs.
- [x] 3.3 Add or document a docs drift check for generated outputs.
- [x] 3.4 Review README, docs, and parity claims for planned-vs-shipped wording.
- [x] 3.5 Produce a source/generated documentation inventory file.

## 4. Release Readiness

- [x] 4.1 Define release evidence required for local/dev, release-candidate, and published states.
- [x] 4.2 Add install smoke checks for built artifacts or local install paths.
- [x] 4.3 Add release checklist items for OpenSpec strict validation, docs truth, harness status, and known gaps.
- [x] 4.4 Ensure release notes include known safety, compatibility, and comparison gaps.
- [x] 4.5 Add a release state matrix covering local/dev, release-candidate, and published evidence requirements.

## 5. Verification

- [x] 5.1 Run OpenSpec strict validation for all active release-bound changes.
- [x] 5.2 Run docs drift or documented equivalent check.
- [x] 5.3 Run install smoke check for the current local build path.
- [x] 5.4 Run `openspec validate openspec-docs-release-hygiene --strict`.

## Evidence

- `python -m pytest harness/test_release_hygiene.py -q`
- `./scripts/check-release-hygiene.sh`
- `python -m pytest harness/ -q`
- `go test ./...`
- `cd docs && npm run build`
- `openspec validate openspec-docs-release-hygiene --strict --json --no-interactive`

## Residual Risks

- `docs/.vitepress/dist` is generated successfully, but generated output should be reviewed and committed only for docs publish/release tasks.
- Publishing docs and release artifacts depends on remote CI/deploy credentials and may fail outside local verification.
