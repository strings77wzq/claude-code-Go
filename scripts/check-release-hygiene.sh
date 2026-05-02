#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT_DIR"

echo "OpenSpec hygiene: checking docs and active changes"
test -f docs/release-hygiene.md
test -f docs/docs-inventory.md
test -f openspec/specs/runtime-safety-gates/spec.md
test -f openspec/specs/extension-runtime-diagnostics/spec.md
test -f openspec/specs/agent-quality-gates/spec.md
test -f openspec/specs/openspec-change-hygiene/spec.md
test -d openspec/changes/archive/2026-05-02-fix-core-runtime-safety
test -d openspec/changes/archive/2026-05-02-productize-extension-boundaries
test -d openspec/changes/archive/2026-05-02-harness-agent-quality-gates
test -d openspec/changes/archive/2026-05-02-openspec-docs-release-hygiene

if grep -R "TBD Purpose" openspec/specs openspec/changes --include='*.md' >/tmp/go-code-hygiene-tbd.txt 2>/dev/null; then
  echo "OpenSpec hygiene: found placeholder purpose text"
  cat /tmp/go-code-hygiene-tbd.txt
  exit 1
fi

echo "docs drift: source/generated inventory documented"
grep -q "docs/.vitepress/dist" docs/docs-inventory.md
grep -q "Source of truth" docs/docs-inventory.md
grep -q "Generated output" docs/docs-inventory.md

echo "install smoke: building local binary"
go build -o bin/go-code ./cmd/go-code
./bin/go-code --help >/tmp/go-code-help.txt
./bin/go-code doctor --offline >/tmp/go-code-doctor.txt || true
grep -q "go-code doctor" /tmp/go-code-doctor.txt
./bin/go-code version >/tmp/go-code-version.txt
grep -q "go-code" /tmp/go-code-version.txt

echo "release hygiene checks passed"
