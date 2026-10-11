#!/usr/bin/env bash
# Builds the browser demo into web/site. Needs Go 1.21+.
set -euo pipefail
cd "$(dirname "$0")/.."

GOOS=js GOARCH=wasm go build -modfile=web/go.mod -ldflags="-s -w" -o web/site/seman.wasm ./cmd/seman-web

# wasm_exec.js must come from the same Go version that built the wasm file.
goroot="$(go env GOROOT)"
for candidate in "$goroot/lib/wasm/wasm_exec.js" "$goroot/misc/wasm/wasm_exec.js"; do
  if [ -f "$candidate" ]; then
    cp "$candidate" web/site/wasm_exec.js
    break
  fi
done

echo "Built web/site. Preview with: python3 -m http.server -d web/site 8000"
