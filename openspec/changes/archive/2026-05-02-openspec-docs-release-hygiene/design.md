## Context

This change keeps project claims aligned with implementation evidence. The current roadmap has a broad umbrella plan and many existing specs; without hygiene gates, future implementation changes can drift, archive incomplete tasks, or publish generated docs that hide meaningful source changes. This slice should run after the urgent runtime safety work but before a public release claim.

## Goals / Non-Goals

**Goals:**
- Make OpenSpec changes small enough to implement and archive with evidence.
- Require clear spec purpose, scenario coverage, task checkboxes, and validation before apply/archive.
- Define which docs are source-of-truth and which are generated artifacts.
- Require release evidence for installability, docs accuracy, harness status, and known gaps.

**Non-Goals:**
- Rewrite all documentation content in this change.
- Implement runtime behavior changes.
- Replace OpenSpec tooling.
- Publish a release automatically.

## Decisions

1. Treat large roadmap changes as umbrellas.
   - Umbrella changes can guide sequencing, but implementation should happen in smaller child changes with concrete specs and tasks.
   - Alternative rejected: applying a 100+ task roadmap directly, because it is hard to verify and archive safely.

2. Require evidence on task completion.
   - Completed tasks must be backed by tests, commands, docs checks, or explicit not-tested notes.
   - Alternative rejected: checkbox-only completion, because future agents cannot distinguish verified work from intent.

3. Separate docs source from generated output.
   - Generated docs must have a declared source command and drift check.
   - Alternative rejected: reviewing generated churn without source provenance, because it obscures real changes.

4. Gate release state transitions.
   - Moving from local/dev to published release requires install verification, OpenSpec validation, docs truth checks, and known-risk notes.
   - Alternative rejected: version bump alone, because it does not prove usability.

## Risks / Trade-offs

- [Risk] Hygiene checks add process overhead. → Mitigation: keep checks scripted and focused on release-critical evidence.
- [Risk] Existing specs may need incremental cleanup. → Mitigation: prioritize active and release-relevant specs first.
- [Risk] Generated docs policy may require CI changes. → Mitigation: start with a documented command and a non-invasive drift check.

## Migration Plan

1. Define OpenSpec hygiene requirements and checklist.
2. Clean high-priority vague spec purposes and active change task structures.
3. Add docs source/generated policy and drift checks.
4. Add release readiness checklist and install smoke evidence requirements.
5. Validate active changes before archive or release.
