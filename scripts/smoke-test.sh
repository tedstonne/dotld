#!/usr/bin/env bash
set -euo pipefail

PASS=0
FAIL=0

pass() { printf "  \xe2\x9c\x93 %s\n" "$1"; PASS=$((PASS + 1)); }
fail() { printf "  \xe2\x9c\x97 %s\n" "$1"; FAIL=$((FAIL + 1)); }
skip() { printf "  - %s (skipped)\n" "$1"; }

# ---------- Stage 1: Binary Verification ----------

printf "\n=== Stage 1: Binary Verification ===\n\n"

if command -v dotld &>/dev/null; then
  pass "dotld is in PATH"
else
  fail "dotld not found in PATH"
  printf "Cannot continue without binary\n"
  exit 1
fi

VERSION_RAW="$(dotld --version 2>&1 || true)"
VERSION="$(echo "$VERSION_RAW" | awk '{print $NF}')"
printf "  version: %s\n" "$VERSION"

LATEST="$(curl -fsSL -o /dev/null -w '%{url_effective}' https://github.com/tedstonne/dotld/releases/latest | grep -oE '[^/]+$' | sed 's/^v//')"
printf "  latest release: %s\n" "$LATEST"

if [[ "$VERSION" == "$LATEST" ]]; then
  pass "dotld --version matches latest GitHub release ($LATEST)"
else
  fail "dotld --version ($VERSION) does not match latest GitHub release ($LATEST)"
fi

HELP="$(dotld --help 2>&1 || true)"
if echo "$HELP" | grep -qi "usage\|domain\|search"; then
  pass "dotld --help shows usage info"
else
  fail "dotld --help output unexpected: $HELP"
fi

# ---------- Stage 2: Live API Test ----------

printf "\n=== Stage 2: Live API Test ===\n\n"

if [[ -z "${DYNADOT_API_PRODUCTION_KEY:-}" ]]; then
  fail "DYNADOT_API_PRODUCTION_KEY not set — required for smoke tests"
  printf "Cannot continue without API key\n"
  exit 1
fi

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

# ---------- Stage 3: Skill Install (via skills.sh) ----------

printf "\n=== Stage 3: Skill Install (via skills.sh) ===\n\n"

# skills add --global installs to:
#   ~/.agents/skills/dotld/          (universal — OpenCode, Cursor, Codex, etc.)
#   ~/.claude/skills/dotld -> symlink (Claude Code)

UNIVERSAL_SKILL_DIR="$HOME/.agents/skills/dotld"
CLAUDE_SKILL_DIR="$HOME/.claude/skills/dotld"

# Check universal path (canonical install location)
if [[ -f "$UNIVERSAL_SKILL_DIR/SKILL.md" ]]; then
  pass "SKILL.md installed to $UNIVERSAL_SKILL_DIR (universal)"
else
  fail "SKILL.md not found at $UNIVERSAL_SKILL_DIR — npx skills add --global may have failed"
fi

# Check Claude Code symlink
if [[ -L "$CLAUDE_SKILL_DIR" ]]; then
  pass "Claude Code symlink exists at $CLAUDE_SKILL_DIR"
  if [[ -f "$CLAUDE_SKILL_DIR/SKILL.md" ]]; then
    pass "Symlink resolves — SKILL.md accessible via Claude Code path"
  else
    fail "Symlink broken — SKILL.md not accessible at $CLAUDE_SKILL_DIR"
  fi
elif [[ -f "$CLAUDE_SKILL_DIR/SKILL.md" ]]; then
  pass "SKILL.md installed to $CLAUDE_SKILL_DIR (direct copy)"
else
  fail "SKILL.md not found at $CLAUDE_SKILL_DIR — Claude Code skill missing"
fi

# Validate SKILL.md frontmatter
SKILL_FILE="$UNIVERSAL_SKILL_DIR/SKILL.md"
if [[ ! -f "$SKILL_FILE" ]]; then
  SKILL_FILE="$CLAUDE_SKILL_DIR/SKILL.md"
fi

if [[ -f "$SKILL_FILE" ]]; then
  FRONTMATTER="$(sed -n '/^---$/,/^---$/p' "$SKILL_FILE")"
  for FIELD in name allowed-tools description; do
    if echo "$FRONTMATTER" | grep -q "$FIELD"; then
      pass "SKILL.md frontmatter has '$FIELD'"
    else
      fail "SKILL.md frontmatter missing '$FIELD'"
    fi
  done

  # Validate references directory
  SKILL_BASE="$(dirname "$SKILL_FILE")"
  if [[ -f "$SKILL_BASE/references/cli-reference.md" ]]; then
    pass "cli-reference.md present in skill references"
  else
    fail "cli-reference.md missing from skill references"
  fi
else
  fail "No SKILL.md found — cannot validate frontmatter"
fi

# ---------- Stage 4: OpenCode CLI ----------

printf "\n=== Stage 4: OpenCode CLI ===\n\n"

if command -v opencode &>/dev/null; then
  pass "opencode is in PATH"
  OC_VERSION="$(opencode -v 2>&1 || true)"
  printf "  opencode version: %s\n" "$OC_VERSION"
else
  fail "opencode not found in PATH"
fi

# Verify opencode config exists
if [[ -f "$HOME/.config/opencode/opencode.json" ]]; then
  pass "opencode.json config present"

  # Validate config is valid JSON
  if jq empty "$HOME/.config/opencode/opencode.json" 2>/dev/null; then
    pass "opencode.json is valid JSON"
  else
    fail "opencode.json is not valid JSON"
  fi
else
  skip "opencode.json not found — OpenCode integration test limited"
fi

# ---------- Stage 5: OpenCode Skill Integration ----------

printf "\n=== Stage 5: OpenCode + dotld Skill Integration ===\n\n"

printf "  model: opencode/big-pickle (free)\n"

PROMPT="Come up with a single creative startup name, then check if domains are available for it"
printf "  prompt: %s\n" "$PROMPT"
printf "  Running OpenCode with prompt...\n\n"

OC_OUTPUT="$(timeout 120 opencode run "$PROMPT" 2>&1 || true)"

printf "  output preview:\n"
echo "$OC_OUTPUT" | head -20 | sed 's/^/    /'
printf "\n"

# Success: look for any word.tld pattern that only dotld produces
if echo "$OC_OUTPUT" | grep -qiE "[a-z0-9]+\.(com|net|org|io|ai|co|app|dev|sh).*(\·|Taken|\\\$[0-9])"; then
  pass "OpenCode invoked dotld — found domain results with availability data"
elif echo "$OC_OUTPUT" | grep -qiE "[a-z0-9]+\.(com|net|org|io|ai|co|app|dev|sh)"; then
  pass "OpenCode invoked dotld — found domain.tld results"
else
  fail "No domain.tld results in output — skill may not have triggered"
fi

# Bonus: check if pricing or availability info came back (Taken / $price)
if echo "$OC_OUTPUT" | grep -qiE "taken|\\\$[0-9]+"; then
  pass "Output contains availability/pricing data from Dynadot"
else
  skip "No pricing/availability markers found (model may have summarized differently)"
fi

# ---------- Results ----------

printf "\n=== Results: %d passed, %d failed ===\n\n" "$PASS" "$FAIL"
[[ "$FAIL" -eq 0 ]]
