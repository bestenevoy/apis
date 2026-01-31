#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

printf '==> Build frontend\n'
cd "$ROOT/frontend"
if [[ -f package-lock.json ]]; then
  npm ci
else
  npm install
fi
npm run build

printf '==> Build Go binary\n'
cd "$ROOT"
mkdir -p dist

go build -trimpath -ldflags "-s -w" -o dist/wrzapi ./cmd/server

echo "Done: $ROOT/dist/wrzapi"
