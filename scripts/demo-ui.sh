#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOCACHE_DIR="$ROOT_DIR/.cache/go-build"
GOMODCACHE_DIR="$ROOT_DIR/.cache/go-mod"
CLIENT_CFG="testdata/local/client-config.json"

cd "$ROOT_DIR"
mkdir -p "$ROOT_DIR/.cache"

echo "Launching TUI client"
echo "Make sure server is already running: ./scripts/run-local.sh"
echo "Demo credentials: demo@example.com / DemoPass123!"

GOCACHE="$GOCACHE_DIR" GOMODCACHE="$GOMODCACHE_DIR" \
  go run ./cmd/client -cfg "$CLIENT_CFG"
