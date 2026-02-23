#!/usr/bin/env bash
set -euo pipefail

BUMP="patch"
DRY_RUN=false
VERSION=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bump)    BUMP="$2";    shift 2 ;;
    --version) VERSION="$2"; shift 2 ;;
    --dry-run) DRY_RUN=true; shift   ;;
    *) printf "Unknown flag: %s\n" "$1" >&2; exit 1 ;;
  esac
done

if [[ -n "$VERSION" ]]; then
  NEXT="$VERSION"
else
  LATEST=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
  MAJOR=$(echo "$LATEST" | sed 's/^v//' | cut -d. -f1)
  MINOR=$(echo "$LATEST" | sed 's/^v//' | cut -d. -f2)
  PATCH=$(echo "$LATEST" | sed 's/^v//' | cut -d. -f3)

  case "$BUMP" in
    major) NEXT="v$((MAJOR + 1)).0.0" ;;
    minor) NEXT="v${MAJOR}.$((MINOR + 1)).0" ;;
    patch) NEXT="v${MAJOR}.${MINOR}.$((PATCH + 1))" ;;
    *) printf "Invalid bump: %s (use major, minor, patch)\n" "$BUMP" >&2; exit 1 ;;
  esac
fi

printf "Releasing %s\n" "$NEXT"

go vet ./...
go test ./...

if [[ "$DRY_RUN" == true ]]; then
  printf "Dry run — skipping tag and release\n"
  goreleaser release --snapshot --clean
  exit 0
fi

git tag "$NEXT"
git push origin "$NEXT"
goreleaser release --clean
