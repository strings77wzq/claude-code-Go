#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT_DIR"

go build -o bin/go-code ./cmd/go-code
python -m pytest harness/ -v --tb=short
