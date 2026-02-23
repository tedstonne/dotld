#!/usr/bin/env bash
set -euo pipefail

PASS=0
FAIL=0

pass() { printf "  ✓ %s\n" "$1"; PASS=$((PASS + 1)); }
fail() { printf "  ✗ %s\n" "$1"; FAIL=$((FAIL + 1)); }

printf "\n=== Stage 1: Binary Install ===\n\n"

curl -fsSL https://raw.githubusercontent.com/tedstonne/dotld/main/scripts/install.sh | bash

if command -v dotld &>/dev/null; then
  pass "dotld is in PATH"
else
  fail "dotld not found in PATH"
  printf "Cannot continue without binary\n"
  exit 1
fi

VERSION="$(dotld --version 2>&1 || true)"
printf "  version: %s\n" "$VERSION"

printf "\n=== Stage 2: Live API Test ===\n\n"

if [[ -z "${DYNADOT_API_PRODUCTION_KEY:-}" ]]; then
  fail "DYNADOT_API_PRODUCTION_KEY not set, skipping API tests"
else
  OUTPUT="$(dotld example.com 2>&1 || true)"
  if echo "$OUTPUT" | grep -qi "taken"; then
    pass "dotld example.com returned Taken"
  else
    fail "dotld example.com unexpected output: $OUTPUT"
  fi

  JSON="$(dotld example.com --json 2>&1 || true)"
  if echo "$JSON" | jq -e '.results[0].domain' &>/dev/null; then
    pass "dotld example.com --json has valid JSON with results[].domain"
  else
    fail "dotld example.com --json invalid: $JSON"
  fi

  if echo "$JSON" | jq -e '.results[0] | has("available", "price", "currency")' &>/dev/null; then
    pass "JSON contains expected fields (available, price, currency)"
  else
    fail "JSON missing expected fields"
  fi
fi

printf "\n=== Stage 3: Skill Install ===\n\n"

SKILL_DIR="$HOME/.claude/skills/dotld"
mkdir -p "$SKILL_DIR"
cp -r /app/skills/dotld/* "$SKILL_DIR/"

if [[ -f "$SKILL_DIR/SKILL.md" ]]; then
  pass "SKILL.md installed to $SKILL_DIR"
else
  fail "SKILL.md not found at $SKILL_DIR"
fi

FRONTMATTER="$(sed -n '/^---$/,/^---$/p' "$SKILL_DIR/SKILL.md")"

for FIELD in name allowed-tools description; do
  if echo "$FRONTMATTER" | grep -q "$FIELD"; then
    pass "SKILL.md frontmatter has '$FIELD'"
  else
    fail "SKILL.md frontmatter missing '$FIELD'"
  fi
done

printf "\n=== Results: %d passed, %d failed ===\n\n" "$PASS" "$FAIL"
[[ "$FAIL" -eq 0 ]]
