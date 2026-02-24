#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

echo "Seeding demo user and notes"
GOCACHE="$ROOT_DIR/.cache/go-build" GOMODCACHE="$ROOT_DIR/.cache/go-mod" \
  go run ./cmd/seed \
    -addr localhost:3200 \
    -username demo-user \
    -email demo@example.com \
    -password 'DemoPass123!'
