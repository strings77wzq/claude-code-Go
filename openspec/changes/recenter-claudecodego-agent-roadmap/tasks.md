## 1. Baseline Truth Audit

- [ ] 1.1 Run `go test ./...` and record the exact pass/fail result in the change notes or task evidence.
- [ ] 1.2 Run `./scripts/run-harness.sh` and record the exact pass/fail result, including failing scenario names if any.
- [ ] 1.3 Run the available docs build or docs lint command and record whether README/docs can be built locally.
- [ ] 1.4 Run `go-code doctor --offline` from a local build and capture gaps between doctor output and documented onboarding.
- [ ] 1.5 Audit README, `PARITY.md`, `docs/`, and `docs/zh/` for claims that are unsupported, stale, or ahead of implementation.
- [ ] 1.6 Produce a short current-state summary listing what is already good, what is broken, and what must be deferred.

## 2. Product Recenter

- [ ] 2.1 Update README positioning to state the Go runtime plus Python harness product identity clearly.
- [ ] 2.2 Update architecture docs to define Go-owned runtime responsibilities and Python-owned harness responsibilities.
- [ ] 2.3 Update roadmap docs to prioritize core prompt/tool/provider/permission/session reliability before extension ecosystem work.
- [ ] 2.4 Update `PARITY.md` into a workflow matrix with status, evidence, known gaps, and next task links.
- [ ] 2.5 Remove or relabel unsupported public claims as planned, partial, or experimental.

## 3. Provider Profile Model

- [ ] 3.1 Refactor provider selection to distinguish user-facing provider profiles from reusable transport implementations.
- [ ] 3.2 Add provider profile metadata for `anthropic`, `openai`, `deepseek`, and `mimo`.
- [ ] 3.3 Add config loading and validation for explicit provider profile names, model IDs, base URLs, and API key environment variable conventions.
- [ ] 3.4 Add unit tests for provider profile detection, explicit provider override, default base URLs, invalid model/provider combinations, and unknown custom models.
- [ ] 3.5 Ensure runtime logs and errors name the selected provider profile rather than only the underlying transport.

## 4. DeepSeek Support

- [ ] 4.1 Update model registry metadata to include `deepseek-v4-flash` and `deepseek-v4-pro` as preferred DeepSeek models.
- [ ] 4.2 Keep `deepseek-chat` and `deepseek-reasoner` only as legacy aliases or reject them with clear migration guidance.
- [ ] 4.3 Add DeepSeek config examples for API key, provider, base URL, model, thinking/reasoning options where supported, and streaming.
- [ ] 4.4 Add unit tests for DeepSeek request construction through the compatible transport.
- [ ] 4.5 Add a Python harness scenario that validates DeepSeek profile selection and mock streaming/tool-call behavior without real API keys.
- [ ] 4.6 Update English and Chinese docs to reflect current DeepSeek model names and deprecation guidance.

## 5. MiMo-V2.5 Support

- [ ] 5.1 Add `mimo-v2.5-pro` model metadata and a MiMo provider profile with documented compatibility status.
- [ ] 5.2 Verify official MiMo API endpoint/auth/streaming assumptions from public docs or user-provided examples before hardcoding endpoint behavior.
- [ ] 5.3 If MiMo uses an OpenAI-compatible API, wire the MiMo profile to the shared compatible transport with MiMo-specific defaults.
- [ ] 5.4 If MiMo requires a different API shape, add a narrow adapter and tests instead of overloading generic OpenAI code.
- [ ] 5.5 Add unit tests for MiMo profile resolution, model validation, config errors, and request construction.
- [ ] 5.6 Add a Python harness scenario for MiMo mock provider streaming and tool-call behavior.
- [ ] 5.7 Add English and Chinese docs showing MiMo-V2.5 setup, known limitations, and verification status.

## 6. Harness Engineering Gate

- [ ] 6.1 Define a harness scenario manifest that maps each supported workflow/provider claim to one or more tests.
- [ ] 6.2 Add deterministic harness coverage for provider streaming, tool calls, file edits, bash permission decisions, session persistence, replay, recoverable provider failures, and model/profile selection.
- [ ] 6.3 Ensure `./scripts/run-harness.sh` builds the binary, starts local mocks, runs scenarios, and exits non-zero on failures.
- [ ] 6.4 Add or update CI jobs so Go tests and Python harness tests run on pull requests.
- [ ] 6.5 Update harness docs with how to add scenarios and how to interpret failures.
- [ ] 6.6 Record verification evidence in `PARITY.md` or a linked docs page after harness scenarios pass.

## 7. Docs And Spec Alignment

- [ ] 7.1 Update docs contribution guidance to require verification evidence for new supported claims.
- [ ] 7.2 Sync English and Chinese provider docs for Anthropic, OpenAI-compatible, DeepSeek, and MiMo profiles.
- [ ] 7.3 Add known-gaps sections to docs pages where implementation is partial.
- [ ] 7.4 Update OpenSpec task discipline so completed tasks include test evidence or an explicit not-tested note.
- [ ] 7.5 Remove stale screenshots, placeholder demos, outdated metrics, or claims that imply unavailable behavior.

## 8. Final Verification

- [ ] 8.1 Run `go test ./...` after all implementation and docs changes.
- [ ] 8.2 Run `./scripts/run-harness.sh` after all implementation and docs changes.
- [ ] 8.3 Run docs build or docs lint after documentation changes.
- [ ] 8.4 Run `openspec status --change recenter-claudecodego-agent-roadmap` and confirm tasks are ready for archive only after implementation evidence is complete.
- [ ] 8.5 Summarize changed files, good existing design points preserved, improvements made, remaining risks, and not-tested gaps.
