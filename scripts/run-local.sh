#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

CERT_PATH="private.pem"
if [[ ! -f "$CERT_PATH" ]]; then
  echo "private.pem not found, generating a local RSA key"
  GOCACHE="$ROOT_DIR/.cache/go-build" GOMODCACHE="$ROOT_DIR/.cache/go-mod" \
    go run ./cmd/certgen -out "$CERT_PATH"
fi

echo "Starting server with testdata/local/server-config.json"
GOCACHE="$ROOT_DIR/.cache/go-build" GOMODCACHE="$ROOT_DIR/.cache/go-mod" \
  go run ./cmd/server -cfg testdata/local/server-config.json
