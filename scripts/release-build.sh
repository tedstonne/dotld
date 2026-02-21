#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${ROOT_DIR}/dist"

mkdir -p "${OUT_DIR}"

bun build --compile "${ROOT_DIR}/src/cli/dotld.ts" --target=bun-linux-x64 --outfile "${OUT_DIR}/dotld-linux-x64"
bun build --compile "${ROOT_DIR}/src/cli/dotld.ts" --target=bun-linux-arm64 --outfile "${OUT_DIR}/dotld-linux-arm64"
bun build --compile "${ROOT_DIR}/src/cli/dotld.ts" --target=bun-darwin-x64 --outfile "${OUT_DIR}/dotld-darwin-x64"
bun build --compile "${ROOT_DIR}/src/cli/dotld.ts" --target=bun-darwin-arm64 --outfile "${OUT_DIR}/dotld-darwin-arm64"

if command -v sha256sum >/dev/null 2>&1; then
  sha256sum "${OUT_DIR}"/dotld-* >"${OUT_DIR}/checksums.txt"
else
  shasum -a 256 "${OUT_DIR}"/dotld-* >"${OUT_DIR}/checksums.txt"
fi

printf "Built release binaries in %s\n" "${OUT_DIR}"
