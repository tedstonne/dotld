#!/usr/bin/env bash
set -euo pipefail

REPO="${DOTLD_REPO:-tedstonne/dotld}"
VERSION="latest"
BIN_DIR="${DOTLD_BIN_DIR:-$HOME/.local/bin}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      REPO="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --bin-dir)
      BIN_DIR="$2"
      shift 2
      ;;
    --)
      shift
      break
      ;;
    *)
      break
      ;;
  esac
done

OS_RAW="$(uname -s)"
ARCH_RAW="$(uname -m)"

case "$OS_RAW" in
  Linux)
    OS="linux"
    ;;
  Darwin)
    OS="darwin"
    ;;
  *)
    printf "Unsupported OS: %s\n" "$OS_RAW" >&2
    exit 1
    ;;
esac

case "$ARCH_RAW" in
  x86_64)
    ARCH="x64"
    ;;
  aarch64 | arm64)
    ARCH="arm64"
    ;;
  *)
    printf "Unsupported architecture: %s\n" "$ARCH_RAW" >&2
    exit 1
    ;;
esac

ASSET="dotld-${OS}-${ARCH}"
if [[ "$VERSION" == "latest" ]]; then
  URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

mkdir -p "$BIN_DIR"
TARGET="${BIN_DIR}/dotld"

printf "Downloading %s\n" "$URL"
curl -fsSL "$URL" -o "$TARGET"
chmod +x "$TARGET"

printf "Installed dotld to %s\n" "$TARGET"

case ":$PATH:" in
  *":$BIN_DIR:"*)
    ;;
  *)
    printf "Add %s to PATH to run dotld globally\n" "$BIN_DIR"
    ;;
esac

if [[ $# -gt 0 ]]; then
  "$TARGET" "$@"
fi
